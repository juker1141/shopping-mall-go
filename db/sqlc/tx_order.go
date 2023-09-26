package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/val"
)

type OrderTxProductParams struct {
	ID  int64 `json:"id"`
	Num int64 `json:"num"`
}

// OrderTxParams contains the input parameters of the role create
type CreateOrderTxParams struct {
	UserID          int64                  `json:"user_id"`
	FullName        string                 `json:"full_name"`
	Email           string                 `json:"email"`
	ShippingAddress string                 `json:"shipping_address"`
	Message         string                 `json:"message"`
	PayMethodID     int64                  `json:"pay_method_id"`
	StatusID        int64                  `json:"status_id"`
	OrderProducts   []OrderTxProductParams `json:"order_products"`
	CouponID        int64                  `json:"coupon_id"`
}

type OrderTxProductResult struct {
	Product
	Num int64 `json:"num"`
}

type OrderTxResult struct {
	Order
	ProductList []OrderTxProductResult `json:"product_list"`
	Status      OrderStatus            `json:"status"`
}

// It creates a order, orderUser, orderProduct, orderCoupon, and get orderStatus within a single database trasaction
func (store *SQLStore) CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		var productList []OrderTxProductResult
		var totalPrice int32
		var finalPrice int32
		for _, orderProduct := range arg.OrderProducts {
			if orderProduct.Num <= 0 {
				err := fmt.Errorf("product num must be positive")
				return err
			}
			// 取得商品
			product, err := q.GetProduct(ctx, orderProduct.ID)
			if err != nil {
				fmt.Println("get product err", err)
				return err
			}
			totalPrice = totalPrice + (int32(orderProduct.Num) * product.OriginPrice)
			finalPrice = finalPrice + (int32(orderProduct.Num) * product.Price)
			productList = append(productList, OrderTxProductResult{
				Product: product,
				Num:     orderProduct.Num,
			})
		}

		if arg.CouponID != 0 {
			coupon, err := q.GetCoupon(context.Background(), arg.CouponID)
			if err != nil {
				return err
			}

			finalPrice = finalPrice * coupon.Percent / 100
		}

		initMessage := ""
		if arg.Message != "" {
			initMessage = arg.Message
		}

		orderArg := CreateOrderParams{
			FullName:        arg.FullName,
			Email:           arg.Email,
			ShippingAddress: arg.ShippingAddress,
			Message: pgtype.Text{
				String: initMessage,
				Valid:  true,
			},
			TotalPrice:  totalPrice,
			FinalPrice:  finalPrice,
			PayMethodID: int32(arg.PayMethodID),
			StatusID:    int32(arg.StatusID),
		}

		// 建立訂單
		result.Order, err = q.CreateOrder(ctx, orderArg)
		if err != nil {
			fmt.Println("create order", err)
			return err
		}

		// 建立訂單跟會員的關聯
		_, err = q.CreateOrderUser(ctx, CreateOrderUserParams{
			OrderID: pgtype.Int4{
				Int32: int32(result.Order.ID),
				Valid: true,
			},
			UserID: pgtype.Int4{
				Int32: int32(arg.UserID),
				Valid: true,
			},
		})
		if err != nil {
			fmt.Println("create order user err", err)
			return err
		}

		// 取得訂單狀態
		result.Status, err = q.GetOrderStatus(ctx, arg.StatusID)
		if err != nil {
			fmt.Println("get order status err", err)
			return err
		}

		for _, orderProduct := range arg.OrderProducts {
			// 建立訂單跟商品的關聯及數量
			_, err = q.CreateOrderProduct(ctx, CreateOrderProductParams{
				OrderID: pgtype.Int4{
					Int32: int32(result.Order.ID),
					Valid: true,
				},
				ProductID: pgtype.Int4{
					Int32: int32(orderProduct.ID),
					Valid: true,
				},
				Num: int32(orderProduct.Num),
			})
			if err != nil {
				fmt.Println("get order product err", err)
				return err
			}
		}

		result.ProductList = productList

		if arg.CouponID != 0 {
			_, err = q.CreateOrderCoupon(ctx, CreateOrderCouponParams{
				OrderID: pgtype.Int4{
					Int32: int32(result.Order.ID),
					Valid: true,
				},
				CouponID: pgtype.Int4{
					Int32: int32(arg.CouponID),
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}

type UpdateOrderTxParams struct {
	ID              int64                  `json:"id"`
	FullName        string                 `json:"full_name"`
	Email           string                 `json:"email"`
	ShippingAddress string                 `json:"shipping_address"`
	Message         string                 `json:"message"`
	TotalPrice      *int32                 `json:"total_price"`
	FinalPrice      *int32                 `json:"final_price"`
	PayMethodID     int64                  `json:"pay_method_id"`
	StatusID        int64                  `json:"status_id"`
	OrderProducts   []OrderTxProductParams `json:"order_products"`
	CouponID        int64                  `json:"coupon_id"`
}

func (store *SQLStore) UpdateOrderTx(ctx context.Context, arg UpdateOrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		orderArg := UpdateOrderParams{
			ID: arg.ID,
		}

		if arg.FullName != "" {
			if err := val.ValidateFullName(arg.FullName); err != nil {
				return err
			}

			orderArg.FullName = pgtype.Text{
				String: arg.FullName,
				Valid:  true,
			}
		}

		if arg.Email != "" {
			orderArg.Email = pgtype.Text{
				String: arg.Email,
				Valid:  true,
			}
		}

		if arg.Message != "" {
			orderArg.Message = pgtype.Text{
				String: arg.Message,
				Valid:  true,
			}
		}

		if arg.ShippingAddress != "" {
			orderArg.ShippingAddress = pgtype.Text{
				String: arg.ShippingAddress,
				Valid:  true,
			}
		}

		if arg.TotalPrice != nil {
			orderArg.TotalPrice = pgtype.Int4{
				Int32: *arg.TotalPrice,
				Valid: true,
			}
		}

		if arg.FinalPrice != nil {
			orderArg.FinalPrice = pgtype.Int4{
				Int32: *arg.FinalPrice,
				Valid: true,
			}
		}

		if arg.PayMethodID != 0 {
			orderArg.PayMethodID = pgtype.Int4{
				Int32: int32(arg.PayMethodID),
				Valid: true,
			}
		}

		if arg.StatusID != 0 {
			orderArg.StatusID = pgtype.Int4{
				Int32: int32(arg.StatusID),
				Valid: true,
			}
		}

		// 更新訂單
		result.Order, err = q.UpdateOrder(ctx, orderArg)
		if err != nil {
			return err
		}

		// 取得訂單狀態
		result.Status, err = q.GetOrderStatus(ctx, int64(result.Order.StatusID))
		if err != nil {
			return err
		}

		var productList []OrderTxProductResult
		if arg.OrderProducts != nil && len(arg.OrderProducts) > 0 {
			// 如果需要更新訂單商品，先把之前建立的關聯移除
			err = q.DeleteOrderProductByOrderId(ctx, pgtype.Int4{
				Int32: int32(result.Order.ID),
				Valid: true,
			})
			if err != nil {
				return err
			}

			for _, orderProduct := range arg.OrderProducts {
				// 取得商品
				product, err := q.GetProduct(ctx, orderProduct.ID)
				if err != nil {
					return err
				}

				// 建立訂單跟商品的關聯及數量
				_, err = q.CreateOrderProduct(ctx, CreateOrderProductParams{
					OrderID: pgtype.Int4{
						Int32: int32(result.Order.ID),
						Valid: true,
					},
					ProductID: pgtype.Int4{
						Int32: int32(product.ID),
						Valid: true,
					},
					Num: int32(orderProduct.Num),
				})
				if err != nil {
					return err
				}

				productList = append(productList, OrderTxProductResult{
					Product: product,
					Num:     orderProduct.Num,
				})
			}
		} else {
			// 如果沒有更新訂單商品，還是要去取得目前的商品
			orderProducts, err := q.ListOrderProductByOrderId(ctx, pgtype.Int4{
				Int32: int32(result.Order.ID),
				Valid: true,
			})
			if err != nil {
				return err
			}

			for _, orderProduct := range orderProducts {
				product, err := q.GetProduct(ctx, int64(orderProduct.ProductID.Int32))
				if err != nil {
					return err
				}
				productList = append(productList, OrderTxProductResult{
					Product: product,
					Num:     int64(orderProduct.Num),
				})
			}
		}

		result.ProductList = productList

		if arg.CouponID != 0 {
			// 如果需要更改，先將之前建立的關聯移除
			err := q.DeleteOrderCouponByOrderId(ctx, pgtype.Int4{
				Int32: int32(result.Order.ID),
				Valid: true,
			})
			if err != nil {
				return err
			}

			_, err = q.CreateOrderCoupon(ctx, CreateOrderCouponParams{
				OrderID: pgtype.Int4{
					Int32: int32(result.Order.ID),
					Valid: true,
				},
				CouponID: pgtype.Int4{
					Int32: int32(arg.CouponID),
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}

type DeleteOrderTxParams struct {
	ID int64 `json:"id"`
}

func (store *SQLStore) DeleteOrderTx(ctx context.Context, arg DeleteOrderTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 建立訂單跟會員的關聯
		err = q.DeleteOrderUserByOrderId(ctx, pgtype.Int4{
			Int32: int32(arg.ID),
			Valid: true,
		})

		if err != nil {
			return err
		}

		err = q.DeleteOrderProductByOrderId(ctx, pgtype.Int4{
			Int32: int32(arg.ID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		err = q.DeleteOrderCouponByOrderId(ctx, pgtype.Int4{
			Int32: int32(arg.ID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		err = q.DeleteOrder(ctx, arg.ID)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
