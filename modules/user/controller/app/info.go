package app

import (
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/user/api/v1"
	"github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

// UserInfoController C端用户信息控制器
type UserInfoController struct {
	*v.Controller
}

func init() {
	c := &UserInfoController{
		&v.Controller{
			Prefix:  "/app/user/info",
			Api:     []string{"Person", "UpdatePerson", "UpdatePassword", "Logoff", "BindPhone", "MiniPhone"},
			Service: service.NewUserInfoService(),
		},
	}
	v.RegisterController(c)
}

// Person 获取个人信息
func (c *UserInfoController) Person(ctx g.Ctx, req *v1.UserInfoPersonReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserInfoService().Person(ctx, service.CurrentUser(ctx).ID)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// UpdatePerson 更新个人信息
func (c *UserInfoController) UpdatePerson(ctx g.Ctx, req *v1.UserInfoUpdatePersonReq) (res *v.BaseRes, err error) {
	err = service.NewUserInfoService().UpdatePerson(ctx, service.CurrentUser(ctx).ID)
	if err != nil {
		return
	}
	res = v.Ok(nil)
	return
}

// UpdatePassword 修改密码
func (c *UserInfoController) UpdatePassword(ctx g.Ctx, req *v1.UserInfoUpdatePasswordReq) (res *v.BaseRes, err error) {
	err = service.NewUserInfoService().UpdatePassword(ctx, service.CurrentUser(ctx).ID, req.Password, req.Code)
	if err != nil {
		return
	}
	res = v.Ok(nil)
	return
}

// Logoff 注销
func (c *UserInfoController) Logoff(ctx g.Ctx, req *v1.UserInfoLogoffReq) (res *v.BaseRes, err error) {
	err = service.NewUserInfoService().Logoff(ctx, service.CurrentUser(ctx).ID)
	if err != nil {
		return
	}
	res = v.Ok(nil)
	return
}

// BindPhone 绑定手机号
func (c *UserInfoController) BindPhone(ctx g.Ctx, req *v1.UserInfoBindPhoneReq) (res *v.BaseRes, err error) {
	err = service.NewUserInfoService().BindPhone(ctx, service.CurrentUser(ctx).ID, req.Phone, req.Code)
	if err != nil {
		return
	}
	res = v.Ok(nil)
	return
}

// MiniPhone 绑定小程序手机号
func (c *UserInfoController) MiniPhone(ctx g.Ctx, req *v1.UserInfoMiniPhoneReq) (res *v.BaseRes, err error) {
	phone, err := service.NewUserInfoService().MiniPhone(ctx, service.CurrentUser(ctx).ID, req.Code, req.EncryptedData, req.Iv)
	if err != nil {
		return
	}
	res = v.Ok(g.Map{"phone": phone})
	return
}
