-- Space模块PostgreSQL数据库回滚迁移文件
-- 创建时间: 2025-01-27
-- 描述: 回滚space模块的空间管理相关表

-- 删除表（按依赖关系逆序删除）
DROP TABLE IF EXISTS space_info;
DROP TABLE IF EXISTS space_type;