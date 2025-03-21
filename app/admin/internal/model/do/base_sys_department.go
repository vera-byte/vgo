// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysDepartment is the golang structure of table base_sys_department for DAO operations like Where/Data.
type BaseSysDepartment struct {
	g.Meta    `orm:"table:base_sys_department, do:true"`
	Id        interface{} // ID
	CreatedAt interface{} // 创建时间
	UpdatedAt interface{} // 更新时间
	TenantId  interface{} // 租户ID
	Name      interface{} // 部门名称
	ParentId  interface{} // 上级部门ID
	OrderNum  interface{} // 排序
	DeletedAt interface{} // 软删除时间
}
