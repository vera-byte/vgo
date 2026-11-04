-- 示例表 demo_sample（模块脚手架自带，可删除）

CREATE SEQUENCE IF NOT EXISTS demo_sample_id_seq;

CREATE TABLE demo_sample (
  "id" integer NOT NULL DEFAULT nextval('demo_sample_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "name" character varying NOT NULL DEFAULT '',
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN demo_sample."id" IS '主键';
COMMENT ON COLUMN demo_sample."createTime" IS '创建时间';
COMMENT ON COLUMN demo_sample."updateTime" IS '更新时间';
COMMENT ON COLUMN demo_sample."name" IS '名称';

CREATE TRIGGER demo_sample_update_updated_time BEFORE UPDATE ON demo_sample
  FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
