package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomOrderUser(t *testing.T) OrderUser {
	user := createRandomUser(t)
	order := createRandomOrder(t)

	arg := CreateOrderUserParams{
		OrderID: pgtype.Int4{
			Int32: int32(order.ID),
			Valid: true,
		},
		UserID: pgtype.Int4{
			Int32: int32(user.ID),
			Valid: true,
		},
	}

	orderUser, err := testStore.CreateOrderUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderUser)

	require.NotZero(t, orderUser.OrderID)
	require.NotZero(t, orderUser.UserID)

	return orderUser
}

func TestCreateOrderUser(t *testing.T) {
	createRandomOrderUser(t)
}

func TestGetOrderUser(t *testing.T) {
	orderUser1 := createRandomOrderUser(t)

	arg := GetOrderUserParams{
		OrderID: pgtype.Int4{
			Int32: orderUser1.OrderID.Int32,
			Valid: true,
		},
		UserID: pgtype.Int4{
			Int32: orderUser1.UserID.Int32,
			Valid: true,
		},
	}

	orderUser2, err := testStore.GetOrderUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderUser2)
}

func TestDeleteOrderUserByOrderId(t *testing.T) {
	orderUser1 := createRandomOrderUser(t)

	err := testStore.DeleteOrderUserByOrderId(context.Background(), orderUser1.OrderID)
	require.NoError(t, err)

	arg := GetOrderUserParams{
		OrderID: pgtype.Int4{
			Int32: orderUser1.OrderID.Int32,
			Valid: true,
		},
		UserID: pgtype.Int4{
			Int32: orderUser1.UserID.Int32,
			Valid: true,
		},
	}

	orderUser2, err := testStore.GetOrderUser(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, orderUser2)
}

func TestDeleteOrderUserByUserId(t *testing.T) {
	orderUser1 := createRandomOrderUser(t)

	err := testStore.DeleteOrderUserByUserId(context.Background(), orderUser1.UserID)
	require.NoError(t, err)

	arg := GetOrderUserParams{
		OrderID: pgtype.Int4{
			Int32: orderUser1.OrderID.Int32,
			Valid: true,
		},
		UserID: pgtype.Int4{
			Int32: orderUser1.UserID.Int32,
			Valid: true,
		},
	}

	orderUser2, err := testStore.GetOrderUser(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, orderUser2)
}

func TestListOrderUserByOrderId(t *testing.T) {
	order := createRandomOrder(t)
	for i := 0; i < 5; i++ {
		user := createRandomUser(t)

		arg := CreateOrderUserParams{
			OrderID: pgtype.Int4{
				Int32: int32(order.ID),
				Valid: true,
			},
			UserID: pgtype.Int4{
				Int32: int32(user.ID),
				Valid: true,
			},
		}
		testStore.CreateOrderUser(context.Background(), arg)
	}

	orderUsers, err := testStore.ListOrderUserByOrderId(context.Background(), pgtype.Int4{
		Int32: int32(order.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderUsers, 5)

	for _, orderUser := range orderUsers {
		require.NotEmpty(t, orderUser)
	}
}

func TestListOrderUserByUserId(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 5; i++ {
		order := createRandomOrder(t)

		arg := CreateOrderUserParams{
			OrderID: pgtype.Int4{
				Int32: int32(order.ID),
				Valid: true,
			},
			UserID: pgtype.Int4{
				Int32: int32(user.ID),
				Valid: true,
			},
		}
		testStore.CreateOrderUser(context.Background(), arg)
	}

	orderUsers, err := testStore.ListOrderUserByUserId(context.Background(), pgtype.Int4{
		Int32: int32(user.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, orderUsers, 5)

	for _, orderUser := range orderUsers {
		require.NotEmpty(t, orderUser)
	}
}
