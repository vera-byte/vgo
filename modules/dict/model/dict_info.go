package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameDictInfo = "dict_info"

// DictInfo mapped from table <dict_info>
type DictInfo struct {
	*v.Model
	TypeID   int32   `gorm:"column:typeId;type:int;not null" json:"typeId"`      // 类型ID
	Name     string  `gorm:"column:name;type:varchar(255);not null" json:"name"` // 名称
	OrderNum int32   `gorm:"column:orderNum;type:int;not null" json:"orderNum"`  // 排序
	Remark   *string `gorm:"column:remark;type:varchar(255)" json:"remark"`      // 备注
	ParentID *int32  `gorm:"column:parentId;type:int" json:"parentId"`           // 父ID
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
