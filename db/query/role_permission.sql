-- name: CreateRolePermission :one
INSERT INTO role_permissions (
  role_id,
  permission_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetRolePermission :one
SELECT * FROM role_permissions
WHERE role_id = $1 AND permission_id = $2
LIMIT 1;

-- name: ListRolePermissions :many
SELECT * FROM role_permissions;

-- name: UpdateRolePermission :one
UPDATE role_permissions
SET role_id = $1, permission_id = $2
WHERE role_id = $3 AND permission_id = $4
RETURNING *;

-- name: DeleteRolePermission :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2;
