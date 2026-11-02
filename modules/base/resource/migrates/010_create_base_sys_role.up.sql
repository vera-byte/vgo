CREATE SEQUENCE IF NOT EXISTS base_sys_role_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_sys_role_id_seq;

CREATE TABLE base_sys_role (
  "id" integer NOT NULL DEFAULT nextval('base_sys_role_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "userId" integer,
  "name" character varying NOT NULL,
  "label" character varying,
  "remark" character varying,
  "relevance" integer DEFAULT 1,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_sys_role."userId" IS '创建者ID';

COMMENT ON COLUMN base_sys_role."name" IS '角色名称';

COMMENT ON COLUMN base_sys_role."label" IS '角色标识';

COMMENT ON COLUMN base_sys_role."remark" IS '备注';

COMMENT ON COLUMN base_sys_role."relevance" IS '关联';

CREATE INDEX idx_base_sys_role_label ON base_sys_role USING btree ("label");
