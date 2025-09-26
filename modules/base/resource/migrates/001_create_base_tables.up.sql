-- Base模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建base模块的所有数据表，包含用户、角色、菜单、部门等核心功能表

-- 1. 部门表
CREATE TABLE IF NOT EXISTS base_sys_department (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    name VARCHAR(255) NOT NULL COMMENT '部门名称',
    "parentId" BIGINT COMMENT '上级部门ID',
    "orderNum" INTEGER NOT NULL DEFAULT 0 COMMENT '排序'
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
    "departmentId" BIGINT COMMENT '部门ID',
    name VARCHAR(255) COMMENT '姓名',
    username VARCHAR(100) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    "passwordV" INTEGER NOT NULL DEFAULT 1 COMMENT '密码版本',
    "nickName" VARCHAR(255) COMMENT '昵称',
    "headImg" VARCHAR(255) COMMENT '头像',
    phone VARCHAR(20) COMMENT '手机',
    email VARCHAR(255) COMMENT '邮箱',
    status INTEGER NOT NULL DEFAULT 1 COMMENT '状态 0:禁用 1:启用',
    remark VARCHAR(255) COMMENT '备注',
    "socketId" VARCHAR(255) COMMENT 'socketId'
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
    "userId" VARCHAR(255) NOT NULL COMMENT '用户ID',
    name VARCHAR(255) NOT NULL COMMENT '名称',
    label VARCHAR(50) COMMENT '角色标签',
    remark VARCHAR(255) COMMENT '备注',
    relevance INTEGER NOT NULL DEFAULT 1 COMMENT '数据权限是否关联上下级'
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
    "parentId" BIGINT COMMENT '父菜单ID',
    name VARCHAR(255) NOT NULL COMMENT '菜单名称',
    router VARCHAR(255) COMMENT '菜单地址',
    perms VARCHAR(255) COMMENT '权限标识',
    type INTEGER NOT NULL COMMENT '类型 0:目录 1:菜单 2:按钮',
    icon VARCHAR(255) COMMENT '图标',
    "orderNum" INTEGER NOT NULL DEFAULT 0 COMMENT '排序',
    "viewPath" VARCHAR(255) COMMENT '视图地址',
    "keepAlive" INTEGER NOT NULL DEFAULT 1 COMMENT '路由缓存',
    "isShow" INTEGER NOT NULL DEFAULT 1 COMMENT '是否显示'
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
    "userId" BIGINT NOT NULL COMMENT '用户ID',
    "roleId" BIGINT NOT NULL COMMENT '角色ID'
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
    "roleId" BIGINT NOT NULL COMMENT '角色ID',
    "menuId" BIGINT NOT NULL COMMENT '菜单ID'
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
    "roleId" BIGINT NOT NULL COMMENT '角色ID',
    "departmentId" BIGINT NOT NULL COMMENT '部门ID'
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
    "cKey" VARCHAR(255) NOT NULL COMMENT '配置键',
    "cValue" TEXT COMMENT '配置值'
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
    "keyName" VARCHAR(255) NOT NULL COMMENT '参数键名',
    name VARCHAR(255) NOT NULL COMMENT '参数名称',
    data TEXT COMMENT '参数数据',
    "dataType" INTEGER NOT NULL DEFAULT 0 COMMENT '数据类型 0:字符串 1:数字 2:布尔 3:JSON',
    remark VARCHAR(255) COMMENT '备注'
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
    "userId" BIGINT COMMENT '用户ID',
    action VARCHAR(255) NOT NULL COMMENT '操作',
    ip VARCHAR(50) COMMENT 'IP地址',
    "ipAddr" VARCHAR(255) COMMENT 'IP归属地',
    params TEXT COMMENT '请求参数'
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
    "group" VARCHAR(255) NOT NULL COMMENT '分组',
    "table" VARCHAR(255) NOT NULL COMMENT '表名'
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
    username VARCHAR(100) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    "nickName" VARCHAR(255) COMMENT '昵称',
    "headImg" VARCHAR(255) COMMENT '头像',
    phone VARCHAR(20) COMMENT '手机',
    email VARCHAR(255) COMMENT '邮箱',
    status INTEGER NOT NULL DEFAULT 1 COMMENT '状态 0:禁用 1:启用',
    remark VARCHAR(255) COMMENT '备注'
);

-- EPS管理员表索引
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_create_time ON base_eps_admin("createTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_update_time ON base_eps_admin("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_deleted_at ON base_eps_admin("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_eps_admin_username ON base_eps_admin(username);

-- 13. EPS应用表
CREATE TABLE IF NOT EXISTS base_eps_app (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    name VARCHAR(255) NOT NULL COMMENT '应用名称',
    "nameSpace" VARCHAR(255) NOT NULL COMMENT '命名空间',
    description TEXT COMMENT '应用描述',
    status INTEGER NOT NULL DEFAULT 1 COMMENT '状态 0:禁用 1:启用'
);

-- EPS应用表索引
CREATE INDEX IF NOT EXISTS idx_base_eps_app_create_time ON base_eps_app("createTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_app_update_time ON base_eps_app("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_app_deleted_at ON base_eps_app("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_eps_app_namespace ON base_eps_app("nameSpace");

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