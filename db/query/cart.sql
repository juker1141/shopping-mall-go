-- name: CreateCart :one
INSERT INTO carts (
  owner,
  total_price,
  final_price
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCart :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1;

-- name: GetCartByOwner :one
SELECT * FROM carts
WHERE owner = $1 LIMIT 1;

-- -- name: ListProductsForCart :many
-- SELECT DISTINCT product.id, product.name, product.created_at
-- FROM carts as c
-- JOIN cart_products AS cp ON c.id = cp.role_id
-- JOIN products AS product ON cp.product_id = product.id
-- WHERE c.id = $1;

-- name: UpdateCart :one
UPDATE carts
SET 
  total_price = $2,
  final_price = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCart :exec
DELETE FROM carts WHERE id = $1;