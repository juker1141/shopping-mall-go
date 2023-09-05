package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomProductCategory(t *testing.T) ProductCategory {
	product := createRandomProduct(t, util.RandomName())
	category := createRandomCategory(t)

	arg := CreateProductCategoryParams{
		ProductID: pgtype.Int4{
			Int32: int32(product.ID),
			Valid: true,
		},
		CategoryID: pgtype.Int4{
			Int32: int32(category.ID),
			Valid: true,
		},
	}

	productCategory, err := testStore.CreateProductCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productCategory)

	require.NotZero(t, productCategory.ProductID)
	require.NotZero(t, productCategory.CategoryID)

	return productCategory
}

func TestCreateProductCategory(t *testing.T) {
	createRandomProductCategory(t)
}

func TestGetProductCategory(t *testing.T) {
	productCategory1 := createRandomProductCategory(t)

	arg := GetProductCategoryParams{
		ProductID: pgtype.Int4{
			Int32: productCategory1.ProductID.Int32,
			Valid: true,
		},
		CategoryID: pgtype.Int4{
			Int32: productCategory1.CategoryID.Int32,
			Valid: true,
		},
	}

	productCategory2, err := testStore.GetProductCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, productCategory2)

	require.Equal(t, productCategory1.ProductID, productCategory2.ProductID)
	require.Equal(t, productCategory1.CategoryID, productCategory2.CategoryID)
}

func TestDeleteProductCategoryByProductId(t *testing.T) {
	productCategory1 := createRandomProductCategory(t)

	err := testStore.DeleteProductCategoryByProductId(context.Background(), productCategory1.ProductID)
	require.NoError(t, err)

	productCategory2, err := testStore.GetProductCategory(context.Background(), GetProductCategoryParams{
		ProductID: pgtype.Int4{
			Int32: productCategory1.ProductID.Int32,
			Valid: true,
		},
		CategoryID: pgtype.Int4{
			Int32: productCategory1.CategoryID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, productCategory2)
}

func TestDeleteProductCategoryByCategoryId(t *testing.T) {
	productCategory1 := createRandomProductCategory(t)

	err := testStore.DeleteProductCategoryByCategoryId(context.Background(), productCategory1.CategoryID)
	require.NoError(t, err)

	productCategory2, err := testStore.GetProductCategory(context.Background(), GetProductCategoryParams{
		ProductID: pgtype.Int4{
			Int32: productCategory1.ProductID.Int32,
			Valid: true,
		},
		CategoryID: pgtype.Int4{
			Int32: productCategory1.CategoryID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, productCategory2)
}

func TestListProductCategoryByProductId(t *testing.T) {
	product := createRandomProduct(t, util.RandomName())
	for i := 0; i < 5; i++ {
		category := createRandomCategory(t)
		testStore.CreateProductCategory(context.Background(), CreateProductCategoryParams{
			ProductID: pgtype.Int4{
				Int32: int32(product.ID),
				Valid: true,
			},
			CategoryID: pgtype.Int4{
				Int32: int32(category.ID),
				Valid: true,
			},
		})
	}

	productCategories, err := testStore.ListProductCategoryByProductId(context.Background(), pgtype.Int4{
		Int32: int32(product.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, productCategories, 5)

	for _, productCategory := range productCategories {
		require.NotEmpty(t, productCategory)
	}
}

func TestListProductCategoryByCategoryId(t *testing.T) {
	category := createRandomCategory(t)
	for i := 0; i < 5; i++ {
		product := createRandomProduct(t, util.RandomName())
		testStore.CreateProductCategory(context.Background(), CreateProductCategoryParams{
			ProductID: pgtype.Int4{
				Int32: int32(product.ID),
				Valid: true,
			},
			CategoryID: pgtype.Int4{
				Int32: int32(category.ID),
				Valid: true,
			},
		})
	}

	productCategories, err := testStore.ListProductCategoryByCategoryId(context.Background(), pgtype.Int4{
		Int32: int32(category.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, productCategories, 5)

	for _, productCategory := range productCategories {
		require.NotEmpty(t, productCategory)
	}
}
