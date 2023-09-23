package db

import (
	"context"
	"testing"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomPayMethod(t *testing.T) PayMethod {
	name := util.RandomName()

	payMethod, err := testStore.CreatePayMethod(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, payMethod)

	require.NotZero(t, payMethod.ID)
	require.Equal(t, name, payMethod.Name)
	return payMethod
}

func TestCreatePayMethod(t *testing.T) {
	createRandomPayMethod(t)
}

func TestGetPayMethod(t *testing.T) {
	payMethod1 := createRandomPayMethod(t)

	payMethod2, err := testStore.GetPayMethod(context.Background(), payMethod1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, payMethod2)

	require.Equal(t, payMethod1.ID, payMethod2.ID)
	require.Equal(t, payMethod1.Name, payMethod2.Name)
}

func TestListPayMethodAndCount(t *testing.T) {
	createRandomPayMethod(t)

	payMethodList, err := testStore.ListPayMethod(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, payMethodList)

	count, err := testStore.GetPayMethodCount(context.Background())
	require.NoError(t, err)

	require.Equal(t, count, int64(len(payMethodList)))
}
