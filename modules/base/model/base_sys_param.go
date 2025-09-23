package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysParam = "base_sys_param"

// BaseSysParam mapped from table <base_sys_param>
type BaseSysParam struct {
	*v.Model
	KeyName  string  `json:"keyName"`  // 键位
	Name     string  `json:"name"`     // 名称
	Data     string  `json:"data"`     // 数据
	DataType int32   `json:"dataType"` // 数据类型 0:字符串 1：数组 2：键值对
	Remark   *string `json:"remark"`   // 备注
}

// TableName BaseSysParam's table name
func (*BaseSysParam) TableName() string {
	return TableNameBaseSysParam
}

// NewBaseSysParam 创建一个新的BaseSysParam
func NewBaseSysParam() *BaseSysParam {
	return &BaseSysParam{
		Model: v.NewModel(),
	}
}

// init 创建表
func init() {
	v.CreateTable(&BaseSysParam{})
}
