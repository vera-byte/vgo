// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-19 16:43:14
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysParam is the golang structure of table base_sys_param for DAO operations like Where/Data.
type BaseSysParam struct {
	g.Meta    `orm:"table:base_sys_param, do:true"`
	Id        interface{} // ID
	CreatedAt interface{} // 创建时间
	UpdatedAt interface{} // 更新时间
	TenantId  interface{} // 租户ID
	KeyName   interface{} // 键
	Name      interface{} // 名称
	Data      interface{} // 数据
	DataType  interface{} // 数据类型 0-字符串 1-富文本 2-文件
	Remark    interface{} // 备注
	DeletedAt interface{} // 软删除时间
}
