// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysRole is the golang structure for table base_sys_role.
type BaseSysRole struct {
	Id               int    `json:"id"               orm:"id"               description:"ID"`          // ID
	CreateTime       string `json:"createTime"       orm:"createTime"       description:"创建时间"`        // 创建时间
	UpdateTime       string `json:"updateTime"       orm:"updateTime"       description:"更新时间"`        // 更新时间
	TenantId         int    `json:"tenantId"         orm:"tenantId"         description:"租户ID"`        // 租户ID
	UserId           string `json:"userId"           orm:"userId"           description:"用户ID"`        // 用户ID
	Name             string `json:"name"             orm:"name"             description:"名称"`          // 名称
	Label            string `json:"label"            orm:"label"            description:"角色标签"`        // 角色标签
	Remark           string `json:"remark"           orm:"remark"           description:"备注"`          // 备注
	Relevance        int    `json:"relevance"        orm:"relevance"        description:"数据权限是否关联上下级"` // 数据权限是否关联上下级
	MenuIdList       string `json:"menuIdList"       orm:"menuIdList"       description:"菜单权限"`        // 菜单权限
	DepartmentIdList string `json:"departmentIdList" orm:"departmentIdList" description:"部门权限"`        // 部门权限
}
