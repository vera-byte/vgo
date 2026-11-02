CREATE SEQUENCE IF NOT EXISTS base_sys_role_menu_id_seq;

CREATE TABLE "base_sys_role_menu" (
  "id" integer NOT NULL DEFAULT nextval('base_sys_role_menu_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "roleId" integer NOT NULL,
  "menuId" integer NOT NULL,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN "base_sys_role_menu"."id" IS 'ID';

COMMENT ON COLUMN "base_sys_role_menu"."createTime" IS '创建时间';

COMMENT ON COLUMN "base_sys_role_menu"."updateTime" IS '更新时间';

COMMENT ON COLUMN "base_sys_role_menu"."tenantId" IS '租户ID';

COMMENT ON COLUMN "base_sys_role_menu"."roleId" IS '角色ID';

COMMENT ON COLUMN "base_sys_role_menu"."menuId" IS '菜单ID';

CREATE INDEX "IDX_3641f81d4201c524a57ce2aa54" ON "base_sys_role_menu" USING btree ("createTime");

CREATE INDEX "IDX_f860298298b26e7a697be36e5b" ON "base_sys_role_menu" USING btree ("updateTime");

CREATE INDEX "IDX_fd2d8bbe13949cfa56b1ed0a5d" ON "base_sys_role_menu" USING btree ("tenantId");