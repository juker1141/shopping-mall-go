-- name: CreateOrderProduct :one
INSERT INTO order_products (
  order_id,
  product_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetOrderProduct :one
SELECT * FROM order_products
WHERE order_id = $1 AND product_id = $2
LIMIT 1;

-- name: ListOrderProductByOrderId :many
SELECT * FROM order_products
WHERE order_id = $1;

-- name: ListOrderProductByProductId :many
SELECT * FROM order_products
WHERE product_id = $1;

-- name: DeleteOrderProductByOrderId :exec
DELETE FROM order_products
WHERE order_id = $1;

-- name: DeleteOrderProductByProductId :exec
DELETE FROM order_products
WHERE product_id = $1;
