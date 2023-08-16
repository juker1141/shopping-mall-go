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


-- name: UpdateRole :one
UPDATE roles
SET
  name = COALESCE(sqlc.narg(name), name)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;