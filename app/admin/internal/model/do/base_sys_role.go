// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysRole is the golang structure of table base_sys_role for DAO operations like Where/Data.
type BaseSysRole struct {
	g.Meta           `orm:"table:base_sys_role, do:true"`
	Id               interface{} // ID
	CreateTime       interface{} // 创建时间
	UpdateTime       interface{} // 更新时间
	TenantId         interface{} // 租户ID
	UserId           interface{} // 用户ID
	Name             interface{} // 名称
	Label            interface{} // 角色标签
	Remark           interface{} // 备注
	Relevance        interface{} // 数据权限是否关联上下级
	MenuIdList       interface{} // 菜单权限
	DepartmentIdList interface{} // 部门权限
}
