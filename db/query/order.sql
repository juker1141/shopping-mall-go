-- name: CreateOrder :one
INSERT INTO orders (
  status_id
) VALUES (
  $1
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- -- name: GetOrderByCode :one
-- SELECT * FROM orders
-- WHERE code = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET 
  is_paid = COALESCE(sqlc.narg(status_id), status_id),
  status_id = COALESCE(sqlc.narg(status_id), status_id),
  updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: GetOrdersCount :one
SELECT COUNT(*) FROM orders;