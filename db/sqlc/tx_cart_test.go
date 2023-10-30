package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func TestUpdateCartTxWithUpdateType(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t, util.RandomString(8))
	num := 5

	cartArg := CreateCartParams{
		Owner: pgtype.Text{
			String: user.Account,
			Valid:  true,
		},
		TotalPrice: 0,
		FinalPrice: 0,
	}

	cart, err := testStore.CreateCart(context.Background(), cartArg)
	require.NoError(t, err)
	require.NotEmpty(t, cart)

	_, err = testStore.CreateCartProduct(context.Background(), CreateCartProductParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
		Num: int32(num),
	})
	require.NoError(t, err)

	arg := UpdateCartTxParams{
		CartID: cart.ID,
	}

	result, err := testStore.UpdateCartTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, result.Cart.ID, cart.ID)
	require.Equal(t, result.Cart.Owner.String, user.Account)
	require.NotEqual(t, result.Cart.TotalPrice, cart.TotalPrice)
	require.NotEqual(t, result.Cart.FinalPrice, cart.FinalPrice)

	require.Equal(t, result.ProductList, []CartTxProductResult{
		{
			Product: product,
			Num:     int64(num),
		},
	})
}

func TestUpdateCartTxWithAddType(t *testing.T) {
	user := createRandomUser(t)
	product := createRandomProduct(t, util.RandomString(8))
	num := 5

	cartArg := CreateCartParams{
		Owner: pgtype.Text{
			String: user.Account,
			Valid:  true,
		},
		TotalPrice: 0,
		FinalPrice: 0,
	}

	cart, err := testStore.CreateCart(context.Background(), cartArg)
	require.NoError(t, err)
	require.NotEmpty(t, cart)

	_, err = testStore.CreateCartProduct(context.Background(), CreateCartProductParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
		Num: int32(num),
	})
	require.NoError(t, err)

	arg := UpdateCartTxParams{
		CartID: cart.ID,
	}

	result1, err := testStore.UpdateCartTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result1)

	require.Equal(t, result1.Cart.ID, cart.ID)
	require.Equal(t, result1.Cart.Owner.String, user.Account)
	require.NotEqual(t, result1.Cart.TotalPrice, cart.TotalPrice)
	require.NotEqual(t, result1.Cart.FinalPrice, cart.FinalPrice)

	require.Equal(t, result1.ProductList, []CartTxProductResult{
		{
			Product: product,
			Num:     int64(num),
		},
	})

	_, err = testStore.UpdateCartProduct(context.Background(), UpdateCartProductParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
		Num: pgtype.Int4{
			Int32: int32(num * 2),
			Valid: true,
		},
	})
	require.NoError(t, err)

	arg = UpdateCartTxParams{
		CartID: cart.ID,
	}

	result2, err := testStore.UpdateCartTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result2)

	require.Equal(t, result2.Cart.ID, cart.ID)
	require.Equal(t, result2.Cart.Owner.String, user.Account)
	require.Equal(t, result1.Cart.TotalPrice*2, result2.Cart.TotalPrice)
	require.Equal(t, result1.Cart.FinalPrice*2, result2.Cart.FinalPrice)

	require.Equal(t, result2.ProductList, []CartTxProductResult{
		{
			Product: product,
			Num:     int64(2 * num),
		},
	})
}
