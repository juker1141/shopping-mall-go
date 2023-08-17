package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type createRoleRequest struct {
	Name          string  `json:"name" binding:"required"`
	PermissionsID []int64 `json:"permissions_id" binding:"required"`
}

func (server *Server) createRole(ctx *gin.Context) {
	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateRoleTxParams{
		Name:          req.Name,
		PermissionsID: req.PermissionsID,
	}

	result, err := server.store.CreateRoleTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type updateRoleRequestQuery struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateRoleRequest struct {
	Name          string  `json:"name"`
	PermissionsID []int64 `json:"permissions_id"`
}

func (server *Server) updateRole(ctx *gin.Context) {
	var query updateRoleRequestQuery
	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateRoleRequest
	ctx.BindJSON(&req)

	arg := db.UpdateRoleTxParams{
		ID: query.ID,
	}

	if req.Name != "" {
		arg.Name = req.Name
	}

	if req.PermissionsID != nil {
		arg.PermissionsID = req.PermissionsID
	}

	result, err := server.store.UpdateRoleTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
