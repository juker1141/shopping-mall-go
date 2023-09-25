-- name: CreateOrder :one
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
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email),
  shipping_address = COALESCE(sqlc.narg(shipping_address), shipping_address),
  message = COALESCE(sqlc.narg(message), message),
  pay_method_id = COALESCE(sqlc.narg(pay_method_id), pay_method_id),
  is_paid = COALESCE(sqlc.narg(is_paid), is_paid),
  total_price = COALESCE(sqlc.narg(total_price), total_price),
  final_price = COALESCE(sqlc.narg(final_price), final_price),
  status_id = COALESCE(sqlc.narg(status_id), status_id),
  updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: GetOrdersCount :one
SELECT COUNT(*) FROM orders;