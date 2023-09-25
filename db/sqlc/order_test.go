package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T) Order {
	payMethod := createRandomPayMethod(t)
	orderStatus := createRandomOrderStatus(t)

	fullName := util.RandomName()
	email := util.RandomEmail()
	shippingAddress := util.RandomAddress()
	message := util.RandomString(10)

	arg := CreateOrderParams{
		FullName:        fullName,
		Email:           email,
		ShippingAddress: shippingAddress,
		Message: pgtype.Text{
			String: message,
			Valid:  true,
		},
		PayMethodID: int32(payMethod.ID),
		StatusID:    int32(orderStatus.ID),
	}

	order, err := testStore.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.NotZero(t, order.ID)
	require.Equal(t, fullName, order.FullName)
	require.Equal(t, email, order.Email)
	require.Equal(t, shippingAddress, order.ShippingAddress)
	require.Equal(t, message, order.Message.String)
	require.False(t, order.IsPaid)
	require.Equal(t, int32(payMethod.ID), order.PayMethodID)
	require.Equal(t, int32(orderStatus.ID), order.StatusID)

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
	require.Equal(t, order1.FullName, order2.FullName)
	require.Equal(t, order1.Email, order2.Email)
	require.Equal(t, order1.ShippingAddress, order2.ShippingAddress)
	require.Equal(t, order1.IsPaid, order2.IsPaid)
	require.Equal(t, order1.PayMethodID, order2.PayMethodID)
	require.Equal(t, order1.StatusID, order2.StatusID)
	require.Equal(t, order1.Message.String, order2.Message.String)

	require.WithinDuration(t, order1.CreatedAt, order2.CreatedAt, time.Second)
	require.WithinDuration(t, order1.UpdatedAt, order2.UpdatedAt, time.Second)
}

func TestUpdateOrder(t *testing.T) {
	oldOrder := createRandomOrder(t)
	newOrderStatus := createRandomOrderStatus(t)
	newPayMethod := createRandomPayMethod(t)
	newBool := true
	newUpdatedTime := time.Now()

	arg := UpdateOrderParams{
		ID: oldOrder.ID,
		FullName: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
		Email: pgtype.Text{
			String: util.RandomEmail(),
			Valid:  true,
		},
		ShippingAddress: pgtype.Text{
			String: util.RandomAddress(),
			Valid:  true,
		},
		Message: pgtype.Text{
			String: util.RandomString(10),
			Valid:  true,
		},
		IsPaid: pgtype.Bool{
			Bool:  newBool,
			Valid: true,
		},
		PayMethodID: pgtype.Int4{
			Int32: int32(newPayMethod.ID),
			Valid: true,
		},
		StatusID: pgtype.Int4{
			Int32: int32(newOrderStatus.ID),
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
	require.NotEqual(t, oldOrder.FullName, newOrder.FullName)
	require.NotEqual(t, oldOrder.Email, newOrder.Email)
	require.NotEqual(t, oldOrder.ShippingAddress, newOrder.ShippingAddress)
	require.NotEqual(t, oldOrder.IsPaid, newOrder.IsPaid)
	require.NotEqual(t, oldOrder.Message.String, newOrder.Message.String)
	require.NotEqual(t, oldOrder.StatusID, newOrder.StatusID)
	require.NotEqual(t, oldOrder.IsPaid, newOrder.IsPaid)

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
