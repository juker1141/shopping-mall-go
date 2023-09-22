-- 刪除外鍵約束
ALTER TABLE "sessions" DROP CONSTRAINT IF EXISTS "sessions_account_fkey";

-- 刪除表格
DROP TABLE IF EXISTS "sessions";