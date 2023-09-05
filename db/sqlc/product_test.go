package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T, name string) Product {
	imageUrl := util.RandomString(10)

	arg := CreateProductParams{
		Title:       name,
		Description: util.RandomString(20),
		Content:     util.RandomString(30),
		OriginPrice: util.RandomPrice(),
		Price:       util.RandomPrice(),
		Unit:        util.RandomName(),
		Status:      1,
		ImageUrl:    imageUrl,
		ImagesUrl: []string{
			imageUrl,
			imageUrl,
			imageUrl,
		},
		CreatedBy: util.RandomName(),
	}

	product, err := testStore.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.NotZero(t, product.ID)
	require.Equal(t, arg.Title, product.Title)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.Content, product.Content)
	require.Equal(t, arg.OriginPrice, product.OriginPrice)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Unit, product.Unit)
	require.Equal(t, arg.Status, product.Status)
	require.Equal(t, arg.ImageUrl, product.ImageUrl)
	require.Equal(t, arg.ImagesUrl, product.ImagesUrl)
	require.Equal(t, arg.CreatedBy, product.CreatedBy)

	require.NotZero(t, product.CreatedAt)
	return product
}

func TestCreateProduct(t *testing.T) {
	createRandomProduct(t, util.RandomName())
}

func TestCreateProductButImagesEmpty(t *testing.T) {
	imageUrl := util.RandomString(10)

	arg := CreateProductParams{
		Title:       util.RandomName(),
		Description: util.RandomString(20),
		Content:     util.RandomString(30),
		OriginPrice: util.RandomPrice(),
		Price:       util.RandomPrice(),
		Unit:        util.RandomName(),
		Status:      1,
		ImageUrl:    imageUrl,
		CreatedBy:   util.RandomName(),
	}

	product, err := testStore.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.NotZero(t, product.ID)
	require.Equal(t, arg.Title, product.Title)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.Content, product.Content)
	require.Equal(t, arg.OriginPrice, product.OriginPrice)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Unit, product.Unit)
	require.Equal(t, arg.Status, product.Status)
	require.Equal(t, arg.ImageUrl, product.ImageUrl)
	require.Equal(t, arg.CreatedBy, product.CreatedBy)

	require.Empty(t, product.ImagesUrl)

	require.NotZero(t, product.CreatedAt)
}

func TestGetProduct(t *testing.T) {
	product1 := createRandomProduct(t, util.RandomName())

	product2, err := testStore.GetProduct(context.Background(), product1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, product1.Title, product2.Title)
	require.Equal(t, product1.Description, product2.Description)
	require.Equal(t, product1.Content, product2.Content)
	require.Equal(t, product1.OriginPrice, product2.OriginPrice)
	require.Equal(t, product1.Price, product2.Price)
	require.Equal(t, product1.Unit, product2.Unit)
	require.Equal(t, product1.Status, product2.Status)
	require.Equal(t, product1.ImageUrl, product2.ImageUrl)
	require.Equal(t, product1.ImagesUrl, product2.ImagesUrl)
	require.Equal(t, product1.CreatedBy, product2.CreatedBy)

	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
}

func TestUpdateProductAllField(t *testing.T) {
	oldProduct := createRandomProduct(t, util.RandomName())

	newImageUrl := util.RandomString(10)

	arg := UpdateProductParams{
		ID: oldProduct.ID,
		Title: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
		Description: pgtype.Text{
			String: util.RandomString(20),
			Valid:  true,
		},
		Content: pgtype.Text{
			String: util.RandomString(30),
			Valid:  true,
		},
		OriginPrice: pgtype.Int4{
			Int32: util.RandomPrice(),
			Valid: true,
		},
		Price: pgtype.Int4{
			Int32: util.RandomPrice(),
			Valid: true,
		},
		Unit: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
		Status: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
		ImageUrl: pgtype.Text{
			String: newImageUrl,
			Valid:  true,
		},
		ImagesUrl: []string{
			newImageUrl,
			newImageUrl,
			newImageUrl,
		},
	}

	newProduct, err := testStore.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newProduct)

	require.Equal(t, oldProduct.ID, newProduct.ID)
	require.Equal(t, oldProduct.CreatedBy, newProduct.CreatedBy)
	require.WithinDuration(t, oldProduct.CreatedAt, newProduct.CreatedAt, time.Second)

	require.NotEqual(t, oldProduct.Title, newProduct.Title)
	require.NotEqual(t, oldProduct.Description, newProduct.Description)
	require.NotEqual(t, oldProduct.Content, newProduct.Content)
	require.NotEqual(t, oldProduct.OriginPrice, newProduct.OriginPrice)
	require.NotEqual(t, oldProduct.Price, newProduct.Price)
	require.NotEqual(t, oldProduct.Unit, newProduct.Unit)
	require.NotEqual(t, oldProduct.Status, newProduct.Status)
	require.NotEqual(t, oldProduct.ImageUrl, newProduct.ImageUrl)
	require.NotEqual(t, oldProduct.ImagesUrl, newProduct.ImagesUrl)
}

func TestUpdateProductOnlyStatus(t *testing.T) {
	oldProduct := createRandomProduct(t, util.RandomName())

	arg := UpdateProductParams{
		ID: oldProduct.ID,
		Status: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
	}

	newProduct, err := testStore.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newProduct)

	require.Equal(t, oldProduct.ID, newProduct.ID)
	require.Equal(t, oldProduct.CreatedBy, newProduct.CreatedBy)
	require.WithinDuration(t, oldProduct.CreatedAt, newProduct.CreatedAt, time.Second)

	require.Equal(t, oldProduct.Title, newProduct.Title)
	require.Equal(t, oldProduct.Description, newProduct.Description)
	require.Equal(t, oldProduct.Content, newProduct.Content)
	require.Equal(t, oldProduct.OriginPrice, newProduct.OriginPrice)
	require.Equal(t, oldProduct.Price, newProduct.Price)
	require.Equal(t, oldProduct.Unit, newProduct.Unit)
	require.Equal(t, oldProduct.ImageUrl, newProduct.ImageUrl)
	require.Equal(t, oldProduct.ImagesUrl, newProduct.ImagesUrl)

	require.NotEqual(t, oldProduct.Status, newProduct.Status)
}

func TestUpdateProductOnlyPrice(t *testing.T) {
	oldProduct := createRandomProduct(t, util.RandomName())

	arg := UpdateProductParams{
		ID: oldProduct.ID,
		OriginPrice: pgtype.Int4{
			Int32: util.RandomPrice(),
			Valid: true,
		},
		Price: pgtype.Int4{
			Int32: util.RandomPrice(),
			Valid: true,
		},
	}

	newProduct, err := testStore.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newProduct)

	require.Equal(t, oldProduct.ID, newProduct.ID)
	require.Equal(t, oldProduct.CreatedBy, newProduct.CreatedBy)
	require.WithinDuration(t, oldProduct.CreatedAt, newProduct.CreatedAt, time.Second)

	require.Equal(t, oldProduct.Title, newProduct.Title)
	require.Equal(t, oldProduct.Description, newProduct.Description)
	require.Equal(t, oldProduct.Content, newProduct.Content)
	require.Equal(t, oldProduct.Unit, newProduct.Unit)
	require.Equal(t, oldProduct.Status, newProduct.Status)
	require.Equal(t, oldProduct.ImageUrl, newProduct.ImageUrl)
	require.Equal(t, oldProduct.ImagesUrl, newProduct.ImagesUrl)

	require.NotEqual(t, oldProduct.OriginPrice, newProduct.OriginPrice)
	require.NotEqual(t, oldProduct.Price, newProduct.Price)
}

func TestDeleteProduct(t *testing.T) {
	product1 := createRandomProduct(t, util.RandomName())

	err := testStore.DeleteProduct(context.Background(), product1.ID)
	require.NoError(t, err)

	product2, err := testStore.GetProduct(context.Background(), product1.ID)
	require.Error(t, err)
	require.Empty(t, product2)
}

func TestListProductsNoSearch(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomProduct(t, util.RandomName())
	}

	arg := ListProductsParams{
		Limit:  5,
		Offset: 5,
	}
	products, err := testStore.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}

func TestListProductsSearchTitle(t *testing.T) {
	n := 3
	name := util.RandomName()
	for i := 0; i < n; i++ {
		createRandomProduct(t, name)
	}
	for i := 0; i < 10-n; i++ {
		createRandomProduct(t, util.RandomName())
	}

	arg := ListProductsParams{
		Key:      KeyTitle,
		KeyValue: name,
		Limit:    10,
		Offset:   0,
	}
	products, err := testStore.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, products, n)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}

func TestListProductsSearchButNoKeyValue(t *testing.T) {
	n := 3
	name := util.RandomName()
	for i := 0; i < n; i++ {
		createRandomProduct(t, name)
	}
	for i := 0; i < 10-n; i++ {
		createRandomProduct(t, util.RandomName())
	}

	arg := ListProductsParams{
		Key:    KeyTitle,
		Limit:  5,
		Offset: 5,
	}
	products, err := testStore.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}

func TestListProductsSearchButNoKey(t *testing.T) {
	n := 3
	name := util.RandomName()
	for i := 0; i < n; i++ {
		createRandomProduct(t, name)
	}
	for i := 0; i < 10-n; i++ {
		createRandomProduct(t, util.RandomName())
	}

	arg := ListProductsParams{
		KeyValue: name,
		Limit:    5,
		Offset:   5,
	}
	products, err := testStore.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}
