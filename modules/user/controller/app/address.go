package app

import (
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/user/api/v1"
	"github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

// UserAddressController C端收货地址控制器（仅本人数据）
type UserAddressController struct {
	*v.Controller
}

func init() {
	c := &UserAddressController{
		&v.Controller{
			Prefix:  "/app/user/address",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page", "Default"},
			Service: service.NewUserAddressService(),
		},
	}
	v.RegisterController(c)
}

// Default 默认地址
func (c *UserAddressController) Default(ctx g.Ctx, req *v1.UserAddressDefaultReq) (res *v.BaseRes, err error) {
	data, err := service.NewUserAddressService().Default(ctx)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}
