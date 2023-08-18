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

-- name: ListPermissionForAdminUser :many
SELECT DISTINCT p.id, p.name, p.created_at
FROM admin_users AS au
JOIN admin_user_roles AS aur ON au.id = aur.admin_user_id
JOIN roles AS r ON aur.role_id = r.id
JOIN role_permissions AS rp ON r.id = rp.role_id
JOIN permissions AS p ON rp.permission_id = p.id
WHERE au.id = $1;

-- name: DeleteAdminUser :exec
DELETE FROM admin_users WHERE id = $1;