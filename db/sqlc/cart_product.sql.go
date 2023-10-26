// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: cart_product.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkCartProductExists = `-- name: CheckCartProductExists :one
SELECT EXISTS (
    SELECT 1
    FROM cart_products
    WHERE cart_id = $1 AND product_id = $2
)
`

type CheckCartProductExistsParams struct {
	CartID    pgtype.Int4 `json:"cart_id"`
	ProductID pgtype.Int4 `json:"product_id"`
}

func (q *Queries) CheckCartProductExists(ctx context.Context, arg CheckCartProductExistsParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkCartProductExists, arg.CartID, arg.ProductID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createCartProduct = `-- name: CreateCartProduct :one
INSERT INTO cart_products (
  cart_id,
  product_id,
  num
) VALUES (
  $1, $2, $3
) RETURNING cart_id, product_id, num
`

type CreateCartProductParams struct {
	CartID    pgtype.Int4 `json:"cart_id"`
	ProductID pgtype.Int4 `json:"product_id"`
	Num       int32       `json:"num"`
}

func (q *Queries) CreateCartProduct(ctx context.Context, arg CreateCartProductParams) (CartProduct, error) {
	row := q.db.QueryRow(ctx, createCartProduct, arg.CartID, arg.ProductID, arg.Num)
	var i CartProduct
	err := row.Scan(&i.CartID, &i.ProductID, &i.Num)
	return i, err
}

const deleteCartProduct = `-- name: DeleteCartProduct :exec
DELETE FROM cart_products
WHERE cart_id = $1 AND product_id = $2
`

type DeleteCartProductParams struct {
	CartID    pgtype.Int4 `json:"cart_id"`
	ProductID pgtype.Int4 `json:"product_id"`
}

func (q *Queries) DeleteCartProduct(ctx context.Context, arg DeleteCartProductParams) error {
	_, err := q.db.Exec(ctx, deleteCartProduct, arg.CartID, arg.ProductID)
	return err
}

const deleteCartProductByCartId = `-- name: DeleteCartProductByCartId :exec
DELETE FROM cart_products
WHERE cart_id = $1
`

func (q *Queries) DeleteCartProductByCartId(ctx context.Context, cartID pgtype.Int4) error {
	_, err := q.db.Exec(ctx, deleteCartProductByCartId, cartID)
	return err
}

const deleteCartProductByProductId = `-- name: DeleteCartProductByProductId :exec
DELETE FROM cart_products
WHERE product_id = $1
`

func (q *Queries) DeleteCartProductByProductId(ctx context.Context, productID pgtype.Int4) error {
	_, err := q.db.Exec(ctx, deleteCartProductByProductId, productID)
	return err
}

const getCartProduct = `-- name: GetCartProduct :one
SELECT cart_id, product_id, num FROM cart_products
WHERE cart_id = $1 AND product_id = $2
LIMIT 1
`

type GetCartProductParams struct {
	CartID    pgtype.Int4 `json:"cart_id"`
	ProductID pgtype.Int4 `json:"product_id"`
}

func (q *Queries) GetCartProduct(ctx context.Context, arg GetCartProductParams) (CartProduct, error) {
	row := q.db.QueryRow(ctx, getCartProduct, arg.CartID, arg.ProductID)
	var i CartProduct
	err := row.Scan(&i.CartID, &i.ProductID, &i.Num)
	return i, err
}

const listCartProductByCartId = `-- name: ListCartProductByCartId :many
SELECT cart_id, product_id, num FROM cart_products
WHERE cart_id = $1
`

func (q *Queries) ListCartProductByCartId(ctx context.Context, cartID pgtype.Int4) ([]CartProduct, error) {
	rows, err := q.db.Query(ctx, listCartProductByCartId, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CartProduct{}
	for rows.Next() {
		var i CartProduct
		if err := rows.Scan(&i.CartID, &i.ProductID, &i.Num); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCartProductByProductId = `-- name: ListCartProductByProductId :many
SELECT cart_id, product_id, num FROM cart_products
WHERE product_id = $1
`

func (q *Queries) ListCartProductByProductId(ctx context.Context, productID pgtype.Int4) ([]CartProduct, error) {
	rows, err := q.db.Query(ctx, listCartProductByProductId, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CartProduct{}
	for rows.Next() {
		var i CartProduct
		if err := rows.Scan(&i.CartID, &i.ProductID, &i.Num); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCartProduct = `-- name: UpdateCartProduct :one
UPDATE cart_products
SET 
  num = COALESCE($1, num)
WHERE
  cart_id = $2 AND product_id = $3
RETURNING cart_id, product_id, num
`

type UpdateCartProductParams struct {
	Num       pgtype.Int4 `json:"num"`
	CartID    pgtype.Int4 `json:"cart_id"`
	ProductID pgtype.Int4 `json:"product_id"`
}

func (q *Queries) UpdateCartProduct(ctx context.Context, arg UpdateCartProductParams) (CartProduct, error) {
	row := q.db.QueryRow(ctx, updateCartProduct, arg.Num, arg.CartID, arg.ProductID)
	var i CartProduct
	err := row.Scan(&i.CartID, &i.ProductID, &i.Num)
	return i, err
}
