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
    WHEN sqlc.arg(key)::varchar = 'title' THEN title ILIKE '%' || sqlc.arg(key_value)::varchar || '%'
    WHEN sqlc.arg(key)::varchar = 'code' THEN code ILIKE '%' || sqlc.arg(key_value)::varchar || '%'
    ELSE TRUE
  END
ORDER BY id
LIMIT sqlc.arg('Limit')
OFFSET sqlc.arg('Offset');

-- name: UpdateCoupon :one
UPDATE coupons
SET 
  title = COALESCE(sqlc.narg(title), title),
  code = COALESCE(sqlc.narg(code), code),
  percent = COALESCE(sqlc.narg(percent), percent),
  start_at = COALESCE(sqlc.narg(start_at), start_at),
  expires_at = COALESCE(sqlc.narg(expires_at), expires_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCoupon :exec
DELETE FROM coupons WHERE id = $1;