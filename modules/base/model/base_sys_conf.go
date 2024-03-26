package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysConf = "base_sys_conf"

// BaseSysConf mapped from table <base_sys_conf>
type BaseSysConf struct {
	*v.Model
	CKey   string `gorm:"column:cKey;type:varchar(255);not null;index" json:"cKey"` // 配置键
	CValue string `gorm:"column:cValue;type:varchar(255);not null" json:"cValue"`   // 配置值
}

// TableName BaseSysConf's table name
func (*BaseSysConf) TableName() string {
	return TableNameBaseSysConf
}

// init 创建表
func init() {
	v.CreateTable(&BaseSysConf{})
}

// NewBaseSysConf 创建实例
func NewBaseSysConf() *BaseSysConf {
	return &BaseSysConf{
		Model: v.NewModel(),
	}
}
