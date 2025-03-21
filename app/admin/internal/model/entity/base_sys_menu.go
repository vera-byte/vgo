// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package entity

import (
	"time"
)

// BaseSysMenu is the golang structure for table base_sys_menu.
type BaseSysMenu struct {
	Id        int       `json:"id"        orm:"id"         description:"ID"`                // ID
	CreatedAt time.Time `json:"createdAt" orm:"created_at" description:"创建时间"`              // 创建时间
	UpdatedAt time.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`              // 更新时间
	TenantId  int       `json:"tenantId"  orm:"tenant_id"  description:"租户ID"`              // 租户ID
	ParentId  int       `json:"parentId"  orm:"parent_id"  description:"父菜单ID"`             // 父菜单ID
	Name      string    `json:"name"      orm:"name"       description:"菜单名称"`              // 菜单名称
	Router    string    `json:"router"    orm:"router"     description:"菜单地址"`              // 菜单地址
	Perms     string    `json:"perms"     orm:"perms"      description:"权限标识"`              // 权限标识
	Type      int       `json:"type"      orm:"type"       description:"类型 0-目录 1-菜单 2-按钮"` // 类型 0-目录 1-菜单 2-按钮
	Icon      string    `json:"icon"      orm:"icon"       description:"图标"`                // 图标
	OrderNum  int       `json:"orderNum"  orm:"order_num"  description:"排序"`                // 排序
	ViewPath  string    `json:"viewPath"  orm:"view_path"  description:"视图地址"`              // 视图地址
	KeepAlive bool      `json:"keepAlive" orm:"keep_alive" description:"路由缓存"`              // 路由缓存
	IsShow    bool      `json:"isShow"    orm:"is_show"    description:"是否显示"`              // 是否显示
	DeletedAt time.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"`             // 软删除时间
}
