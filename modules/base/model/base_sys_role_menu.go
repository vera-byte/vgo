package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysRoleMenu = "base_sys_role_menu"

// BaseSysRoleMenu mapped from table <base_sys_role_menu>
type BaseSysRoleMenu struct {
	*v.Model
	RoleID uint `json:"roleId"` // 角色ID
	MenuID uint `json:"menuId"` // 菜单ID
}

// TableName BaseSysRoleMenu's table name
func (*BaseSysRoleMenu) TableName() string {
	return TableNameBaseSysRoleMenu
}

// NewBaseSysRoleMenu create a new BaseSysRoleMenu
func NewBaseSysRoleMenu() *BaseSysRoleMenu {
	return &BaseSysRoleMenu{
		Model: &v.Model{},
	}
}
