
CREATE TRIGGER update_space_info_updated_time BEFORE UPDATE ON space_info FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_space_type_updated_time BEFORE UPDATE ON space_type FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
