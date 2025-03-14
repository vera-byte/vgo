// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysParam is the golang structure for table base_sys_param.
type BaseSysParam struct {
	Id         int    `json:"id"         orm:"id"         description:"ID"`                    // ID
	CreateTime string `json:"createTime" orm:"createTime" description:"创建时间"`                  // 创建时间
	UpdateTime string `json:"updateTime" orm:"updateTime" description:"更新时间"`                  // 更新时间
	TenantId   int    `json:"tenantId"   orm:"tenantId"   description:"租户ID"`                  // 租户ID
	KeyName    string `json:"keyName"    orm:"keyName"    description:"键"`                     // 键
	Name       string `json:"name"       orm:"name"       description:"名称"`                    // 名称
	Data       string `json:"data"       orm:"data"       description:"数据"`                    // 数据
	DataType   int    `json:"dataType"   orm:"dataType"   description:"数据类型 0-字符串 1-富文本 2-文件"` // 数据类型 0-字符串 1-富文本 2-文件
	Remark     string `json:"remark"     orm:"remark"     description:"备注"`                    // 备注
}
