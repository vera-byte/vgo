CREATE SEQUENCE IF NOT EXISTS base_sys_param_id_seq;

CREATE TABLE "base_sys_param" (
  "id" integer NOT NULL DEFAULT nextval('base_sys_param_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "keyName" character varying NOT NULL,
  "name" character varying NOT NULL,
  "data" text NOT NULL,
  "dataType" integer NOT NULL DEFAULT 0,
  "remark" character varying,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "base_sys_param"."id" IS 'ID';

COMMENT ON COLUMN "base_sys_param"."createTime" IS '创建时间';

COMMENT ON COLUMN "base_sys_param"."updateTime" IS '更新时间';

COMMENT ON COLUMN "base_sys_param"."tenantId" IS '租户ID';

COMMENT ON COLUMN "base_sys_param"."keyName" IS '键';

COMMENT ON COLUMN "base_sys_param"."name" IS '名称';

COMMENT ON COLUMN "base_sys_param"."data" IS '数据';

COMMENT ON COLUMN "base_sys_param"."dataType" IS '数据类型 0-字符串 1-富文本 2-文件 ';

COMMENT ON COLUMN "base_sys_param"."remark" IS '备注';

CREATE INDEX "IDX_479122e3bf464112f7a7253dac" ON "base_sys_param" USING btree ("updateTime");

CREATE INDEX "IDX_7bcb57371b481d8e2d66ddeaea" ON "base_sys_param" USING btree ("createTime");

CREATE INDEX "IDX_8a0ab598ca7d63475356ca1157" ON "base_sys_param" USING btree ("tenantId");

CREATE UNIQUE INDEX "IDX_cf19b5e52d8c71caa9c4534454" ON "base_sys_param" USING btree ("keyName");