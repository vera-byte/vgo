package admin

import (
	"github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

// UserInfoController 管理端-用户管理控制器
type UserInfoController struct {
	*v.Controller
}

func init() {
	c := &UserInfoController{
		&v.Controller{
			Prefix:  "/admin/user/info",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewUserInfoService(),
		},
	}
	v.RegisterController(c)
}
