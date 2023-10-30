-- name: CreateCartCoupon :one
INSERT INTO cart_coupons (
  cart_id,
  coupon_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetCartCoupon :one
SELECT * FROM cart_coupons
WHERE cart_id = $1 AND coupon_id = $2
LIMIT 1;

-- name: CheckCartCouponExists :one
SELECT EXISTS (
  SELECT 1
  FROM cart_coupons
  WHERE cart_id = $1
);

-- name: ListCartCouponByCartId :many
SELECT * FROM cart_coupons
WHERE cart_id = $1;

-- name: ListCartCouponByCouponId :many
SELECT * FROM cart_coupons
WHERE coupon_id = $1;

-- name: DeleteCartCouponByCartId :exec
DELETE FROM cart_coupons
WHERE cart_id = $1;

-- name: DeleteCartCouponByCouponId :exec
DELETE FROM cart_coupons
WHERE coupon_id = $1;
