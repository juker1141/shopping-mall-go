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

	"github.com/golang/mock/gomock"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func TestCreateRole(t *testing.T) {
	role := randomRole()

	n := 5

	permissionList, permissionsID := randomPermissionList(n)
	result := db.RoleTxResult{
		Role:           role,
		PermissionList: permissionList,
	}

	testCases := []struct {
		name          string
		body          db.CreateRoleTxParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: db.CreateRoleTxParams{
				Name:          role.Name,
				PermissionsID: permissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateRoleTxParams{
					Name:          role.Name,
					PermissionsID: permissionsID,
				}

				store.EXPECT().
					CreateRoleTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoleTx(t, recorder.Body, result)
			},
		},
		{
			name: "InternalError",
			body: db.CreateRoleTxParams{
				Name:          role.Name,
				PermissionsID: permissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateRoleTxParams{
					Name:          role.Name,
					PermissionsID: permissionsID,
				}

				store.EXPECT().
					CreateRoleTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.RoleTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidJSONName",
			body: db.CreateRoleTxParams{
				Name:          "",
				PermissionsID: permissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateRoleTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidJSONPermissionsID",
			body: db.CreateRoleTxParams{
				Name:          role.Name,
				PermissionsID: nil,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateRoleTx(gomock.Any(), gomock.Any()).
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

			url := "/admin/roles"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateRole(t *testing.T) {
	role := randomRole()
	updateRoleName := util.RandomRole()
	n := 5

	updatePermissionList, updatedPermissionsID := randomPermissionList(n)

	arg := db.UpdateRoleTxParams{
		ID:            role.ID,
		Name:          updateRoleName,
		PermissionsID: updatedPermissionsID,
	}

	result := db.RoleTxResult{
		Role:           role,
		PermissionList: updatePermissionList,
	}

	testCases := []struct {
		name          string
		body          db.UpdateRoleTxParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: db.UpdateRoleTxParams{
				ID:            role.ID,
				Name:          updateRoleName,
				PermissionsID: updatedPermissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoleTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoleTx(t, recorder.Body, result)
			},
		},
		{
			name: "OnlyUpdateRoleName",
			body: db.UpdateRoleTxParams{
				ID:   role.ID,
				Name: updateRoleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateRoleTxParams{
					ID:   role.ID,
					Name: updateRoleName,
				}

				store.EXPECT().
					UpdateRoleTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoleTx(t, recorder.Body, result)
			},
		},
		{
			name: "OnlyUpdateRolePermission",
			body: db.UpdateRoleTxParams{
				ID:            role.ID,
				PermissionsID: updatedPermissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateRoleTxParams{
					ID:            role.ID,
					PermissionsID: updatedPermissionsID,
				}

				store.EXPECT().
					UpdateRoleTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoleTx(t, recorder.Body, result)
			},
		},
		{
			name: "InternalError",
			body: db.UpdateRoleTxParams{
				ID:            role.ID,
				Name:          updateRoleName,
				PermissionsID: updatedPermissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoleTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RoleTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "ErrorID",
			body: db.UpdateRoleTxParams{
				ID:            -1,
				Name:          updateRoleName,
				PermissionsID: updatedPermissionsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoleTx(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/role/%d", tc.body.ID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomRole() db.Role {
	return db.Role{
		ID:   util.RandomInt(1, 1000),
		Name: util.RandomName(),
	}
}

func randomPermissionList(size int) ([]db.Permission, []int64) {
	var permissionList []db.Permission
	var permissionsID []int64

	for i := 0; i < size; i++ {
		permission := randomPermission()
		permissionList = append(permissionList, permission)
		permissionsID = append(permissionsID, permission.ID)
	}

	return permissionList, permissionsID
}

func requireBodyMatchRoleTx(t *testing.T, body *bytes.Buffer, roleTxResult db.RoleTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotRoleTx db.RoleTxResult
	err = json.Unmarshal(data, &gotRoleTx)
	require.NoError(t, err)
	require.Equal(t, roleTxResult, gotRoleTx)
}
