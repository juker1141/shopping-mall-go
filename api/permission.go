package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type createPermissionRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createPermission(ctx *gin.Context) {
	var req createPermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	permission, err := server.store.CreatePermission(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, permission)
}

type listPermissionRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listPermission(ctx *gin.Context) {
	var req listPermissionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListPermissionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	permission, err := server.store.ListPermissions(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, permission)
}

type getPermissionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPermission(ctx *gin.Context) {
	var req getPermissionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	permission, err := server.store.GetPermission(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, permission)
}

type updatePermissionQuery struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updatePermission(ctx *gin.Context) {
	var req updatePermissionRequest
	var query updatePermissionQuery

	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePermissionParams{
		ID:   query.ID,
		Name: req.Name,
	}

	permission, err := server.store.UpdatePermission(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, permission)
}
