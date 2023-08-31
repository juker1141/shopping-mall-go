package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomCart(t *testing.T) Cart {
	user := createRandomUser(t)

	arg := CreateCartParams{
		Owner: pgtype.Text{
			String: user.Account,
			Valid:  true,
		},
		TotalPrice: util.RandomPrice(),
		FinalPrice: util.RandomPrice(),
	}

	cart, err := testStore.CreateCart(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cart)

	require.Equal(t, arg.Owner, cart.Owner)
	require.Equal(t, arg.TotalPrice, cart.TotalPrice)
	require.Equal(t, arg.FinalPrice, cart.FinalPrice)

	require.NotZero(t, cart.CreatedAt)

	return cart
}

func TestCreateCart(t *testing.T) {
	createRandomCart(t)
}

func TestGetCart(t *testing.T) {
	cart1 := createRandomCart(t)

	cart2, err := testStore.GetCart(context.Background(), cart1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cart2)

	require.Equal(t, cart1.ID, cart2.ID)
	require.Equal(t, cart1.Owner, cart2.Owner)
	require.Equal(t, cart1.TotalPrice, cart2.TotalPrice)
	require.Equal(t, cart1.FinalPrice, cart2.FinalPrice)

	require.WithinDuration(t, cart1.CreatedAt, cart2.CreatedAt, time.Second)
}

func TestGetCartByOwner(t *testing.T) {
	cart1 := createRandomCart(t)

	cart2, err := testStore.GetCartByOwner(context.Background(), cart1.Owner)
	require.NoError(t, err)
	require.NotEmpty(t, cart2)

	require.Equal(t, cart1.ID, cart2.ID)
	require.Equal(t, cart1.Owner, cart2.Owner)
	require.Equal(t, cart1.TotalPrice, cart2.TotalPrice)
	require.Equal(t, cart1.FinalPrice, cart2.FinalPrice)

	require.WithinDuration(t, cart1.CreatedAt, cart2.CreatedAt, time.Second)
}

func TestUpdateCart(t *testing.T) {
	oldCart := createRandomCart(t)

	arg := UpdateCartParams{
		ID:         oldCart.ID,
		TotalPrice: util.RandomPrice(),
		FinalPrice: util.RandomPrice(),
	}
	newCart, err := testStore.UpdateCart(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCart)

	require.Equal(t, oldCart.ID, newCart.ID)
	require.Equal(t, oldCart.Owner, newCart.Owner)
	require.WithinDuration(t, oldCart.CreatedAt, newCart.CreatedAt, time.Second)

	require.NotEqual(t, oldCart.TotalPrice, newCart.TotalPrice)
	require.NotEqual(t, oldCart.FinalPrice, newCart.FinalPrice)
}

func TestDeleteCart(t *testing.T) {
	cart1 := createRandomCart(t)

	err := testStore.DeleteCart(context.Background(), cart1.ID)
	require.NoError(t, err)

	cart2, err := testStore.GetCart(context.Background(), cart1.ID)
	require.Error(t, err)
	require.Empty(t, cart2)
}
