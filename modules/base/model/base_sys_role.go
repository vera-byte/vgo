package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysRole = "base_sys_role"

// BaseSysRole mapped from table <base_sys_role>
type BaseSysRole struct {
	*v.Model
	UserID    string  `json:"userId"`    // 用户ID
	Name      string  `json:"name"`      // 名称
	Label     *string `json:"label"`     // 角色标签
	Remark    *string `json:"remark"`    // 备注
	Relevance *int32  `json:"relevance"` // 数据权限是否关联上下级
}

// TableName BaseSysRole's table name
func (*BaseSysRole) TableName() string {
	return TableNameBaseSysRole
}

// NewBaseSysRole create a new BaseSysRole
func NewBaseSysRole() *BaseSysRole {
	return &BaseSysRole{
		Model: v.NewModel(),
	}
}
