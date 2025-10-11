package model

const TableNameBaseSysInit = "base_sys_init"

// BaseSysInit mapped from table <base_sys_init>
type BaseSysInit struct {
	Id    uint   `json:"id"`
	Table string `json:"table"`
	Group string `json:"group"`
}

// TableName BaseSysInit's table namer
func (*BaseSysInit) TableName() string {
	return TableNameBaseSysInit
}

// TableGroup BaseSysInit's table group
func (*BaseSysInit) GroupName() string {
	return "default"
}

// GetStruct BaseSysInit's struct
func (m *BaseSysInit) GetStruct() interface{} {
	return m
}
