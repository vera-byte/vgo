-- 时间列默认值与更新触发器对齐为 'YYYY-MM-DD HH24:MI:SS'（now()::text 会带微秒与时区），
-- 并归一化存量脏格式数据
ALTER TABLE "recycle_data" ALTER COLUMN "createTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
ALTER TABLE "recycle_data" ALTER COLUMN "updateTime" SET DEFAULT to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
UPDATE "recycle_data" SET "createTime" = left("createTime", 19) WHERE length("createTime") > 19;
UPDATE "recycle_data" SET "updateTime" = left("updateTime", 19) WHERE length("updateTime") > 19;

