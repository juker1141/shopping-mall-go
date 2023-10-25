-- name: CreateCartProduct :one
INSERT INTO cart_products (
  cart_id,
  product_id,
  num
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCartProduct :one
SELECT * FROM cart_products
WHERE cart_id = $1 AND product_id = $2
LIMIT 1;

-- name: UpdateCartProduct :one
UPDATE cart_products
SET 
  num = COALESCE(sqlc.narg(num), num)
WHERE
  cart_id = sqlc.arg(cart_id) AND product_id = sqlc.arg(product_id)
RETURNING *;

-- name: ListCartProductByCartId :many
SELECT * FROM cart_products
WHERE cart_id = $1;

-- name: ListCartProductByProductId :many
SELECT * FROM cart_products
WHERE product_id = $1;

-- name: DeleteCartProduct :exec
DELETE FROM cart_products
WHERE cart_id = $1 AND product_id = $2;

-- name: DeleteCartProductByCartId :exec
DELETE FROM cart_products
WHERE cart_id = $1;

-- name: DeleteCartProductByProductId :exec
DELETE FROM cart_products
WHERE product_id = $1;
