// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package entity

import (
	"time"
)

// BaseSysLog is the golang structure for table base_sys_log.
type BaseSysLog struct {
	Id        int       `json:"id"        orm:"id"         description:"ID"`    // ID
	CreatedAt time.Time `json:"createdAt" orm:"created_at" description:"创建时间"`  // 创建时间
	UpdatedAt time.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`  // 更新时间
	TenantId  int       `json:"tenantId"  orm:"tenant_id"  description:"租户ID"`  // 租户ID
	UserId    int       `json:"userId"    orm:"user_id"    description:"用户ID"`  // 用户ID
	Action    string    `json:"action"    orm:"action"     description:"行为"`    // 行为
	Ip        string    `json:"ip"        orm:"ip"         description:"ip"`    // ip
	Params    string    `json:"params"    orm:"params"     description:"参数"`    // 参数
	DeletedAt time.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"` // 软删除时间
}
