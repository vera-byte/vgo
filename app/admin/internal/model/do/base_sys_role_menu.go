// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysRoleMenu is the golang structure of table base_sys_role_menu for DAO operations like Where/Data.
type BaseSysRoleMenu struct {
	g.Meta    `orm:"table:base_sys_role_menu, do:true"`
	Id        interface{} // ID
	CreatedAt interface{} // 创建时间
	UpdatedAt interface{} // 更新时间
	TenantId  interface{} // 租户ID
	RoleId    interface{} // 角色ID
	MenuId    interface{} // 菜单ID
	DeletedAt interface{} // 软删除时间
}
