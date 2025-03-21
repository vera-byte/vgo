// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysLog is the golang structure of table base_sys_log for DAO operations like Where/Data.
type BaseSysLog struct {
	g.Meta    `orm:"table:base_sys_log, do:true"`
	Id        interface{} // ID
	CreatedAt interface{} // 创建时间
	UpdatedAt interface{} // 更新时间
	TenantId  interface{} // 租户ID
	UserId    interface{} // 用户ID
	Action    interface{} // 行为
	Ip        interface{} // ip
	Params    interface{} // 参数
	DeletedAt interface{} // 软删除时间
}
