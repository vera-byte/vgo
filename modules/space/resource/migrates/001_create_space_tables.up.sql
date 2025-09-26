-- Space模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建space模块的空间管理相关表

-- 1. 空间类型表
CREATE TABLE IF NOT EXISTS space_type (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    name VARCHAR(255) NOT NULL COMMENT '类别名称',
    "parentId" INTEGER COMMENT '父分类ID'
);

-- 空间类型表索引
CREATE INDEX IF NOT EXISTS idx_space_type_create_time ON space_type("createTime");
CREATE INDEX IF NOT EXISTS idx_space_type_update_time ON space_type("updateTime");
CREATE INDEX IF NOT EXISTS idx_space_type_deleted_at ON space_type("deletedAt");
CREATE INDEX IF NOT EXISTS idx_space_type_parent_id ON space_type("parentId");
CREATE INDEX IF NOT EXISTS idx_space_type_name ON space_type(name);

-- 2. 空间信息表
CREATE TABLE IF NOT EXISTS space_info (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    url VARCHAR(255) NOT NULL COMMENT '地址',
    type VARCHAR(255) NOT NULL COMMENT '类型',
    "classifyId" BIGINT COMMENT '分类ID'
);

-- 空间信息表索引
CREATE INDEX IF NOT EXISTS idx_space_info_create_time ON space_info("createTime");
CREATE INDEX IF NOT EXISTS idx_space_info_update_time ON space_info("updateTime");
CREATE INDEX IF NOT EXISTS idx_space_info_deleted_at ON space_info("deletedAt");
CREATE INDEX IF NOT EXISTS idx_space_info_classify_id ON space_info("classifyId");
CREATE INDEX IF NOT EXISTS idx_space_info_type ON space_info(type);
CREATE INDEX IF NOT EXISTS idx_space_info_url ON space_info(url);

-- 空间信息表外键约束
ALTER TABLE space_info ADD CONSTRAINT fk_space_info_classify 
    FOREIGN KEY ("classifyId") REFERENCES space_type(id) ON DELETE SET NULL;

-- 创建触发器函数用于自动更新updateTime字段
CREATE OR REPLACE FUNCTION update_space_updated_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updateTime" = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表创建更新时间触发器
CREATE TRIGGER update_space_type_updated_time BEFORE UPDATE ON space_type FOR EACH ROW EXECUTE FUNCTION update_space_updated_time_column();
CREATE TRIGGER update_space_info_updated_time BEFORE UPDATE ON space_info FOR EACH ROW EXECUTE FUNCTION update_space_updated_time_column();