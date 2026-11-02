package v1

import "github.com/gogf/gf/v2/frame/g"

// UserInfoPersonReq 获取个人信息
type UserInfoPersonReq struct {
	g.Meta `path:"/person" method:"GET" summary:"获取个人信息" tags:"用户信息"`
}

// UserInfoUpdatePersonReq 更新个人信息
type UserInfoUpdatePersonReq struct {
	g.Meta   `path:"/updatePerson" method:"POST" summary:"更新个人信息" tags:"用户信息"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	Gender    *int   `json:"gender"`
	Description string `json:"description"`
}

// UserInfoUpdatePasswordReq 修改密码
type UserInfoUpdatePasswordReq struct {
	g.Meta   `path:"/updatePassword" method:"POST" summary:"修改密码" tags:"用户信息"`
	Password string `json:"password" v:"required#密码不能为空"`
	Code     string `json:"code" v:"required#验证码不能为空"`
}

// UserInfoLogoffReq 注销
type UserInfoLogoffReq struct {
	g.Meta `path:"/logoff" method:"POST" summary:"注销" tags:"用户信息"`
}

// UserInfoBindPhoneReq 绑定手机号
type UserInfoBindPhoneReq struct {
	g.Meta `path:"/bindPhone" method:"POST" summary:"绑定手机号" tags:"用户信息"`
	Phone  string `json:"phone" v:"required#手机号不能为空"`
	Code   string `json:"code" v:"required#验证码不能为空"`
}

// UserInfoMiniPhoneReq 绑定小程序手机号
type UserInfoMiniPhoneReq struct {
	g.Meta        `path:"/miniPhone" method:"POST" summary:"绑定小程序手机号" tags:"用户信息"`
	Code          string `json:"code" v:"required#code不能为空"`
	EncryptedData string `json:"encryptedData" v:"required#encryptedData不能为空"`
	Iv            string `json:"iv" v:"required#iv不能为空"`
}

// UserAddressDefaultReq 默认地址
type UserAddressDefaultReq struct {
	g.Meta `path:"/default" method:"GET" summary:"默认地址" tags:"用户地址"`
}

// UserWxMpConfigReq 微信公众号配置
type UserWxMpConfigReq struct {
	g.Meta `path:"/wxMpConfig" method:"POST" summary:"获取微信公众号配置" tags:"用户通用"`
	Url    string `json:"url" v:"required#url不能为空"`
}
