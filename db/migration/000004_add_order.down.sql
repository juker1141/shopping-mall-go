-- 移除外鍵約束
ALTER TABLE "order_coupons" DROP CONSTRAINT "order_coupons_coupon_id_fkey";
ALTER TABLE "order_coupons" DROP CONSTRAINT "order_coupons_order_id_fkey";
ALTER TABLE "order_products" DROP CONSTRAINT "order_products_product_id_fkey";
ALTER TABLE "order_products" DROP CONSTRAINT "order_products_order_id_fkey";
ALTER TABLE "order_users" DROP CONSTRAINT "order_users_user_id_fkey";
ALTER TABLE "order_users" DROP CONSTRAINT "order_users_order_id_fkey";

-- 移除索引
DROP INDEX "order_coupons_order_id_coupon_id_idx";
DROP INDEX "order_products_order_id_product_id_idx";
DROP INDEX "order_users_order_id_user_id_idx";

-- 移除表格
DROP TABLE "order_coupons";
DROP TABLE "order_products";
DROP TABLE "order_users";
DROP TABLE "orders";
