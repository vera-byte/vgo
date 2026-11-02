CREATE SEQUENCE IF NOT EXISTS base_sys_department_id_seq;

CREATE TABLE "base_sys_department" (
  "id" integer NOT NULL DEFAULT nextval('base_sys_department_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "name" character varying NOT NULL,
  "userId" integer,
  "parentId" integer,
  "orderNum" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "base_sys_department"."id" IS 'ID';

COMMENT ON COLUMN "base_sys_department"."createTime" IS '创建时间';

COMMENT ON COLUMN "base_sys_department"."updateTime" IS '更新时间';

COMMENT ON COLUMN "base_sys_department"."tenantId" IS '租户ID';

COMMENT ON COLUMN "base_sys_department"."name" IS '部门名称';

COMMENT ON COLUMN "base_sys_department"."userId" IS '创建者ID';

COMMENT ON COLUMN "base_sys_department"."parentId" IS '上级部门ID';

COMMENT ON COLUMN "base_sys_department"."orderNum" IS '排序';

CREATE INDEX "IDX_be4c53cd671384fa588ca9470a" ON "base_sys_department" USING btree ("createTime");

CREATE INDEX "IDX_ca1473a793961ec55bc0c8d268" ON "base_sys_department" USING btree ("updateTime");

CREATE INDEX "IDX_f19e8ffd9c62ddb17e76c8b9d7" ON "base_sys_department" USING btree ("tenantId");

CREATE INDEX "IDX_f5fb5ac30b3609c27af3517727" ON "base_sys_department" USING btree ("userId");