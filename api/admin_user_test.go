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
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

type eqCreateAdminUserParamsMatcher struct {
	arg      db.CreateAdminUserTxParams
	password string
}

func (e eqCreateAdminUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAdminUserTxParams)
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

func EqCreateAdminUserParams(arg db.CreateAdminUserTxParams, password string) gomock.Matcher {
	return eqCreateAdminUserParamsMatcher{arg, password}
}

func TestCreateAdminUserAPI(t *testing.T) {
	adminUser, password := randomAdminUser(t, int32(1))

	n := 5

	roleList, rolesID := randomRoleList(n)

	result := db.AdminUserTxResult{
		AdminUser: adminUser,
		RoleList:  roleList,
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
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"roles_id":  []int64{1},
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.CreateAdminUserTxParams{
					Account:  adminUser.Account,
					FullName: adminUser.FullName,
					Status:   adminUser.Status,
					RolesID:  []int64{1},
				}

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), EqCreateAdminUserParams(arg, password)).
					Times(1).
					Return(result, nil)
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUserTxResult{}, sql.ErrConnDone)
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUserTxResult{}, db.ErrUniqueViolation)
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidRolesIDLength",
			body: gin.H{
				"account":   adminUser.Account,
				"full_name": adminUser.FullName,
				"password":  password,
				"status":    adminUser.Status,
				"roles_id":  nil,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
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
	adminUser, _ := randomAdminUser(t, int32(1))

	n := 5

	roleList, _ := randomRoleList(n)

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
					ListRolesForAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(roleList, nil)
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

	adminUserList := make([]db.AdminUser, n)
	for i := 0; i < n; i++ {
		adminUser, _ := randomAdminUser(t, int32(1))
		adminUserList[i] = adminUser
	}

	roleList, _ := randomRoleList(n)

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
						ListRolesForAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
						Times(1).
						Return(roleList, nil)
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
	arg      db.UpdateAdminUserTxParams
	password string
}

func (e eqUpdateAdminUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateAdminUserTxParams)

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

func (e eqUpdateAdminUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateAdminUserParams(arg db.UpdateAdminUserTxParams, password string) gomock.Matcher {
	return eqUpdateAdminUserParamsMatcher{arg, password}
}

func TestUpdateAdminUserAPI(t *testing.T) {
	adminUser, password := randomAdminUser(t, int32(1))
	newPassword := util.RandomString(8)
	newName := util.RandomString(8)
	newStatus := int32(0)

	invalidStatus := int32(3)

	n := 5
	roleList, rolesID := randomRoleList(n)

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
				"status":       &newStatus,
				"roles_id":     rolesID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserTxParams{
					ID:       adminUser.ID,
					FullName: newName,
					Status:   &newStatus,
					RolesID:  rolesID,
				}

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), EqUpdateAdminUserParams(arg, newPassword)).
					Times(1).
					Return(db.AdminUserTxResult{
						AdminUser: adminUser,
						RoleList:  roleList,
					}, nil)
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
				"roles_id":     rolesID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":     rolesID,
				"old_password": password,
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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

				arg := db.UpdateAdminUserTxParams{
					ID:       adminUser.ID,
					FullName: newName,
				}

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.AdminUserTxResult{
						AdminUser: adminUser,
						RoleList:  roleList,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OnlyUpdateStatus",
			ID:   adminUser.ID,
			body: gin.H{
				"status": &newStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserTxParams{
					ID:     adminUser.ID,
					Status: &newStatus,
				}

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.AdminUserTxResult{
						AdminUser: adminUser,
						RoleList:  roleList,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OnlyUpdateRolesID",
			ID:   adminUser.ID,
			body: gin.H{
				"roles_id": rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", accountPermissions)

				arg := db.UpdateAdminUserTxParams{
					ID:      adminUser.ID,
					RolesID: rolesID,
				}

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.AdminUserTxResult{
						AdminUser: adminUser,
						RoleList:  roleList,
					}, nil)
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

				arg := db.UpdateAdminUserTxParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), EqUpdateAdminUserParams(arg, newPassword)).
					Times(1).
					Return(db.AdminUserTxResult{
						AdminUser: adminUser,
						RoleList:  roleList,
					}, nil)
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

				arg := db.UpdateAdminUserTxParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.AdminUserTxResult{}, db.ErrRecordNotFound)
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
				"roles_id":     rolesID,
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
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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
				"roles_id":     rolesID,
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
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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

				arg := db.UpdateAdminUserTxParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.AdminUserTxResult{
						AdminUser: adminUser,
						RoleList:  roleList,
					}, nil)
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
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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
					UpdateAdminUserTx(gomock.Any(), gomock.Any()).
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

				arg := db.UpdateAdminUserTxParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					GetAdminUser(gomock.Any(), gomock.Eq(adminUser.ID)).
					Times(1).
					Return(adminUser, nil)

				store.EXPECT().
					UpdateAdminUserTx(gomock.Any(), EqUpdateAdminUserParams(arg, newPassword)).
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
	adminUser, _ := randomAdminUser(t, int32(0))

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

				arg := db.DeleteAdminUserTxParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					DeleteAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.DeleteAdminUserTxResult{
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
					DeleteAdminUserTx(gomock.Any(), gomock.Any()).
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
					DeleteAdminUserTx(gomock.Any(), gomock.Any()).
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
					DeleteAdminUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.DeleteAdminUserTxResult{}, sql.ErrConnDone)
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

				arg := db.DeleteAdminUserTxParams{
					ID: adminUser.ID,
				}

				store.EXPECT().
					DeleteAdminUserTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.DeleteAdminUserTxResult{}, db.ErrRecordNotFound)
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
					DeleteAdminUserTx(gomock.Any(), gomock.Any()).
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
	adminUser, password := randomAdminUser(t, int32(1))

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

func randomAdminUser(t *testing.T, status int32) (adminUser db.AdminUser, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	adminUser = db.AdminUser{
		ID:             util.RandomID(),
		Account:        util.RandomAccount(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		Status:         status,
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
