CREATE SEQUENCE IF NOT EXISTS space_type_id_seq;

CREATE TABLE "space_type" (
  "id" integer NOT NULL DEFAULT nextval('space_type_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "name" character varying NOT NULL,
  "parentId" integer,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "space_type"."id" IS 'ID';

COMMENT ON COLUMN "space_type"."createTime" IS '创建时间';

COMMENT ON COLUMN "space_type"."updateTime" IS '更新时间';

COMMENT ON COLUMN "space_type"."tenantId" IS '租户ID';

COMMENT ON COLUMN "space_type"."name" IS '类别名称';

COMMENT ON COLUMN "space_type"."parentId" IS '父分类ID';

CREATE INDEX "IDX_0749b509b68488caecd4cc2bbc" ON "space_type" USING btree ("updateTime");

CREATE INDEX "IDX_5e7f846b8cdabbceba95ed3314" ON "space_type" USING btree ("tenantId");

CREATE INDEX "IDX_6669449501d275f367ca295472" ON "space_type" USING btree ("createTime");