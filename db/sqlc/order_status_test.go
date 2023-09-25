package db

import (
	"context"
	"testing"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomOrderStatus(t *testing.T) OrderStatus {
	name := util.RandomName()
	description := util.RandomString(20)

	arg := CreateOrderStatusParams{
		Name:        name,
		Description: description,
	}

	orderStatus, err := testStore.CreateOrderStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, orderStatus)

	require.NotZero(t, orderStatus.ID)
	require.Equal(t, name, orderStatus.Name)
	require.Equal(t, description, orderStatus.Description)
	return orderStatus
}

func TestCreateOrderStatus(t *testing.T) {
	createRandomOrderStatus(t)
}

func TestGetOrderStatus(t *testing.T) {
	orderStatus1 := createRandomOrderStatus(t)

	orderStatus2, err := testStore.GetOrderStatus(context.Background(), orderStatus1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, orderStatus2)

	require.Equal(t, orderStatus1.ID, orderStatus2.ID)
	require.Equal(t, orderStatus1.Name, orderStatus2.Name)
	require.Equal(t, orderStatus1.Description, orderStatus2.Description)
}

func TestListOrderStatusAndCount(t *testing.T) {
	createRandomOrderStatus(t)

	orderStatusList, err := testStore.ListOrderStatus(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, orderStatusList)

	count, err := testStore.GetOrderStatusCount(context.Background())
	require.NoError(t, err)

	require.Equal(t, count, int64(len(orderStatusList)))
}
