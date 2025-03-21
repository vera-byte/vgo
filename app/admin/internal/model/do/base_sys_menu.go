// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysMenu is the golang structure of table base_sys_menu for DAO operations like Where/Data.
type BaseSysMenu struct {
	g.Meta    `orm:"table:base_sys_menu, do:true"`
	Id        interface{} // ID
	CreatedAt interface{} // 创建时间
	UpdatedAt interface{} // 更新时间
	TenantId  interface{} // 租户ID
	ParentId  interface{} // 父菜单ID
	Name      interface{} // 菜单名称
	Router    interface{} // 菜单地址
	Perms     interface{} // 权限标识
	Type      interface{} // 类型 0-目录 1-菜单 2-按钮
	Icon      interface{} // 图标
	OrderNum  interface{} // 排序
	ViewPath  interface{} // 视图地址
	KeepAlive interface{} // 路由缓存
	IsShow    interface{} // 是否显示
	DeletedAt interface{} // 软删除时间
}
