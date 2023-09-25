package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
)

type createOrderRequest struct {
	FullName        string                  `json:"full_name" binding:"required,fullName"`
	Email           string                  `json:"email" binding:"required"`
	ShippingAddress string                  `json:"shipping_address" binding:"required"`
	Message         string                  `json:"message"`
	PayMethodID     int64                   `json:"pay_method_id" binding:"required"`
	OrderProducts   []db.OrderProductParams `json:"order_products" binding:"required,min=1"`
	CouponID        int64                   `json:"coupon_id"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if len(req.OrderProducts) <= 0 {
		err := fmt.Errorf("at least one item must be present in the order")
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

	user, err := server.store.GetUserByAccount(ctx, typePayload.Account)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	defaultStatus := int64(1)

	arg := db.CreateOrderTxParams{
		UserID:          user.ID,
		FullName:        req.FullName,
		Email:           req.Email,
		ShippingAddress: req.ShippingAddress,
		PayMethodID:     req.PayMethodID,
		StatusID:        defaultStatus,
		OrderProducts:   req.OrderProducts,
	}

	if req.Message != "" {
		arg.Message = req.Message
	}

	if req.CouponID != 0 {
		arg.CouponID = req.CouponID
	}

	result, err := server.store.CreateOrderTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
