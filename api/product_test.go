package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
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

func TestCreateProductAPI(t *testing.T) {
	product := randomProduct()

	templateBody := gin.H{
		"title":        product.Title,
		"category":     product.Category,
		"origin_price": product.OriginPrice,
		"price":        product.Price,
		"unit":         product.Unit,
		"description":  product.Description,
		"content":      product.Content,
		"status":       1,
		"image_url":    product.ImageUrl,
	}

	testCases := []struct {
		name           string
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
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_product.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.CreateProductParams{
					Title:       product.Title,
					Category:    product.Category,
					Description: product.Description,
					Content:     product.Content,
					Price:       product.Price,
					OriginPrice: product.OriginPrice,
					Unit:        product.Unit,
					Status:      product.Status,
					ImageUrl:    product.ImageUrl,
					ImagesUrl:   []string{},
					CreatedBy:   "user",
				}

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "OK_ButNoUploadImage",
			isUploadAvatar: false,
			fileType:       "image",
			fileName:       "fake_product.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.CreateProductParams{
					Title:       product.Title,
					Category:    product.Category,
					Description: product.Description,
					Content:     product.Content,
					Price:       product.Price,
					OriginPrice: product.OriginPrice,
					Unit:        product.Unit,
					Status:      product.Status,
					ImageUrl:    product.ImageUrl,
					ImagesUrl:   []string{},
					CreatedBy:   "user",
				}

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "NoAuthorization",
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_product.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:           "NoRequiredPermission",
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_product.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", emptyPermission)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:           "InternalError",
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_product.png",
			body:           templateBody,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				arg := db.CreateProductParams{
					Title:       product.Title,
					Category:    product.Category,
					Description: product.Description,
					Content:     product.Content,
					Price:       product.Price,
					OriginPrice: product.OriginPrice,
					Unit:        product.Unit,
					Status:      product.Status,
					ImageUrl:    product.ImageUrl,
					ImagesUrl:   []string{},
					CreatedBy:   "user",
				}

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:           "TooShortTitle",
			isUploadAvatar: true,
			fileType:       "image",
			fileName:       "fake_product.png",
			body: gin.H{
				"title":        "aa",
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
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
			fileName:       "fake_product.gif",
			body: gin.H{
				"title":        "aa",
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
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
			fileName:       "fake_product.pdf",
			body: gin.H{
				"title":        "aa",
				"category":     product.Category,
				"origin_price": product.OriginPrice,
				"price":        product.Price,
				"unit":         product.Unit,
				"description":  product.Description,
				"content":      product.Content,
				"status":       1,
				"image_url":    product.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				addPermissionMiddleware(store, "user", productPermissions)

				store.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
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
					fileWriter, err = CreateImageFormFile(writer, "products", tc.fileName)
				} else {
					fileWriter, err = writer.CreateFormFile("products", tc.fileName)
				}
				require.NoError(t, err)

				fakeFileContent = append(fakeFileContent, []byte("Fake image data...")...)
				_, err = fileWriter.Write(fakeFileContent)
				require.NoError(t, err)
			}

			// 結束FormData
			writer.Close()

			url := "/admin/product"
			request, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomProduct() db.Product {
	targetPath := filepath.Join("static", "products", "fake_product.png")

	return db.Product{
		ID:          util.RandomID(),
		Title:       util.RandomName(),
		Category:    util.RandomName(),
		OriginPrice: util.RandomPrice(),
		Price:       util.RandomPrice(),
		Unit:        util.RandomName(),
		Description: util.RandomString(20),
		Content:     util.RandomString(20),
		Status:      1,
		ImageUrl:    targetPath,
		CreatedBy:   "user",
	}
}
