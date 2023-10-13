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
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/val"
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

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		err := errors.New("can not get token payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	typePayload, ok := payload.(*token.Payload)
	if !ok {
		err := errors.New("payload is not of type payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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

type listProductsQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listProductsResponse struct {
	Count int32        `json:"count"`
	Data  []db.Product `json:"data"`
}

func (server *Server) listProducts(ctx *gin.Context) {
	var query listProductsQuery
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

	rsp := listProductsResponse{
		Count: int32(counts),
		Data:  products,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type productRoutesUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateProductRequest struct {
	Title       string                `form:"title"`
	Category    string                `form:"category"`
	Description string                `form:"description"`
	Content     string                `form:"content"`
	OriginPrice int32                 `form:"origin_price"`
	Price       int32                 `form:"price"`
	Unit        string                `form:"unit"`
	Status      *int32                `form:"status"`
	ImageFile   *multipart.FileHeader `form:"image_file"`
	ImageUrl    string                `form:"image_url"`
	ImagesUrl   *[]string             `form:"images_url"`
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var uri productRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateProductParams{
		ID: uri.ID,
	}

	if req.Title != "" {
		if err := val.ValidateTitle(req.Title); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		arg.Title = pgtype.Text{
			String: req.Title,
			Valid:  true,
		}
	}

	if req.Category != "" {
		arg.Category = pgtype.Text{
			String: req.Category,
			Valid:  true,
		}
	}

	if req.Content != "" {
		arg.Content = pgtype.Text{
			String: req.Content,
			Valid:  true,
		}
	}

	if req.Description != "" {
		arg.Description = pgtype.Text{
			String: req.Description,
			Valid:  true,
		}
	}

	if req.Unit != "" {
		arg.Unit = pgtype.Text{
			String: req.Unit,
			Valid:  true,
		}
	}

	if req.Price >= 0 {
		arg.Price = pgtype.Int4{
			Int32: req.Price,
			Valid: true,
		}
	}

	if req.OriginPrice >= 0 {
		arg.OriginPrice = pgtype.Int4{
			Int32: req.OriginPrice,
			Valid: true,
		}
	}

	if req.Status != nil {
		if val.IsValidStatus(int(*req.Status)) {
			arg.Status = pgtype.Int4{
				Int32: *req.Status,
				Valid: true,
			}
		} else {
			err := fmt.Errorf("invalid status input: %d", *req.Status)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
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

		arg.ImageUrl = pgtype.Text{
			String: targetPath,
			Valid:  true,
		}
	} else if req.ImageUrl != "" {
		arg.ImageUrl = pgtype.Text{
			String: req.ImageUrl,
			Valid:  true,
		}
	}

	if req.ImagesUrl != nil && len(*req.ImagesUrl) > 0 {
		arg.ImagesUrl = *req.ImagesUrl
	}

	product, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) getProduct(ctx *gin.Context) {
	var uri productRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, uri.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	var uri productRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteProduct(ctx, uri.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := deleteResult{
		Message: "Delete product success.",
	}

	ctx.JSON(http.StatusOK, result)
}
