-- 微信用户表

CREATE SEQUENCE IF NOT EXISTS user_wx_id_seq;

CREATE TABLE user_wx (
  "id" integer NOT NULL DEFAULT nextval('user_wx_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "unionid" character varying,
  "openid" character varying NOT NULL,
  "avatarUrl" character varying,
  "nickName" character varying,
  "gender" integer NOT NULL DEFAULT 0,
  "language" character varying,
  "city" character varying,
  "province" character varying,
  "country" character varying,
  "type" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN user_wx."unionid" IS '微信unionid';
COMMENT ON COLUMN user_wx."openid" IS '微信openid';
COMMENT ON COLUMN user_wx."avatarUrl" IS '头像';
COMMENT ON COLUMN user_wx."nickName" IS '昵称';
COMMENT ON COLUMN user_wx."gender" IS '性别 0-未知 1-男 2-女';
COMMENT ON COLUMN user_wx."language" IS '语言';
COMMENT ON COLUMN user_wx."city" IS '城市';
COMMENT ON COLUMN user_wx."province" IS '省份';
COMMENT ON COLUMN user_wx."country" IS '国家';
COMMENT ON COLUMN user_wx."type" IS '类型 0-小程序 1-公众号 2-H5 3-APP';

CREATE INDEX idx_user_wx_unionid ON user_wx USING btree ("unionid");
CREATE UNIQUE INDEX uk_user_wx_openid ON user_wx("openid");
