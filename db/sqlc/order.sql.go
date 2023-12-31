// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: order.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
  full_name,
  email,
  shipping_address,
  message,
  total_price,
  final_price,
  pay_method_id,
  status_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, full_name, email, shipping_address, message, is_paid, total_price, final_price, pay_method_id, status_id, created_at, updated_at
`

type CreateOrderParams struct {
	FullName        string      `json:"full_name"`
	Email           string      `json:"email"`
	ShippingAddress string      `json:"shipping_address"`
	Message         pgtype.Text `json:"message"`
	TotalPrice      int32       `json:"total_price"`
	FinalPrice      int32       `json:"final_price"`
	PayMethodID     int32       `json:"pay_method_id"`
	StatusID        int32       `json:"status_id"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder,
		arg.FullName,
		arg.Email,
		arg.ShippingAddress,
		arg.Message,
		arg.TotalPrice,
		arg.FinalPrice,
		arg.PayMethodID,
		arg.StatusID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.ShippingAddress,
		&i.Message,
		&i.IsPaid,
		&i.TotalPrice,
		&i.FinalPrice,
		&i.PayMethodID,
		&i.StatusID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteOrder, id)
	return err
}

const getOrder = `-- name: GetOrder :one
SELECT id, full_name, email, shipping_address, message, is_paid, total_price, final_price, pay_method_id, status_id, created_at, updated_at FROM orders
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetOrder(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRow(ctx, getOrder, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.ShippingAddress,
		&i.Message,
		&i.IsPaid,
		&i.TotalPrice,
		&i.FinalPrice,
		&i.PayMethodID,
		&i.StatusID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrdersCount = `-- name: GetOrdersCount :one
SELECT COUNT(*) FROM orders
`

func (q *Queries) GetOrdersCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getOrdersCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listOrders = `-- name: ListOrders :many

SELECT id, full_name, email, shipping_address, message, is_paid, total_price, final_price, pay_method_id, status_id, created_at, updated_at FROM orders
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListOrdersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// -- name: GetOrderByCode :one
// SELECT * FROM orders
// WHERE code = $1 LIMIT 1;
func (q *Queries) ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error) {
	rows, err := q.db.Query(ctx, listOrders, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.ShippingAddress,
			&i.Message,
			&i.IsPaid,
			&i.TotalPrice,
			&i.FinalPrice,
			&i.PayMethodID,
			&i.StatusID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders
SET 
  full_name = COALESCE($1, full_name),
  email = COALESCE($2, email),
  shipping_address = COALESCE($3, shipping_address),
  message = COALESCE($4, message),
  pay_method_id = COALESCE($5, pay_method_id),
  is_paid = COALESCE($6, is_paid),
  total_price = COALESCE($7, total_price),
  final_price = COALESCE($8, final_price),
  status_id = COALESCE($9, status_id),
  updated_at = COALESCE($10, updated_at)
WHERE id = $11
RETURNING id, full_name, email, shipping_address, message, is_paid, total_price, final_price, pay_method_id, status_id, created_at, updated_at
`

type UpdateOrderParams struct {
	FullName        pgtype.Text        `json:"full_name"`
	Email           pgtype.Text        `json:"email"`
	ShippingAddress pgtype.Text        `json:"shipping_address"`
	Message         pgtype.Text        `json:"message"`
	PayMethodID     pgtype.Int4        `json:"pay_method_id"`
	IsPaid          pgtype.Bool        `json:"is_paid"`
	TotalPrice      pgtype.Int4        `json:"total_price"`
	FinalPrice      pgtype.Int4        `json:"final_price"`
	StatusID        pgtype.Int4        `json:"status_id"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	ID              int64              `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, updateOrder,
		arg.FullName,
		arg.Email,
		arg.ShippingAddress,
		arg.Message,
		arg.PayMethodID,
		arg.IsPaid,
		arg.TotalPrice,
		arg.FinalPrice,
		arg.StatusID,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.ShippingAddress,
		&i.Message,
		&i.IsPaid,
		&i.TotalPrice,
		&i.FinalPrice,
		&i.PayMethodID,
		&i.StatusID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
