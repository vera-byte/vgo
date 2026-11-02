-- 维护 updateTime 触发器

CREATE TRIGGER update_user_info_updated_time BEFORE UPDATE ON user_info FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_user_address_updated_time BEFORE UPDATE ON user_address FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
CREATE TRIGGER update_user_wx_updated_time BEFORE UPDATE ON user_wx FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
