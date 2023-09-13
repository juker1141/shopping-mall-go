
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
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "genders" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

INSERT INTO "genders" ("name") VALUES
  ('男生'),
  ('女生'),
  ('其他');

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
  "code" varchar NOT NULL,
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


CREATE INDEX ON "users" ("account");

CREATE INDEX ON "carts" ("owner");

CREATE INDEX ON "coupons" ("title");

CREATE INDEX ON "coupons" ("code");

CREATE INDEX ON "coupons" ("start_at");

CREATE INDEX ON "coupons" ("expires_at");


CREATE UNIQUE INDEX ON "cart_products" ("cart_id", "product_id");

CREATE UNIQUE INDEX ON "cart_coupons" ("cart_id", "coupon_id");


COMMENT ON COLUMN "users"."status" IS 'must be either 0 or 1';

COMMENT ON COLUMN "carts"."total_price" IS 'must be positive';

COMMENT ON COLUMN "carts"."final_price" IS 'must be positive';

COMMENT ON COLUMN "products"."status" IS 'must be either 0 or 1';


ALTER TABLE "users" ADD FOREIGN KEY ("gender_id") REFERENCES "genders" ("id");

ALTER TABLE "carts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("account");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("cart_id") REFERENCES "carts" ("id");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "cart_coupons" ADD FOREIGN KEY ("cart_id") REFERENCES "carts" ("id");

ALTER TABLE "cart_coupons" ADD FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id");