ALTER TABLE "plugin_info" ADD COLUMN "devPath" character varying;
COMMENT ON COLUMN "plugin_info"."devPath" IS '开发模式源码路径（go run 目标，仅 RunMode=dev 生效）';
