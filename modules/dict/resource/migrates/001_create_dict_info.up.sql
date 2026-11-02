CREATE SEQUENCE IF NOT EXISTS dict_info_id_seq;

CREATE TABLE "dict_info" (
  "id" integer NOT NULL DEFAULT nextval('dict_info_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "typeId" integer NOT NULL,
  "name" character varying NOT NULL,
  "value" character varying,
  "orderNum" integer NOT NULL DEFAULT 0,
  "remark" character varying,
  "parentId" integer,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "dict_info"."id" IS 'ID';

COMMENT ON COLUMN "dict_info"."createTime" IS '创建时间';

COMMENT ON COLUMN "dict_info"."updateTime" IS '更新时间';

COMMENT ON COLUMN "dict_info"."tenantId" IS '租户ID';

COMMENT ON COLUMN "dict_info"."typeId" IS '类型ID';

COMMENT ON COLUMN "dict_info"."name" IS '名称';

COMMENT ON COLUMN "dict_info"."value" IS '值';

COMMENT ON COLUMN "dict_info"."orderNum" IS '排序';

COMMENT ON COLUMN "dict_info"."remark" IS '备注';

COMMENT ON COLUMN "dict_info"."parentId" IS '父ID';

CREATE INDEX "IDX_10362a62adbf120821fff209d8" ON "dict_info" USING btree ("updateTime");

CREATE INDEX "IDX_5c311a4af30de1181a5d7a7cc2" ON "dict_info" USING btree ("createTime");

CREATE INDEX "IDX_c26dc4b1ccb26e642191995edd" ON "dict_info" USING btree ("tenantId");