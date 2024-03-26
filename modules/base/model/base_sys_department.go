package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysDepartment = "base_sys_department"

// BaseSysDepartment mapped from table <base_sys_department>
type BaseSysDepartment struct {
	*v.Model
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"` // 部门名称
	ParentID uint   `gorm:"column:parentId;type:bigint" json:"parentId"`        // 上级部门ID
	OrderNum int32  `gorm:"column:orderNum;type:int;not null" json:"orderNum"`  // 排序
}

// TableName BaseSysDepartment's table name
func (*BaseSysDepartment) TableName() string {
	return TableNameBaseSysDepartment
}

// NewBaseSysDepartment 创建一个BaseSysDepartment实例
func NewBaseSysDepartment() *BaseSysDepartment {
	return &BaseSysDepartment{
		Model: v.NewModel(),
	}
}

// init 创建表
func init() {
	v.CreateTable(&BaseSysDepartment{})
}
