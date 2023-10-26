package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type updateCartRequest struct {
	ProductID  int64  `json:"product_id" binding:"required"`
	Num        int32  `json:"num" binding:"required,gt=0"`
	CouponCode string `json:"coupon_code"`
}

type updateCartQuery struct {
	Type string `form:"type" binding:"required"`
}

type cartResponse struct {
	db.Cart
	ProductList []db.OrderTxProductResult `json:"product_list"`
	Coupon      db.Coupon                 `json:"coupon"`
}

func (server *Server) updateCart(ctx *gin.Context) {
	// var req updateCartRequest
	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }

	// var query updateCartQuery
	// if err := ctx.ShouldBindQuery(&query); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }

	// payload, exists := ctx.Get(authorizationPayloadKey)
	// if !exists {
	// 	err := errors.New("can not get token payload")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }

	// typePayload, ok := payload.(*token.Payload)
	// if !ok {
	// 	err := errors.New("payload is not of type payload")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }

	// rsp := cartResponse{
	// 	Cart:        updateCart,
	// 	ProductList: productList,
	// 	Coupon:      coupon,
	// }

	// ctx.JSON(http.StatusOK, rsp)
}
