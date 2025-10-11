package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysUser = "base_sys_user"

// BaseSysUser mapped from table <base_sys_user>
type BaseSysUser struct {
	*v.Model
	DepartmentID uint    `json:"departmentId"` // 部门ID
	Name         *string `json:"name"`         // 姓名
	Username     string  `json:"username"`     // 用户名
	Password     string  `json:"password"`     // 密码
	PasswordV    *int32  `json:"passwordV"`    // 密码版本, 作用是改完密码，让原来的token失效
	NickName     *string `json:"nickName"`     // 昵称
	HeadImg      *string `json:"headImg"`      // 头像
	Phone        *string `json:"phone"`        // 手机
	Email        *string `json:"email"`        // 邮箱
	Status       *int32  `json:"status"`       // 状态 0:禁用 1：启用
	Remark       *string `json:"remark"`       // 备注
	SocketID     *string `json:"socketId"`     // socketId
}

// TableName BaseSysUser's table name
func (*BaseSysUser) TableName() string {
	return TableNameBaseSysUser
}

// NewBaseSysUser 创建一个新的BaseSysUser
func NewBaseSysUser() *BaseSysUser {
	return &BaseSysUser{
		Model: v.NewModel(),
	}
}
