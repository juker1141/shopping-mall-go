-- name: CreateOrderCoupon :one
INSERT INTO order_coupons (
  order_id,
  coupon_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetOrderCoupon :one
SELECT * FROM order_coupons
WHERE order_id = $1 AND coupon_id = $2
LIMIT 1;

-- name: ListOrderCouponByOrderId :many
SELECT * FROM order_coupons
WHERE order_id = $1;

-- name: ListOrderCouponByCouponId :many
SELECT * FROM order_coupons
WHERE coupon_id = $1;

-- name: UpdateOrderCouponByOrderId :one
UPDATE order_coupons
SET 
  coupon_id = $2
WHERE order_id = $1
RETURNING *;

-- name: DeleteOrderCouponByOrderId :exec
DELETE FROM order_coupons
WHERE order_id = $1;

-- name: DeleteOrderCouponByCouponId :exec
DELETE FROM order_coupons
WHERE coupon_id = $1;
