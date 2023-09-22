-- name: CreateRole :one
INSERT INTO roles (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListRolesOption :many
SELECT * FROM roles
ORDER BY id;

-- name: GetRolesCount :one
SELECT COUNT(*) FROM roles;

-- name: ListPermissionsForRole :many
SELECT DISTINCT p.id, p.name, p.created_at
FROM roles AS r
JOIN role_permissions AS rp ON r.id = rp.role_id
JOIN permissions AS p ON rp.permission_id = p.id
WHERE r.id = $1;

-- name: UpdateRole :one
UPDATE roles
SET
  name = COALESCE(sqlc.narg(name), name)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;