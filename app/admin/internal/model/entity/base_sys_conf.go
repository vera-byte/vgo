// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysConf is the golang structure for table base_sys_conf.
type BaseSysConf struct {
	Id         int    `json:"id"         orm:"id"         description:"ID"`   // ID
	CreateTime string `json:"createTime" orm:"createTime" description:"创建时间"` // 创建时间
	UpdateTime string `json:"updateTime" orm:"updateTime" description:"更新时间"` // 更新时间
	TenantId   int    `json:"tenantId"   orm:"tenantId"   description:"租户ID"` // 租户ID
	CKey       string `json:"cKey"       orm:"cKey"       description:"配置键"`  // 配置键
	CValue     string `json:"cValue"     orm:"cValue"     description:"配置值"`  // 配置值
}
