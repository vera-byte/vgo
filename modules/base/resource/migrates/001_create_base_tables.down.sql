-- Base模块PostgreSQL数据库回滚迁移文件
-- 创建时间: 2025-01-27
-- 描述: 回滚base模块的所有数据表

-- 删除表（按依赖关系逆序删除）
DROP TABLE IF EXISTS base_sys_user_role;
DROP TABLE IF EXISTS base_sys_role_menu;
DROP TABLE IF EXISTS base_sys_role_department;
DROP TABLE IF EXISTS base_sys_user;
DROP TABLE IF EXISTS base_sys_role;
DROP TABLE IF EXISTS base_sys_menu;
DROP TABLE IF EXISTS base_sys_department;
DROP TABLE IF EXISTS base_sys_param;
DROP TABLE IF EXISTS base_sys_log;
DROP TABLE IF EXISTS base_sys_conf;