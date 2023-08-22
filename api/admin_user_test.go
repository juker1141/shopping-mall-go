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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateAdminUserTxParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
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

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateAdminUserTxParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateAdminUser(t *testing.T) {
	adminUser, password := randomAdminUser(t)

	n := 5

	roleList, rolesID := randomRoleList(n)

	result := db.AdminUserTxResult{
		AdminUser: adminUser,
		RoleList:  roleList,
	}

	testCases := []struct {
		name          string
		body          gin.H
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
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAdminUserTxParams{
					Account:  adminUser.Account,
					FullName: adminUser.FullName,
					Status:   adminUser.Status,
					RolesID:  rolesID,
				}

				store.EXPECT().
					CreateAdminUserTx(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			jsonData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/admin/admin_user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAdminUser(t *testing.T) (adminUser db.AdminUser, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	adminUser = db.AdminUser{
		Account:        util.RandomAccount(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		Status:         1,
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
