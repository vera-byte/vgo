// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysDepartment is the golang structure for table base_sys_department.
type BaseSysDepartment struct {
	Id         int    `json:"id"         orm:"id"         description:"ID"`     // ID
	CreateTime string `json:"createTime" orm:"createTime" description:"创建时间"`   // 创建时间
	UpdateTime string `json:"updateTime" orm:"updateTime" description:"更新时间"`   // 更新时间
	TenantId   int    `json:"tenantId"   orm:"tenantId"   description:"租户ID"`   // 租户ID
	Name       string `json:"name"       orm:"name"       description:"部门名称"`   // 部门名称
	ParentId   int    `json:"parentId"   orm:"parentId"   description:"上级部门ID"` // 上级部门ID
	OrderNum   int    `json:"orderNum"   orm:"orderNum"   description:"排序"`     // 排序
}
