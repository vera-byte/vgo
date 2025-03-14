// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysMenu is the golang structure for table base_sys_menu.
type BaseSysMenu struct {
	Id         int    `json:"id"         orm:"id"         description:"ID"`                // ID
	CreateTime string `json:"createTime" orm:"createTime" description:"创建时间"`              // 创建时间
	UpdateTime string `json:"updateTime" orm:"updateTime" description:"更新时间"`              // 更新时间
	TenantId   int    `json:"tenantId"   orm:"tenantId"   description:"租户ID"`              // 租户ID
	ParentId   int    `json:"parentId"   orm:"parentId"   description:"父菜单ID"`             // 父菜单ID
	Name       string `json:"name"       orm:"name"       description:"菜单名称"`              // 菜单名称
	Router     string `json:"router"     orm:"router"     description:"菜单地址"`              // 菜单地址
	Perms      string `json:"perms"      orm:"perms"      description:"权限标识"`              // 权限标识
	Type       int    `json:"type"       orm:"type"       description:"类型 0-目录 1-菜单 2-按钮"` // 类型 0-目录 1-菜单 2-按钮
	Icon       string `json:"icon"       orm:"icon"       description:"图标"`                // 图标
	OrderNum   int    `json:"orderNum"   orm:"orderNum"   description:"排序"`                // 排序
	ViewPath   string `json:"viewPath"   orm:"viewPath"   description:"视图地址"`              // 视图地址
	KeepAlive  int    `json:"keepAlive"  orm:"keepAlive"  description:"路由缓存"`              // 路由缓存
	IsShow     int    `json:"isShow"     orm:"isShow"     description:"是否显示"`              // 是否显示
}
