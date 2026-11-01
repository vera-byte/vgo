-- Dict模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建dict模块的数据字典相关表

-- 1. 字典类型表
CREATE TABLE IF NOT EXISTS dict_type (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL
);

-- 字典类型表索引
CREATE INDEX IF NOT EXISTS idx_dict_type_create_time ON dict_type("createTime");
CREATE INDEX IF NOT EXISTS idx_dict_type_update_time ON dict_type("updateTime");
CREATE INDEX IF NOT EXISTS idx_dict_type_deleted_at ON dict_type("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_dict_type_key ON dict_type(key);

-- 2. 字典信息表
CREATE TABLE IF NOT EXISTS dict_info (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "typeId" INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    "orderNum" INTEGER NOT NULL DEFAULT 0,
    remark VARCHAR(255),
    "parentId" INTEGER
);

-- 字典信息表索引
CREATE INDEX IF NOT EXISTS idx_dict_info_create_time ON dict_info("createTime");
CREATE INDEX IF NOT EXISTS idx_dict_info_update_time ON dict_info("updateTime");
CREATE INDEX IF NOT EXISTS idx_dict_info_deleted_at ON dict_info("deletedAt");
CREATE INDEX IF NOT EXISTS idx_dict_info_type_id ON dict_info("typeId");
CREATE INDEX IF NOT EXISTS idx_dict_info_parent_id ON dict_info("parentId");
CREATE INDEX IF NOT EXISTS idx_dict_info_order_num ON dict_info("orderNum");

-- 字典信息表外键约束
ALTER TABLE dict_info ADD CONSTRAINT fk_dict_info_type 
    FOREIGN KEY ("typeId") REFERENCES dict_type(id) ON DELETE CASCADE;

-- 创建触发器函数用于自动更新updateTime字段
CREATE OR REPLACE FUNCTION update_dict_updated_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updateTime" = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表创建更新时间触发器
CREATE TRIGGER update_dict_type_updated_time BEFORE UPDATE ON dict_type FOR EACH ROW EXECUTE FUNCTION update_dict_updated_time_column();
CREATE TRIGGER update_dict_info_updated_time BEFORE UPDATE ON dict_info FOR EACH ROW EXECUTE FUNCTION update_dict_updated_time_column();

-- 字段注释
COMMENT ON COLUMN dict_type."name" IS '名称';
COMMENT ON COLUMN dict_type."key" IS '标识';
COMMENT ON COLUMN dict_info."typeId" IS '类型ID';
COMMENT ON COLUMN dict_info."name" IS '名称';
COMMENT ON COLUMN dict_info."orderNum" IS '排序';
COMMENT ON COLUMN dict_info."remark" IS '备注';
COMMENT ON COLUMN dict_info."parentId" IS '父ID';
