// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package entity

import (
	"time"
)

// BaseSysParam is the golang structure for table base_sys_param.
type BaseSysParam struct {
	Id        int       `json:"id"        orm:"id"         description:"ID"`                    // ID
	CreatedAt time.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                  // 创建时间
	UpdatedAt time.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                  // 更新时间
	TenantId  int       `json:"tenantId"  orm:"tenant_id"  description:"租户ID"`                  // 租户ID
	KeyName   string    `json:"keyName"   orm:"key_name"   description:"键"`                     // 键
	Name      string    `json:"name"      orm:"name"       description:"名称"`                    // 名称
	Data      string    `json:"data"      orm:"data"       description:"数据"`                    // 数据
	DataType  int       `json:"dataType"  orm:"data_type"  description:"数据类型 0-字符串 1-富文本 2-文件"` // 数据类型 0-字符串 1-富文本 2-文件
	Remark    string    `json:"remark"    orm:"remark"     description:"备注"`                    // 备注
	DeletedAt time.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"`                 // 软删除时间
}
