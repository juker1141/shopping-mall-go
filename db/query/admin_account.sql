-- name: CreateAdminUser :one
INSERT INTO admin_users (
  account,
  full_name,
  hashed_password,
  status
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAdminUser :one
SELECT * FROM admin_users
WHERE id = $1 LIMIT 1;

-- name: ListAdminUsers :many
SELECT * FROM admin_users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAdminUser :one
UPDATE admin_users
SET 
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  status = COALESCE(sqlc.narg(status), status)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAdminUser :exec
DELETE FROM admin_users WHERE id = $1;