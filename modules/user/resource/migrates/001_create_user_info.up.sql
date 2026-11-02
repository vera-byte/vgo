-- 用户信息表（参照 cool-admin-midway user 模块）

CREATE SEQUENCE IF NOT EXISTS user_info_id_seq;

CREATE TABLE user_info (
  "id" integer NOT NULL DEFAULT nextval('user_info_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "unionid" character varying,
  "avatarUrl" character varying,
  "nickName" character varying,
  "phone" character varying,
  "gender" integer NOT NULL DEFAULT 0,
  "status" integer NOT NULL DEFAULT 1,
  "loginType" integer NOT NULL DEFAULT 0,
  "password" character varying,
  "description" text,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN user_info."unionid" IS '登录唯一ID';
COMMENT ON COLUMN user_info."avatarUrl" IS '头像';
COMMENT ON COLUMN user_info."nickName" IS '昵称';
COMMENT ON COLUMN user_info."phone" IS '手机号';
COMMENT ON COLUMN user_info."gender" IS '性别 0-未知 1-男 2-女';
COMMENT ON COLUMN user_info."status" IS '状态 0-禁用 1-正常 2-已注销';
COMMENT ON COLUMN user_info."loginType" IS '登录方式 0-小程序 1-公众号 2-H5';
COMMENT ON COLUMN user_info."password" IS '密码（bcrypt）';
COMMENT ON COLUMN user_info."description" IS '介绍';

CREATE UNIQUE INDEX uk_user_info_unionid ON user_info("unionid");
CREATE UNIQUE INDEX uk_user_info_phone ON user_info("phone");
CREATE INDEX idx_user_info_status ON user_info USING btree (status);
