package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
)

type createOrderRequest struct {
	FullName        string                    `json:"fullName" binding:"required,fullName"`
	Email           string                    `json:"email" binding:"required"`
	ShippingAddress string                    `json:"shippingAddress" binding:"required"`
	Message         string                    `json:"message"`
	PayMethodID     int64                     `json:"payMethodId" binding:"required"`
	OrderProducts   []db.OrderTxProductParams `json:"orderProducts" binding:"required,min=1"`
	CouponID        int64                     `json:"couponId"`
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

type orderRoutesUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrder(ctx *gin.Context) {
	var uri orderRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, uri.ID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	status, err := server.store.GetOrderStatus(ctx, int64(order.StatusID))
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderProducts, err := server.store.ListOrderProductByOrderId(ctx, pgtype.Int4{
		Int32: int32(order.ID),
		Valid: true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var productList []db.OrderTxProductResult
	for _, orderProduct := range orderProducts {
		product, err := server.store.GetProduct(ctx, int64(orderProduct.ProductID.Int32))
		if err != nil {
			if err == db.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		productList = append(productList, db.OrderTxProductResult{
			Product: product,
			Num:     int64(orderProduct.Num),
		})
	}

	rsp := db.OrderTxResult{
		Order:       order,
		ProductList: productList,
		Status:      status,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listOrdersQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=5,max=10"`
}

type listOrdersResponse struct {
	Count int32              `json:"count"`
	Data  []db.OrderTxResult `json:"data"`
}

func (server *Server) listOrders(ctx *gin.Context) {
	var query listOrdersQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListOrdersParams{
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
	}

	orders, err := server.store.ListOrders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var data []db.OrderTxResult
	for _, order := range orders {
		status, err := server.store.GetOrderStatus(ctx, int64(order.StatusID))
		if err != nil {
			if err == db.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		orderProducts, err := server.store.ListOrderProductByOrderId(ctx, pgtype.Int4{
			Int32: int32(order.ID),
			Valid: true,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var productList []db.OrderTxProductResult
		for _, orderProduct := range orderProducts {
			product, err := server.store.GetProduct(ctx, int64(orderProduct.ProductID.Int32))
			if err != nil {
				if err == db.ErrRecordNotFound {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			productList = append(productList, db.OrderTxProductResult{
				Product: product,
				Num:     int64(orderProduct.Num),
			})
		}

		data = append(data, db.OrderTxResult{
			Order:       order,
			ProductList: productList,
			Status:      status,
		})
	}

	counts, err := server.store.GetOrdersCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listOrdersResponse{
		Count: int32(counts),
		Data:  data,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type updateOrderRequest struct {
	FullName        string                    `json:"fullName" binding:"required,fullName"`
	Email           string                    `json:"email" binding:"required"`
	ShippingAddress string                    `json:"shippingAddress" binding:"required"`
	Message         string                    `json:"message"`
	PayMethodID     int64                     `json:"payMethodId" binding:"required"`
	OrderProducts   []db.OrderTxProductParams `json:"orderProducts" binding:"required,min=1"`
	CouponID        int64                     `json:"couponId"`
}

func (server *Server) updateOrder(ctx *gin.Context) {
	var uri orderRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateOrderRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateOrderTxParams{
		ID: uri.ID,
	}

	if req.FullName != "" {
		arg.FullName = req.FullName
	}

	if req.Email != "" {
		arg.Email = req.Email
	}

	if req.ShippingAddress != "" {
		arg.ShippingAddress = req.ShippingAddress
	}

	if req.Message != "" {
		arg.Message = req.Message
	}

	if req.PayMethodID != 0 {
		arg.PayMethodID = req.PayMethodID
	}

	if req.OrderProducts != nil && len(req.OrderProducts) != 0 {
		arg.OrderProducts = req.OrderProducts
	}

	if req.CouponID != 0 {
		arg.CouponID = req.CouponID
	}

	result, err := server.store.UpdateOrderTx(ctx, arg)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
