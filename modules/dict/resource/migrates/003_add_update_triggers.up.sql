
CREATE TRIGGER update_dict_info_updated_time BEFORE UPDATE ON dict_info FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();

CREATE TRIGGER update_dict_type_updated_time BEFORE UPDATE ON dict_type FOR EACH ROW EXECUTE FUNCTION update_updated_time_column();
