package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

type couponRoutesUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateCouponRequest struct {
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	Percent   int32     `json:"percent"`
	StartAt   time.Time `json:"start_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (server *Server) updateCoupon(ctx *gin.Context) {
	var uri couponRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateCouponRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCouponParams{
		ID: uri.ID,
	}

	if req.Title != "" {
		arg.Title = pgtype.Text{
			String: req.Title,
			Valid:  true,
		}
	}

	if req.Code != "" {
		arg.Code = pgtype.Text{
			String: req.Code,
			Valid:  true,
		}
	}

	if req.Percent > 0 {
		arg.Percent = pgtype.Int4{
			Int32: req.Percent,
			Valid: true,
		}
	}

	if !req.StartAt.IsZero() {
		arg.StartAt = pgtype.Timestamptz{
			Time:  req.StartAt,
			Valid: true,
		}
	}

	if !req.ExpiresAt.IsZero() {
		arg.ExpiresAt = pgtype.Timestamptz{
			Time:  req.ExpiresAt,
			Valid: true,
		}
	}

	coupon, err := server.store.UpdateCoupon(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, coupon)
}
