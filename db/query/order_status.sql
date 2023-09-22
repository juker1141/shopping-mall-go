-- name: CreateOrderStatus :one
INSERT INTO order_status (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetOrderStatus :one
SELECT * FROM order_status
WHERE id = $1 LIMIT 1;

-- name: ListOrderStatus :many
SELECT * FROM order_status
ORDER BY id;