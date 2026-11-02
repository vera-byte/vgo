CREATE SEQUENCE IF NOT EXISTS base_sys_user_role_id_seq;

    CREATE TABLE "base_sys_user_role" (
  "id" integer NOT NULL DEFAULT nextval('base_sys_user_role_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "userId" integer NOT NULL,
  "roleId" integer NOT NULL,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "base_sys_user_role"."id" IS 'ID';

COMMENT ON COLUMN "base_sys_user_role"."createTime" IS '创建时间';

COMMENT ON COLUMN "base_sys_user_role"."updateTime" IS '更新时间';

COMMENT ON COLUMN "base_sys_user_role"."tenantId" IS '租户ID';

COMMENT ON COLUMN "base_sys_user_role"."userId" IS '用户ID';

COMMENT ON COLUMN "base_sys_user_role"."roleId" IS '角色ID';

CREATE INDEX "IDX_2f1dc0b6aad5604a2ddf37fba6" ON "base_sys_user_role" USING btree ("tenantId");

CREATE INDEX "IDX_3e36c0d2b1a4c659c6b4fc64b3" ON "base_sys_user_role" USING btree ("updateTime");

CREATE INDEX "IDX_fa9555e03e42fce748c9046b1c" ON "base_sys_user_role" USING btree ("createTime");