-- 刪除外鍵約束
ALTER TABLE "role_permissions" DROP CONSTRAINT IF EXISTS "role_permissions_role_id_fkey";
ALTER TABLE "role_permissions" DROP CONSTRAINT IF EXISTS "role_permissions_permission_id_fkey";
ALTER TABLE "admin_user_roles" DROP CONSTRAINT IF EXISTS "admin_user_roles_admin_user_id_fkey";
ALTER TABLE "admin_user_roles" DROP CONSTRAINT IF EXISTS "admin_user_roles_role_id_fkey";

-- 刪除索引
DROP INDEX IF EXISTS "admin_users_account_idx";
DROP INDEX IF EXISTS "role_permissions_role_id_permission_id_idx";
DROP INDEX IF EXISTS "admin_user_roles_admin_user_id_role_id_idx";

-- 刪除關聯表
DROP TABLE IF EXISTS "role_permissions";
DROP TABLE IF EXISTS "admin_user_roles";

-- 刪除主要表
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "admin_users";
