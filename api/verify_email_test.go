package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func TestVerifyEmailAPI(t *testing.T) {
	user, _ := randomUser(t, 1)
	verifyEmail := randomVerifyEmail(user)

	type Query struct {
		emailId    int64
		secretCode string
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				emailId:    verifyEmail.ID,
				secretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().VerifyEmailTx(gomock.Any(), db.VerifyEmailTxParams{
					EmailId:    verifyEmail.ID,
					SecretCode: verifyEmail.SecretCode,
				}).Times(1).Return(db.VerifyEmailTxResult{
					User:        user,
					VerifyEmail: verifyEmail,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				emailId:    verifyEmail.ID,
				secretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().VerifyEmailTx(gomock.Any(), db.VerifyEmailTxParams{
					EmailId:    verifyEmail.ID,
					SecretCode: verifyEmail.SecretCode,
				}).Times(1).Return(db.VerifyEmailTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidEmailId",
			query: Query{
				emailId:    -1,
				secretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().VerifyEmailTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidLengthSecretCode",
			query: Query{
				emailId:    verifyEmail.ID,
				secretCode: util.RandomString(20),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().VerifyEmailTx(gomock.Any(), gomock.Any()).Times(0)
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

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			url := "/user/verifyEmail"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("emailId", fmt.Sprintf("%d", tc.query.emailId))
			q.Add("secretCode", tc.query.secretCode)
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomVerifyEmail(user db.User) db.VerifyEmail {
	return db.VerifyEmail{
		ID: util.RandomID(),
		UserID: pgtype.Int4{
			Int32: int32(user.ID),
			Valid: true,
		},
		Email:      user.Email,
		SecretCode: util.RandomString(32),
		IsUsed:     false,
	}
}
