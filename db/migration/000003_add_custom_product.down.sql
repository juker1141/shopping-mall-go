-- Drop foreign keys
ALTER TABLE "product_categorys" DROP CONSTRAINT "product_categorys_category_id_fkey";
ALTER TABLE "product_categorys" DROP CONSTRAINT "product_categorys_product_id_fkey";
ALTER TABLE "cart_coupons" DROP CONSTRAINT "cart_coupons_coupon_id_fkey";
ALTER TABLE "cart_coupons" DROP CONSTRAINT "cart_coupons_cart_id_fkey";
ALTER TABLE "cart_products" DROP CONSTRAINT "cart_products_product_id_fkey";
ALTER TABLE "cart_products" DROP CONSTRAINT "cart_products_cart_id_fkey";
ALTER TABLE "carts" DROP CONSTRAINT "carts_owner_fkey";
ALTER TABLE "users" DROP CONSTRAINT "users_gender_id_fkey";

-- Drop tables
DROP TABLE "product_categorys";
DROP TABLE "cart_coupons";
DROP TABLE "cart_products";
DROP TABLE "categorys";
DROP TABLE "products";
DROP TABLE "coupons";
DROP TABLE "carts";
DROP TABLE "genders";
DROP TABLE "users";
