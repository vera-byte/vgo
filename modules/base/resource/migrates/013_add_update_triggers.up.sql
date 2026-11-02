-- 维护 updateTime 触发器函数（新导出 schema 缺失，恢复旧设计）
CREATE OR REPLACE FUNCTION update_updated_time_column() RETURNS trigger AS $$
BEGIN
    NEW."updateTime" := to_char(now(), 'YYYY-MM-DD HH24:MI:SS');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_base_sys_conf_updated_time BEFORE UPDATE ON base_sys_conf FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_department_updated_time BEFORE UPDATE ON base_sys_department FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_log_updated_time BEFORE UPDATE ON base_sys_log FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_menu_updated_time BEFORE UPDATE ON base_sys_menu FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_param_updated_time BEFORE UPDATE ON base_sys_param FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_role_updated_time BEFORE UPDATE ON base_sys_role FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_role_department_updated_time BEFORE UPDATE ON base_sys_role_department FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_role_menu_updated_time BEFORE UPDATE ON base_sys_role_menu FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_user_updated_time BEFORE UPDATE ON base_sys_user FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_sys_user_role_updated_time BEFORE UPDATE ON base_sys_user_role FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_eps_admin_updated_time BEFORE UPDATE ON base_eps_admin FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_base_eps_app_updated_time BEFORE UPDATE ON base_eps_app FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
