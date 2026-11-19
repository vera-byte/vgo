DROP INDEX IF EXISTS "idx_base_sys_user_username";
CREATE UNIQUE INDEX "UQ_base_sys_user_username" ON "base_sys_user" ("username");
