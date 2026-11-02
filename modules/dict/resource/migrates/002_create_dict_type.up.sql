CREATE SEQUENCE IF NOT EXISTS dict_type_id_seq;

CREATE TABLE "dict_type" (
  "id" integer NOT NULL DEFAULT nextval('dict_type_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "name" character varying NOT NULL,
  "key" character varying NOT NULL,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "dict_type"."id" IS 'ID';

COMMENT ON COLUMN "dict_type"."createTime" IS '创建时间';

COMMENT ON COLUMN "dict_type"."updateTime" IS '更新时间';

COMMENT ON COLUMN "dict_type"."tenantId" IS '租户ID';

COMMENT ON COLUMN "dict_type"."name" IS '名称';

COMMENT ON COLUMN "dict_type"."key" IS '标识';

CREATE INDEX "IDX_69734e5c2d29cc2139d5078f2c" ON "dict_type" USING btree ("createTime");

CREATE INDEX "IDX_6cccb2e33846cd354e8dc0e0ef" ON "dict_type" USING btree ("updateTime");

CREATE INDEX "IDX_7d4f3d2336e1afdda38278a07e" ON "dict_type" USING btree ("tenantId");