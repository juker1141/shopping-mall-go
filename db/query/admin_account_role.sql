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

-- name: ListAdminUserRoles :many
SELECT * FROM admin_user_roles;

-- name: UpdateAdminUserRole :one
UPDATE admin_user_roles
SET admin_user_id = $1, role_id = $2
WHERE admin_user_id = $3 AND role_id = $4
RETURNING *;

-- name: DeleteAdminUserRole :exec
DELETE FROM admin_user_roles
WHERE admin_user_id = $1 AND role_id = $2;