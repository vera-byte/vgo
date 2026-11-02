package app

import (
	"github.com/gogf/gf/v2/frame/g"
	baseApi "github.com/vera-byte/vgo/modules/base/api/v1"
	baseService "github.com/vera-byte/vgo/modules/base/service"
	v1 "github.com/vera-byte/vgo/modules/user/api/v1"
	"github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

// UserLoginController C端登录控制器（全部免登录接口）
type UserLoginController struct {
	*v.Controller
}

func init() {
	c := &UserLoginController{
		&v.Controller{
			Prefix:  "/app/user/login",
			Api:     []string{"Mini", "Mp", "WxApp", "Phone", "UniPhone", "MiniPhone", "Captcha", "SmsCode", "RefreshToken", "Password"},
			Service: service.NewUserInfoService(),
		},
	}
	v.RegisterController(c)
}

// Mini 小程序登录
func (c *UserLoginController) Mini(ctx g.Ctx, req *v1.UserLoginMiniReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().Mini(ctx, req.Code, req.EncryptedData, req.Iv)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// Mp 公众号登录
func (c *UserLoginController) Mp(ctx g.Ctx, req *v1.UserLoginMpReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().Mp(ctx, req.Code)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// WxApp 微信APP授权登录
func (c *UserLoginController) WxApp(ctx g.Ctx, req *v1.UserLoginWxAppReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().WxApp(ctx, req.Code)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// Phone 手机验证码登录
func (c *UserLoginController) Phone(ctx g.Ctx, req *v1.UserLoginPhoneReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().PhoneVerifyCode(ctx, req.Phone, req.SmsCode)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// UniPhone 一键手机号登录
func (c *UserLoginController) UniPhone(ctx g.Ctx, req *v1.UserLoginUniPhoneReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().UniPhone(ctx, req.AccessToken, req.Openid, req.AppId)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// MiniPhone 绑定小程序手机号（登录）
func (c *UserLoginController) MiniPhone(ctx g.Ctx, req *v1.UserLoginMiniPhoneReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().MiniPhone(ctx, req.Code, req.EncryptedData, req.Iv)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// Captcha 图片验证码（复用基础模块能力）
func (c *UserLoginController) Captcha(ctx g.Ctx, req *v1.UserLoginCaptchaReq) (res *v.BaseRes, err error) {
	data, err := baseService.NewBaseSysLoginService().Captcha(&baseApi.BaseOpenCaptchaReq{
		Width:  req.Width,
		Height: req.Height,
		Color:  req.Color,
	})
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// SmsCode 发送短信验证码
func (c *UserLoginController) SmsCode(ctx g.Ctx, req *v1.UserLoginSmsCodeReq) (res *v.BaseRes, err error) {
	if err = service.NewUserLoginService().SmsCode(ctx, req.Phone, req.CaptchaId, req.Code); err != nil {
		return
	}
	res = v.Ok(nil)
	return
}

// RefreshToken 刷新token
func (c *UserLoginController) RefreshToken(ctx g.Ctx, req *v1.UserLoginRefreshReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// Password 密码登录
func (c *UserLoginController) Password(ctx g.Ctx, req *v1.UserLoginPasswordReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserLoginService().Password(ctx, req.Phone, req.Password)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}
