-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-10-02T06:01:05.767Z

CREATE TABLE "permissions" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "admin_users" (
  "id" bigserial PRIMARY KEY,
  "account" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "role_id" int,
  "status" int NOT NULL DEFAULT 1,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "role_permissions" (
  "role_id" int,
  "permission_id" int
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "account" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "account" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "gender_id" int,
  "cellphone" varchar NOT NULL,
  "address" varchar NOT NULL,
  "shipping_address" varchar NOT NULL,
  "post_code" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "status" int NOT NULL DEFAULT 1,
  "avatar_url" varchar NOT NULL,
  "is_email_verified" bool NOT NULL DEFAULT false,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "user_id" int,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "genders" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "carts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar,
  "total_price" int NOT NULL,
  "final_price" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "coupons" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "code" varchar UNIQUE NOT NULL,
  "percent" int NOT NULL,
  "created_by" varchar NOT NULL,
  "start_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL DEFAULT '2100-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "category" varchar NOT NULL,
  "origin_price" int NOT NULL,
  "price" int NOT NULL,
  "unit" varchar NOT NULL,
  "description" varchar NOT NULL,
  "content" varchar NOT NULL,
  "status" int NOT NULL DEFAULT 1,
  "image_url" varchar NOT NULL,
  "images_url" varchar[],
  "created_by" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cart_products" (
  "cart_id" int,
  "product_id" int,
  "num" int NOT NULL
);

CREATE TABLE "cart_coupons" (
  "cart_id" int,
  "coupon_id" int
);

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
);

CREATE TABLE "order_coupons" (
  "order_id" int,
  "coupon_id" int
);

CREATE INDEX ON "admin_users" ("account");

CREATE UNIQUE INDEX ON "role_permissions" ("role_id", "permission_id");

CREATE INDEX ON "users" ("account");

CREATE INDEX ON "carts" ("owner");

CREATE INDEX ON "coupons" ("title");

CREATE INDEX ON "coupons" ("code");

CREATE INDEX ON "coupons" ("start_at");

CREATE INDEX ON "coupons" ("expires_at");

CREATE UNIQUE INDEX ON "cart_products" ("cart_id", "product_id");

CREATE UNIQUE INDEX ON "cart_coupons" ("cart_id", "coupon_id");

CREATE UNIQUE INDEX ON "order_users" ("order_id", "user_id");

CREATE UNIQUE INDEX ON "order_products" ("order_id", "product_id");

CREATE UNIQUE INDEX ON "order_coupons" ("order_id", "coupon_id");

COMMENT ON COLUMN "admin_users"."status" IS 'must be either 0 or 1';

COMMENT ON COLUMN "users"."status" IS 'must be either 0 or 1';

COMMENT ON COLUMN "carts"."total_price" IS 'must be positive';

COMMENT ON COLUMN "carts"."final_price" IS 'must be positive';

COMMENT ON COLUMN "products"."status" IS 'must be either 0 or 1';

COMMENT ON COLUMN "orders"."total_price" IS 'must be positive';

COMMENT ON COLUMN "orders"."final_price" IS 'must be positive';

ALTER TABLE "admin_users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("account") REFERENCES "admin_users" ("account");

ALTER TABLE "users" ADD FOREIGN KEY ("gender_id") REFERENCES "genders" ("id");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("account");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("cart_id") REFERENCES "carts" ("id");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "cart_coupons" ADD FOREIGN KEY ("cart_id") REFERENCES "carts" ("id");

ALTER TABLE "cart_coupons" ADD FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("pay_method_id") REFERENCES "pay_methods" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("status_id") REFERENCES "order_status" ("id");

ALTER TABLE "order_users" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_coupons" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_coupons" ADD FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id");
