CREATE SEQUENCE IF NOT EXISTS base_sys_conf_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_sys_conf_id_seq;

CREATE TABLE base_sys_conf (
  "id" integer NOT NULL DEFAULT nextval('base_sys_conf_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "cKey" character varying NOT NULL,
  "cValue" text,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_sys_conf."cKey" IS '配置键';

COMMENT ON COLUMN base_sys_conf."cValue" IS '配置值';

CREATE INDEX idx_base_sys_conf_cKey ON base_sys_conf USING btree ("cKey");
