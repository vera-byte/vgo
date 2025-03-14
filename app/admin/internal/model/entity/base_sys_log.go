// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysLog is the golang structure for table base_sys_log.
type BaseSysLog struct {
	Id         int    `json:"id"         orm:"id"         description:"ID"`   // ID
	CreateTime string `json:"createTime" orm:"createTime" description:"创建时间"` // 创建时间
	UpdateTime string `json:"updateTime" orm:"updateTime" description:"更新时间"` // 更新时间
	TenantId   int    `json:"tenantId"   orm:"tenantId"   description:"租户ID"` // 租户ID
	UserId     int    `json:"userId"     orm:"userId"     description:"用户ID"` // 用户ID
	Action     string `json:"action"     orm:"action"     description:"行为"`   // 行为
	Ip         string `json:"ip"         orm:"ip"         description:"ip"`   // ip
	Params     string `json:"params"     orm:"params"     description:"参数"`   // 参数
}
