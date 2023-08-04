-- name: CreatePermission :one
INSERT INTO permissions (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetPermission :one
SELECT * FROM permissions
WHERE id = $1 LIMIT 1;

-- name: ListPermissions :many
SELECT * FROM permissions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePermission :one
UPDATE permissions
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;