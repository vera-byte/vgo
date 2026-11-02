-- 插件信息表（参照 cool-admin-midway plugin 模块）

CREATE SEQUENCE IF NOT EXISTS plugin_info_id_seq;

CREATE TABLE plugin_info (
  "id" integer NOT NULL DEFAULT nextval('plugin_info_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "name" character varying NOT NULL,
  "description" character varying,
  "keyName" character varying NOT NULL,
  "hook" character varying,
  "readme" text,
  "version" character varying,
  "logo" text,
  "author" character varying,
  "status" integer NOT NULL DEFAULT 0,
  "content" text,
  "tsContent" text,
  "pluginJson" text,
  "config" text,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN plugin_info."name" IS '名称';
COMMENT ON COLUMN plugin_info."description" IS '简介';
COMMENT ON COLUMN plugin_info."keyName" IS 'Key名（唯一标识）';
COMMENT ON COLUMN plugin_info."hook" IS 'Hook 名（如 upload），同 Hook 同时只能启用一个';
COMMENT ON COLUMN plugin_info."readme" IS '使用说明';
COMMENT ON COLUMN plugin_info."version" IS '版本';
COMMENT ON COLUMN plugin_info."logo" IS 'Logo（base64）';
COMMENT ON COLUMN plugin_info."author" IS '作者';
COMMENT ON COLUMN plugin_info."status" IS '状态 0-禁用 1-启用';
COMMENT ON COLUMN plugin_info."content" IS '插件内容（JSON {type,data}）';
COMMENT ON COLUMN plugin_info."tsContent" IS 'TS 源码（JSON {type,data}，留档）';
COMMENT ON COLUMN plugin_info."pluginJson" IS '插件 plugin.json（JSON）';
COMMENT ON COLUMN plugin_info."config" IS '配置（JSON，支持 @env 多环境）';

CREATE UNIQUE INDEX uk_plugin_info_key_name ON plugin_info("keyName");
CREATE INDEX idx_plugin_info_status ON plugin_info USING btree (status);
CREATE INDEX idx_plugin_info_hook ON plugin_info USING btree (hook);

CREATE TRIGGER update_plugin_info_updated_time BEFORE UPDATE ON plugin_info FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
