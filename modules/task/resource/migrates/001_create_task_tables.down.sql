-- Task模块PostgreSQL数据库回滚迁移文件
-- 创建时间: 2025-01-27
-- 描述: 回滚task模块的任务管理相关表

-- 删除表（按依赖关系逆序删除）
DROP TABLE IF EXISTS task_log;
DROP TABLE IF EXISTS task_info;