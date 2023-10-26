package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func TestUpdateCartTx(t *testing.T) {
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

	_, err = testStore.GetCartProduct(context.Background(), GetCartProductParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
	})
	fmt.Println(err)

	arg := UpdateCartTxParams{
		Account:   user.Account,
		Type:      "update",
		ProductID: product.ID,
		Num:       int32(num),
	}

	result, err := testStore.UpdateCartTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, result.Cart.ID, cart.ID)
	require.NotEqual(t, result.Cart.Owner, user.Account)
	require.NotEqual(t, result.Cart.TotalPrice, cart.TotalPrice)
	require.NotEqual(t, result.Cart.FinalPrice, cart.FinalPrice)

	require.Equal(t, result.ProductList, CartTxProductResult{
		Product: product,
		Num:     int64(num),
	})
}
