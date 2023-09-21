
CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "is_paid" bool DEFAULT false,
  "status" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_users" (
  "order_id" int,
  "user_id" int
);

CREATE TABLE "order_products" (
  "order_id" int,
  "product_id" int,
  "num" int NOT NULL DEFAULT 1

  CONSTRAINT check_nums CHECK (num > 0)
);

CREATE TABLE "order_coupons" (
  "order_id" int,
  "coupon_id" int
);

CREATE UNIQUE INDEX ON "order_users" ("order_id", "user_id");

CREATE UNIQUE INDEX ON "order_products" ("order_id", "product_id");

CREATE UNIQUE INDEX ON "order_coupons" ("order_id", "coupon_id");

ALTER TABLE "order_users" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_coupons" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_coupons" ADD FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id");
