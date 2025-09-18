package v1

import "github.com/gogf/gf/v2/frame/g"

// login 接口请求
type BaseOpenLoginReq struct {
	g.Meta     `path:"/login" method:"POST" summary:"用户登录" tags:"开放接口"`
	Username   string `json:"username" p:"username" v:"required"`
	Password   string `json:"password" p:"password" v:"required"`
	CaptchaId  string `json:"captchaId" p:"captchaId" v:"required"`
	VerifyCode string `json:"verifyCode" p:"verifyCode" v:"required"`
}

// captcha 验证码接口

type BaseOpenCaptchaReq struct {
	g.Meta `path:"/captcha" method:"GET" summary:"获取验证码" tags:"开放接口"`
	Height int    `json:"height" in:"query" default:"40"`
	Width  int    `json:"width"  in:"query" default:"150"`
	Color  string `json:"color" in:"query" default:"#2c3142"`
}
