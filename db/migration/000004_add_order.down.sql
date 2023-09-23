-- 刪除外鍵約束
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_status_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "pay_methods_id_fkey";
ALTER TABLE "order_users" DROP CONSTRAINT IF EXISTS "order_users_order_id_fkey";
ALTER TABLE "order_users" DROP CONSTRAINT IF EXISTS "order_users_user_id_fkey";
ALTER TABLE "order_products" DROP CONSTRAINT IF EXISTS "order_products_order_id_fkey";
ALTER TABLE "order_products" DROP CONSTRAINT IF EXISTS "order_products_product_id_fkey";
ALTER TABLE "order_coupons" DROP CONSTRAINT IF EXISTS "order_coupons_order_id_fkey";
ALTER TABLE "order_coupons" DROP CONSTRAINT IF EXISTS "order_coupons_coupon_id_fkey";

-- 刪除索引

DROP INDEX IF EXISTS "order_users_order_id_user_id_idx";
DROP INDEX IF EXISTS "order_products_order_id_product_id_idx";
DROP INDEX IF EXISTS "order_coupons_order_id_coupon_id_idx";

-- 刪除表格
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "pay_methods";
DROP TABLE IF EXISTS "order_status";
DROP TABLE IF EXISTS "order_users";
DROP TABLE IF EXISTS "order_products";
DROP TABLE IF EXISTS "order_coupons";
