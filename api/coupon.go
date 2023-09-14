package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
)

type createCouponRequest struct {
	Title     string    `json:"title" binding:"required"`
	Code      string    `json:"code" binding:"required"`
	Percent   int32     `json:"percent" binding:"required"`
	StartAt   time.Time `json:"start_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (server *Server) createCoupon(ctx *gin.Context) {
	var req createCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	arg := db.CreateCouponParams{
		Title:     req.Title,
		Code:      req.Code,
		Percent:   req.Percent,
		CreatedBy: typePayload.Account,
	}

	if req.StartAt.IsZero() {
		arg.StartAt = time.Now()
	} else {
		arg.StartAt = req.StartAt
	}

	if req.ExpiresAt.IsZero() {
		arg.ExpiresAt = time.Time{}
	} else {
		arg.ExpiresAt = req.ExpiresAt
	}

	coupon, err := server.store.CreateCoupon(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, coupon)
}
