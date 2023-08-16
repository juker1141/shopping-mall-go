package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/val"
)

type createRoleRequest struct {
	Name          string  `json:"name" binding:"required"`
	Status        int32   `json:"status" binding:"required"`
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
		Status:        req.Status,
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
	Status        int32   `json:"status"`
	PermissionsID []int64 `json:"permissions_id"`
}

func (server *Server) updateRole(ctx *gin.Context) {
	var query updateRoleRequestQuery
	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateRoleTxParams{
		ID: query.ID,
	}

	if len(req.Name) != 0 {
		arg.Name = req.Name
	}
	if val.IsValidStatus(int(req.Status)) {
		arg.Status = req.Status
	}
	if len(req.PermissionsID) != 0 || req.PermissionsID != nil {
		arg.PermissionsID = req.PermissionsID
	}

	result, err := server.store.UpdateRoleTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
