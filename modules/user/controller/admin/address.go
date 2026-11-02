package admin

import (
	userService "github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

// UserAddressController 管理端-收货地址控制器
type UserAddressController struct {
	*v.Controller
}

func init() {
	c := &UserAddressController{
		&v.Controller{
			Prefix:  "/admin/user/address",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: userService.NewUserAddressService(),
		},
	}
	v.RegisterController(c)
}
