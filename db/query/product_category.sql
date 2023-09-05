-- name: CreateProductCategory :one
INSERT INTO product_categories (
  product_id,
  category_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetProductCategory :one
SELECT * FROM product_categories
WHERE product_id = $1 AND category_id = $2
LIMIT 1;

-- name: ListProductCategoryByProductId :many
SELECT * FROM product_categories
WHERE product_id = $1;

-- name: ListProductCategoryByCategoryId :many
SELECT * FROM product_categories
WHERE category_id = $1;

-- name: DeleteProductCategoryByProductId :exec
DELETE FROM product_categories
WHERE product_id = $1;

-- name: DeleteProductCategoryByCategoryId :exec
DELETE FROM product_categories
WHERE category_id = $1;
