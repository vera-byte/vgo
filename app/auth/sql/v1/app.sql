CREATE TABLE apps (
    id              SERIAL PRIMARY KEY,
    app_id          VARCHAR(64) NOT NULL UNIQUE,         -- 应用唯一标识
    name            VARCHAR(128) NOT NULL,               -- 应用名称
    public_key      TEXT NOT NULL,                       -- RSA 公钥 PEM
    platform        VARCHAR(32) NOT NULL,                -- 平台标识：android / ios / web / jsapi
    android_md5     VARCHAR(64),                         -- Android 专属：keystore 的 md5
    disabled        BOOLEAN DEFAULT FALSE,               -- 禁用标志
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

TRUNCATE TABLE apps RESTART IDENTITY CASCADE;

-- 添加应用（假设一个 Android 应用）
INSERT INTO apps (app_id, name, public_key, platform, android_md5, disabled)
VALUES (
    'com.demo.app',
    '演示应用',
    '-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A...\n-----END PUBLIC KEY-----',
    'android',
    '9f5c6fe7d09c3f9ed443c1a32c77ee3c', -- Android keystore 的 MD5
    false
);