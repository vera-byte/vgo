// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysConf is the golang structure of table base_sys_conf for DAO operations like Where/Data.
type BaseSysConf struct {
	g.Meta     `orm:"table:base_sys_conf, do:true"`
	Id         interface{} // ID
	CreateTime interface{} // 创建时间
	UpdateTime interface{} // 更新时间
	TenantId   interface{} // 租户ID
	CKey       interface{} // 配置键
	CValue     interface{} // 配置值
}
