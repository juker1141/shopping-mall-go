-- Drop foreign keys
ALTER TABLE "cart_coupons" DROP CONSTRAINT IF EXISTS "cart_coupons_coupon_id_fkey";
ALTER TABLE "cart_coupons" DROP CONSTRAINT IF EXISTS "cart_coupons_cart_id_fkey";
ALTER TABLE "cart_products" DROP CONSTRAINT IF EXISTS "cart_products_product_id_fkey";
ALTER TABLE "cart_products" DROP CONSTRAINT IF EXISTS "cart_products_cart_id_fkey";
ALTER TABLE "carts" DROP CONSTRAINT IF EXISTS "carts_owner_fkey";
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_gender_id_fkey";
ALTER TABLE "coupons" DROP CONSTRAINT "check_dates";

-- Drop indexes
DROP INDEX IF EXISTS "cart_coupons_cart_id_coupon_id_idx";
DROP INDEX IF EXISTS "cart_products_cart_id_product_id_idx";
DROP INDEX IF EXISTS "coupons_expires_at_idx";
DROP INDEX IF EXISTS "coupons_start_at_idx";
DROP INDEX IF EXISTS "coupons_code_idx";
DROP INDEX IF EXISTS "coupons_title_idx";
DROP INDEX IF EXISTS "carts_owner_idx";
DROP INDEX IF EXISTS "users_account_idx";

-- Drop tables
DROP TABLE IF EXISTS "cart_coupons";
DROP TABLE IF EXISTS "cart_products";
DROP TABLE IF EXISTS "products";
DROP TABLE IF EXISTS "coupons";
DROP TABLE IF EXISTS "carts";
DROP TABLE IF EXISTS "genders";
DROP TABLE IF EXISTS "users";

