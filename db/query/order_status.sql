-- name: CreateOrderStatus :one
INSERT INTO order_status (
  name,
  description
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetOrderStatus :one
SELECT * FROM order_status
WHERE id = $1 LIMIT 1;

-- name: ListOrderStatus :many
SELECT * FROM order_status
ORDER BY id;

-- name: GetOrderStatusCount :one
SELECT COUNT(*) FROM order_status;