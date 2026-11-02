package v1

import "github.com/gogf/gf/v2/frame/g"

// UserLoginMiniReq 小程序登录
type UserLoginMiniReq struct {
	g.Meta        `path:"/mini" method:"POST" summary:"小程序登录" tags:"用户登录"`
	Code          string `json:"code" v:"required#code不能为空"`
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

// UserLoginMpReq 公众号登录
type UserLoginMpReq struct {
	g.Meta `path:"/mp" method:"POST" summary:"公众号登录" tags:"用户登录"`
	Code   string `json:"code" v:"required#code不能为空"`
}

// UserLoginWxAppReq 微信APP授权登录
type UserLoginWxAppReq struct {
	g.Meta `path:"/wxApp" method:"POST" summary:"微信APP授权登录" tags:"用户登录"`
	Code   string `json:"code" v:"required#code不能为空"`
}

// UserLoginPhoneReq 手机验证码登录
type UserLoginPhoneReq struct {
	g.Meta  `path:"/phone" method:"POST" summary:"手机号登录" tags:"用户登录"`
	Phone   string `json:"phone" v:"required#手机号不能为空"`
	SmsCode string `json:"smsCode" v:"required#验证码不能为空"`
}

// UserLoginUniPhoneReq 一键手机号登录
type UserLoginUniPhoneReq struct {
	g.Meta      `path:"/uniPhone" method:"POST" summary:"一键手机号登录" tags:"用户登录"`
	AccessToken string `json:"access_token" v:"required#access_token不能为空"`
	Openid      string `json:"openid" v:"required#openid不能为空"`
	AppId       string `json:"appId"`
}

// UserLoginMiniPhoneReq 绑定小程序手机号
type UserLoginMiniPhoneReq struct {
	g.Meta        `path:"/miniPhone" method:"POST" summary:"绑定小程序手机号" tags:"用户登录"`
	Code          string `json:"code" v:"required#code不能为空"`
	EncryptedData string `json:"encryptedData" v:"required#encryptedData不能为空"`
	Iv            string `json:"iv" v:"required#iv不能为空"`
}

// UserLoginCaptchaReq 图片验证码
type UserLoginCaptchaReq struct {
	g.Meta `path:"/captcha" method:"GET" summary:"图片验证码" tags:"用户登录"`
	Width  int `json:"width" d:"150"`
	Height int `json:"height" d:"40"`
	Color  string `json:"color"`
}

// UserLoginSmsCodeReq 发送短信验证码
type UserLoginSmsCodeReq struct {
	g.Meta    `path:"/smsCode" method:"POST" summary:"发送短信验证码" tags:"用户登录"`
	Phone     string `json:"phone" v:"required#手机号不能为空"`
	CaptchaId string `json:"captchaId" v:"required#captchaId不能为空"`
	Code      string `json:"code" v:"required#验证码不能为空"`
}

// UserLoginRefreshReq 刷新token
type UserLoginRefreshReq struct {
	g.Meta       `path:"/refreshToken" method:"POST" summary:"刷新token" tags:"用户登录"`
	RefreshToken string `json:"refreshToken" v:"required#refreshToken不能为空"`
}

// UserLoginPasswordReq 密码登录
type UserLoginPasswordReq struct {
	g.Meta   `path:"/password" method:"POST" summary:"密码登录" tags:"用户登录"`
	Phone    string `json:"phone" v:"required#手机号不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}
