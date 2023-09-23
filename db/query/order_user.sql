-- name: CreateOrderUser :one
INSERT INTO order_users (
  order_id,
  user_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetOrderUser :one
SELECT * FROM order_users
WHERE order_id = $1 AND user_id = $2
LIMIT 1;

-- name: ListOrderUserByOrderId :many
SELECT * FROM order_users
WHERE order_id = $1;

-- name: ListOrderUserByUserId :many
SELECT * FROM order_users
WHERE user_id = $1;

-- name: DeleteOrderUserByOrderId :exec
DELETE FROM order_users
WHERE order_id = $1;

-- name: DeleteOrderUserByUserId :exec
DELETE FROM order_users
WHERE user_id = $1;
