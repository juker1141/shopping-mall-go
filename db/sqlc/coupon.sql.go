// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: coupon.sql

package db

import (
	"context"
	"time"
)

const createCoupon = `-- name: CreateCoupon :one
INSERT INTO coupons (
  title,
  code,
  percent,
  created_by,
  start_at,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id, title, code, percent, created_by, start_at, expires_at, created_at
`

type CreateCouponParams struct {
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	Percent   int32     `json:"percent"`
	CreatedBy string    `json:"created_by"`
	StartAt   time.Time `json:"start_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateCoupon(ctx context.Context, arg CreateCouponParams) (Coupon, error) {
	row := q.db.QueryRow(ctx, createCoupon,
		arg.Title,
		arg.Code,
		arg.Percent,
		arg.CreatedBy,
		arg.StartAt,
		arg.ExpiresAt,
	)
	var i Coupon
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Code,
		&i.Percent,
		&i.CreatedBy,
		&i.StartAt,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteCoupon = `-- name: DeleteCoupon :exec
DELETE FROM coupons WHERE id = $1
`

func (q *Queries) DeleteCoupon(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteCoupon, id)
	return err
}

const getCoupon = `-- name: GetCoupon :one
SELECT id, title, code, percent, created_by, start_at, expires_at, created_at FROM coupons
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCoupon(ctx context.Context, id int64) (Coupon, error) {
	row := q.db.QueryRow(ctx, getCoupon, id)
	var i Coupon
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Code,
		&i.Percent,
		&i.CreatedBy,
		&i.StartAt,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const listCoupons = `-- name: ListCoupons :many
SELECT id, title, code, percent, created_by, start_at, expires_at, created_at FROM coupons
WHERE
  CASE
    WHEN $1::varchar = 'title' THEN title ILIKE '%' || $2::varchar || '%'
    WHEN $1::varchar = 'code' THEN code ILIKE '%' || $2::varchar || '%'
    WHEN $1::varchar = 'start_at' THEN start_at ILIKE '%' || $2::varchar || '%'
    WHEN $1::varchar = 'expires_at' THEN expires_at ILIKE '%' || $2::varchar || '%'
    ELSE TRUE
  END
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListCouponsParams struct {
	Column1 string `json:"column_1"`
	Column2 string `json:"column_2"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

func (q *Queries) ListCoupons(ctx context.Context, arg ListCouponsParams) ([]Coupon, error) {
	rows, err := q.db.Query(ctx, listCoupons,
		arg.Column1,
		arg.Column2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Coupon{}
	for rows.Next() {
		var i Coupon
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Code,
			&i.Percent,
			&i.CreatedBy,
			&i.StartAt,
			&i.ExpiresAt,
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

const updateCoupon = `-- name: UpdateCoupon :one
UPDATE coupons
SET 
  title = $2,
  code = $3,
  percent = $4,
  start_at = $5,
  expires_at = $6
WHERE id = $1
RETURNING id, title, code, percent, created_by, start_at, expires_at, created_at
`

type UpdateCouponParams struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	Percent   int32     `json:"percent"`
	StartAt   time.Time `json:"start_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) UpdateCoupon(ctx context.Context, arg UpdateCouponParams) (Coupon, error) {
	row := q.db.QueryRow(ctx, updateCoupon,
		arg.ID,
		arg.Title,
		arg.Code,
		arg.Percent,
		arg.StartAt,
		arg.ExpiresAt,
	)
	var i Coupon
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Code,
		&i.Percent,
		&i.CreatedBy,
		&i.StartAt,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}