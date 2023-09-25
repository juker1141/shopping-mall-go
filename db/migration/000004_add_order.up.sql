
CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "shipping_address" varchar NOT NULL,
  "message" varchar,
  "is_paid" bool NOT NULL DEFAULT false,
  "total_price" int NOT NULL,
  "final_price" int NOT NULL,
  "pay_method_id" int NOT NULL,
  "status_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "pay_methods" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "order_status" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "description" varchar NOT NULL
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

COMMENT ON COLUMN "orders"."total_price" IS 'must be positive';

COMMENT ON COLUMN "orders"."final_price" IS 'must be positive';

CREATE UNIQUE INDEX ON "order_users" ("order_id", "user_id");

CREATE UNIQUE INDEX ON "order_products" ("order_id", "product_id");

CREATE UNIQUE INDEX ON "order_coupons" ("order_id", "coupon_id");

ALTER TABLE "order_users" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_coupons" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_coupons" ADD FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("status_id") REFERENCES "order_status" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("pay_method_id") REFERENCES "pay_methods" ("id");

