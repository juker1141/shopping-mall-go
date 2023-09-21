package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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

type eqCreateAdminUserParamsMatcher struct {
	arg      db.CreateAdminUserParams
	password string
}

func (e eqCreateAdminUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAdminUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateAdminUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateAdminUserParams(arg db.CreateAdminUserParams, password string) gomock.Matcher {
	return eqCreateAdminUserParamsMatcher{arg, password}
}

func TestCreateAdminUserAPI(t *testing.T) {
	role := randomRole()
	adminUser, password := randomAdminUser(t, int32(1), role.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(role, nil)

				arg := db.CreateAdminUserParams{
					Account:  adminUser.Account,
					FullName: adminUser.FullName,
					Status:   adminUser.Status,
					RoleID: pgtype.Int4{
						Int32: adminUser.RoleID.Int32,
						Valid: true,
					},
				}

				store.EXPECT().
					CreateAdminUser(gomock.Any(), EqCreateAdminUserParams(arg, password)).
					Times(1).
					Return(adminUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(db.Role{}, sql.ErrConnDone)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateAccount",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(role, nil)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUser{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"account":   "invalid-user#1",
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidFullName",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": "invalid_FullName#@",
				"password":  password,
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  "psw",
				"status":    adminUser.Status,
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidStatusInput",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    "1",
				"role_id":   role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "RoleIDNotFound",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"roles_id":  role.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)

				store.EXPECT().
					CreateAdminUser(gomock.Any(), gomock.Any()).
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

			url := "/admin/manager-user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAdminUserAPI(t *testing.T) {
	role := randomRole()
	adminUser, _ := randomAdminUser(t, int32(1), role.ID)

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(int64(adminUser.RoleID.Int32))).
					Times(1).
					Return(role, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(db.AdminUser{}, db.ErrRecordNotFound)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(db.AdminUser{}, sql.ErrConnDone)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)
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
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/manager-user/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAdminUsersAPI(t *testing.T) {
	n := 5

	role := randomRole()
	adminUserList := make([]db.AdminUser, n)
	for i := 0; i < n; i++ {
		adminUser, _ := randomAdminUser(t, int32(1), role.ID)
		adminUserList[i] = adminUser
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
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.ListAdminUsersParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListAdminUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(adminUserList, nil)

				for _, adminUser := range adminUserList {
					store.EXPECT().
						GetRole(gomock.Any(), gomock.Eq(int64(adminUser.RoleID.Int32))).
						Times(1).
						Return(role, nil)
				}

				store.EXPECT().GetAdminUsersCount(gomock.Any()).Times(1).Return(int64(n), nil)
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
					ListAdminUsers(gomock.Any(), gomock.Any()).
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
					ListAdminUsers(gomock.Any(), gomock.Any()).
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
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					ListAdminUsers(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.AdminUser{}, sql.ErrConnDone)
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
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					ListAdminUsers(gomock.Any(), gomock.Any()).
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
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					ListAdminUsers(gomock.Any(), gomock.Any()).
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

			url := "/admin/manager-users"
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

type eqUpdateAdminUserParamsMatcher struct {
	arg      db.UpdateAdminUserParams
	password string
}

func (e eqUpdateAdminUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateAdminUserParams)

	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword.String)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = pgtype.Text{
		String: arg.HashedPassword.String,
		Valid:  true,
	}
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateAdminUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateAdminUserParams(arg db.UpdateAdminUserParams, password string) gomock.Matcher {
	return eqUpdateAdminUserParamsMatcher{arg, password}
}

func TestUpdateAdminUserAPI(t *testing.T) {
	role := randomRole()
	adminUser, password := randomAdminUser(t, int32(1), role.ID)
	newPassword := util.RandomString(8)
	newName := util.RandomString(8)
	newStatus := int32(0)

	invalidStatus := int32(3)

	newRole := randomRole()

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
			ID:   adminUser.ID,
			body: gin.H{
				"full_name":    newName,
				"status":       newStatus,
				"role_id":      int64(adminUser.RoleID.Int32),
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: newStatus,
						Valid: true,
					},
					RoleID: pgtype.Int4{
						Int32: adminUser.RoleID.Int32,
						Valid: true,
					},
				}

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), EqUpdateAdminUserParams(arg, newPassword)).
					Times(1).
					Return(adminUser, nil)

				fmt.Println(newRole.ID, newRole)
				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(int64(adminUser.RoleID.Int32))).
					Times(1).
					Return(role, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   adminUser.ID,
			body: gin.H{
				"full_name":    newName,
				"status":       &newStatus,
				"role_id":      newRole.ID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   adminUser.ID,
			body: gin.H{
				"full_name":    newName,
				"status":       &newStatus,
				"role_id":      newRole.ID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "OnlyUpdateFullName",
			ID:   adminUser.ID,
			body: gin.H{
				"full_name": newName,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(role, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OnlyUpdateStatus",
			ID:   adminUser.ID,
			body: gin.H{
				"status": newStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
					Status: pgtype.Int4{
						Int32: newStatus,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), EqUpdateAdminUserParams(arg, password)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(role, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OnlyUpdateRolesID",
			ID:   adminUser.ID,
			body: gin.H{
				"role_id": newRole.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
					RoleID: pgtype.Int4{
						Int32: int32(newRole.ID),
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(newRole.ID)).
					Times(1).
					Return(newRole, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OnlyUpdatePassword",
			ID:   adminUser.ID,
			body: gin.H{
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), EqUpdateAdminUserParams(arg, newPassword)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(role, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   adminUser.ID,
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.AdminUser{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NotFoundWithPassword",
			ID:   adminUser.ID,
			body: gin.H{
				"full_name":    newName,
				"status":       &newStatus,
				"role_id":      newRole.ID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(db.AdminUser{}, db.ErrRecordNotFound)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalErrror",
			ID:   adminUser.ID,
			body: gin.H{
				"full_name":    newName,
				"status":       &newStatus,
				"role_id":      role.ID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUser{}, sql.ErrConnDone)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "AllFieldEmpty",
			ID:   adminUser.ID,
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(1).
					Return(role, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidFullName",
			ID:   adminUser.ID,
			body: gin.H{
				"full_name": "Invalid$#123",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidStatus",
			ID:   adminUser.ID,
			body: gin.H{
				"status": &invalidStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidRolesIDLength",
			ID:   adminUser.ID,
			body: gin.H{
				"roles_id": []int64{},
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "WrongPassword",
			ID:   adminUser.ID,
			body: gin.H{
				"old_password": "12345678",
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					UpdateAdminUser(gomock.Any(), EqUpdateAdminUserParams(arg, newPassword)).
					Times(0)

				store.EXPECT().
					GetRole(gomock.Any(), gomock.Eq(role.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			fmt.Println(tc.body)
			jsonData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/admin/manager-user/%d", tc.ID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteAdminUserAPI(t *testing.T) {
	role := randomRole()
	adminUser, _ := randomAdminUser(t, int32(0), role.ID)

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					DeleteAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(deleteResult{
						Message: "Delete adminUser success.",
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					DeleteAdminUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					DeleteAdminUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(deleteResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   adminUser.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					DeleteAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(deleteResult{}, db.ErrRecordNotFound)
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
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					DeleteAdminUser(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/manager-user/%d", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginAdminUserAPI(t *testing.T) {
	role := randomRole()
	adminUser, password := randomAdminUser(t, int32(1), role.ID)

	n := 5
	permissionList, _ := randomPermissionList(n)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account":  adminUser.Account,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAdminUserByAccount(gomock.Any(), gomock.Eq(adminUser.Account)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					ListPermissionsForAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(permissionList, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "AccountNotFound",
			body: gin.H{
				"account":  "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAdminUserByAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUser{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"account":  adminUser.Account,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAdminUserByAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUser{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"account":  adminUser.Account,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAdminUserByAccount(gomock.Any(), gomock.Eq(adminUser.Account)).
					Times(1).
					Return(db.AdminUser{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := "/admin/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAdminUser(t *testing.T, status int32, role_id int64) (adminUser db.AdminUser, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	adminUser = db.AdminUser{
		ID:             util.RandomID(),
		Account:        util.RandomAccount(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		Status:         status,
		RoleID: pgtype.Int4{
			Int32: int32(role_id),
			Valid: true,
		},
	}
	return
}

func randomRoleList(size int) ([]db.Role, []int64) {
	var roleList []db.Role
	var rolesID []int64

	for i := 0; i < size; i++ {
		role := randomRole()
		roleList = append(roleList, role)
		rolesID = append(rolesID, role.ID)
	}

	return roleList, rolesID
}

// func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotUser db.User
// 	err = json.Unmarshal(data, &gotUser)

// 	require.NoError(t, err)
// 	require.Equal(t, user.ID, gotUser.ID)
// 	require.Equal(t, user.Account, gotUser.Account)
// 	require.Equal(t, user.Email, gotUser.Email)
// 	require.Equal(t, user.FullName, gotUser.FullName)
// 	require.Equal(t, user.AvatarUrl, gotUser.AvatarUrl)
// }
