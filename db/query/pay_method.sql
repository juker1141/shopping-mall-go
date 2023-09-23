-- name: CreatePayMethod :one
INSERT INTO pay_methods (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetPayMethod :one
SELECT * FROM pay_methods
WHERE id = $1 LIMIT 1;

-- name: ListPayMethod :many
SELECT * FROM pay_methods
ORDER BY id;

-- name: GetPayMethodCount :one
SELECT COUNT(*) FROM pay_methods;