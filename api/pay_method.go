package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type listPayMethodsOptionResponse struct {
	PayMethods []db.PayMethod `json:"payMethods"`
}

func (server *Server) listPayMethodsOption(ctx *gin.Context) {
	paymethods, err := server.store.ListPayMethod(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listPayMethodsOptionResponse{
		PayMethods: paymethods,
	}

	ctx.JSON(http.StatusOK, rsp)
}
