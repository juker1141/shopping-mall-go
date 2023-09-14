package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
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

type eqCreateProductParamsMatcher struct {
	arg db.CreateProductParams
}

func (e eqCreateProductParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateProductParams)
	if !ok {
		return false
	}

	arg.ImageUrl = e.arg.ImageUrl

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateProductParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqCreateProductParams(arg db.CreateProductParams) gomock.Matcher {
	return eqCreateProductParamsMatcher{arg}
}

func TestCreateProductAPI(t *testing.T) {
	product := randomProduct()

	templateBody := gin.H{
		"title":        product.Title,
		"category":     product.Category,
		"origin_price": product.OriginPrice,
		"price":        product.Price,
		"unit":         product.Unit,
		"description":  product.Description,
		"content":      product.Content,
		"status":       1,
		"image_url":    product.ImageUrl,
	}

	testCases := []struct {
		name          string
		isUploadImage bool
		fileType      string
		fileName      string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.CreateProductParams{
					Title:       product.Title,
					Category:    product.Category,
					Description: product.Description,
					Content:     product.Content,
					Price:       product.Price,
					OriginPrice: product.OriginPrice,
					Unit:        product.Unit,
					Status:      product.Status,
					ImageUrl:    product.ImageUrl,
					ImagesUrl:   []string{},
					CreatedBy:   "user",
				}

				store.EXPECT().
					CreateProduct(gomock.Any(), EqCreateProductParams(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProduct(t, recorder.Body, product)
			},
		},
		{
			name:          "OK_ButNoUploadImage",
			isUploadImage: false,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.CreateProductParams{
					Title:       product.Title,
					Category:    product.Category,
					Description: product.Description,
					Content:     product.Content,
					Price:       product.Price,
					OriginPrice: product.OriginPrice,
					Unit:        product.Unit,
					Status:      product.Status,
					ImageUrl:    product.ImageUrl,
					ImagesUrl:   []string{},
					CreatedBy:   "user",
				}

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProduct(t, recorder.Body, product)
			},
		},
		{
			name:          "NoAuthorization",
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:          "NoRequiredPermission",
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:          "InternalError",
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.CreateProductParams{
					Title:       product.Title,
					Category:    product.Category,
					Description: product.Description,
					Content:     product.Content,
					Price:       product.Price,
					OriginPrice: product.OriginPrice,
					Unit:        product.Unit,
					Status:      product.Status,
					ImageUrl:    product.ImageUrl,
					ImagesUrl:   []string{},
					CreatedBy:   "user",
				}

				store.EXPECT().
					CreateProduct(gomock.Any(), EqCreateProductParams(arg)).
					Times(1).
					Return(db.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "TooShortTitle",
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body: gin.H{
				"title":        "aa",
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "InvalidImageType",
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.gif",
			body: gin.H{
				"title":        product.Title,
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "InvalidImageContentType",
			isUploadImage: true,
			fileType:      "other",
			fileName:      "fake_product.pdf",
			body: gin.H{
				"title":        product.Title,
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
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

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range tc.body {
				fieldWriter, err := writer.CreateFormField(key)
				require.NoError(t, err)
				fieldValue := fmt.Sprintf("%v", value)
				_, err = fieldWriter.Write([]byte(fieldValue))
				require.NoError(t, err)
			}

			if tc.isUploadImage {
				var fakeFileContent []byte
				var fileWriter io.Writer
				var err error

				if tc.fileType == "image" {
					fileWriter, err = CreateImageFormFile(writer, "image_file", tc.fileName)
				} else {
					fileWriter, err = writer.CreateFormFile("image_file", tc.fileName)
				}
				require.NoError(t, err)

				fakeFileContent = append(fakeFileContent, []byte("Fake image data...")...)
				_, err = fileWriter.Write(fakeFileContent)
				require.NoError(t, err)
			}

			// 結束FormData
			writer.Close()

			url := "/admin/product"
			request, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

type eqUpdateProductParamsMatcher struct {
	arg db.UpdateProductParams
}

func (e eqUpdateProductParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateProductParams)
	if !ok {
		return false
	}

	arg.ImageUrl = e.arg.ImageUrl

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateProductParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqUpdateProductParams(arg db.UpdateProductParams) gomock.Matcher {
	return eqUpdateProductParamsMatcher{arg}
}

func TestUpdateProductAPI(t *testing.T) {
	product := randomProduct()

	templateBody := gin.H{
		"title":        product.Title,
		"category":     product.Category,
		"origin_price": product.OriginPrice,
		"price":        product.Price,
		"unit":         product.Unit,
		"description":  product.Description,
		"content":      product.Content,
		"status":       1,
		"image_url":    product.ImageUrl,
	}

	invalidStatus := int32(3)
	images_url := []string{
		product.ImageUrl,
		product.ImageUrl,
	}

	testCases := []struct {
		name          string
		ID            int64
		isUploadImage bool
		fileType      string
		fileName      string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.UpdateProductParams{
					ID: product.ID,
					Title: pgtype.Text{
						String: product.Title,
						Valid:  true,
					},
					Category: pgtype.Text{
						String: product.Category,
						Valid:  true,
					},
					Description: pgtype.Text{
						String: product.Description,
						Valid:  true,
					},
					Content: pgtype.Text{
						String: product.Content,
						Valid:  true,
					},
					Price: pgtype.Int4{
						Int32: product.Price,
						Valid: true,
					},
					OriginPrice: pgtype.Int4{
						Int32: product.OriginPrice,
						Valid: true,
					},
					Unit: pgtype.Text{
						String: product.Unit,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: product.Status,
						Valid: true,
					},
					ImageUrl: pgtype.Text{
						String: product.ImageUrl,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateProduct(gomock.Any(), EqUpdateProductParams(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "OK_NoUploadImageFile",
			ID:            product.ID,
			isUploadImage: false,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.UpdateProductParams{
					ID: product.ID,
					Title: pgtype.Text{
						String: product.Title,
						Valid:  true,
					},
					Category: pgtype.Text{
						String: product.Category,
						Valid:  true,
					},
					Description: pgtype.Text{
						String: product.Description,
						Valid:  true,
					},
					Content: pgtype.Text{
						String: product.Content,
						Valid:  true,
					},
					Price: pgtype.Int4{
						Int32: product.Price,
						Valid: true,
					},
					OriginPrice: pgtype.Int4{
						Int32: product.OriginPrice,
						Valid: true,
					},
					Unit: pgtype.Text{
						String: product.Unit,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: product.Status,
						Valid: true,
					},
					ImageUrl: pgtype.Text{
						String: product.ImageUrl,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "OK_WithImagesUrl",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body: gin.H{
				"title":        product.Title,
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
				"images_url":   images_url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.UpdateProductParams{
					ID: product.ID,
					Title: pgtype.Text{
						String: product.Title,
						Valid:  true,
					},
					Category: pgtype.Text{
						String: product.Category,
						Valid:  true,
					},
					Description: pgtype.Text{
						String: product.Description,
						Valid:  true,
					},
					Content: pgtype.Text{
						String: product.Content,
						Valid:  true,
					},
					Price: pgtype.Int4{
						Int32: product.Price,
						Valid: true,
					},
					OriginPrice: pgtype.Int4{
						Int32: product.OriginPrice,
						Valid: true,
					},
					Unit: pgtype.Text{
						String: product.Unit,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: product.Status,
						Valid: true,
					},
					ImageUrl: pgtype.Text{
						String: product.ImageUrl,
						Valid:  true,
					},
					ImagesUrl: images_url,
				}

				store.EXPECT().
					UpdateProduct(gomock.Any(), EqUpdateProductParams(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "NoAuthorization",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:          "NoRequiredPermission",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:          "InternalError",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body:          templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.UpdateProductParams{
					ID: product.ID,
					Title: pgtype.Text{
						String: product.Title,
						Valid:  true,
					},
					Category: pgtype.Text{
						String: product.Category,
						Valid:  true,
					},
					Description: pgtype.Text{
						String: product.Description,
						Valid:  true,
					},
					Content: pgtype.Text{
						String: product.Content,
						Valid:  true,
					},
					Price: pgtype.Int4{
						Int32: product.Price,
						Valid: true,
					},
					OriginPrice: pgtype.Int4{
						Int32: product.OriginPrice,
						Valid: true,
					},
					Unit: pgtype.Text{
						String: product.Unit,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: product.Status,
						Valid: true,
					},
					ImageUrl: pgtype.Text{
						String: product.ImageUrl,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateProduct(gomock.Any(), EqUpdateProductParams(arg)).
					Times(1).
					Return(db.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "InvalidTitle",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body: gin.H{
				"title": "aa",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:          "InvalidStatus",
			ID:            product.ID,
			isUploadImage: true,
			fileType:      "image",
			fileName:      "fake_product.png",
			body: gin.H{
				"status": &invalidStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
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

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range tc.body {
				switch key {
				case "images_url":
					if values, ok := value.([]string); ok {
						for _, v := range values {
							fieldWriter, err := writer.CreateFormField("images_url")
							require.NoError(t, err)
							_, err = fieldWriter.Write([]byte(v))
							require.NoError(t, err)
						}
					}
				default:
					fieldWriter, err := writer.CreateFormField(key)
					require.NoError(t, err)
					fieldValue := fmt.Sprintf("%v", value)
					_, err = fieldWriter.Write([]byte(fieldValue))
					require.NoError(t, err)
				}
			}

			if tc.isUploadImage {
				var fakeFileContent []byte
				var fileWriter io.Writer
				var err error

				if tc.fileType == "image" {
					fileWriter, err = CreateImageFormFile(writer, "image_file", tc.fileName)
				} else {
					fileWriter, err = writer.CreateFormFile("image_file", tc.fileName)
				}
				require.NoError(t, err)

				fakeFileContent = append(fakeFileContent, []byte("Fake image data...")...)
				_, err = fileWriter.Write(fakeFileContent)
				require.NoError(t, err)
			}

			// 結束FormData
			writer.Close()

			url := fmt.Sprintf("/admin/product/%d", tc.ID)
			request, err := http.NewRequest(http.MethodPatch, url, body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListProductsAPI(t *testing.T) {
	n := 5

	productList := make([]db.Product, n)
	for i := 0; i < n; i++ {
		product := randomProduct()
		productList = append(productList, product)
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
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.ListProductsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(productList, nil)

				store.EXPECT().
					GetProductsCount(gomock.Any()).
					Times(1).
					Return(int64(n), nil)
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
					ListProducts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetProductsCount(gomock.Any()).
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
					ListProducts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetProductsCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.ListProductsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Product{}, sql.ErrConnDone)

				store.EXPECT().
					GetProductsCount(gomock.Any()).
					Times(0)
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
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetProductsCount(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPage",
			query: Query{
				page:     1,
				pageSize: 10000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetProductsCount(gomock.Any()).
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

			url := "/admin/products"
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

func TestGetProductAPI(t *testing.T) {
	product := randomProduct()

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(db.Product{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(db.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/product/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteProductAPI(t *testing.T) {
	product := randomProduct()

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   product.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/product/%d", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomProduct() db.Product {
	targetPath := filepath.Join("static", "products", "fake_product.png")

	return db.Product{
		ID:          util.RandomID(),
		Title:       util.RandomName(),
		Category:    util.RandomName(),
		OriginPrice: util.RandomPrice(),
		Price:       util.RandomPrice(),
		Unit:        util.RandomName(),
		Description: util.RandomString(20),
		Content:     util.RandomString(20),
		Status:      1,
		ImageUrl:    targetPath,
		CreatedBy:   "user",
	}
}

func requireBodyMatchProduct(t *testing.T, body *bytes.Buffer, product db.Product) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProduct db.Product
	err = json.Unmarshal(data, &gotProduct)

	fmt.Printf("%+v", gotProduct)
	require.NoError(t, err)
	require.Equal(t, product.ID, gotProduct.ID)
	require.Equal(t, product.Title, gotProduct.Title)
	require.Equal(t, product.Category, gotProduct.Category)
	require.Equal(t, product.OriginPrice, gotProduct.OriginPrice)
	require.Equal(t, product.Price, gotProduct.Price)
	require.Equal(t, product.Unit, gotProduct.Unit)
	require.Equal(t, product.Description, gotProduct.Description)
	require.Equal(t, product.Content, gotProduct.Content)
	require.Equal(t, product.Status, gotProduct.Status)
	require.Equal(t, product.ImageUrl, gotProduct.ImageUrl)
	require.Equal(t, product.ImagesUrl, gotProduct.ImagesUrl)
	require.Equal(t, product.CreatedBy, gotProduct.CreatedBy)
}
