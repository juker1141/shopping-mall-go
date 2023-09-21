-- name: CreateAdminUser :one
INSERT INTO admin_users (
  account,
  full_name,
  hashed_password,
  status,
  role_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetAdminUser :one
SELECT * FROM admin_users
WHERE id = $1 LIMIT 1;

-- name: GetAdminUserByAccount :one
SELECT * FROM admin_users
WHERE account = $1 LIMIT 1;

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
  status = COALESCE(sqlc.narg(status), status),
  role_id = COALESCE(sqlc.narg(role_id), role_id)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: ListPermissionsForAdminUser :many
SELECT DISTINCT p.id, p.name, p.created_at
FROM admin_users AS au
JOIN roles AS r ON au.role_id = r.id
JOIN role_permissions AS rp ON r.id = rp.role_id
JOIN permissions AS p ON rp.permission_id = p.id
WHERE au.id = $1;

-- name: ListPermissionsIDByAccount :many
SELECT p.id
FROM admin_users AS au
JOIN roles AS r ON au.role_id = r.id
JOIN role_permissions AS rp ON r.id = rp.role_id
JOIN permissions AS p ON rp.permission_id = p.id
WHERE au.account = $1;

-- name: GetAdminUsersCount :one
SELECT COUNT(*) FROM admin_users;

-- name: DeleteAdminUser :exec
DELETE FROM admin_users WHERE id = $1;