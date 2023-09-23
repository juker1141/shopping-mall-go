package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderProduct(t *testing.T) OrderProduct {
	order := createRandomOrder(t)
	product := createRandomProduct(t, util.RandomName())

	arg := CreateOrderProductParams{
		OrderID: pgtype.Int4{
			Int32: int32(order.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
	}

	orderProduct, err := testStore.CreateOrderProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderProduct)

	require.NotZero(t, orderProduct.OrderID)
	require.NotZero(t, orderProduct.ProductID)

	return orderProduct
}

func TestCreateOrderProduct(t *testing.T) {
	createRandomOrderProduct(t)
}

func TestGetOrderProduct(t *testing.T) {
	orderProduct1 := createRandomOrderProduct(t)

	arg := GetOrderProductParams{
		OrderID: pgtype.Int4{
			Int32: orderProduct1.OrderID.Int32,
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: orderProduct1.ProductID.Int32,
			Valid: true,
		},
	}

	orderProduct2, err := testStore.GetOrderProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderProduct2)
}

func TestDeleteOrderProductByOrderId(t *testing.T) {
	orderProduct1 := createRandomOrderProduct(t)

	err := testStore.DeleteOrderProductByOrderId(context.Background(), orderProduct1.OrderID)
	require.NoError(t, err)

	arg := GetOrderProductParams{
		OrderID: pgtype.Int4{
			Int32: orderProduct1.OrderID.Int32,
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: orderProduct1.ProductID.Int32,
			Valid: true,
		},
	}

	orderProduct2, err := testStore.GetOrderProduct(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, orderProduct2)
}

func TestDeleteOrderProductByProductId(t *testing.T) {
	orderProduct1 := createRandomOrderProduct(t)

	err := testStore.DeleteOrderProductByProductId(context.Background(), orderProduct1.ProductID)
	require.NoError(t, err)

	arg := GetOrderProductParams{
		OrderID: pgtype.Int4{
			Int32: orderProduct1.OrderID.Int32,
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: orderProduct1.ProductID.Int32,
			Valid: true,
		},
	}

	orderProduct2, err := testStore.GetOrderProduct(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, orderProduct2)
}

func TestListOrderProductByOrderId(t *testing.T) {
	order := createRandomOrder(t)
	for i := 0; i < 5; i++ {
		product := createRandomProduct(t, util.RandomName())

		arg := CreateOrderProductParams{
			OrderID: pgtype.Int4{
				Int32: int32(order.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(product.ID),
				Valid: true,
			},
		}
		testStore.CreateOrderProduct(context.Background(), arg)
	}

	orderProducts, err := testStore.ListOrderProductByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(order.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderProducts, 5)

	for _, orderProduct := range orderProducts {
		require.NotEmpty(t, orderProduct)
	}
}

func TestListOrderProductByProductId(t *testing.T) {
	product := createRandomProduct(t, util.RandomName())
	for i := 0; i < 5; i++ {
		order := createRandomOrder(t)

		arg := CreateOrderProductParams{
			OrderID: pgtype.Int4{
				Int32: int32(order.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(product.ID),
				Valid: true,
			},
		}
		testStore.CreateOrderProduct(context.Background(), arg)
	}

	orderProducts, err := testStore.ListOrderProductByProductId(context.Background(), pgtype.Int4{
		Int32: int32(product.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderProducts, 5)

	for _, orderProduct := range orderProducts {
		require.NotEmpty(t, orderProduct)
	}
}
