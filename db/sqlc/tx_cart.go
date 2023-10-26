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
	Account    string `json:"account" binding:"required"`
	Type       string `json:"type"`
	ProductID  int64  `json:"product_id" binding:"required"`
	Num        int32  `json:"num" binding:"required,gt=0"`
	CouponCode string `json:"coupon_code"`
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

		cart, err := q.GetCartByOwner(ctx, pgtype.Text{
			String: arg.Account,
			Valid:  true,
		})
		if err != nil {
			fmt.Println(err, "getCart")
			return err
		}

		err = q.DeleteCartProduct(context.Background(), DeleteCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(arg.ProductID),
				Valid: true,
			},
		})

		_, err = q.CreateCartProduct(ctx, CreateCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(arg.ProductID),
				Valid: true,
			},
			Num: arg.Num,
		})

		if err != nil {
			fmt.Println(err, "Create")
			return err
		}

		// if reflect.DeepEqual(cartProduct, CartProduct{}) {

		// 	_, err = q.CreateCartProduct(ctx, CreateCartProductParams{
		// 		CartID: pgtype.Int4{
		// 			Int32: int32(cart.ID),
		// 			Valid: true,
		// 		},
		// 		ProductID: pgtype.Int4{
		// 			Int32: int32(arg.ProductID),
		// 			Valid: true,
		// 		},
		// 		Num: arg.Num,
		// 	})
		// } else {

		// 	_, err = q.UpdateCartProduct(ctx, UpdateCartProductParams{
		// 		CartID: pgtype.Int4{
		// 			Int32: int32(cart.ID),
		// 			Valid: true,
		// 		},
		// 		ProductID: pgtype.Int4{
		// 			Int32: int32(arg.ProductID),
		// 			Valid: true,
		// 		},
		// 		Num: pgtype.Int4{
		// 			Int32: updateNum,
		// 			Valid: true,
		// 		},
		// 	})
		// 	if err != nil {
		// 		fmt.Println(err, "updateCartP")
		// 		return err
		// 	}
		// }

		cartProducts, err := q.ListCartProductByCartId(ctx, pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		})
		if err != nil {
			fmt.Println(err, "ListCartProduct")
			return err
		}

		var productList []CartTxProductResult
		var totalPrice int32
		var finalPrice int32
		for _, cartProduct := range cartProducts {
			if cartProduct.Num <= 0 {
				err := fmt.Errorf("product num must be positive")
				return err
			}
			// 取得商品
			product, err := q.GetProduct(ctx, int64(cartProduct.ProductID.Int32))
			if err != nil {
				fmt.Println(err, "getProduct")
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

		var coupon Coupon
		if arg.CouponCode != "" {
			coupon, err = q.GetCouponByCode(ctx, arg.CouponCode)
			if err != nil {
				fmt.Println(err, "get Coupon")
				return err
			}
			finalPrice = finalPrice * (100 - coupon.Percent) / 100

			err = q.DeleteCartCouponByCartId(ctx, pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			})
			if err != nil {
				fmt.Println(err, "del Coupon")
				return err
			}

			_, err = q.CreateCartCoupon(ctx, CreateCartCouponParams{
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
				fmt.Println(err, "cartCoupon")
				return err
			}
		}

		cartArg := UpdateCartParams{
			ID:         cart.ID,
			TotalPrice: totalPrice,
			FinalPrice: finalPrice,
		}

		result.Cart, err = q.UpdateCart(ctx, cartArg)
		if err != nil {
			fmt.Println(err, "updateCart")
			return err
		}

		return nil
	})
	return result, err
}
