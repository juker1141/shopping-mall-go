package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomCartCoupon(t *testing.T) CartCoupon {
	cart := createRandomCart(t)
	coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	arg := CreateCartCouponParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: int32(coupon.ID),
			Valid: true,
		},
	}

	cartCoupon, err := testStore.CreateCartCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cartCoupon)

	require.NotZero(t, cartCoupon.CartID)
	require.NotZero(t, cartCoupon.CouponID)

	return cartCoupon
}

func TestCreateCartCoupon(t *testing.T) {
	createRandomCartCoupon(t)
}

func TestGetCartCoupon(t *testing.T) {
	cartCoupon1 := createRandomCartCoupon(t)

	arg := GetCartCouponParams{
		CartID: pgtype.Int4{
			Int32: cartCoupon1.CartID.Int32,
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: cartCoupon1.CouponID.Int32,
			Valid: true,
		},
	}

	cartCoupon2, err := testStore.GetCartCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cartCoupon2)

	require.Equal(t, cartCoupon1.CartID, cartCoupon2.CartID)
	require.Equal(t, cartCoupon1.CouponID, cartCoupon2.CouponID)
}

func TestDeleteCartCouponByCartId(t *testing.T) {
	cartCoupon1 := createRandomCartCoupon(t)

	err := testStore.DeleteCartCouponByCartId(context.Background(), cartCoupon1.CartID)
	require.NoError(t, err)

	cartCoupon2, err := testStore.GetCartCoupon(context.Background(), GetCartCouponParams{
		CartID: pgtype.Int4{
			Int32: cartCoupon1.CartID.Int32,
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: cartCoupon1.CouponID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.Empty(t, cartCoupon2)
}

func TestDeleteCartCouponByCouponId(t *testing.T) {
	cartCoupon1 := createRandomCartCoupon(t)

	err := testStore.DeleteCartCouponByCouponId(context.Background(), cartCoupon1.CouponID)
	require.NoError(t, err)

	cartCoupon2, err := testStore.GetCartCoupon(context.Background(), GetCartCouponParams{
		CartID: pgtype.Int4{
			Int32: cartCoupon1.CartID.Int32,
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: cartCoupon1.CouponID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.Empty(t, cartCoupon2)
}

func TestListCartCouponByCartId(t *testing.T) {
	cart := createRandomCart(t)
	for i := 0; i < 5; i++ {
		coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

		testStore.CreateCartCoupon(context.Background(), CreateCartCouponParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			CouponID: pgtype.Int4{
				Int32: int32(coupon.ID),
				Valid: true,
			},
		})
	}

	cartCoupons, err := testStore.ListCartCouponByCartId(context.Background(), pgtype.Int4{
		Int32: int32(cart.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, cartCoupons, 5)

	for _, cartCoupon := range cartCoupons {
		require.NotEmpty(t, cartCoupon)
	}
}

func TestListCartCouponByCouponId(t *testing.T) {
	coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	for i := 0; i < 5; i++ {
		cart := createRandomCart(t)

		testStore.CreateCartCoupon(context.Background(), CreateCartCouponParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			CouponID: pgtype.Int4{
				Int32: int32(coupon.ID),
				Valid: true,
			},
		})
	}

	cartCoupons, err := testStore.ListCartCouponByCouponId(context.Background(), pgtype.Int4{
		Int32: int32(coupon.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, cartCoupons, 5)

	for _, cartCoupon := range cartCoupons {
		require.NotEmpty(t, cartCoupon)
	}
}
