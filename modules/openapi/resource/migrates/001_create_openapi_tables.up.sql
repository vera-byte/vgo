-- OpenAPI模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建openapi模块的开放平台相关表

-- 1. 开放平台应用表
CREATE TABLE IF NOT EXISTS openapi_app (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    app_id VARCHAR(64) NOT NULL COMMENT '应用ID',
    app_name VARCHAR(255) NOT NULL COMMENT '应用名称',
    app_secret VARCHAR(255) NOT NULL COMMENT '应用密钥',
    public_key TEXT NOT NULL COMMENT 'RSA公钥',
    private_key TEXT NOT NULL COMMENT 'RSA私钥',
    status INTEGER NOT NULL DEFAULT 1 COMMENT '状态 0:禁用 1:启用',
    description VARCHAR(500) COMMENT '应用描述',
    remark VARCHAR(255) COMMENT '备注'
);

-- 开放平台应用表索引
CREATE INDEX IF NOT EXISTS idx_openapi_app_create_time ON openapi_app("createTime");
CREATE INDEX IF NOT EXISTS idx_openapi_app_update_time ON openapi_app("updateTime");
CREATE INDEX IF NOT EXISTS idx_openapi_app_deleted_at ON openapi_app("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_openapi_app_app_id ON openapi_app(app_id);
CREATE INDEX IF NOT EXISTS idx_openapi_app_status ON openapi_app(status);

-- 2. 开放平台签名日志表
CREATE TABLE IF NOT EXISTS openapi_sign_log (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    app_id VARCHAR(64) NOT NULL COMMENT '应用ID',
    request_id VARCHAR(64) NOT NULL COMMENT '请求ID',
    timestamp BIGINT NOT NULL COMMENT '时间戳',
    nonce VARCHAR(64) NOT NULL COMMENT '随机数',
    request_body TEXT COMMENT '请求体',
    signature TEXT NOT NULL COMMENT '生成的签名',
    client_ip VARCHAR(45) COMMENT '客户端IP',
    user_agent VARCHAR(500) COMMENT '用户代理',
    status INTEGER NOT NULL DEFAULT 1 COMMENT '状态 0:失败 1:成功',
    error_msg VARCHAR(500) COMMENT '错误信息'
);

-- 开放平台签名日志表索引
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_create_time ON openapi_sign_log("createTime");
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_update_time ON openapi_sign_log("updateTime");
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_deleted_at ON openapi_sign_log("deletedAt");
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_app_id ON openapi_sign_log(app_id);
CREATE UNIQUE INDEX IF NOT EXISTS uk_openapi_sign_log_request_id ON openapi_sign_log(request_id);
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_timestamp ON openapi_sign_log(timestamp);
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_client_ip ON openapi_sign_log(client_ip);
CREATE INDEX IF NOT EXISTS idx_openapi_sign_log_status ON openapi_sign_log(status);

-- 签名日志表外键约束
ALTER TABLE openapi_sign_log ADD CONSTRAINT fk_openapi_sign_log_app 
    FOREIGN KEY (app_id) REFERENCES openapi_app(app_id) ON DELETE CASCADE;

-- 创建触发器函数用于自动更新updateTime字段
CREATE OR REPLACE FUNCTION update_openapi_updated_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updateTime" = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表创建更新时间触发器
CREATE TRIGGER update_openapi_app_updated_time BEFORE UPDATE ON openapi_app FOR EACH ROW EXECUTE FUNCTION update_openapi_updated_time_column();
CREATE TRIGGER update_openapi_sign_log_updated_time BEFORE UPDATE ON openapi_sign_log FOR EACH ROW EXECUTE FUNCTION update_openapi_updated_time_column();