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

-- name: ListRolePermissionByRoleId :many
SELECT * FROM role_permissions
WHERE role_id = $1;

-- name: ListRolePermissionByPermissionId :many
SELECT * FROM role_permissions
WHERE permission_id = $1;

-- name: DeleteRolePermissionByRoleId :exec
DELETE FROM role_permissions
WHERE role_id = $1;

-- name: DeleteRolePermissionByPermissionId :exec
DELETE FROM role_permissions
WHERE permission_id = $1;
