-- 数据回收站：软删除/删除时数据快照回收到该表（参照 cool-admin-midway recycle 模块）

CREATE SEQUENCE IF NOT EXISTS recycle_data_id_seq;

CREATE TABLE recycle_data (
  "id" integer NOT NULL DEFAULT nextval('recycle_data_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "entityInfo" character varying NOT NULL,
  "userId" integer,
  "data" text NOT NULL,
  "url" character varying,
  "params" text,
  "count" integer NOT NULL DEFAULT 1,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN recycle_data."entityInfo" IS '来源信息，格式 {"group":"default","table":"表名"}';
COMMENT ON COLUMN recycle_data."userId" IS '操作人ID';
COMMENT ON COLUMN recycle_data."data" IS '被删除的数据（JSON 数组快照）';
COMMENT ON COLUMN recycle_data."url" IS '请求的接口';
COMMENT ON COLUMN recycle_data."params" IS '请求参数（JSON）';
COMMENT ON COLUMN recycle_data."count" IS '删除数据条数';

CREATE INDEX idx_recycle_data_user_id ON recycle_data USING btree ("userId");
CREATE INDEX idx_recycle_data_create_time ON recycle_data USING btree ("createTime");

CREATE TRIGGER update_recycle_data_updated_time BEFORE UPDATE ON recycle_data FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
