-- -- name: CreateOrder :one
-- INSERT INTO orders (
-- ) VALUES (
-- ) RETURNING *;

-- -- name: GetOrder :one
-- SELECT * FROM orders
-- WHERE id = $1 LIMIT 1;

-- -- name: GetOrderByCode :one
-- SELECT * FROM orders
-- WHERE code = $1 LIMIT 1;

-- -- name: ListOrders :many
-- SELECT * FROM orders
-- WHERE
--   CASE
--     WHEN sqlc.arg(key)::varchar = 'title' THEN title ILIKE '%' || sqlc.arg(key_value)::varchar || '%'
--     WHEN sqlc.arg(key)::varchar = 'code' THEN code ILIKE '%' || sqlc.arg(key_value)::varchar || '%'
--     ELSE TRUE
--   END
-- ORDER BY id
-- LIMIT sqlc.arg('Limit')
-- OFFSET sqlc.arg('Offset');

-- -- name: UpdateOrder :one
-- UPDATE orders
-- SET 
--   title = COALESCE(sqlc.narg(title), title),
--   code = COALESCE(sqlc.narg(code), code),
--   percent = COALESCE(sqlc.narg(percent), percent),
--   start_at = COALESCE(sqlc.narg(start_at), start_at),
--   expires_at = COALESCE(sqlc.narg(expires_at), expires_at)
-- WHERE id = sqlc.arg(id)
-- RETURNING *;

-- -- name: DeleteOrder :exec
-- DELETE FROM orders WHERE id = $1;

-- -- name: GetOrdersCount :one
-- SELECT COUNT(*) FROM orders;