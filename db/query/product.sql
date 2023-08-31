-- name: CreateProduct :one
INSERT INTO products (
  title,
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
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT p.*
FROM products AS p
JOIN product_categories AS pc ON p.id = pc.product_id
JOIN categories AS c ON pc.category_id = c.id
WHERE
  CASE
    WHEN $1::varchar = 'title' THEN p.title ILIKE '%' || $2::varchar || '%'
    WHEN $1::varchar = 'category' THEN c.name ILIKE '%' || $2::varchar || '%'
    ELSE TRUE
  END
ORDER BY p.id
LIMIT $3
OFFSET $4;

-- name: UpdateProduct :one
UPDATE products
SET
  title = COALESCE(sqlc.narg(title), title),
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