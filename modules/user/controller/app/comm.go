package app

import (
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/user/api/v1"
	"github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

// UserCommController C端通用控制器
type UserCommController struct {
	*v.Controller
}

func init() {
	c := &UserCommController{
		&v.Controller{
			Prefix:  "/app/user/comm",
			Api:     []string{"WxMpConfig"},
			Service: service.NewUserInfoService(),
		},
	}
	v.RegisterController(c)
}

// WxMpConfig 获取微信公众号配置
func (c *UserCommController) WxMpConfig(ctx g.Ctx, req *v1.UserWxMpConfigReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserWxService().WxMpConfig(ctx, req.Url)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}
