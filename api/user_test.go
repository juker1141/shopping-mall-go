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
	"net/textproto"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/juker1141/shopping-mall-go/db/mock"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	e.arg.AvatarUrl = ""
	arg.AvatarUrl = ""

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func CreateImageFormFile(w *multipart.Writer, fieldname string, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", "image/*")

	return w.CreatePart(h)
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t, int32(1))
	cart := randomCart(user.Account)
	genderId := util.RandomGender()

	templateBody := gin.H{
		"account":          user.Account,
		"email":            user.Email,
		"full_name":        user.FullName,
		"password":         password,
		"gender_id":        genderId,
		"phone":            user.Phone,
		"address":          user.Address,
		"shipping_address": user.ShippingAddress,
		"post_code":        user.PostCode,
		"status":           user.Status,
	}

	testCases := []struct {
		name           string
		isUploadAvatar bool
		body           gin.H
		buildStubs     func(store *mockdb.MockStore)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:           "OK",
			isUploadAvatar: true,
			body:           templateBody,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Account:  user.Account,
					Email:    user.Email,
					FullName: user.FullName,
					GenderID: pgtype.Int4{
						Int32: genderId,
						Valid: true,
					},
					Phone:           user.Phone,
					Address:         user.Address,
					ShippingAddress: user.ShippingAddress,
					PostCode:        user.PostCode,
					Status:          user.Status,
					AvatarUrl:       user.AvatarUrl,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)

				cartArg := db.CreateCartParams{
					Owner: pgtype.Text{
						String: user.Account,
						Valid:  true,
					},
					TotalPrice: 0,
					FinalPrice: 0,
				}

				store.EXPECT().CreateCart(gomock.Any(), gomock.Eq(cartArg)).Times(1).Return(cart, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "InternalError",
			body:           templateBody,
			isUploadAvatar: false,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:           "DuplicateAccount",
			body:           templateBody,
			isUploadAvatar: false,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"account":          "invalid-user#1",
				"email":            user.Email,
				"full_name":        user.FullName,
				"password":         password,
				"gender_id":        genderId,
				"phone":            user.Phone,
				"address":          user.Address,
				"shipping_address": user.ShippingAddress,
				"post_code":        user.PostCode,
				"status":           user.Status,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidFullName",
			body: gin.H{
				"account":          user.Account,
				"email":            user.Email,
				"full_name":        "invalid_FullName#@",
				"password":         password,
				"gender_id":        genderId,
				"phone":            user.Phone,
				"address":          user.Address,
				"shipping_address": user.ShippingAddress,
				"post_code":        user.PostCode,
				"status":           user.Status,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"account":          user.Account,
				"email":            user.Email,
				"full_name":        user.FullName,
				"password":         "psw",
				"gender_id":        genderId,
				"phone":            user.Phone,
				"address":          user.Address,
				"shipping_address": user.ShippingAddress,
				"post_code":        user.PostCode,
				"status":           user.Status,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmailAddress",
			body: gin.H{
				"account":          user.Account,
				"email":            "invalidEmail",
				"full_name":        user.FullName,
				"password":         password,
				"gender_id":        genderId,
				"phone":            user.Phone,
				"address":          user.Address,
				"shipping_address": user.ShippingAddress,
				"post_code":        user.PostCode,
				"status":           user.Status,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmailPhoneNumber",
			body: gin.H{
				"account":          user.Account,
				"email":            user.Email,
				"full_name":        user.FullName,
				"password":         password,
				"gender_id":        genderId,
				"phone":            "123456789",
				"address":          user.Address,
				"shipping_address": user.ShippingAddress,
				"post_code":        user.PostCode,
				"status":           user.Status,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidStatusInput",
			body: gin.H{
				"account":          user.Account,
				"email":            user.Email,
				"full_name":        user.FullName,
				"password":         password,
				"gender_id":        genderId,
				"phone":            user.Phone,
				"address":          user.Address,
				"shipping_address": user.ShippingAddress,
				"post_code":        user.PostCode,
				"status":           "0",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
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

			if tc.isUploadAvatar {
				var fakeFileContent []byte

				fileWriter, err := CreateImageFormFile(writer, "avatar_file", "fake_avatar.png")
				require.NoError(t, err)

				fakeFileContent = append(fakeFileContent, []byte("Fake image data...")...)
				_, err = fileWriter.Write(fakeFileContent)
				require.NoError(t, err)
			}

			// 結束FormData
			writer.Close()

			// url := "/user"
			request, err := http.NewRequest(http.MethodPost, "/user", body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t, int32(1))

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account":  user.Account,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq(user.Account)).
					Times(1).
					Return(user, nil)
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
					GetUserByAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"account":  user.Account,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"account":  user.Account,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByAccount(gomock.Any(), gomock.Eq(user.Account)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
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

			url := "/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T, status int32) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	targetPath := filepath.Join("static", "avatar_images", "fake_avatar.png")

	user = db.User{
		ID:             util.RandomInt(1, 100),
		Account:        util.RandomAccount(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		GenderID: pgtype.Int4{
			Int32: util.RandomGender(),
			Valid: true,
		},
		Phone:           util.RandomPhone(),
		Address:         util.RandomString(20),
		ShippingAddress: util.RandomString(20),
		PostCode:        util.RandomPostCode(),
		Status:          status,
		AvatarUrl:       targetPath,
	}
	return
}

func randomCart(account string) db.Cart {
	return db.Cart{
		Owner: pgtype.Text{
			String: account,
			Valid:  true,
		},
		TotalPrice: 0,
		FinalPrice: 0,
	}
}
