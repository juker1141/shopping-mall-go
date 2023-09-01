package db

import (
	"context"
	"testing"
	"time"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	name := util.RandomName()

	category, err := testStore.CreateCategory(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, name, category.Name)
	require.NotZero(t, category.CreatedAt)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategory(t *testing.T) {
	category1 := createRandomCategory(t)

	category2, err := testStore.GetCategory(context.Background(), category1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, category1.Name, category2.Name)
	require.WithinDuration(t, category1.CreatedAt, category2.CreatedAt, time.Second)
}

func TestUpdateCategory(t *testing.T) {
	oldCategory := createRandomCategory(t)

	arg := UpdateCategoryParams{
		ID:   oldCategory.ID,
		Name: util.RandomName(),
	}
	newCategory, err := testStore.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCategory)

	require.Equal(t, oldCategory.ID, newCategory.ID)
	require.NotEqual(t, oldCategory.Name, newCategory.Name)
	require.WithinDuration(t, oldCategory.CreatedAt, newCategory.CreatedAt, time.Second)
}

func TestDeleteCategory(t *testing.T) {
	category1 := createRandomCategory(t)

	err := testStore.DeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	category2, err := testStore.GetCategory(context.Background(), category1.ID)
	require.Error(t, err)
	require.Empty(t, category2)
}

func TestListCategory(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomCategory(t)
	}

	categories, err := testStore.ListCategories(context.Background())
	require.NoError(t, err)
	require.NotZero(t, categories)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}
