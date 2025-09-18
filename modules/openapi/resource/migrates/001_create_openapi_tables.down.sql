-- OpenAPI模块PostgreSQL数据库回滚迁移文件
-- 创建时间: 2025-01-27
-- 描述: 回滚openapi模块的开放平台相关表

-- 删除表（按依赖关系逆序删除）
DROP TABLE IF EXISTS openapi_sign_log;
DROP TABLE IF EXISTS openapi_app;