package model

import "github.com/vera-byte/vgo/v"

const TableNameUserWx = "user_wx"

// UserWx 微信用户（参照 cool-admin-midway user 模块 UserWxEntity）
type UserWx struct {
	*v.Model
	Unionid   *string `json:"unionid"`   // 微信unionid
	Openid    string  `json:"openid"`    // 微信openid
	AvatarUrl *string `json:"avatarUrl"` // 头像
	NickName  *string `json:"nickName"`  // 昵称
	Gender    int     `json:"gender"`    // 性别 0-未知 1-男 2-女
	Language  *string `json:"language"`  // 语言
	City      *string `json:"city"`      // 城市
	Province  *string `json:"province"`  // 省份
	Country   *string `json:"country"`   // 国家
	Type      int     `json:"type"`      // 类型 0-小程序 1-公众号 2-H5 3-APP
}

// TableName UserWx 的表名
func (*UserWx) TableName() string {
	return TableNameUserWx
}

// NewUserWx 创建一个 UserWx 模型实例
func NewUserWx() *UserWx {
	return &UserWx{Model: v.NewModel()}
}
