-- name: CreateUser :one
INSERT INTO users (
  account,
  email,
  full_name,
  gender_id,
  phone,
  address,
  shipping_address,
  post_code,
  hashed_password,
  status,
  avatar_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByAccount :one
SELECT * FROM users
WHERE account = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  phone = COALESCE(sqlc.narg(phone), phone),
  address = COALESCE(sqlc.narg(address), address),
  shipping_address = COALESCE(sqlc.narg(shipping_address), shipping_address),
  post_code = COALESCE(sqlc.narg(post_code), post_code),
  avatar_url = COALESCE(sqlc.narg(avatar_url), avatar_url),
  status = COALESCE(sqlc.narg(status), status)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: GetUsersCount :one
SELECT COUNT(*) FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;