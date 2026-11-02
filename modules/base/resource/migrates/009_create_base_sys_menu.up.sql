CREATE SEQUENCE IF NOT EXISTS base_sys_menu_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_sys_menu_id_seq;

CREATE TABLE base_sys_menu (
  "id" integer NOT NULL DEFAULT nextval('base_sys_menu_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "parentId" integer,
  "name" character varying NOT NULL,
  "router" character varying,
  "perms" character varying,
  "type" integer,
  "icon" character varying,
  "orderNum" integer DEFAULT 0,
  "viewPath" character varying,
  "keepAlive" boolean,
  "isShow" boolean,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_sys_menu."parentId" IS '上级菜单ID';

COMMENT ON COLUMN base_sys_menu."name" IS '菜单名称';

COMMENT ON COLUMN base_sys_menu."router" IS '路由';

COMMENT ON COLUMN base_sys_menu."perms" IS '权限标识';

COMMENT ON COLUMN base_sys_menu."type" IS '类型';

COMMENT ON COLUMN base_sys_menu."icon" IS '图标';

COMMENT ON COLUMN base_sys_menu."orderNum" IS '排序';

COMMENT ON COLUMN base_sys_menu."viewPath" IS '视图路径';

COMMENT ON COLUMN base_sys_menu."keepAlive" IS '缓存';

COMMENT ON COLUMN base_sys_menu."isShow" IS '显示';

CREATE INDEX idx_base_sys_menu_parentId ON base_sys_menu USING btree ("parentId");

CREATE INDEX idx_base_sys_menu_perms ON base_sys_menu USING btree ("perms");
