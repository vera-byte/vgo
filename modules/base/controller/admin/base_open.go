package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseOpen struct {
	*v.ControllerSimple
	baseSysLoginService *service.BaseSysLoginService
	baseOpenService     *service.BaseOpenService
}

func init() {
	var open = &BaseOpen{
		ControllerSimple:    &v.ControllerSimple{Perfix: "/admin/base/open"},
		baseSysLoginService: service.NewBaseSysLoginService(),
		baseOpenService:     service.NewBaseOpenService(),
	}
	// 注册路由
	v.RegisterControllerSimple(open)
}

// 验证码接口
func (c *BaseOpen) BaseOpenCaptcha(ctx context.Context, req *v1.BaseOpenCaptchaReq) (res *v.BaseRes, err error) {
	data, err := c.baseSysLoginService.Captcha(req)
	res = v.Ok(data)
	return
}

// eps 接口请求
type BaseOpenEpsReq struct {
	g.Meta `path:"/eps" method:"GET"`
}

// eps 接口
func (c *BaseOpen) Eps(ctx context.Context, req *BaseOpenEpsReq) (res *v.BaseRes, err error) {
	if !v.Config.Eps {
		g.Log().Error(ctx, "eps is not open")
		res = v.Ok(nil)
		return
	}
	data, err := c.baseOpenService.AdminEPS(ctx)
	if err != nil {
		g.Log().Error(ctx, "eps error", err)
		return v.Fail(err.Error()), err
	}
	res = v.Ok(data)
	return
}

// login 接口
func (c *BaseOpen) Login(ctx context.Context, req *v1.BaseOpenLoginReq) (res *v.BaseRes, err error) {
	data, err := c.baseSysLoginService.Login(ctx, req)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// RefreshTokenReq 刷新token请求
type RefreshTokenReq struct {
	g.Meta       `path:"/refreshToken" method:"GET"`
	RefreshToken string `json:"refreshToken" v:"required#refreshToken不能为空"`
}

// RefreshToken 刷新token
func (c *BaseOpen) RefreshToken(ctx context.Context, req *RefreshTokenReq) (res *v.BaseRes, err error) {
	data, err := c.baseSysLoginService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}
