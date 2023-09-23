package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderCoupon(t *testing.T) OrderCoupon {
	coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	order := createRandomOrder(t)

	arg := CreateOrderCouponParams{
		OrderID: pgtype.Int4{
			Int32: int32(order.ID),
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: int32(coupon.ID),
			Valid: true,
		},
	}

	orderCoupon, err := testStore.CreateOrderCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderCoupon)

	require.NotZero(t, orderCoupon.OrderID)
	require.NotZero(t, orderCoupon.CouponID)

	return orderCoupon
}

func TestCreateOrderCoupon(t *testing.T) {
	createRandomOrderCoupon(t)
}

func TestGetOrderCoupon(t *testing.T) {
	orderCoupon1 := createRandomOrderCoupon(t)

	arg := GetOrderCouponParams{
		OrderID: pgtype.Int4{
			Int32: orderCoupon1.OrderID.Int32,
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: orderCoupon1.CouponID.Int32,
			Valid: true,
		},
	}

	orderCoupon2, err := testStore.GetOrderCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderCoupon2)
}

func TestDeleteOrderCouponByOrderId(t *testing.T) {
	orderCoupon1 := createRandomOrderCoupon(t)

	err := testStore.DeleteOrderCouponByOrderId(context.Background(), orderCoupon1.OrderID)
	require.NoError(t, err)

	arg := GetOrderCouponParams{
		OrderID: pgtype.Int4{
			Int32: orderCoupon1.OrderID.Int32,
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: orderCoupon1.CouponID.Int32,
			Valid: true,
		},
	}

	orderCoupon2, err := testStore.GetOrderCoupon(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, orderCoupon2)
}

func TestDeleteOrderCouponByCouponId(t *testing.T) {
	orderCoupon1 := createRandomOrderCoupon(t)

	err := testStore.DeleteOrderCouponByCouponId(context.Background(), orderCoupon1.CouponID)
	require.NoError(t, err)

	arg := GetOrderCouponParams{
		OrderID: pgtype.Int4{
			Int32: orderCoupon1.OrderID.Int32,
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: orderCoupon1.CouponID.Int32,
			Valid: true,
		},
	}

	orderCoupon2, err := testStore.GetOrderCoupon(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, orderCoupon2)
}

func TestListOrderCouponByOrderId(t *testing.T) {
	order := createRandomOrder(t)
	for i := 0; i < 5; i++ {
		coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

		arg := CreateOrderCouponParams{
			OrderID: pgtype.Int4{
				Int32: int32(order.ID),
				Valid: true,
			},
			CouponID: pgtype.Int4{
				Int32: int32(coupon.ID),
				Valid: true,
			},
		}
		testStore.CreateOrderCoupon(context.Background(), arg)
	}

	orderCoupons, err := testStore.ListOrderCouponByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(order.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderCoupons, 5)

	for _, orderCoupon := range orderCoupons {
		require.NotEmpty(t, orderCoupon)
	}
}

func TestListOrderCouponByCouponId(t *testing.T) {
	coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	for i := 0; i < 5; i++ {
		order := createRandomOrder(t)

		arg := CreateOrderCouponParams{
			OrderID: pgtype.Int4{
				Int32: int32(order.ID),
				Valid: true,
			},
			CouponID: pgtype.Int4{
				Int32: int32(coupon.ID),
				Valid: true,
			},
		}
		testStore.CreateOrderCoupon(context.Background(), arg)
	}

	orderCoupons, err := testStore.ListOrderCouponByCouponId(context.Background(), pgtype.Int4{
		Int32: int32(coupon.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderCoupons, 5)

	for _, orderCoupon := range orderCoupons {
		require.NotEmpty(t, orderCoupon)
	}
}
