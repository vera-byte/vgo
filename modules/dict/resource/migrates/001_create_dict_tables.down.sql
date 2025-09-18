-- Dict模块PostgreSQL数据库回滚迁移文件
-- 创建时间: 2025-01-27
-- 描述: 回滚dict模块的数据字典相关表

-- 删除表（按依赖关系逆序删除）
DROP TABLE IF EXISTS dict_info;
DROP TABLE IF EXISTS dict_type;