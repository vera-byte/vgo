// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package entity

import (
	"time"
)

// BaseSysRoleDepartment is the golang structure for table base_sys_role_department.
type BaseSysRoleDepartment struct {
	Id           int       `json:"id"           orm:"id"            description:"ID"`    // ID
	CreatedAt    time.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`  // 创建时间
	UpdatedAt    time.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`  // 更新时间
	TenantId     int       `json:"tenantId"     orm:"tenant_id"     description:"租户ID"`  // 租户ID
	RoleId       int       `json:"roleId"       orm:"role_id"       description:"角色ID"`  // 角色ID
	DepartmentId int       `json:"departmentId" orm:"department_id" description:"部门ID"`  // 部门ID
	DeletedAt    time.Time `json:"deletedAt"    orm:"deleted_at"    description:"软删除时间"` // 软删除时间
}
