-- 用户收货地址表

CREATE SEQUENCE IF NOT EXISTS user_address_id_seq;

CREATE TABLE user_address (
  "id" integer NOT NULL DEFAULT nextval('user_address_id_seq'::regclass),
  "createTime" character varying NOT NULL DEFAULT (now()::text),
  "updateTime" character varying NOT NULL DEFAULT (now()::text),
  "tenantId" integer,
  "userId" integer NOT NULL,
  "contact" character varying NOT NULL,
  "phone" character varying NOT NULL,
  "province" character varying NOT NULL,
  "city" character varying NOT NULL,
  "district" character varying NOT NULL,
  "address" character varying NOT NULL,
  "isDefault" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);

COMMENT ON COLUMN user_address."userId" IS '用户ID';
COMMENT ON COLUMN user_address."contact" IS '联系人';
COMMENT ON COLUMN user_address."phone" IS '手机号';
COMMENT ON COLUMN user_address."province" IS '省';
COMMENT ON COLUMN user_address."city" IS '市';
COMMENT ON COLUMN user_address."district" IS '区';
COMMENT ON COLUMN user_address."address" IS '地址';
COMMENT ON COLUMN user_address."isDefault" IS '是否默认';

CREATE INDEX idx_user_address_user_id ON user_address USING btree ("userId");
CREATE INDEX idx_user_address_phone ON user_address USING btree ("phone");
