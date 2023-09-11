-- name: CreateProduct :one
INSERT INTO products (
  title,
  category,
  description,
  content,
  origin_price,
  price,
  unit,
  status,
  image_url,
  images_url,
  created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT p.*
FROM products AS p
WHERE
  CASE
    WHEN sqlc.arg(key)::varchar = 'title' THEN p.title ILIKE '%' || sqlc.arg(key_value)::varchar || '%'
    ELSE TRUE
  END
ORDER BY p.id
LIMIT sqlc.arg('Limit')
OFFSET sqlc.arg('Offset');

-- name: UpdateProduct :one
UPDATE products
SET
  title = COALESCE(sqlc.narg(title), title),
  category = COALESCE(sqlc.narg(category), category),
  description = COALESCE(sqlc.narg(description), description),
  content = COALESCE(sqlc.narg(content), content),
  origin_price = COALESCE(sqlc.narg(origin_price), origin_price),
  price = COALESCE(sqlc.narg(price), price),
  unit = COALESCE(sqlc.narg(unit), unit),
  status = COALESCE(sqlc.narg(status), status),
  image_url = COALESCE(sqlc.narg(image_url), image_url),
  images_url = COALESCE(sqlc.narg(images_url), images_url)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProductsCount :one
SELECT COUNT(*) FROM products;