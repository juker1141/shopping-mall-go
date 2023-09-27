package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type listOrderStatusesOptionResponse struct {
	Statuses []db.OrderStatus `json:"statuses"`
}

func (server *Server) listOrderStatusesOption(ctx *gin.Context) {
	statuses, err := server.store.ListOrderStatus(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listOrderStatusesOptionResponse{
		Statuses: statuses,
	}

	ctx.JSON(http.StatusOK, rsp)
}
