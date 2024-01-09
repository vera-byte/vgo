package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameSpaceType = "space_type"

// SpaceType mapped from table <space_type>
type SpaceType struct {
	*v.Model
	Name     string `gorm:"column:name;type:varchar(255);not null;comment:类别名称 " json:"name"` // 类别名称
	ParentID *int32 `gorm:"column:parentId;comment:父分类ID" json:"parentId"`                    // 父分类ID
}

// TableName SpaceType's table name
func (*SpaceType) TableName() string {
	return TableNameSpaceType
}

// GroupName SpaceType's table group
func (*SpaceType) GroupName() string {
	return "default"
}

// NewSpaceType create a new SpaceType
func NewSpaceType() *SpaceType {
	return &SpaceType{
		Model: v.NewModel(),
	}
}

// init 创建表
func init() {
	v.CreateTable(&SpaceType{})
}
