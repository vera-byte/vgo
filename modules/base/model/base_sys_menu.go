package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysMenu = "base_sys_menu"

// BaseSysMenu mapped from table <base_sys_menu>
type BaseSysMenu struct {
	*v.Model
	ParentID  uint    `json:"parentId"`  // 父菜单ID
	Name      string  `json:"name"`      // 菜单名称
	Router    *string `json:"router"`    // 菜单地址
	Perms     *string `json:"perms"`     // 权限标识
	Type      int32   `json:"type"`      // 类型 0：目录 1：菜单 2：按钮
	Icon      *string `json:"icon"`      // 图标
	OrderNum  int32   `json:"orderNum"`  // 排序
	ViewPath  *string `json:"viewPath"`  // 视图地址
	KeepAlive *int32  `json:"keepAlive"` // 路由缓存
	IsShow    *int32  `json:"isShow"`    // 是否显示
}

// TableName BaseSysMenu's table name
func (*BaseSysMenu) TableName() string {
	return TableNameBaseSysMenu
}

// NewBaseSysMenu create a new BaseSysMenu
func NewBaseSysMenu() *BaseSysMenu {
	return &BaseSysMenu{
		Model: v.NewModel(),
	}
}

// init 创建表
func init() {
	v.CreateTable(&BaseSysMenu{})
}
