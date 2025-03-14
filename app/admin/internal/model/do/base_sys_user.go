// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysUser is the golang structure of table base_sys_user for DAO operations like Where/Data.
type BaseSysUser struct {
	g.Meta       `orm:"table:base_sys_user, do:true"`
	Id           interface{} // ID
	CreateTime   interface{} // 创建时间
	UpdateTime   interface{} // 更新时间
	TenantId     interface{} // 租户ID
	DepartmentId interface{} // 部门ID
	Name         interface{} // 姓名
	Username     interface{} // 用户名
	Password     interface{} // 密码
	PasswordV    interface{} // 密码版本, 作用是改完密码，让原来的token失效
	NickName     interface{} // 昵称
	HeadImg      interface{} // 头像
	Phone        interface{} // 手机
	Email        interface{} // 邮箱
	Remark       interface{} // 备注
	Status       interface{} // 状态 0-禁用 1-启用
	SocketId     interface{} // socketId
}
