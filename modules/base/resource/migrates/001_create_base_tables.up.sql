-- Base模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建base模块的所有数据表，包含用户、角色、菜单、部门等核心功能表

-- 1. 部门表
CREATE TABLE IF NOT EXISTS base_sys_department (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    name VARCHAR(255) NOT NULL,
    "parentId" BIGINT,
    "orderNum" INTEGER NOT NULL DEFAULT 0
);

-- 部门表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_department_create_time ON base_sys_department("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_department_update_time ON base_sys_department("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_department_deleted_at ON base_sys_department("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_department_parent_id ON base_sys_department("parentId");

-- 2. 用户表
CREATE TABLE IF NOT EXISTS base_sys_user (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "departmentId" BIGINT,
    name VARCHAR(255),
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    "passwordV" INTEGER NOT NULL DEFAULT 1,
    "nickName" VARCHAR(255),
    "headImg" VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(255),
    status INTEGER NOT NULL DEFAULT 1,
    remark VARCHAR(255),
    "socketId" VARCHAR(255)
);

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_user_create_time ON base_sys_user("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_update_time ON base_sys_user("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_deleted_at ON base_sys_user("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_department_id ON base_sys_user("departmentId");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_username ON base_sys_user(username);
CREATE INDEX IF NOT EXISTS idx_base_sys_user_phone ON base_sys_user(phone);

-- 用户表外键约束
ALTER TABLE base_sys_user ADD CONSTRAINT fk_base_sys_user_department 
    FOREIGN KEY ("departmentId") REFERENCES base_sys_department(id) ON DELETE SET NULL;

-- 3. 角色表
CREATE TABLE IF NOT EXISTS base_sys_role (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "userId" VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(50),
    remark VARCHAR(255),
    relevance INTEGER NOT NULL DEFAULT 1
);

-- 角色表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_role_create_time ON base_sys_role("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_update_time ON base_sys_role("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_deleted_at ON base_sys_role("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_name ON base_sys_role(name);
CREATE INDEX IF NOT EXISTS idx_base_sys_role_label ON base_sys_role(label);

-- 4. 菜单表
CREATE TABLE IF NOT EXISTS base_sys_menu (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "parentId" BIGINT,
    name VARCHAR(255) NOT NULL,
    router VARCHAR(255),
    perms VARCHAR(255),
    type INTEGER NOT NULL,
    icon VARCHAR(255),
    "orderNum" INTEGER NOT NULL DEFAULT 0,
    "viewPath" VARCHAR(255),
    "keepAlive" INTEGER NOT NULL DEFAULT 1,
    "isShow" INTEGER NOT NULL DEFAULT 1
);

-- 菜单表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_menu_create_time ON base_sys_menu("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_menu_update_time ON base_sys_menu("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_menu_deleted_at ON base_sys_menu("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_menu_parent_id ON base_sys_menu("parentId");
CREATE INDEX IF NOT EXISTS idx_base_sys_menu_type ON base_sys_menu(type);

-- 5. 用户角色关联表
CREATE TABLE IF NOT EXISTS base_sys_user_role (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "userId" BIGINT NOT NULL,
    "roleId" BIGINT NOT NULL
);

-- 用户角色关联表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_user_role_create_time ON base_sys_user_role("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_role_update_time ON base_sys_user_role("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_role_deleted_at ON base_sys_user_role("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_role_user_id ON base_sys_user_role("userId");
CREATE INDEX IF NOT EXISTS idx_base_sys_user_role_role_id ON base_sys_user_role("roleId");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_sys_user_role ON base_sys_user_role("userId", "roleId");

-- 用户角色关联表外键约束
ALTER TABLE base_sys_user_role ADD CONSTRAINT fk_base_sys_user_role_user 
    FOREIGN KEY ("userId") REFERENCES base_sys_user(id) ON DELETE CASCADE;
ALTER TABLE base_sys_user_role ADD CONSTRAINT fk_base_sys_user_role_role 
    FOREIGN KEY ("roleId") REFERENCES base_sys_role(id) ON DELETE CASCADE;

-- 6. 角色菜单关联表
CREATE TABLE IF NOT EXISTS base_sys_role_menu (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "roleId" BIGINT NOT NULL,
    "menuId" BIGINT NOT NULL
);

-- 角色菜单关联表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_role_menu_create_time ON base_sys_role_menu("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_menu_update_time ON base_sys_role_menu("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_menu_deleted_at ON base_sys_role_menu("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_menu_role_id ON base_sys_role_menu("roleId");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_menu_menu_id ON base_sys_role_menu("menuId");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_sys_role_menu ON base_sys_role_menu("roleId", "menuId");

-- 角色菜单关联表外键约束
ALTER TABLE base_sys_role_menu ADD CONSTRAINT fk_base_sys_role_menu_role 
    FOREIGN KEY ("roleId") REFERENCES base_sys_role(id) ON DELETE CASCADE;
ALTER TABLE base_sys_role_menu ADD CONSTRAINT fk_base_sys_role_menu_menu 
    FOREIGN KEY ("menuId") REFERENCES base_sys_menu(id) ON DELETE CASCADE;

-- 7. 角色部门关联表
CREATE TABLE IF NOT EXISTS base_sys_role_department (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "roleId" BIGINT NOT NULL,
    "departmentId" BIGINT NOT NULL
);

-- 角色部门关联表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_role_department_create_time ON base_sys_role_department("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_department_update_time ON base_sys_role_department("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_department_deleted_at ON base_sys_role_department("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_department_role_id ON base_sys_role_department("roleId");
CREATE INDEX IF NOT EXISTS idx_base_sys_role_department_department_id ON base_sys_role_department("departmentId");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_sys_role_department ON base_sys_role_department("roleId", "departmentId");

-- 角色部门关联表外键约束
ALTER TABLE base_sys_role_department ADD CONSTRAINT fk_base_sys_role_department_role 
    FOREIGN KEY ("roleId") REFERENCES base_sys_role(id) ON DELETE CASCADE;
ALTER TABLE base_sys_role_department ADD CONSTRAINT fk_base_sys_role_department_department 
    FOREIGN KEY ("departmentId") REFERENCES base_sys_department(id) ON DELETE CASCADE;

-- 8. 系统配置表
CREATE TABLE IF NOT EXISTS base_sys_conf (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "cKey" VARCHAR(255) NOT NULL,
    "cValue" TEXT
);

-- 系统配置表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_conf_create_time ON base_sys_conf("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_conf_update_time ON base_sys_conf("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_conf_deleted_at ON base_sys_conf("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_sys_conf_key ON base_sys_conf("cKey");

-- 9. 系统参数表
CREATE TABLE IF NOT EXISTS base_sys_param (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "keyName" VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    data TEXT,
    "dataType" INTEGER NOT NULL DEFAULT 0,
    remark VARCHAR(255)
);

-- 系统参数表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_param_create_time ON base_sys_param("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_param_update_time ON base_sys_param("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_param_deleted_at ON base_sys_param("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_sys_param_key_name ON base_sys_param("keyName");

-- 10. 系统日志表
CREATE TABLE IF NOT EXISTS base_sys_log (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "userId" BIGINT,
    action VARCHAR(255) NOT NULL,
    ip VARCHAR(50),
    "ipAddr" VARCHAR(255),
    params TEXT
);

-- 系统日志表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_log_create_time ON base_sys_log("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_log_update_time ON base_sys_log("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_log_deleted_at ON base_sys_log("deletedAt");
CREATE INDEX IF NOT EXISTS idx_base_sys_log_user_id ON base_sys_log("userId");
CREATE INDEX IF NOT EXISTS idx_base_sys_log_action ON base_sys_log(action);
CREATE INDEX IF NOT EXISTS idx_base_sys_log_ip ON base_sys_log(ip);

-- 11. 系统初始化表
CREATE TABLE IF NOT EXISTS base_sys_init (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "group" VARCHAR(255) NOT NULL,
    "table" VARCHAR(255) NOT NULL
);

-- 系统初始化表索引
CREATE INDEX IF NOT EXISTS idx_base_sys_init_create_time ON base_sys_init("createTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_init_update_time ON base_sys_init("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_sys_init_deleted_at ON base_sys_init("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_sys_init_group_table ON base_sys_init("group", "table");

-- 12. EPS管理员表
CREATE TABLE IF NOT EXISTS base_eps_admin (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    module VARCHAR(100) NOT NULL,
    method VARCHAR(16) NOT NULL,
    path VARCHAR(255) NOT NULL,
    prefix VARCHAR(255) NOT NULL,
    summary VARCHAR(255),
    tag VARCHAR(100),
    dts TEXT
);

-- EPS管理员表索引
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_create_time ON base_eps_admin("createTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_update_time ON base_eps_admin("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_deleted_at ON base_eps_admin("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_eps_admin_route ON base_eps_admin(module, prefix, path, method);

-- 13. EPS应用表
CREATE TABLE IF NOT EXISTS base_eps_app (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    module VARCHAR(100) NOT NULL,
    method VARCHAR(16) NOT NULL,
    path VARCHAR(255) NOT NULL,
    prefix VARCHAR(255) NOT NULL,
    summary VARCHAR(255),
    tag VARCHAR(100),
    dts TEXT
);

-- EPS应用表索引
CREATE INDEX IF NOT EXISTS idx_base_eps_app_create_time ON base_eps_app("createTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_app_update_time ON base_eps_app("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_app_deleted_at ON base_eps_app("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_eps_app_route ON base_eps_app(module, prefix, path, method);

-- 创建触发器函数用于自动更新updateTime字段
CREATE OR REPLACE FUNCTION update_updated_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updateTime" = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表创建更新时间触发器
CREATE TRIGGER update_base_sys_department_updated_time BEFORE UPDATE ON base_sys_department FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_user_updated_time BEFORE UPDATE ON base_sys_user FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_role_updated_time BEFORE UPDATE ON base_sys_role FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_menu_updated_time BEFORE UPDATE ON base_sys_menu FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_user_role_updated_time BEFORE UPDATE ON base_sys_user_role FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_role_menu_updated_time BEFORE UPDATE ON base_sys_role_menu FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_role_department_updated_time BEFORE UPDATE ON base_sys_role_department FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_conf_updated_time BEFORE UPDATE ON base_sys_conf FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_param_updated_time BEFORE UPDATE ON base_sys_param FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_log_updated_time BEFORE UPDATE ON base_sys_log FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_sys_init_updated_time BEFORE UPDATE ON base_sys_init FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_eps_admin_updated_time BEFORE UPDATE ON base_eps_admin FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_eps_app_updated_time BEFORE UPDATE ON base_eps_app FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

-- 字段注释
COMMENT ON COLUMN base_sys_department."name" IS '部门名称';
COMMENT ON COLUMN base_sys_department."parentId" IS '上级部门ID';
COMMENT ON COLUMN base_sys_department."orderNum" IS '排序';
COMMENT ON COLUMN base_sys_user."departmentId" IS '部门ID';
COMMENT ON COLUMN base_sys_user."name" IS '姓名';
COMMENT ON COLUMN base_sys_user."username" IS '用户名';
COMMENT ON COLUMN base_sys_user."password" IS '密码';
COMMENT ON COLUMN base_sys_user."passwordV" IS '密码版本';
COMMENT ON COLUMN base_sys_user."nickName" IS '昵称';
COMMENT ON COLUMN base_sys_user."headImg" IS '头像';
COMMENT ON COLUMN base_sys_user."phone" IS '手机';
COMMENT ON COLUMN base_sys_user."email" IS '邮箱';
COMMENT ON COLUMN base_sys_user."status" IS '状态 0:禁用 1:启用';
COMMENT ON COLUMN base_sys_user."remark" IS '备注';
COMMENT ON COLUMN base_sys_user."socketId" IS 'socketId';
COMMENT ON COLUMN base_sys_role."userId" IS '用户ID';
COMMENT ON COLUMN base_sys_role."name" IS '名称';
COMMENT ON COLUMN base_sys_role."label" IS '角色标签';
COMMENT ON COLUMN base_sys_role."remark" IS '备注';
COMMENT ON COLUMN base_sys_role."relevance" IS '数据权限是否关联上下级';
COMMENT ON COLUMN base_sys_menu."parentId" IS '父菜单ID';
COMMENT ON COLUMN base_sys_menu."name" IS '菜单名称';
COMMENT ON COLUMN base_sys_menu."router" IS '菜单地址';
COMMENT ON COLUMN base_sys_menu."perms" IS '权限标识';
COMMENT ON COLUMN base_sys_menu."type" IS '类型 0:目录 1:菜单 2:按钮';
COMMENT ON COLUMN base_sys_menu."icon" IS '图标';
COMMENT ON COLUMN base_sys_menu."orderNum" IS '排序';
COMMENT ON COLUMN base_sys_menu."viewPath" IS '视图地址';
COMMENT ON COLUMN base_sys_menu."keepAlive" IS '路由缓存';
COMMENT ON COLUMN base_sys_menu."isShow" IS '是否显示';
COMMENT ON COLUMN base_sys_user_role."userId" IS '用户ID';
COMMENT ON COLUMN base_sys_user_role."roleId" IS '角色ID';
COMMENT ON COLUMN base_sys_role_menu."roleId" IS '角色ID';
COMMENT ON COLUMN base_sys_role_menu."menuId" IS '菜单ID';
COMMENT ON COLUMN base_sys_role_department."roleId" IS '角色ID';
COMMENT ON COLUMN base_sys_role_department."departmentId" IS '部门ID';
COMMENT ON COLUMN base_sys_conf."cKey" IS '配置键';
COMMENT ON COLUMN base_sys_conf."cValue" IS '配置值';
COMMENT ON COLUMN base_sys_param."keyName" IS '参数键名';
COMMENT ON COLUMN base_sys_param."name" IS '参数名称';
COMMENT ON COLUMN base_sys_param."data" IS '参数数据';
COMMENT ON COLUMN base_sys_param."dataType" IS '数据类型 0:字符串 1:数字 2:布尔 3:JSON';
COMMENT ON COLUMN base_sys_param."remark" IS '备注';
COMMENT ON COLUMN base_sys_log."userId" IS '用户ID';
COMMENT ON COLUMN base_sys_log."action" IS '操作';
COMMENT ON COLUMN base_sys_log."ip" IS 'IP地址';
COMMENT ON COLUMN base_sys_log."ipAddr" IS 'IP归属地';
COMMENT ON COLUMN base_sys_log."params" IS '请求参数';
COMMENT ON COLUMN base_sys_init."group" IS '分组';
COMMENT ON COLUMN base_sys_init."table" IS '表名';
COMMENT ON COLUMN base_eps_admin.module IS '所属模块名称';
COMMENT ON COLUMN base_eps_admin.method IS '请求方法';
COMMENT ON COLUMN base_eps_admin.path IS '请求路径';
COMMENT ON COLUMN base_eps_admin.prefix IS '路由前缀';
COMMENT ON COLUMN base_eps_admin.summary IS '描述';
COMMENT ON COLUMN base_eps_admin.tag IS '标签';
COMMENT ON COLUMN base_eps_admin.dts IS '附加数据';
COMMENT ON COLUMN base_eps_app.module IS '所属模块名称';
COMMENT ON COLUMN base_eps_app.method IS '请求方法';
COMMENT ON COLUMN base_eps_app.path IS '请求路径';
COMMENT ON COLUMN base_eps_app.prefix IS '路由前缀';
COMMENT ON COLUMN base_eps_app.summary IS '描述';
COMMENT ON COLUMN base_eps_app.tag IS '标签';
COMMENT ON COLUMN base_eps_app.dts IS '附加数据';
