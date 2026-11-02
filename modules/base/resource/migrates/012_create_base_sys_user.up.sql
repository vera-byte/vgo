CREATE SEQUENCE IF NOT EXISTS base_sys_user_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_sys_user_id_seq;

CREATE TABLE base_sys_user (
  "id" integer NOT NULL DEFAULT nextval('base_sys_user_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "departmentId" integer,
  "userId" integer,
  "name" character varying,
  "username" character varying NOT NULL,
  "password" character varying NOT NULL,
  "passwordV" integer DEFAULT 1,
  "nickName" character varying,
  "headImg" character varying,
  "phone" character varying,
  "email" character varying,
  "remark" character varying,
  "status" integer DEFAULT 1,
  "socketId" character varying,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_sys_user."departmentId" IS '部门ID';

COMMENT ON COLUMN base_sys_user."userId" IS '创建者ID';

COMMENT ON COLUMN base_sys_user."name" IS '姓名';

COMMENT ON COLUMN base_sys_user."username" IS '用户名';

COMMENT ON COLUMN base_sys_user."password" IS '密码';

COMMENT ON COLUMN base_sys_user."passwordV" IS '密码版本';

COMMENT ON COLUMN base_sys_user."nickName" IS '昵称';

COMMENT ON COLUMN base_sys_user."headImg" IS '头像';

COMMENT ON COLUMN base_sys_user."phone" IS '手机';

COMMENT ON COLUMN base_sys_user."email" IS '邮箱';

COMMENT ON COLUMN base_sys_user."remark" IS '备注';

COMMENT ON COLUMN base_sys_user."status" IS '状态 0:禁用 1:启用';

COMMENT ON COLUMN base_sys_user."socketId" IS 'socketId';

CREATE INDEX idx_base_sys_user_username ON base_sys_user USING btree ("username");

CREATE INDEX idx_base_sys_user_departmentId ON base_sys_user USING btree ("departmentId");
