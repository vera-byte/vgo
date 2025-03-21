// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysConf is the golang structure of table base_sys_conf for DAO operations like Where/Data.
type BaseSysConf struct {
	g.Meta    `orm:"table:base_sys_conf, do:true"`
	Id        interface{} // ID
	CreatedAt interface{} // 创建时间
	UpdatedAt interface{} // 更新时间
	TenantId  interface{} // 租户ID
	CKey      interface{} // 配置键
	CValue    interface{} // 配置值
	DeletedAt interface{} // 软删除时间
}
