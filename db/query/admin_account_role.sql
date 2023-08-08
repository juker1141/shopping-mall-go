-- name: CreateAdminUserRole :one
INSERT INTO admin_user_roles (
  admin_user_id,
  role_id
)
VALUES (
  $1, $2
) RETURNING *;

-- name: GetAdminUserRole :one
SELECT * FROM admin_user_roles
WHERE admin_user_id = $1 AND role_id = $2
LIMIT 1;

-- name: ListAdminUserRoleByRoleId :many
SELECT * FROM admin_user_roles
WHERE role_id = $1;

-- name: ListAdminUserRoleByAdminUserId :many
SELECT * FROM admin_user_roles
WHERE admin_user_id = $1;

-- name: DeleteAdminUserRoleByRoleId :exec
DELETE FROM admin_user_roles
WHERE role_id = $1;

-- name: DeleteAdminUserRoleByAdminUserId :exec
DELETE FROM admin_user_roles
WHERE admin_user_id = $1;