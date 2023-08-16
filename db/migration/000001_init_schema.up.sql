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
  "status" int NOT NULL DEFAULT 1,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "role_permissions" (
  "role_id" int,
  "permission_id" int
);

CREATE TABLE "admin_user_roles" (
  "admin_user_id" int,
  "role_id" int
);

CREATE INDEX ON "admin_users" ("account");

CREATE UNIQUE INDEX ON "role_permissions" ("role_id", "permission_id");

CREATE UNIQUE INDEX ON "admin_user_roles" ("admin_user_id", "role_id");


COMMENT ON COLUMN "admin_users"."status" IS 'must be either 0 or 1';

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "admin_user_roles" ADD FOREIGN KEY ("admin_user_id") REFERENCES "admin_users" ("id");

ALTER TABLE "admin_user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");
