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
	"github.com/jackc/pgx/v5/pgtype"
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

func TestCreateCouponAPI(t *testing.T) {
	coupon := randomCoupon()

	centerTime := time.Now()

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
		{
			name: "CheckDateIsFailed",
			body: gin.H{
				"title":      coupon.Title,
				"code":       coupon.Code,
				"percent":    coupon.Percent,
				"start_at":   centerTime,
				"expires_at": centerTime.Add(-5 * time.Minute),
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
					StartAt:   centerTime,
					ExpiresAt: centerTime.Add(-5 * time.Minute),
				}

				store.EXPECT().
					CreateCoupon(gomock.Any(), EqCreateCouponParams(arg)).
					Times(1).
					Return(db.Coupon{}, db.ErrCheckDateFailed)
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

type eqUpdateCouponParamsMatcher struct {
	arg db.UpdateCouponParams
}

func (e eqUpdateCouponParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateCouponParams)
	if !ok {
		return false
	}

	if !e.arg.StartAt.Time.IsZero() && !e.arg.StartAt.Time.Equal(arg.StartAt.Time) {
		return false
	}

	if !e.arg.ExpiresAt.Time.IsZero() && !e.arg.ExpiresAt.Time.Equal(arg.ExpiresAt.Time) {
		return false
	}

	e.arg.StartAt = arg.StartAt
	e.arg.ExpiresAt = arg.ExpiresAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateCouponParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqUpdateCouponParams(arg db.UpdateCouponParams) gomock.Matcher {
	return eqUpdateCouponParamsMatcher{arg}
}

func TestUpdateCouponAPI(t *testing.T) {
	coupon := randomCoupon()

	centerTime := time.Now()

	templateBody := gin.H{
		"title":      coupon.Title,
		"code":       coupon.Code,
		"percent":    coupon.Percent,
		"start_at":   coupon.StartAt,
		"expires_at": coupon.ExpiresAt,
	}

	testCases := []struct {
		name          string
		ID            int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   coupon.ID,
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
					Title: pgtype.Text{
						String: coupon.Title,
						Valid:  true,
					},
					Code: pgtype.Text{
						String: coupon.Code,
						Valid:  true,
					},
					Percent: pgtype.Int4{
						Int32: coupon.Percent,
						Valid: true,
					},
					StartAt: pgtype.Timestamptz{
						Time:  coupon.StartAt,
						Valid: true,
					},
					ExpiresAt: pgtype.Timestamptz{
						Time:  coupon.ExpiresAt,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), EqUpdateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_NoInputTime",
			ID:   coupon.ID,
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

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
					Title: pgtype.Text{
						String: coupon.Title,
						Valid:  true,
					},
					Code: pgtype.Text{
						String: coupon.Code,
						Valid:  true,
					},
					Percent: pgtype.Int4{
						Int32: coupon.Percent,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), EqUpdateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_OnlyInputStartTime",
			ID:   coupon.ID,
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

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
					Title: pgtype.Text{
						String: coupon.Title,
						Valid:  true,
					},
					Code: pgtype.Text{
						String: coupon.Code,
						Valid:  true,
					},
					Percent: pgtype.Int4{
						Int32: coupon.Percent,
						Valid: true,
					},
					StartAt: pgtype.Timestamptz{
						Time:  coupon.StartAt,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), EqUpdateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK_OnlyInputExpiresTime",
			ID:   coupon.ID,
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

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
					Title: pgtype.Text{
						String: coupon.Title,
						Valid:  true,
					},
					Code: pgtype.Text{
						String: coupon.Code,
						Valid:  true,
					},
					Percent: pgtype.Int4{
						Int32: coupon.Percent,
						Valid: true,
					},
					ExpiresAt: pgtype.Timestamptz{
						Time:  coupon.ExpiresAt,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), EqUpdateCouponParams(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   coupon.ID,
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   coupon.ID,
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					UpdateCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   coupon.ID,
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
					Title: pgtype.Text{
						String: coupon.Title,
						Valid:  true,
					},
					Code: pgtype.Text{
						String: coupon.Code,
						Valid:  true,
					},
					Percent: pgtype.Int4{
						Int32: coupon.Percent,
						Valid: true,
					},
					StartAt: pgtype.Timestamptz{
						Time:  coupon.StartAt,
						Valid: true,
					},
					ExpiresAt: pgtype.Timestamptz{
						Time:  coupon.ExpiresAt,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), EqUpdateCouponParams(arg)).
					Times(1).
					Return(db.Coupon{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FieldEmpty",
			ID:   coupon.ID,
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   coupon.ID,
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Coupon{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   -1,
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					UpdateCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "CheckDateIsFailed",
			ID:   coupon.ID,
			body: gin.H{
				"start_at":   centerTime,
				"expires_at": centerTime.Add(-5 * time.Minute),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.UpdateCouponParams{
					ID: coupon.ID,
					StartAt: pgtype.Timestamptz{
						Time:  centerTime,
						Valid: true,
					},
					ExpiresAt: pgtype.Timestamptz{
						Time:  centerTime.Add(-5 * time.Minute),
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateCoupon(gomock.Any(), EqUpdateCouponParams(arg)).
					Times(1).
					Return(db.Coupon{}, db.ErrCheckDateFailed)
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

			url := fmt.Sprintf("/admin/coupon/%d", tc.ID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetCouponAPI(t *testing.T) {
	coupon := randomCoupon()

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCoupon(gomock.Any(), gomock.Eq(coupon.ID)).
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
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCoupon(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(1).
					Return(db.Coupon{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCoupon(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(1).
					Return(db.Coupon{}, db.ErrRecordNotFound)
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
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCoupon(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/coupon/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCheckCouponAPI(t *testing.T) {
	coupon := randomCoupon()

	notStartTime := time.Now().Add(5 * time.Minute)

	notStartCoupon := db.Coupon{
		ID:        util.RandomID(),
		Title:     util.RandomString(6),
		Code:      util.RandomString(6),
		Percent:   int32(util.RandomInt(1, 100)),
		CreatedBy: "user",
		StartAt:   notStartTime,
		ExpiresAt: notStartTime.Add(time.Minute),
	}

	expiredStartTime := time.Now().Add(-5 * time.Minute)

	exoiredCoupon := db.Coupon{
		ID:        util.RandomID(),
		Title:     util.RandomString(6),
		Code:      util.RandomString(6),
		Percent:   int32(util.RandomInt(1, 100)),
		CreatedBy: "user",
		StartAt:   expiredStartTime,
		ExpiresAt: expiredStartTime.Add(time.Minute),
	}

	templateBody := gin.H{
		"code": coupon.Code,
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

				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Eq(coupon.Code)).
					Times(1).
					Return(coupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				// requireBodyMatchCoupon(t, recorder.Body, coupon)
			},
		},
		{
			name: "NoAuthorization",
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Eq(coupon.Code)).
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
					GetCouponByCode(gomock.Any(), gomock.Eq(coupon.Code)).
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

				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Eq(coupon.Code)).
					Times(1).
					Return(db.Coupon{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Eq(coupon.Code)).
					Times(1).
					Return(db.Coupon{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "EmptyCode",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotStartCoupon",
			body: gin.H{
				"code": notStartCoupon.Code,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Eq(notStartCoupon.Code)).
					Times(1).Return(notStartCoupon, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ExpiredCoupon",
			body: gin.H{
				"code": exoiredCoupon.Code,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					GetCouponByCode(gomock.Any(), gomock.Eq(exoiredCoupon.Code)).
					Times(1).Return(exoiredCoupon, nil)
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

			url := "/admin/coupon/check"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListCouponsAPI(t *testing.T) {
	n := 5

	couponList := make([]db.Coupon, n)
	for i := 0; i < n; i++ {
		coupon := randomCoupon()
		couponList = append(couponList, coupon)
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
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.ListCouponsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListCoupons(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(couponList, nil)

				store.EXPECT().
					GetCouponsCount(gomock.Any()).
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
					ListCoupons(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetCouponsCount(gomock.Any()).
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
					ListCoupons(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetCouponsCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", couponPermissions)

				arg := db.ListCouponsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListCoupons(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Coupon{}, sql.ErrConnDone)

				store.EXPECT().
					GetCouponsCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					ListCoupons(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetCouponsCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					ListCoupons(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetCouponsCount(gomock.Any()).
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

			url := "/admin/coupons"
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

func TestDeleteCouponAPI(t *testing.T) {
	coupon := randomCoupon()

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					DeleteCoupon(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					DeleteCoupon(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					DeleteCoupon(gomock.Any(), gomock.Eq(coupon.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   coupon.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					DeleteCoupon(gomock.Any(), gomock.Eq(coupon.ID)).
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
				addPermissionMiddleware(store, "user", couponPermissions)

				store.EXPECT().
					DeleteCoupon(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/coupon/%d", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
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
	require.NotEmpty(t, gotCoupon.StartAt)
	require.NotEmpty(t, gotCoupon.ExpiresAt)
}
