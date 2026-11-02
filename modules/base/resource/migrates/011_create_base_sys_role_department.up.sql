CREATE SEQUENCE IF NOT EXISTS base_sys_role_department_id_seq;

CREATE SEQUENCE IF NOT EXISTS base_sys_role_department_id_seq;

CREATE TABLE base_sys_role_department (
  "id" integer NOT NULL DEFAULT nextval('base_sys_role_department_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "roleId" integer NOT NULL,
  "departmentId" integer NOT NULL,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN base_sys_role_department."roleId" IS '角色ID';

COMMENT ON COLUMN base_sys_role_department."departmentId" IS '部门ID';

CREATE INDEX idx_base_sys_role_department_roleId ON base_sys_role_department USING btree ("roleId");

CREATE INDEX idx_base_sys_role_department_departmentId ON base_sys_role_department USING btree ("departmentId");
