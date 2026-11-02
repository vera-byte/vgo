CREATE SEQUENCE IF NOT EXISTS task_log_id_seq;

CREATE SEQUENCE IF NOT EXISTS task_log_id_seq;

CREATE TABLE task_log (
  "id" integer NOT NULL DEFAULT nextval('task_log_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "taskId" integer,
  "status" integer NOT NULL,
  "detail" text,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN task_log."taskId" IS '任务ID';

COMMENT ON COLUMN task_log."status" IS '状态 1:成功 0:失败';

COMMENT ON COLUMN task_log."detail" IS '执行详情';

CREATE INDEX idx_task_log_taskId ON task_log USING btree ("taskId");

CREATE INDEX idx_task_log_status ON task_log USING btree ("status");
