
CREATE TRIGGER update_task_info_updated_time BEFORE UPDATE ON task_info FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_task_log_updated_time BEFORE UPDATE ON task_log FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
