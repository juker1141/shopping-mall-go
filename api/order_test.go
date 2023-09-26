package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderAPI(t *testing.T) {
	user, _ := randomUser(t, 1)
	paymethod := randomPayMethod()
	status := defaultOrderStatus()
	coupon := randomCoupon()
	totalPrice := int32(0)
	finalPrice := int32(0)

	n := 3
	productlist := make([]db.OrderTxProductResult, n)
	orderProducts := make([]db.OrderTxProductParams, n)
	invalidOrderProducts := make([]db.OrderTxProductParams, n)
	for i := 0; i < n; i++ {
		product := randomProduct()
		num := util.RandomInt(1, 10)

		productlist[i] = db.OrderTxProductResult{
			Product: product,
			Num:     num,
		}

		orderProducts[i] = db.OrderTxProductParams{
			ID:  product.ID,
			Num: num,
		}
		invalidOrderProducts[i] = db.OrderTxProductParams{
			ID:  product.ID,
			Num: 0,
		}
		totalPrice = totalPrice + int32(num*int64(product.OriginPrice))
		finalPrice = finalPrice + int32(num*int64(product.Price))
	}

	order := randomOrder(paymethod.ID, status.ID, totalPrice, finalPrice)

	result := db.OrderTxResult{
		Order:       order,
		Status:      status,
		ProductList: productlist,
	}

	message := util.RandomString(20)

	templateBody := gin.H{
		"full_name":        order.FullName,
		"email":            order.Email,
		"shipping_address": order.ShippingAddress,
		"message":          message,
		"pay_method_id":    paymethod.ID,
		"order_products":   orderProducts,
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq("user")).
					Times(1).
					Return(user, nil)

				arg := db.CreateOrderTxParams{
					UserID:          user.ID,
					FullName:        order.FullName,
					Email:           order.Email,
					ShippingAddress: order.ShippingAddress,
					Message:         message,
					PayMethodID:     paymethod.ID,
					StatusID:        status.ID,
					OrderProducts:   orderProducts,
				}

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderTx(t, recorder.Body, result)
			},
		},
		{
			name: "OKWithCoupon",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"coupon_id":        coupon.ID,
				"message":          message,
				"pay_method_id":    paymethod.ID,
				"order_products":   orderProducts,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq("user")).
					Times(1).
					Return(user, nil)

				arg := db.CreateOrderTxParams{
					UserID:          user.ID,
					FullName:        order.FullName,
					Email:           order.Email,
					ShippingAddress: order.ShippingAddress,
					Message:         message,
					PayMethodID:     paymethod.ID,
					StatusID:        status.ID,
					OrderProducts:   orderProducts,
					CouponID:        coupon.ID,
				}

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderTx(t, recorder.Body, result)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"coupon_id":        coupon.ID,
				"message":          message,
				"pay_method_id":    paymethod.ID,
				"order_products":   orderProducts,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"coupon_id":        coupon.ID,
				"message":          message,
				"pay_method_id":    paymethod.ID,
				"order_products":   orderProducts,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"message":          message,
				"pay_method_id":    paymethod.ID,
				"order_products":   invalidOrderProducts,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq("user")).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidProductNum",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"message":          message,
				"pay_method_id":    paymethod.ID,
				"order_products":   invalidOrderProducts,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq("user")).
					Times(1).
					Return(user, nil)

				arg := db.CreateOrderTxParams{
					UserID:          user.ID,
					FullName:        order.FullName,
					Email:           order.Email,
					ShippingAddress: order.ShippingAddress,
					Message:         message,
					PayMethodID:     paymethod.ID,
					StatusID:        status.ID,
					OrderProducts:   invalidOrderProducts,
				}

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.OrderTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPayMethodID",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"message":          message,
				"pay_method_id":    0,
				"order_products":   invalidOrderProducts,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq("user")).
					Times(0)

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidOrderProductLength",
			body: gin.H{
				"full_name":        order.FullName,
				"email":            order.Email,
				"shipping_address": order.ShippingAddress,
				"message":          message,
				"pay_method_id":    paymethod.ID,
				"order_products":   []db.OrderProduct{},
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateOrderTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			jsonData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/admin/order"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListOrdersAPI(t *testing.T) {
	paymethod := randomPayMethod()
	status := defaultOrderStatus()
	totalPrice := int32(0)
	finalPrice := int32(0)

	n := 5

	productList := make([]db.Product, n)
	for i := range productList {
		product := randomProduct()
		productList[i] = product
	}

	orderList := make([]db.Order, n)
	var orderProducts []db.OrderProduct
	for i := range orderList {
		num := util.RandomInt(1, 10)
		totalPrice = totalPrice + int32(num*int64(productList[i].OriginPrice))
		finalPrice = finalPrice + int32(num*int64(productList[i].Price))

		order := randomOrder(paymethod.ID, status.ID, totalPrice, finalPrice)
		orderList[i] = order

		for productIndex := range productList {
			orderProducts = append(orderProducts, db.OrderProduct{
				OrderID: pgtype.Int4{
					Int32: int32(order.ID),
					Valid: true,
				},
				ProductID: pgtype.Int4{
					Int32: int32(productList[productIndex].ID),
					Valid: true,
				},
				Num: int32(num),
			})
		}
	}

	type Query struct {
		page     int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				page:     1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				arg := db.ListOrdersParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListOrders(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(orderList, nil)

				for _, order := range orderList {
					store.EXPECT().
						GetOrderStatus(gomock.Any(), gomock.Eq(int64(order.StatusID))).
						Times(1).
						Return(status, nil)

					store.EXPECT().
						ListOrderProductByOrderId(gomock.Any(), gomock.Eq(pgtype.Int4{
							Int32: int32(order.ID),
							Valid: true,
						})).
						Times(1).
						Return(
							getOrderProductByOrderId(order.ID, orderProducts),
							nil,
						)

					for i, orderProduct := range getOrderProductByOrderId(order.ID, orderProducts) {
						store.EXPECT().
							GetProduct(gomock.Any(), gomock.Eq(int64(orderProduct.ProductID.Int32))).
							Times(1).
							Return(productList[i], nil)
					}
				}

				store.EXPECT().GetOrdersCount(gomock.Any()).Times(1).Return(int64(n), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				page:     1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListOrders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			query: Query{
				page:     1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					ListOrders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				page:     1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					ListOrders(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Order{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPage",
			query: Query{
				page:     -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					ListOrders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				page:     1,
				pageSize: 10000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", orderPermissions)

				store.EXPECT().
					ListOrders(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/admin/orders"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.page))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomPayMethod() db.PayMethod {
	return db.PayMethod{
		ID:   util.RandomID(),
		Name: util.RandomName(),
	}
}

func defaultOrderStatus() db.OrderStatus {
	return db.OrderStatus{
		ID:          1,
		Name:        util.RandomName(),
		Description: util.RandomString(20),
	}
}

func randomOrder(payMethodID, orderStatus int64, totalPrice, finalPrice int32) db.Order {
	return db.Order{
		ID:              util.RandomID(),
		FullName:        util.RandomName(),
		Email:           util.RandomEmail(),
		ShippingAddress: util.RandomAddress(),
		IsPaid:          false,
		TotalPrice:      util.RandomPrice(),
		FinalPrice:      util.RandomPrice(),
		PayMethodID:     int32(payMethodID),
		StatusID:        int32(orderStatus),
	}
}

func getOrderProductByOrderId(orderID int64, orderProducts []db.OrderProduct) []db.OrderProduct {
	var result []db.OrderProduct
	for _, orderProduct := range orderProducts {
		if orderProduct.OrderID.Int32 == int32(orderID) {
			result = append(result, orderProduct)
		}
	}
	return result
}

func requireBodyMatchOrderTx(t *testing.T, body *bytes.Buffer, result db.OrderTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.OrderTxResult
	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)
	require.Equal(t, gotResult, result)
}
