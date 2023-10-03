package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomCartProduct(t *testing.T) CartProduct {
	cart := createRandomCart(t)
	product := createRandomProduct(t, util.RandomName())
	num := int32(util.RandomInt(1, 100))

	arg := CreateCartProductParams{
		CartID: pgtype.Int4{
			Int32: int32(cart.ID),
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
		Num: num,
	}

	cartProduct, err := testStore.CreateCartProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cartProduct)

	require.NotZero(t, cartProduct.CartID)
	require.NotZero(t, cartProduct.ProductID)
	require.Equal(t, num, cartProduct.Num)

	return cartProduct
}

func TestCreateCartProduct(t *testing.T) {
	createRandomCartProduct(t)
}

func TestGetCartProduct(t *testing.T) {
	cartProduct1 := createRandomCartProduct(t)

	arg := GetCartProductParams{
		CartID:    cartProduct1.CartID,
		ProductID: cartProduct1.ProductID,
	}

	cartProduct2, err := testStore.GetCartProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cartProduct2)

	require.Equal(t, cartProduct1.CartID, cartProduct2.CartID)
	require.Equal(t, cartProduct1.ProductID, cartProduct2.ProductID)
	require.Equal(t, cartProduct1.Num, cartProduct2.Num)
}

func TestUpdateCartProduct(t *testing.T) {
	oldCartProduct := createRandomCartProduct(t)

	arg := UpdateCartProductParams{
		CartID: pgtype.Int4{
			Int32: oldCartProduct.CartID.Int32,
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: oldCartProduct.ProductID.Int32,
			Valid: true,
		},
		Num: pgtype.Int4{
			Int32: oldCartProduct.Num + 10,
			Valid: true,
		},
	}

	newCartProduct, err := testStore.UpdateCartProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCartProduct)

	require.Equal(t, oldCartProduct.CartID, newCartProduct.CartID)
	require.Equal(t, oldCartProduct.ProductID, newCartProduct.ProductID)
	require.NotEqual(t, oldCartProduct.Num, newCartProduct.Num)
}

func TestDeleteCartProductByCartId(t *testing.T) {
	cartProduct1 := createRandomCartProduct(t)

	err := testStore.DeleteCartProductByCartId(context.Background(), cartProduct1.CartID)
	require.NoError(t, err)

	cartProduct2, err := testStore.GetCartProduct(context.Background(), GetCartProductParams{
		CartID: pgtype.Int4{
			Int32: cartProduct1.CartID.Int32,
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: cartProduct1.ProductID.Int32,
			Valid: true,
		},
	})

	require.Error(t, err)
	require.Empty(t, cartProduct2)
}

func TestDeleteCartProductByProductId(t *testing.T) {
	cartProduct1 := createRandomCartProduct(t)

	err := testStore.DeleteCartProductByProductId(context.Background(), cartProduct1.ProductID)
	require.NoError(t, err)

	cartProduct2, err := testStore.GetCartProduct(context.Background(), GetCartProductParams{
		CartID: pgtype.Int4{
			Int32: cartProduct1.CartID.Int32,
			Valid: true,
		},
		ProductID: pgtype.Int4{
			Int32: cartProduct1.ProductID.Int32,
			Valid: true,
		},
	})

	require.Error(t, err)
	require.Empty(t, cartProduct2)
}

func TestListCartProductByCartId(t *testing.T) {
	cart := createRandomCart(t)
	for i := 0; i < 5; i++ {
		product := createRandomProduct(t, util.RandomName())

		testStore.CreateCartProduct(context.Background(), CreateCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(product.ID),
				Valid: true,
			},
			Num: int32(util.RandomInt(1, 100)),
		})
	}

	cartProducts, err := testStore.ListCartProductByCartId(context.Background(), pgtype.Int4{
		Int32: int32(cart.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, cartProducts, 5)

	for _, cartProduct := range cartProducts {
		require.NotEmpty(t, cartProduct)
	}
}

func TestListCartProductByProductId(t *testing.T) {
	product := createRandomProduct(t, util.RandomName())
	for i := 0; i < 5; i++ {
		cart := createRandomCart(t)

		testStore.CreateCartProduct(context.Background(), CreateCartProductParams{
			CartID: pgtype.Int4{
				Int32: int32(cart.ID),
				Valid: true,
			},
			ProductID: pgtype.Int4{
				Int32: int32(product.ID),
				Valid: true,
			},
			Num: int32(util.RandomInt(1, 100)),
		})
	}

	cartProducts, err := testStore.ListCartProductByProductId(context.Background(), pgtype.Int4{
		Int32: int32(product.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, cartProducts, 5)

	for _, cartProduct := range cartProducts {
		require.NotEmpty(t, cartProduct)
	}
}
