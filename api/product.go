package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
)

type createProductRequest struct {
	Title       string                `form:"title" binding:"required,min=3"`
	Category    string                `form:"category" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Content     string                `form:"content" binding:"required"`
	OriginPrice int32                 `form:"origin_price" binding:"required"`
	Price       int32                 `form:"price" binding:"required"`
	Unit        string                `form:"unit" binding:"required"`
	Status      int32                 `form:"status" binding:"required"`
	ImageFile   *multipart.FileHeader `form:"image_file"`
	ImageUrl    string                `form:"image_url"`
	ImagesUrl   *[]string             `form:"images_url"`
}

// type productResponse struct {
// 	Title       string   `json:"title"`
// 	Category    string   `json:"category"`
// 	Description string   `json:"description"`
// 	Content     string   `json:"content"`
// 	OriginPrice int32    `json:"origin_price"`
// 	Price       int32    `json:"price"`
// 	Unit        string   `json:"unit"`
// 	Status      int32    `json:"status"`
// 	ImageUrl    string   `json:"image_url"`
// 	ImagesUrl   []string `json:"images_url"`
// 	CreatedBy   string   `json:"created_by"`
// }

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		err := errors.New("can not get permission")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	typePayload, ok := payload.(*token.Payload)
	if !ok {
		err := errors.New("payload is not of type Payload")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.CreateProductParams{
		Title:       req.Title,
		Category:    req.Category,
		Description: req.Description,
		Content:     req.Content,
		OriginPrice: req.OriginPrice,
		Price:       req.Price,
		Unit:        req.Unit,
		Status:      req.Status,
		CreatedBy:   typePayload.Account,
	}

	if req.ImageFile == nil && req.ImageUrl == "" {
		err := fmt.Errorf("an image file or an image URL is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.ImageFile != nil {
		file, err := ctx.FormFile("image_file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") || strings.HasSuffix(file.Filename, ".gif") {
			err := fmt.Errorf("only non-GIF image files are allowed for upload")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		const maxSize = 5 << 20 // 5MB
		if file.Size > maxSize {
			err := fmt.Errorf("file size exceeds 5 MB")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
		targetPath := filepath.Join("static", "products", filename)

		err = ctx.SaveUploadedFile(req.ImageFile, targetPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		arg.ImageUrl = targetPath
	} else if req.ImageUrl != "" {
		arg.ImageUrl = req.ImageUrl
	}

	if req.ImagesUrl != nil && len(*req.ImagesUrl) > 0 {
		arg.ImagesUrl = *req.ImagesUrl
	} else {
		arg.ImagesUrl = []string{}
	}

	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

type listProductQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listProductResponse struct {
	Count int32        `json:"count"`
	Data  []db.Product `json:"data"`
}

func (server *Server) listProduct(ctx *gin.Context) {
	var query listProductQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductsParams{
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
	}

	products, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	counts, err := server.store.GetProductsCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listProductResponse{
		Count: int32(counts),
		Data:  products,
	}

	ctx.JSON(http.StatusOK, rsp)
}