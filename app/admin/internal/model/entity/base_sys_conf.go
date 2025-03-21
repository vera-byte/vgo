// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package entity

import (
	"time"
)

// BaseSysConf is the golang structure for table base_sys_conf.
type BaseSysConf struct {
	Id        int       `json:"id"        orm:"id"         description:"ID"`    // ID
	CreatedAt time.Time `json:"createdAt" orm:"created_at" description:"创建时间"`  // 创建时间
	UpdatedAt time.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`  // 更新时间
	TenantId  int       `json:"tenantId"  orm:"tenant_id"  description:"租户ID"`  // 租户ID
	CKey      string    `json:"cKey"      orm:"c_key"      description:"配置键"`   // 配置键
	CValue    string    `json:"cValue"    orm:"c_value"    description:"配置值"`   // 配置值
	DeletedAt time.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"` // 软删除时间
}
