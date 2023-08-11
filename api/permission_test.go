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
	"github.com/jackc/pgx/v5"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func TestGetPermissionAPI(t *testing.T) {
	permission := randomPermission()

	testCases := []struct {
		name          string
		permissionID  int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			permissionID: permission.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPermission(gomock.Any(), gomock.Eq(permission.ID)).
					Times(1).
					Return(permission, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPermission(t, recorder.Body, permission)
			},
		},
		{
			name:         "NotFound",
			permissionID: permission.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPermission(gomock.Any(), gomock.Eq(permission.ID)).
					Times(1).
					Return(db.Permission{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:         "InternalError",
			permissionID: permission.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPermission(gomock.Any(), gomock.Eq(permission.ID)).
					Times(1).
					Return(db.Permission{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "InvalidID",
			permissionID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPermission(gomock.Any(), gomock.Any()).
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

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/admin/permission/%d", tc.permissionID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomPermission() db.Permission {
	return db.Permission{
		ID:   util.RandomInt(1, 1000),
		Name: util.RandomName(),
	}
}

func requireBodyMatchPermission(t *testing.T, body *bytes.Buffer, permission db.Permission) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotPermission db.Permission
	err = json.Unmarshal(data, &gotPermission)
	require.NoError(t, err)
	require.Equal(t, permission, gotPermission)
}
