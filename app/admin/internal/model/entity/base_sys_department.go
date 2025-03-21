// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package entity

import (
	"time"
)

// BaseSysDepartment is the golang structure for table base_sys_department.
type BaseSysDepartment struct {
	Id        int       `json:"id"        orm:"id"         description:"ID"`     // ID
	CreatedAt time.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
	UpdatedAt time.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`   // 更新时间
	TenantId  int       `json:"tenantId"  orm:"tenant_id"  description:"租户ID"`   // 租户ID
	Name      string    `json:"name"      orm:"name"       description:"部门名称"`   // 部门名称
	ParentId  int       `json:"parentId"  orm:"parent_id"  description:"上级部门ID"` // 上级部门ID
	OrderNum  int       `json:"orderNum"  orm:"order_num"  description:"排序"`     // 排序
	DeletedAt time.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"`  // 软删除时间
}
