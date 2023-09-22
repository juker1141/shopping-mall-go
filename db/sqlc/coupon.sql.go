// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: coupon.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

const getCouponByCode = `-- name: GetCouponByCode :one
SELECT id, title, code, percent, created_by, start_at, expires_at, created_at FROM coupons
WHERE code = $1 LIMIT 1
`

func (q *Queries) GetCouponByCode(ctx context.Context, code string) (Coupon, error) {
	row := q.db.QueryRow(ctx, getCouponByCode, code)
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

const getCouponsCount = `-- name: GetCouponsCount :one
SELECT COUNT(*) FROM coupons
`

func (q *Queries) GetCouponsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getCouponsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listCoupons = `-- name: ListCoupons :many
SELECT id, title, code, percent, created_by, start_at, expires_at, created_at FROM coupons
WHERE
  CASE
    WHEN $1::varchar = 'title' THEN title ILIKE '%' || $2::varchar || '%'
    WHEN $1::varchar = 'code' THEN code ILIKE '%' || $2::varchar || '%'
    WHEN $1::varchar = 'start_time' THEN start_at >= $3::timestamptz
    WHEN $1::varchar = 'expires_time' THEN expires_at <= $3::timestamptz
    ELSE TRUE
  END
ORDER BY id
LIMIT $5
OFFSET $4
`

type ListCouponsParams struct {
	Key          string    `json:"key"`
	KeyValue     string    `json:"key_value"`
	KeyTimeValue time.Time `json:"key_time_value"`
	Offset       int32     `json:"Offset"`
	Limit        int32     `json:"Limit"`
}

func (q *Queries) ListCoupons(ctx context.Context, arg ListCouponsParams) ([]Coupon, error) {
	rows, err := q.db.Query(ctx, listCoupons,
		arg.Key,
		arg.KeyValue,
		arg.KeyTimeValue,
		arg.Offset,
		arg.Limit,
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
  title = COALESCE($1, title),
  code = COALESCE($2, code),
  percent = COALESCE($3, percent),
  start_at = COALESCE($4, start_at),
  expires_at = COALESCE($5, expires_at)
WHERE id = $6
RETURNING id, title, code, percent, created_by, start_at, expires_at, created_at
`

type UpdateCouponParams struct {
	Title     pgtype.Text        `json:"title"`
	Code      pgtype.Text        `json:"code"`
	Percent   pgtype.Int4        `json:"percent"`
	StartAt   pgtype.Timestamptz `json:"start_at"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
	ID        int64              `json:"id"`
}

func (q *Queries) UpdateCoupon(ctx context.Context, arg UpdateCouponParams) (Coupon, error) {
	row := q.db.QueryRow(ctx, updateCoupon,
		arg.Title,
		arg.Code,
		arg.Percent,
		arg.StartAt,
		arg.ExpiresAt,
		arg.ID,
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
