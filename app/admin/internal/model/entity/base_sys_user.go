// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2025-03-14 15:21:22
// =================================================================================

package entity

// BaseSysUser is the golang structure for table base_sys_user.
type BaseSysUser struct {
	Id           int    `json:"id"           orm:"id"           description:"ID"`                        // ID
	CreateTime   string `json:"createTime"   orm:"createTime"   description:"创建时间"`                      // 创建时间
	UpdateTime   string `json:"updateTime"   orm:"updateTime"   description:"更新时间"`                      // 更新时间
	TenantId     int    `json:"tenantId"     orm:"tenantId"     description:"租户ID"`                      // 租户ID
	DepartmentId int    `json:"departmentId" orm:"departmentId" description:"部门ID"`                      // 部门ID
	Name         string `json:"name"         orm:"name"         description:"姓名"`                        // 姓名
	Username     string `json:"username"     orm:"username"     description:"用户名"`                       // 用户名
	Password     string `json:"password"     orm:"password"     description:"密码"`                        // 密码
	PasswordV    int    `json:"passwordV"    orm:"passwordV"    description:"密码版本, 作用是改完密码，让原来的token失效"` // 密码版本, 作用是改完密码，让原来的token失效
	NickName     string `json:"nickName"     orm:"nickName"     description:"昵称"`                        // 昵称
	HeadImg      string `json:"headImg"      orm:"headImg"      description:"头像"`                        // 头像
	Phone        string `json:"phone"        orm:"phone"        description:"手机"`                        // 手机
	Email        string `json:"email"        orm:"email"        description:"邮箱"`                        // 邮箱
	Remark       string `json:"remark"       orm:"remark"       description:"备注"`                        // 备注
	Status       int    `json:"status"       orm:"status"       description:"状态 0-禁用 1-启用"`              // 状态 0-禁用 1-启用
	SocketId     string `json:"socketId"     orm:"socketId"     description:"socketId"`                  // socketId
}
