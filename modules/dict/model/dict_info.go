package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameDictInfo = "dict_info"

// DictInfo mapped from table <dict_info>
type DictInfo struct {
	*v.Model
	TypeID   int32   `json:"typeId"`   // 类型ID
	Name     string  `json:"name"`     // 名称
	OrderNum int32   `json:"orderNum"` // 排序
	Remark   *string `json:"remark"`   // 备注
	ParentID *int32  `json:"parentId"` // 父ID
}

// TableName DictInfo's table name
func (*DictInfo) TableName() string {
	return TableNameDictInfo
}

// GroupName DictInfo's table group
func (*DictInfo) GroupName() string {
	return "default"
}

// NewDictInfo create a new DictInfo
func NewDictInfo() *DictInfo {
	return &DictInfo{
		Model: v.NewModel(),
	}
}

// init 创建表
func init() {
	v.CreateTable(&DictInfo{})
}
