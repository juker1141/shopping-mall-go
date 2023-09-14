-- name: CreateGender :one
INSERT INTO genders (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetGender :one
SELECT * FROM genders
WHERE id = $1 LIMIT 1;

-- name: ListGenders :many
SELECT * FROM genders
ORDER BY id;