CREATE SEQUENCE IF NOT EXISTS space_info_id_seq;

CREATE TABLE "space_info" (
  "id" integer NOT NULL DEFAULT nextval('space_info_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "url" character varying NOT NULL,
  "type" character varying NOT NULL,
  "classifyId" integer,
  "fileId" character varying NOT NULL,
  "name" character varying NOT NULL,
  "size" integer NOT NULL,
  "version" integer NOT NULL DEFAULT 1,
  "key" character varying NOT NULL,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "space_info"."id" IS 'ID';

COMMENT ON COLUMN "space_info"."createTime" IS '创建时间';

COMMENT ON COLUMN "space_info"."updateTime" IS '更新时间';

COMMENT ON COLUMN "space_info"."tenantId" IS '租户ID';

COMMENT ON COLUMN "space_info"."url" IS '地址';

COMMENT ON COLUMN "space_info"."type" IS '类型';

COMMENT ON COLUMN "space_info"."classifyId" IS '分类ID';

COMMENT ON COLUMN "space_info"."fileId" IS '文件id';

COMMENT ON COLUMN "space_info"."name" IS '文件名';

COMMENT ON COLUMN "space_info"."size" IS '文件大小';

COMMENT ON COLUMN "space_info"."version" IS '文档版本';

COMMENT ON COLUMN "space_info"."key" IS '文件位置';

CREATE INDEX "IDX_0975633032bfe6574468b3a4ae" ON "space_info" USING btree ("fileId");

CREATE INDEX "IDX_6001c5ed2088b893c0d69bb244" ON "space_info" USING btree ("tenantId");

CREATE INDEX "IDX_d7a2539961e9aacba8b353f3c9" ON "space_info" USING btree ("updateTime");

CREATE INDEX "IDX_eb1da2f304c760846b5add09b3" ON "space_info" USING btree ("createTime");