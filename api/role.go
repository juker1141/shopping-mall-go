package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
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
