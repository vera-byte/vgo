CREATE SEQUENCE IF NOT EXISTS task_info_id_seq;

CREATE TABLE "task_info" (
  "id" integer NOT NULL DEFAULT nextval('task_info_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "jobId" character varying,
  "repeatConf" character varying(1000),
  "name" character varying NOT NULL,
  "cron" character varying,
  "limit" integer,
  "every" integer,
  "remark" character varying,
  "status" integer NOT NULL DEFAULT 1,
  "startDate" timestamp without time zone,
  "endDate" timestamp without time zone,
  "data" character varying,
  "service" character varying,
  "type" integer NOT NULL DEFAULT 0,
  "nextRunTime" timestamp without time zone,
  "taskType" integer NOT NULL DEFAULT 0,
  "lastExecuteTime" timestamp without time zone,
  "lockExpireTime" timestamp without time zone,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "task_info"."id" IS 'ID';

COMMENT ON COLUMN "task_info"."createTime" IS '创建时间';

COMMENT ON COLUMN "task_info"."updateTime" IS '更新时间';

COMMENT ON COLUMN "task_info"."tenantId" IS '租户ID';

COMMENT ON COLUMN "task_info"."jobId" IS '任务ID';

COMMENT ON COLUMN "task_info"."repeatConf" IS '任务配置';

COMMENT ON COLUMN "task_info"."name" IS '名称';

COMMENT ON COLUMN "task_info"."cron" IS 'cron';

COMMENT ON COLUMN "task_info"."limit" IS '最大执行次数 不传为无限次';

COMMENT ON COLUMN "task_info"."every" IS '每间隔多少毫秒执行一次 如果cron设置了 这项设置就无效';

COMMENT ON COLUMN "task_info"."remark" IS '备注';

COMMENT ON COLUMN "task_info"."status" IS '状态 0-停止 1-运行';

COMMENT ON COLUMN "task_info"."startDate" IS '开始时间';

COMMENT ON COLUMN "task_info"."endDate" IS '结束时间';

COMMENT ON COLUMN "task_info"."data" IS '数据';

COMMENT ON COLUMN "task_info"."service" IS '执行的service实例ID';

COMMENT ON COLUMN "task_info"."type" IS '状态 0-系统 1-用户';

COMMENT ON COLUMN "task_info"."nextRunTime" IS '下一次执行时间';

COMMENT ON COLUMN "task_info"."taskType" IS '状态 0-cron 1-时间间隔';

CREATE INDEX "IDX_11b991dc4a7a5585c636008d3a" ON "task_info" USING btree ("tenantId");

CREATE INDEX "IDX_2adc6f9c241391126f27dac145" ON "task_info" USING btree ("updateTime");

CREATE INDEX "IDX_6ced02f467e59bd6306b549bb0" ON "task_info" USING btree ("createTime");