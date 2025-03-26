package v1

import "github.com/gogf/gf/v2/frame/g"

type CaptchaReq struct {
	g.Meta `path:"captcha" method:"get" sm:"图形验证码" tags:"开放"`
	Height int32 `json:"height" v:"required#请输入验证码高度" dc:"验证码高度"`
	Width  int32 `json:"width" v:"required#请输入验证码宽度" dc:"验证码宽度"`
}

type CaptchaRes struct {
	CaptchaId string `json:"captchaId" dc:"验证码ID"`
	Data      string `json:"data" dc:"验证码base64"`
}

type LoginReq struct {
	g.Meta     `path:"login" method:"post" sm:"账号密码登录" tags:"开放"`
	CaptchaId  string `json:"captchaId" v:"required#请输入验证码ID" dc:"验证码ID"`
	VerifyCode string `json:"verifyCode" v:"required#请输入验证码" dc:"验证码"`
	Password   string `json:"password" v:"required#请输入密码" dc:"密码"`
	Username   string `json:"username" v:"required#请输入用户名" dc:"用户名"`
}

type LoginRes struct {
	Expire        int64  `json:"expire"  dc:"过期时间"`
	RefreshExpire int64  `json:"refreshExpire" dc:"刷新时效"`
	RefreshToken  string `json:"refreshToken"  dc:"刷新token"`
	Token         string `json:"token"  dc:"访问token"`
}

type RefreshTokenReq struct {
	g.Meta       `path:"refreshToken" method:"get" sm:"刷新token" tags:"开放"`
	RefreshToken string `json:"refreshToken" v:"required#请输入token" dc:"token"`
}

type RefreshTokenRes struct {
	Expire        int64  `json:"expire"  dc:"过期时间"`
	RefreshExpire int64  `json:"refreshExpire" dc:"刷新时效"`
	RefreshToken  string `json:"refreshToken"  dc:"刷新token"`
	Token         string `json:"token"  dc:"访问token"`
}
