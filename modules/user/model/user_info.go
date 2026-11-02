package model

import "github.com/vera-byte/vgo/v"

const TableNameUserInfo = "user_info"

// UserInfo 用户信息（参照 cool-admin-midway user 模块 UserInfoEntity）
type UserInfo struct {
	*v.Model
	Unionid   *string `json:"unionid"`   // 登录唯一ID
	AvatarUrl *string `json:"avatarUrl"` // 头像
	NickName  *string `json:"nickName"`  // 昵称
	Phone     *string `json:"phone"`     // 手机号
	Gender    int     `json:"gender"`    // 性别 0-未知 1-男 2-女
	Status    int     `json:"status"`    // 状态 0-禁用 1-正常 2-已注销
	LoginType int     `json:"loginType"` // 登录方式 0-小程序 1-公众号 2-H5
	Password  *string `json:"password"`  // 密码
	Description *string `json:"description"` // 介绍
}

// TableName UserInfo 的表名
func (*UserInfo) TableName() string {
	return TableNameUserInfo
}

// NewUserInfo 创建一个 UserInfo 模型实例
func NewUserInfo() *UserInfo {
	return &UserInfo{Model: v.NewModel()}
}
