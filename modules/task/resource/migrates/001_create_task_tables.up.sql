-- Task模块PostgreSQL数据库迁移文件
-- 创建时间: 2025-01-27
-- 描述: 创建task模块的任务管理相关表

-- 1. 任务信息表
CREATE TABLE IF NOT EXISTS task_info (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "jobId" VARCHAR(255),
    "repeatConf" TEXT,
    name VARCHAR(255),
    cron VARCHAR(255),
    "limit" INTEGER,
    every INTEGER,
    remark VARCHAR(255),
    status INTEGER,
    "startDate" TIMESTAMP,
    "endDate" TIMESTAMP,
    data VARCHAR(255),
    service VARCHAR(255),
    type INTEGER,
    "nextRunTime" TIMESTAMP,
    "taskType" INTEGER
);

-- 任务信息表索引
CREATE INDEX IF NOT EXISTS idx_task_info_create_time ON task_info("createTime");
CREATE INDEX IF NOT EXISTS idx_task_info_update_time ON task_info("updateTime");
CREATE INDEX IF NOT EXISTS idx_task_info_deleted_at ON task_info("deletedAt");
CREATE UNIQUE INDEX IF NOT EXISTS uk_task_info_job_id ON task_info("jobId");
CREATE INDEX IF NOT EXISTS idx_task_info_status ON task_info(status);
CREATE INDEX IF NOT EXISTS idx_task_info_type ON task_info(type);
CREATE INDEX IF NOT EXISTS idx_task_info_task_type ON task_info("taskType");
CREATE INDEX IF NOT EXISTS idx_task_info_next_run_time ON task_info("nextRunTime");
CREATE INDEX IF NOT EXISTS idx_task_info_start_date ON task_info("startDate");
CREATE INDEX IF NOT EXISTS idx_task_info_end_date ON task_info("endDate");
CREATE INDEX IF NOT EXISTS idx_task_info_service ON task_info(service);

-- 2. 任务日志表
CREATE TABLE IF NOT EXISTS task_log (
    id BIGSERIAL PRIMARY KEY,
    "createTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updateTime" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "taskId" BIGINT,
    status SMALLINT NOT NULL,
    detail TEXT
);

-- 任务日志表索引
CREATE INDEX IF NOT EXISTS idx_task_log_create_time ON task_log("createTime");
CREATE INDEX IF NOT EXISTS idx_task_log_update_time ON task_log("updateTime");
CREATE INDEX IF NOT EXISTS idx_task_log_deleted_at ON task_log("deletedAt");
CREATE INDEX IF NOT EXISTS idx_task_log_task_id ON task_log("taskId");
CREATE INDEX IF NOT EXISTS idx_task_log_status ON task_log(status);

-- 任务日志表外键约束
ALTER TABLE task_log ADD CONSTRAINT fk_task_log_task 
    FOREIGN KEY ("taskId") REFERENCES task_info(id) ON DELETE CASCADE;

-- 创建触发器函数用于自动更新updateTime字段
CREATE OR REPLACE FUNCTION update_task_updated_time_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updateTime" = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表创建更新时间触发器
CREATE TRIGGER update_task_info_updated_time BEFORE UPDATE ON task_info FOR EACH ROW EXECUTE FUNCTION update_task_updated_time_column();
CREATE TRIGGER update_task_log_updated_time BEFORE UPDATE ON task_log FOR EACH ROW EXECUTE FUNCTION update_task_updated_time_column();

-- 字段注释
COMMENT ON COLUMN task_info."jobId" IS '任务ID';
COMMENT ON COLUMN task_info."repeatConf" IS '重复配置';
COMMENT ON COLUMN task_info."name" IS '任务名称';
COMMENT ON COLUMN task_info."cron" IS 'cron表达式';
COMMENT ON COLUMN task_info."limit" IS '限制次数 不传为不限制';
COMMENT ON COLUMN task_info."every" IS '间隔时间 单位秒';
COMMENT ON COLUMN task_info."remark" IS '备注';
COMMENT ON COLUMN task_info."status" IS '状态 0:关闭 1:开启';
COMMENT ON COLUMN task_info."startDate" IS '开始时间';
COMMENT ON COLUMN task_info."endDate" IS '结束时间';
COMMENT ON COLUMN task_info."data" IS '数据';
COMMENT ON COLUMN task_info."service" IS '执行的服务';
COMMENT ON COLUMN task_info."type" IS '类型 0:系统 1:用户';
COMMENT ON COLUMN task_info."nextRunTime" IS '下次执行时间';
COMMENT ON COLUMN task_info."taskType" IS '任务类型 0:cron 1:时间间隔';
COMMENT ON COLUMN task_log."taskId" IS '任务ID';
COMMENT ON COLUMN task_log."status" IS '状态 0:失败 1:成功';
COMMENT ON COLUMN task_log."detail" IS '详情';
