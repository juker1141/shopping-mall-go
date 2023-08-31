-- name: CreateCoupon :one
INSERT INTO coupons (
  title,
  code,
  percent,
  created_by,
  start_at,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetCoupon :one
SELECT * FROM coupons
WHERE id = $1 LIMIT 1;

-- name: ListCoupons :many
SELECT * FROM coupons
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
OFFSET $4;

-- name: UpdateCoupon :one
UPDATE coupons
SET 
  title = $2,
  code = $3,
  percent = $4,
  start_at = $5,
  expires_at = $6
WHERE id = $1
RETURNING *;

-- name: DeleteCoupon :exec
DELETE FROM coupons WHERE id = $1;