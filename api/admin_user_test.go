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

func TestCreateAdminUser(t *testing.T) {
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
				"roles_id":  rolesID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAdminUserTxParams{
					Account:  adminUser.Account,
					FullName: adminUser.FullName,
					Status:   adminUser.Status,
					RolesID:  rolesID,
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
			name: "DuplicateUsername",
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
				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AdminUserTxResult{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			url := "/admin/user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// func TestGetAdminUser(t *testing.T) {
// 	adminUser, password := randomAdminUser(t)

// 	n := 5

// 	roleList, rolesID := randomRoleList(n)

// 	result := db.AdminUserTxResult{
// 		AdminUser: adminUser,
// 		RoleList:  roleList,
// 	}

// 	testCases := []struct {
// 		name          string
// 		ID            int64
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			ID:   adminUser.ID,
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.CreateAdminUserTxParams{
// 					Account:  adminUser.Account,
// 					FullName: adminUser.FullName,
// 					Status:   adminUser.Status,
// 					RolesID:  rolesID,
// 				}

// 				store.EXPECT().
// 					GetAdminUser(gomock.Any(), EqCreateAdminUserParams(arg, password)).
// 					Times(1).
// 					Return(result, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoAuthorization",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": adminUser.FullName,
// 				"password":  password,
// 				"status":    adminUser.Status,
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": adminUser.FullName,
// 				"password":  password,
// 				"status":    adminUser.Status,
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.AdminUserTxResult{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "DuplicateUsername",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": adminUser.FullName,
// 				"password":  password,
// 				"status":    adminUser.Status,
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.AdminUserTxResult{}, db.ErrUniqueViolation)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusForbidden, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidUsername",
// 			body: gin.H{
// 				"account":   "invalid-user#1",
// 				"full_name": adminUser.FullName,
// 				"password":  password,
// 				"status":    adminUser.Status,
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidFullName",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": "invalid_FullName#@",
// 				"password":  password,
// 				"status":    adminUser.Status,
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "TooShortPassword",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": adminUser.FullName,
// 				"password":  "psw",
// 				"status":    adminUser.Status,
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidStatusInput",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": adminUser.FullName,
// 				"password":  password,
// 				"status":    "1",
// 				"roles_id":  rolesID,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidRolesIDLength",
// 			body: gin.H{
// 				"account":   adminUser.Account,
// 				"full_name": adminUser.FullName,
// 				"password":  password,
// 				"status":    adminUser.Status,
// 				"roles_id":  nil,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAdminUserTx(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			jsonData, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/admin/user"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.tokenMaker)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

func TestListAdminUser(t *testing.T) {
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
	fmt.Println(e, x, arg, "arg")

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

func TestUpdateAdminUser(t *testing.T) {
}

func TestDeleteAdminUser(t *testing.T) {
}

func TestLoginAdminUser(t *testing.T) {
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
		ID:             util.RandomInt(1, 100),
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
