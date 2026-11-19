-- 时间列默认值与更新触发器对齐为 'YYYY-MM-DD HH24:MI:SS'（now()::text 会带微秒与时区），
-- 并归一化存量脏格式数据
ALTER TABLE "base_eps_admin" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_eps_admin" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_eps_admin" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_eps_admin" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_eps_app" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_eps_app" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_eps_app" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_eps_app" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_conf" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_conf" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_conf" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_conf" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_department" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_department" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_department" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_department" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_log" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_log" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_log" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_log" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_menu" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_menu" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_menu" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_menu" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_param" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_param" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_param" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_param" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_role" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_role" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_role" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_role" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_role_department" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_role_department" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_role_department" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_role_department" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_role_menu" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_role_menu" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_role_menu" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_role_menu" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_user" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_user" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_user" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_user" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "base_sys_user_role" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "base_sys_user_role" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "base_sys_user_role" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "base_sys_user_role" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;
