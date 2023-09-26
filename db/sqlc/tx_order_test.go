package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderTx(t *testing.T) OrderTxResult {
	user := createRandomUser(t)
	payMethod := createRandomPayMethod(t)
	status := createRandomOrderStatus(t)
	coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	n := 3
	orderProducts := make([]OrderTxProductParams, n)
	productList := make([]OrderTxProductResult, n)
	for i := 0; i < n; i++ {
		product := createRandomProduct(t, util.RandomName())
		num := util.RandomInt(1, 10)
		orderProducts[i] = OrderTxProductParams{
			ID:  product.ID,
			Num: num,
		}
		productList[i] = OrderTxProductResult{
			Product: product,
			Num:     num,
		}
	}

	arg := CreateOrderTxParams{
		UserID:          user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		ShippingAddress: user.ShippingAddress,
		Message:         util.RandomString(20),
		PayMethodID:     payMethod.ID,
		StatusID:        status.ID,
		OrderProducts:   orderProducts,
		CouponID:        coupon.ID,
	}

	result, err := testStore.CreateOrderTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.FullName, result.FullName)
	require.Equal(t, user.Email, result.Email)
	require.Equal(t, user.ShippingAddress, result.ShippingAddress)
	require.NotZero(t, result.TotalPrice)
	require.NotZero(t, result.FinalPrice)
	require.Equal(t, arg.Message, result.Message.String)
	require.Equal(t, productList, result.ProductList)
	require.Equal(t, int32(payMethod.ID), result.PayMethodID)
	require.Equal(t, int32(status.ID), result.StatusID)

	orderUser, err := testStore.ListOrderUserByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderUser, 1)
	require.Equal(t, result.ID, int64(orderUser[0].OrderID.Int32))
	require.Equal(t, user.ID, int64(orderUser[0].UserID.Int32))

	orderProduct, err := testStore.ListOrderProductByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderProduct, n)

	orderCoupon, err := testStore.GetOrderCoupon(context.Background(), GetOrderCouponParams{
		OrderID: pgtype.Int4{
			Int32: int32(result.ID),
			Valid: true,
		},
		CouponID: pgtype.Int4{
			Int32: int32(coupon.ID),
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, orderCoupon)

	return result
}

func TestCreateOrderTx(t *testing.T) {
	createRandomOrderTx(t)
}

func TestCreateOrderTxEmptyMessageEmptyCoupon(t *testing.T) {
	user := createRandomUser(t)
	payMethod := createRandomPayMethod(t)
	status := createRandomOrderStatus(t)

	n := 3
	orderProducts := make([]OrderTxProductParams, n)
	productList := make([]OrderTxProductResult, n)
	for i := 0; i < n; i++ {
		product := createRandomProduct(t, util.RandomName())
		num := util.RandomInt(1, 10)

		orderProducts[i] = OrderTxProductParams{
			ID:  product.ID,
			Num: num,
		}
		productList[i] = OrderTxProductResult{
			Product: product,
			Num:     num,
		}
	}

	arg := CreateOrderTxParams{
		UserID:          user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		ShippingAddress: user.ShippingAddress,
		PayMethodID:     payMethod.ID,
		StatusID:        status.ID,
		OrderProducts:   orderProducts,
	}

	result, err := testStore.CreateOrderTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.FullName, result.FullName)
	require.Equal(t, user.Email, result.Email)
	require.Equal(t, user.ShippingAddress, result.ShippingAddress)
	require.NotZero(t, result.TotalPrice)
	require.NotZero(t, result.FinalPrice)
	require.Empty(t, result.Message.String)
	require.Equal(t, productList, result.ProductList)
	require.Equal(t, int32(payMethod.ID), result.PayMethodID)
	require.Equal(t, int32(status.ID), result.StatusID)

	orderUser, err := testStore.ListOrderUserByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderUser, 1)
	require.Equal(t, result.ID, int64(orderUser[0].OrderID.Int32))
	require.Equal(t, user.ID, int64(orderUser[0].UserID.Int32))

	orderProduct, err := testStore.ListOrderProductByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderProduct, n)
}

func testUpdateOrderTxResult(t *testing.T, oldResult, newResult OrderTxResult) {
	orderUser, err := testStore.ListOrderUserByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(newResult.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderUser, 1)
	require.Equal(t, int32(oldResult.ID), orderUser[0].OrderID.Int32)
	require.NotZero(t, orderUser[0].UserID.Int32)

	orderProduct, err := testStore.ListOrderProductByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(newResult.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Equal(t, len(newResult.ProductList), len(orderProduct))

	orderCoupon, err := testStore.ListOrderCouponByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(newResult.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.NotEmpty(t, orderCoupon)
	require.Len(t, orderCoupon, 1)
}

func TestUpdateOrderTx(t *testing.T) {
	oldResult := createRandomOrderTx(t)
	newPayMethod := createRandomPayMethod(t)
	newStatus := createRandomOrderStatus(t)
	coupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	n := 3
	orderProducts := make([]OrderTxProductParams, n)
	productList := make([]OrderTxProductResult, n)
	for i := 0; i < n; i++ {
		product := createRandomProduct(t, util.RandomName())
		num := util.RandomInt(1, 10)

		orderProducts[i] = OrderTxProductParams{
			ID:  product.ID,
			Num: num,
		}
		productList[i] = OrderTxProductResult{
			Product: product,
			Num:     num,
		}
	}

	arg := UpdateOrderTxParams{
		ID:              oldResult.ID,
		FullName:        util.RandomName(),
		Email:           util.RandomEmail(),
		ShippingAddress: util.RandomAddress(),
		Message:         util.RandomString(10),
		PayMethodID:     newPayMethod.ID,
		StatusID:        newStatus.ID,
		OrderProducts:   orderProducts,
		CouponID:        coupon.ID,
	}

	newResult, err := testStore.UpdateOrderTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newResult)

	require.NotEqual(t, oldResult.FullName, newResult.FullName)
	require.NotEqual(t, oldResult.Email, newResult.Email)
	require.NotEqual(t, oldResult.ShippingAddress, newResult.ShippingAddress)
	require.NotEqual(t, oldResult.TotalPrice, newResult.TotalPrice)
	require.NotEqual(t, oldResult.FinalPrice, newResult.FinalPrice)
	require.NotEqual(t, oldResult.Message.String, newResult.Message.String)
	require.NotEqual(t, oldResult.ProductList, newResult.ProductList)
	require.Equal(t, int32(newPayMethod.ID), newResult.PayMethodID)
	require.Equal(t, int32(newStatus.ID), newResult.StatusID)

	testUpdateOrderTxResult(t, oldResult, newResult)
}

func TestUpdateOrderTxAllFieldEmpty(t *testing.T) {
	oldResult := createRandomOrderTx(t)

	arg := UpdateOrderTxParams{
		ID: oldResult.ID,
	}

	newResult, err := testStore.UpdateOrderTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newResult)

	require.Equal(t, oldResult.FullName, newResult.FullName)
	require.Equal(t, oldResult.Email, newResult.Email)
	require.Equal(t, oldResult.ShippingAddress, newResult.ShippingAddress)
	require.Equal(t, oldResult.TotalPrice, newResult.TotalPrice)
	require.Equal(t, oldResult.FinalPrice, newResult.FinalPrice)
	require.Equal(t, oldResult.Message.String, newResult.Message.String)
	require.Equal(t, oldResult.ProductList, newResult.ProductList)
	require.Equal(t, oldResult.PayMethodID, newResult.PayMethodID)
	require.Equal(t, oldResult.StatusID, newResult.StatusID)

	testUpdateOrderTxResult(t, oldResult, newResult)
}

func TestDeleteOrderTx(t *testing.T) {
	result := createRandomOrderTx(t)

	arg := DeleteOrderTxParams{
		ID: result.ID,
	}

	err := testStore.DeleteOrderTx(context.Background(), arg)
	require.NoError(t, err)

	order, err := testStore.GetOrder(context.Background(), result.ID)
	require.Error(t, err)
	require.Empty(t, order)

	orderUser, err := testStore.ListOrderUserByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderUser, 0)

	orderProduct, err := testStore.ListOrderProductByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderProduct, 0)

	orderCoupon, err := testStore.ListOrderCouponByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(result.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderCoupon, 0)
}
