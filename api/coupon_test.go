package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

type eqCreateCouponParamsMatcher struct {
	arg db.CreateCouponParams
}

func (e eqCreateCouponParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateCouponParams)
	if !ok {
		return false
	}

	if !e.arg.StartAt.IsZero() && !e.arg.StartAt.Equal(arg.StartAt) {
		return false
	}

	if !e.arg.ExpiresAt.IsZero() && !e.arg.ExpiresAt.Equal(arg.ExpiresAt) {
		return false
	}

	e.arg.StartAt = arg.StartAt
	e.arg.ExpiresAt = arg.ExpiresAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateCouponParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqCreateCouponParams(arg db.CreateCouponParams) gomock.Matcher {
	return eqCreateCouponParamsMatcher{arg}
}

func TestCreateCoupon(t *testing.T) {
	coupon := randomCoupon()

	templateBody := gin.H{
		"title":      coupon.Title,
		"code":       coupon.Code,
		"percent":    coupon.Percent,
		"start_at":   coupon.StartAt,
		"expires_at": coupon.ExpiresAt,
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
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.CreateCouponParams{
					Title:     coupon.Title,
					Code:      coupon.Code,
					Percent:   coupon.Percent,
					CreatedBy: coupon.CreatedBy,
					StartAt:   coupon.StartAt,
					ExpiresAt: coupon.ExpiresAt,
				}

				store.EXPECT().
					CreateCoupon(gomock.Any(), EqCreateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCoupon(t, recorder.Body, coupon)
			},
		},
		{
			name: "OK_NoInputTime",
			body: gin.H{
				"title":   coupon.Title,
				"code":    coupon.Code,
				"percent": coupon.Percent,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.CreateCouponParams{
					Title:     coupon.Title,
					Code:      coupon.Code,
					Percent:   coupon.Percent,
					CreatedBy: coupon.CreatedBy,
				}

				store.EXPECT().
					CreateCoupon(gomock.Any(), EqCreateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCoupon(t, recorder.Body, coupon)
			},
		},
		{
			name: "OK_OnlyInputStartTime",
			body: gin.H{
				"title":    coupon.Title,
				"code":     coupon.Code,
				"percent":  coupon.Percent,
				"start_at": coupon.StartAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.CreateCouponParams{
					Title:     coupon.Title,
					Code:      coupon.Code,
					Percent:   coupon.Percent,
					CreatedBy: coupon.CreatedBy,
					StartAt:   coupon.StartAt,
				}

				store.EXPECT().
					CreateCoupon(gomock.Any(), EqCreateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCoupon(t, recorder.Body, coupon)
			},
		},
		{
			name: "OK_OnlyInputExpiresTime",
			body: gin.H{
				"title":      coupon.Title,
				"code":       coupon.Code,
				"percent":    coupon.Percent,
				"expires_at": coupon.ExpiresAt,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.CreateCouponParams{
					Title:     coupon.Title,
					Code:      coupon.Code,
					Percent:   coupon.Percent,
					CreatedBy: coupon.CreatedBy,
					ExpiresAt: coupon.ExpiresAt,
				}

				store.EXPECT().
					CreateCoupon(gomock.Any(), EqCreateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCoupon(t, recorder.Body, coupon)
			},
		},
		{
			name: "NoAuthorization",
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					CreateCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.CreateCouponParams{
					Title:     coupon.Title,
					Code:      coupon.Code,
					Percent:   coupon.Percent,
					CreatedBy: coupon.CreatedBy,
					StartAt:   coupon.StartAt,
					ExpiresAt: coupon.ExpiresAt,
				}

				store.EXPECT().
					CreateCoupon(gomock.Any(), EqCreateCouponParams(arg)).
					Times(1).
					Return(db.Coupon{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FieldEmpty",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					CreateCoupon(gomock.Any(), gomock.Any()).
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

			url := "/admin/coupon"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomCoupon() db.Coupon {
	return db.Coupon{
		ID:        util.RandomID(),
		Title:     util.RandomString(6),
		Code:      util.RandomString(6),
		Percent:   int32(util.RandomInt(1, 100)),
		CreatedBy: "user",
		StartAt:   time.Now(),
		ExpiresAt: time.Now().Add(time.Minute),
	}
}

func requireBodyMatchCoupon(t *testing.T, body *bytes.Buffer, coupon db.Coupon) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCoupon db.Coupon
	err = json.Unmarshal(data, &gotCoupon)
	require.NoError(t, err)
	require.Equal(t, coupon.Title, gotCoupon.Title)
	require.Equal(t, coupon.Code, gotCoupon.Code)
	require.Equal(t, coupon.Percent, gotCoupon.Percent)
	require.Equal(t, coupon.CreatedBy, gotCoupon.CreatedBy)
	require.WithinDuration(t, coupon.StartAt, gotCoupon.StartAt, time.Second)
	require.WithinDuration(t, coupon.ExpiresAt, gotCoupon.ExpiresAt, time.Second)
}
