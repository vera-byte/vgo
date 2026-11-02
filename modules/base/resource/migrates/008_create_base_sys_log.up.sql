CREATE SEQUENCE IF NOT EXISTS base_sys_log_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_sys_log_id_seq;

CREATE TABLE base_sys_log (
  "id" integer NOT NULL DEFAULT nextval('base_sys_log_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "userId" integer,
  "action" character varying,
  "ip" character varying,
  "ipAddr" character varying,
  "params" text,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_sys_log."userId" IS '用户ID';

COMMENT ON COLUMN base_sys_log."action" IS '操作';

COMMENT ON COLUMN base_sys_log."ip" IS 'IP地址';

COMMENT ON COLUMN base_sys_log."ipAddr" IS 'IP归属地';

COMMENT ON COLUMN base_sys_log."params" IS '请求参数';

CREATE INDEX idx_base_sys_log_userId ON base_sys_log USING btree ("userId");

CREATE INDEX idx_base_sys_log_action ON base_sys_log USING btree ("action");
