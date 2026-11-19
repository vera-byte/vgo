-- 时间列默认值与更新触发器对齐为 'YYYY-MM-DD HH24:MI:SS'（now()::text 会带微秒与时区），
-- 并归一化存量脏格式数据
ALTER TABLE "task_info" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "task_info" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "task_info" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "task_info" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

ALTER TABLE "task_log" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "task_log" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "task_log" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "task_log" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

