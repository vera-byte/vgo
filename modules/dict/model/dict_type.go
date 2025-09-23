package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameDictType = "dict_type"

// DictType mapped from table <dict_type>
type DictType struct {
	*v.Model
	Name string `json:"name"` // 名称
	Key  string `json:"key"`  // 标识
}

// TableName DictType's table name
func (*DictType) TableName() string {
	return TableNameDictType
}

// GroupName DictType's table group
func (*DictType) GroupName() string {
	return "default"
}

// NewDictType create a new DictType
func NewDictType() *DictType {
	return &DictType{
		Model: v.NewModel(),
	}
}

// init 创建表
func init() {
	v.CreateTable(&DictType{})
}
