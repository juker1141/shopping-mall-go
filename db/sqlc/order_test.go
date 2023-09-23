package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T) Order {
	orderStatus := createRandomOrderStatus(t)
	statusID := pgtype.Int4{
		Int32: int32(orderStatus.ID),
		Valid: true,
	}
	order, err := testStore.CreateOrder(context.Background(), statusID)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.NotZero(t, order.ID)
	require.Equal(t, statusID.Int32, order.StatusID.Int32)
	require.False(t, order.IsPaid.Bool)

	require.NotZero(t, order.CreatedAt)
	require.NotZero(t, order.UpdatedAt)

	return order
}

func TestCreateOrder(t *testing.T) {
	createRandomOrder(t)
}

func TestGetOrder(t *testing.T) {
	order1 := createRandomOrder(t)

	order2, err := testStore.GetOrder(context.Background(), order1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.ID, order2.ID)
	require.Equal(t, order1.StatusID.Int32, order2.StatusID.Int32)
	require.Equal(t, order1.IsPaid.Bool, order2.IsPaid.Bool)

	require.WithinDuration(t, order1.CreatedAt, order2.CreatedAt, time.Second)
	require.WithinDuration(t, order1.UpdatedAt, order2.UpdatedAt, time.Second)
}

func TestUpdateOrder(t *testing.T) {
	oldOrder := createRandomOrder(t)
	newOrderStatus := createRandomOrderStatus(t)
	newBool := true
	newUpdatedTime := time.Now()

	arg := UpdateOrderParams{
		ID: oldOrder.ID,
		StatusID: pgtype.Int4{
			Int32: int32(newOrderStatus.ID),
			Valid: true,
		},
		IsPaid: pgtype.Bool{
			Bool:  newBool,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  newUpdatedTime,
			Valid: true,
		},
	}

	newOrder, err := testStore.UpdateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newOrder)

	require.Equal(t, oldOrder.ID, newOrder.ID)
	require.NotEqual(t, oldOrder.StatusID.Int32, newOrder.StatusID.Int32)
	require.NotEqual(t, oldOrder.IsPaid.Bool, newOrder.IsPaid.Bool)

	require.WithinDuration(t, newUpdatedTime, newOrder.UpdatedAt, time.Second)
	require.WithinDuration(t, oldOrder.CreatedAt, newOrder.CreatedAt, time.Second)
}

func TestDeleteOrder(t *testing.T) {
	order1 := createRandomOrder(t)

	err := testStore.DeleteOrder(context.Background(), order1.ID)
	require.NoError(t, err)

	order2, err := testStore.GetOrder(context.Background(), order1.ID)
	require.Error(t, err)
	require.Empty(t, order2)
}

func TestListOrder(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		createRandomOrder(t)
	}

	arg := ListOrdersParams{
		Limit:  5,
		Offset: 5,
	}

	orders, err := testStore.ListOrders(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orders)

	for _, order := range orders {
		require.NotEmpty(t, order)
	}
}
