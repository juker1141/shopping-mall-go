// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: product.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
  title,
  category,
  description,
  content,
  origin_price,
  price,
  unit,
  status,
  image_url,
  images_url,
  created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING id, title, category, origin_price, price, unit, description, content, status, image_url, images_url, created_by, created_at
`

type CreateProductParams struct {
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	OriginPrice int32    `json:"origin_price"`
	Price       int32    `json:"price"`
	Unit        string   `json:"unit"`
	Status      int32    `json:"status"`
	ImageUrl    string   `json:"image_url"`
	ImagesUrl   []string `json:"images_url"`
	CreatedBy   string   `json:"created_by"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.Title,
		arg.Category,
		arg.Description,
		arg.Content,
		arg.OriginPrice,
		arg.Price,
		arg.Unit,
		arg.Status,
		arg.ImageUrl,
		arg.ImagesUrl,
		arg.CreatedBy,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.OriginPrice,
		&i.Price,
		&i.Unit,
		&i.Description,
		&i.Content,
		&i.Status,
		&i.ImageUrl,
		&i.ImagesUrl,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, title, category, origin_price, price, unit, description, content, status, image_url, images_url, created_by, created_at FROM products
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRow(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.OriginPrice,
		&i.Price,
		&i.Unit,
		&i.Description,
		&i.Content,
		&i.Status,
		&i.ImageUrl,
		&i.ImagesUrl,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}

const getProductsCount = `-- name: GetProductsCount :one
SELECT COUNT(*) FROM products
`

func (q *Queries) GetProductsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getProductsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listProducts = `-- name: ListProducts :many
SELECT p.id, p.title, p.category, p.origin_price, p.price, p.unit, p.description, p.content, p.status, p.image_url, p.images_url, p.created_by, p.created_at
FROM products AS p
WHERE
  CASE
    WHEN $1::varchar = 'title' THEN p.title ILIKE '%' || $2::varchar || '%'
    ELSE TRUE
  END
ORDER BY p.id
LIMIT $4
OFFSET $3
`

type ListProductsParams struct {
	Key      string `json:"key"`
	KeyValue string `json:"key_value"`
	Offset   int32  `json:"Offset"`
	Limit    int32  `json:"Limit"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.Query(ctx, listProducts,
		arg.Key,
		arg.KeyValue,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Category,
			&i.OriginPrice,
			&i.Price,
			&i.Unit,
			&i.Description,
			&i.Content,
			&i.Status,
			&i.ImageUrl,
			&i.ImagesUrl,
			&i.CreatedBy,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE products
SET
  title = COALESCE($1, title),
  category = COALESCE($2, category),
  description = COALESCE($3, description),
  content = COALESCE($4, content),
  origin_price = COALESCE($5, origin_price),
  price = COALESCE($6, price),
  unit = COALESCE($7, unit),
  status = COALESCE($8, status),
  image_url = COALESCE($9, image_url),
  images_url = COALESCE($10, images_url)
WHERE
  id = $11
RETURNING id, title, category, origin_price, price, unit, description, content, status, image_url, images_url, created_by, created_at
`

type UpdateProductParams struct {
	Title       pgtype.Text `json:"title"`
	Category    pgtype.Text `json:"category"`
	Description pgtype.Text `json:"description"`
	Content     pgtype.Text `json:"content"`
	OriginPrice pgtype.Int4 `json:"origin_price"`
	Price       pgtype.Int4 `json:"price"`
	Unit        pgtype.Text `json:"unit"`
	Status      pgtype.Int4 `json:"status"`
	ImageUrl    pgtype.Text `json:"image_url"`
	ImagesUrl   []string    `json:"images_url"`
	ID          int64       `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.Title,
		arg.Category,
		arg.Description,
		arg.Content,
		arg.OriginPrice,
		arg.Price,
		arg.Unit,
		arg.Status,
		arg.ImageUrl,
		arg.ImagesUrl,
		arg.ID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Category,
		&i.OriginPrice,
		&i.Price,
		&i.Unit,
		&i.Description,
		&i.Content,
		&i.Status,
		&i.ImageUrl,
		&i.ImagesUrl,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}
