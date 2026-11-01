-- OpenAPI模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建openapi模块的开放平台相关表

-- 1. 开放平台应用表
CREATE TABLE IF NOT EXISTS openapi_app (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    app_id VARCHAR(64) NOT NULL,
    app_name VARCHAR(255) NOT NULL,
    app_secret VARCHAR(255) NOT NULL,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    description VARCHAR(500),
    remark VARCHAR(255)
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
    app_id VARCHAR(64) NOT NULL,
    request_id VARCHAR(64) NOT NULL,
    timestamp BIGINT NOT NULL,
    nonce VARCHAR(64) NOT NULL,
    request_body TEXT,
    signature TEXT NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    status INTEGER NOT NULL DEFAULT 1,
    error_msg VARCHAR(500)
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

-- 字段注释
COMMENT ON COLUMN openapi_app."app_id" IS '应用ID';
COMMENT ON COLUMN openapi_app."app_name" IS '应用名称';
COMMENT ON COLUMN openapi_app."app_secret" IS '应用密钥';
COMMENT ON COLUMN openapi_app."public_key" IS 'RSA公钥';
COMMENT ON COLUMN openapi_app."private_key" IS 'RSA私钥';
COMMENT ON COLUMN openapi_app."status" IS '状态 0:禁用 1:启用';
COMMENT ON COLUMN openapi_app."description" IS '应用描述';
COMMENT ON COLUMN openapi_app."remark" IS '备注';
COMMENT ON COLUMN openapi_sign_log."app_id" IS '应用ID';
COMMENT ON COLUMN openapi_sign_log."request_id" IS '请求ID';
COMMENT ON COLUMN openapi_sign_log."timestamp" IS '时间戳';
COMMENT ON COLUMN openapi_sign_log."nonce" IS '随机数';
COMMENT ON COLUMN openapi_sign_log."request_body" IS '请求体';
COMMENT ON COLUMN openapi_sign_log."signature" IS '生成的签名';
COMMENT ON COLUMN openapi_sign_log."client_ip" IS '客户端IP';
COMMENT ON COLUMN openapi_sign_log."user_agent" IS '用户代理';
COMMENT ON COLUMN openapi_sign_log."status" IS '状态 0:失败 1:成功';
COMMENT ON COLUMN openapi_sign_log."error_msg" IS '错误信息';
