package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
)

type updateCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Num       int32 `json:"num" binding:"required,gt=0"`
}

type updateCartQuery struct {
	Type string `form:"type" binding:"required"`
}

type cartResponse struct {
	db.Cart
	ProductList []db.CartTxProductResult `json:"product_list"`
	Coupon      db.Coupon                `json:"coupon"`
}

func (server *Server) updateCartProduct(ctx *gin.Context) {
	var req updateCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var query updateCartQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
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

	cart, err := server.store.GetCartByOwner(ctx, pgtype.Text{
		String: typePayload.Account,
		Valid:  true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cartProductisExists, err := server.store.CheckCartProductExists(ctx, db.CheckCartProductExistsParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(req.ProductID),
			Valid: true,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !cartProductisExists {
		_, err = server.store.CreateCartProduct(ctx, db.CreateCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(req.ProductID),
				Valid: true,
			},
			Num: req.Num,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		cartProduct, err := server.store.GetCartProduct(ctx, db.GetCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(req.ProductID),
				Valid: true,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var updateNum int32
		if strings.ToLower(query.Type) == "add" {
			updateNum = cartProduct.Num + req.Num
		} else if strings.ToLower(query.Type) == "update" {
			updateNum = req.Num
		}

		_, err = server.store.UpdateCartProduct(ctx, db.UpdateCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(req.ProductID),
				Valid: true,
			},
			Num: pgtype.Int4{
				Int32: updateNum,
				Valid: true,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	result, err := server.store.UpdateCartTx(ctx, db.UpdateCartTxParams{
		CartID: cart.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := cartResponse{
		Cart:        result.Cart,
		ProductList: result.ProductList,
		Coupon:      result.Coupon,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type addCartCouponRequest struct {
	Code string `json:"code" binding:"required"`
}

func (server *Server) addCartCoupon(ctx *gin.Context) {
	var req addCartCouponRequest
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

	cart, err := server.store.GetCartByOwner(context.Background(), pgtype.Text{
		String: typePayload.Account,
		Valid:  true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	coupon, err := server.store.GetCouponByCode(ctx, req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteCartCouponByCartId(ctx, pgtype.Int4{
		Int32: int32(cart.ID),
		Valid: true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateCartCoupon(ctx, db.CreateCartCouponParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: int32(coupon.ID),
			Valid: true,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result, err := server.store.UpdateCartTx(ctx, db.UpdateCartTxParams{
		CartID: cart.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := cartResponse{
		Cart:        result.Cart,
		ProductList: result.ProductList,
		Coupon:      result.Coupon,
	}

	ctx.JSON(http.StatusOK, rsp)
}
