CREATE SEQUENCE IF NOT EXISTS base_eps_admin_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_eps_admin_id_seq;

CREATE TABLE base_eps_admin (
  "id" integer NOT NULL DEFAULT nextval('base_eps_admin_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "module" character varying NOT NULL,
  "method" character varying NOT NULL,
  "path" character varying NOT NULL,
  "prefix" character varying NOT NULL,
  "summary" character varying,
  "tag" character varying,
  "dts" text,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_eps_admin."module" IS '所属模块名称';

COMMENT ON COLUMN base_eps_admin."method" IS '请求方法';

COMMENT ON COLUMN base_eps_admin."path" IS '请求路径';

COMMENT ON COLUMN base_eps_admin."prefix" IS '路由前缀';

COMMENT ON COLUMN base_eps_admin."summary" IS '描述';

COMMENT ON COLUMN base_eps_admin."tag" IS '标签';

COMMENT ON COLUMN base_eps_admin."dts" IS '附加数据';

CREATE INDEX idx_base_eps_admin_module ON base_eps_admin USING btree ("module");

CREATE INDEX idx_base_eps_admin_prefix ON base_eps_admin USING btree ("prefix");

CREATE INDEX idx_base_eps_admin_path ON base_eps_admin USING btree ("path");
