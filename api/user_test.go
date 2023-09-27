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
	e.arg.AvatarUrl = arg.AvatarUrl

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

	templateBody := gin.H{
		"account":         user.Account,
		"email":           user.Email,
		"fullName":        user.FullName,
		"password":        password,
		"genderId":        user.GenderID.Int32,
		"cellphone":       user.Cellphone,
		"address":         user.Address,
		"shippingAddress": user.ShippingAddress,
		"postCode":        user.PostCode,
		"status":          user.Status,
	}

	testCases := []struct {
		name           string
		isUploadAvatar bool
		fileType       string
		fileName       string
		body           gin.H
		buildStubs     func(store *mockdb.MockStore)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:           "OK",
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           templateBody,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Account:  user.Account,
					Email:    user.Email,
					FullName: user.FullName,
					GenderID: pgtype.Int4{
						Int32: user.GenderID.Int32,
						Valid: true,
					},
					Cellphone:       user.Cellphone,
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
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:           "OK_NoImage",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           templateBody,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Account:  user.Account,
					Email:    user.Email,
					FullName: user.FullName,
					GenderID: pgtype.Int4{
						Int32: user.GenderID.Int32,
						Valid: true,
					},
					Cellphone:       user.Cellphone,
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
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:           "InternalError",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           templateBody,
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
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           templateBody,
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
			name:           "InvalidUsername",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         "invalid-user#1",
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          user.Status,
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
			name:           "InvalidFullName",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        "invalid_FullName#@",
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          user.Status,
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
			name:           "TooShortPassword",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        "psw",
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          user.Status,
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
			name:           "InvalidEmailAddress",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         user.Account,
				"email":           "invalidEmail",
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          user.Status,
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
			name:           "InvalidGenderID",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        4,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          user.Status,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Account:  user.Account,
					Email:    user.Email,
					FullName: user.FullName,
					GenderID: pgtype.Int4{
						Int32: 4,
						Valid: true,
					},
					Cellphone:       user.Cellphone,
					Address:         user.Address,
					ShippingAddress: user.ShippingAddress,
					PostCode:        user.PostCode,
					Status:          user.Status,
					AvatarUrl:       user.AvatarUrl,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:           "InvalidEmailPhoneNumber",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       "123456789",
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          user.Status,
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
			name:           "InvalidStatusInput",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          "0",
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
			name:           "InvalidImageType",
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.gif",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          1,
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
			name:           "InvalidImageContentType",
			isUploadAvatar: true,
			fileType:       "other",
			fileName:       "fake_avatar.pdf",
			body: gin.H{
				"account":         user.Account,
				"email":           user.Email,
				"fullName":        user.FullName,
				"password":        password,
				"genderId":        user.GenderID.Int32,
				"cellphone":       user.Cellphone,
				"address":         user.Address,
				"shippingAddress": user.ShippingAddress,
				"postCode":        user.PostCode,
				"status":          1,
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
				var fileWriter io.Writer
				var err error

				if tc.fileType == "image" {
					fileWriter, err = CreateImageFormFile(writer, "avatarFile", tc.fileName)
				} else {
					fileWriter, err = writer.CreateFormFile("avatarFile", tc.fileName)
				}

				require.NoError(t, err)

				fakeFileContent = append(fakeFileContent, []byte("Fake image data...")...)
				_, err = fileWriter.Write(fakeFileContent)
				require.NoError(t, err)
			}

			// 結束FormData
			writer.Close()

			url := "/user"
			request, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
}

func (e eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	if arg.HashedPassword.String != "" {
		err := util.CheckPassword(e.password, arg.HashedPassword.String)
		if err != nil {
			return false
		}
		e.arg.HashedPassword = pgtype.Text{
			String: arg.HashedPassword.String,
			Valid:  true,
		}
		e.arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  arg.PasswordChangedAt.Time,
			Valid: true,
		}
	}

	if arg.AvatarUrl.String != "" {
		e.arg.AvatarUrl = pgtype.Text{
			String: arg.AvatarUrl.String,
			Valid:  true,
		}
	}

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateUserParams(arg db.UpdateUserParams, password string) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg, password}
}

func TestUpdateUserByAdminAPI(t *testing.T) {
	user, password := randomUser(t, 1)
	newPassword := util.RandomString(8)
	newName := util.RandomName()
	newPhone := util.RandomCellPhone()
	newAddress := util.RandomAddress()
	newPostCode := util.RandomPostCode()
	newStatus := int32(0)

	invalidStatus := int32(3)

	templateBody := gin.H{
		"fullName":        newName,
		"cellphone":       newPhone,
		"address":         newAddress,
		"shippingAddress": newAddress,
		"postCode":        newPostCode,
		"status":          newStatus,
		"oldPassword":     password,
		"newPassword":     newPassword,
	}

	testCases := []struct {
		name           string
		ID             int64
		isUploadAvatar bool
		fileType       string
		fileName       string
		body           gin.H
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs     func(store *mockdb.MockStore)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:           "OK",
			ID:             user.ID,
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)

				arg := db.UpdateUserParams{
					ID: user.ID,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Cellphone: pgtype.Text{
						String: newPhone,
						Valid:  true,
					},
					Address: pgtype.Text{
						String: newAddress,
						Valid:  true,
					},
					ShippingAddress: pgtype.Text{
						String: newAddress,
						Valid:  true,
					},
					PostCode: pgtype.Text{
						String: newPostCode,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: newStatus,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg, newPassword)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "OK_NoUpdatePassword",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"fullName":        newName,
				"cellphone":       newPhone,
				"address":         newAddress,
				"shippingAddress": newAddress,
				"postCode":        newPostCode,
				"status":          newStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				arg := db.UpdateUserParams{
					ID: user.ID,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Cellphone: pgtype.Text{
						String: newPhone,
						Valid:  true,
					},
					Address: pgtype.Text{
						String: newAddress,
						Valid:  true,
					},
					ShippingAddress: pgtype.Text{
						String: newAddress,
						Valid:  true,
					},
					PostCode: pgtype.Text{
						String: newPostCode,
						Valid:  true,
					},
					Status: pgtype.Int4{
						Int32: newStatus,
						Valid: true,
					},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "OK_OnlyPassword",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"oldPassword": password,
				"newPassword": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)

				arg := db.UpdateUserParams{
					ID: user.ID,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg, newPassword)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "OK_AllFieldEmpty",
			ID:             user.ID,
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				arg := db.UpdateUserParams{
					ID: user.ID,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "NoAuthorization",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"fullName":        newName,
				"cellphone":       newPhone,
				"address":         newAddress,
				"shippingAddress": newAddress,
				"postCode":        newPostCode,
				"status":          newStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:           "NoRequiredPermission",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"fullName":        newName,
				"cellphone":       newPhone,
				"address":         newAddress,
				"shippingAddress": newAddress,
				"postCode":        newPostCode,
				"status":          newStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:           "InternalError",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:           "InvalidFullName",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"fullName": "invalid-name",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "InvalidStatus",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"status": invalidStatus,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "InvalidPhone",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"cellphone": "1234567890",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "InvalidImageType",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_avatar.gif",
			body: gin.H{
				"cellphone": "1234567890",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "InvalidImageType",
			ID:             user.ID,
			isUploadAvatar: true,
			fileType:       "other",
			fileName:       "fake_avatar.pdf",
			body: gin.H{
				"cellphone": "1234567890",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "InvalidUserID",
			ID:             -1,
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body:           gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "WrongPassword",
			ID:             user.ID,
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_avatar.png",
			body: gin.H{
				"newPassword": "aa345678",
				"oldPassword": "za345678",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
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
				var fileWriter io.Writer
				var err error

				if tc.fileType == "image" {
					fileWriter, err = CreateImageFormFile(writer, "avatarFile", tc.fileName)
				} else {
					fileWriter, err = writer.CreateFormFile("avatarFile", tc.fileName)
				}

				require.NoError(t, err)

				fakeFileContent = append(fakeFileContent, []byte("Fake image data...")...)
				_, err = fileWriter.Write(fakeFileContent)
				require.NoError(t, err)
			}

			// 結束FormData
			writer.Close()

			url := fmt.Sprintf("/admin/memberUser/%d", tc.ID)
			request, err := http.NewRequest(http.MethodPatch, url, body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetUserByAdminAPI(t *testing.T) {
	user, _ := randomUser(t, int32(1))

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/memberUser/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListUsersByAdminAPI(t *testing.T) {
	n := 5

	userList := make([]db.User, n)
	for i := 0; i < n; i++ {
		user, _ := randomUser(t, int32(1))
		userList[i] = user
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
				addPermissionMiddleware(store, "user", memberPermissions)

				arg := db.ListUsersParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(userList, nil)

				store.EXPECT().
					GetUsersCount(gomock.Any()).
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
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetUsersCount(gomock.Any()).
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
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetUsersCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", memberPermissions)

				arg := db.ListUsersParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.User{}, sql.ErrConnDone)

				store.EXPECT().
					GetUsersCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetUsersCount(gomock.Any()).
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
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetUsersCount(gomock.Any()).
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

			url := "/admin/memberUsers"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.page))
			q.Add("pageSize", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteUserByAdminAPI(t *testing.T) {
	user, _ := randomUser(t, int32(1))

	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoRequiredPermission",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   user.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
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
				addPermissionMiddleware(store, "user", memberPermissions)

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/admin/memberUser/%d", tc.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

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
		ID:             util.RandomID(),
		Account:        util.RandomAccount(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		GenderID: pgtype.Int4{
			Int32: util.RandomGender(),
			Valid: true,
		},
		Cellphone:       util.RandomCellPhone(),
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

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser userResponse
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Account, gotUser.Account)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.AvatarUrl, gotUser.AvatarUrl)
}
