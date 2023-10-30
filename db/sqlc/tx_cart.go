package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type CartTxProductResult struct {
	Product
	Num int64 `json:"num"`
}

type UpdateCartTxParams struct {
	CartID int64 `json:"cart_id" binding:"required,gt=0"`
}

type CartTxResult struct {
	Cart
	ProductList []CartTxProductResult `json:"product_list"`
	Coupon      Coupon                `json:"coupon"`
}

func (store *SQLStore) UpdateCartTx(ctx context.Context, arg UpdateCartTxParams) (CartTxResult, error) {
	var result CartTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var productList []CartTxProductResult
		var totalPrice int32
		var finalPrice int32

		cartProducts, err := q.ListCartProductByCartId(ctx, pgtype.Int4{
			Int32: int32(arg.CartID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		for _, cartProduct := range cartProducts {
			if cartProduct.Num <= 0 {
				err := fmt.Errorf("product num must be positive")
				return err
			}
			// 取得商品
			product, err := q.GetProduct(ctx, int64(cartProduct.ProductID.Int32))
			if err != nil {
				return err
			}
			totalPrice = totalPrice + (int32(cartProduct.Num) * product.OriginPrice)
			finalPrice = finalPrice + (int32(cartProduct.Num) * product.Price)
			productList = append(productList, CartTxProductResult{
				Product: product,
				Num:     int64(cartProduct.Num),
			})
		}

		result.ProductList = productList

		cartCouponisExists, err := q.CheckCartCouponExists(context.Background(), pgtype.Int4{
			Int32: int32(arg.CartID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		var coupon Coupon
		if cartCouponisExists {
			cartCoupons, err := q.ListCartCouponByCartId(context.Background(),
				pgtype.Int4{
					Int32: int32(arg.CartID),
				})
			if err != nil {
				return err
			}

			cartCoupon := cartCoupons[0]

			coupon, err = q.GetCoupon(context.Background(), int64(cartCoupon.CouponID.Int32))
			if err != nil {
				return err
			}
			finalPrice = finalPrice * (100 - coupon.Percent) / 100
		}
		result.Coupon = coupon

		cartArg := UpdateCartParams{
			ID:         arg.CartID,
			TotalPrice: totalPrice,
			FinalPrice: finalPrice,
		}

		result.Cart, err = q.UpdateCart(ctx, cartArg)
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
