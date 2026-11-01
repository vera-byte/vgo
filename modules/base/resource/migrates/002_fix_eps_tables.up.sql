-- 修复 base_eps_admin / base_eps_app 建表错误
-- 背景：001 迁移误将 EPS 接口表建成了用户表/应用表结构（复制粘贴错误），
--       导致接口元数据无法写入（input data match no fields in table）。
--       此迁移重建两表为正确的路由条目结构；旧表中不可能存在有效数据（插入从未成功）。

DROP TABLE IF EXISTS base_eps_admin;
DROP TABLE IF EXISTS base_eps_app;

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
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_create_time ON base_eps_admin("createTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_update_time ON base_eps_admin("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_admin_deleted_at ON base_eps_admin("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_eps_admin_route ON base_eps_admin(module, prefix, path, method);

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
CREATE INDEX IF NOT EXISTS idx_base_eps_app_create_time ON base_eps_app("createTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_app_update_time ON base_eps_app("updateTime");
CREATE INDEX IF NOT EXISTS idx_base_eps_app_deleted_at ON base_eps_app("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_base_eps_app_route ON base_eps_app(module, prefix, path, method);

CREATE TRIGGER update_base_eps_admin_updated_time BEFORE UPDATE ON base_eps_admin FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_base_eps_app_updated_time BEFORE UPDATE ON base_eps_app FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

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
