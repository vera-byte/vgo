package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysUserRole = "base_sys_user_role"

// BaseSysUserRole mapped from table <base_sys_user_role>
type BaseSysUserRole struct {
	*v.Model
	UserID uint `json:"userId"` // 用户ID
	RoleID uint `json:"roleId"` // 角色ID
}

// TableName BaseSysUserRole's table name
func (*BaseSysUserRole) TableName() string {
	return TableNameBaseSysUserRole
}

// NewBaseSysUserRole create a new BaseSysUserRole
func NewBaseSysUserRole() *BaseSysUserRole {
	return &BaseSysUserRole{
		Model: v.NewModel(),
	}
}
