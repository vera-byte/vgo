package admin

import (
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysMenuController struct {
	*v.Controller
}

func init() {
	var base_sys_menu_controller = &BaseSysMenuController{
		&v.Controller{
			Perfix:  "/admin/base/sys/menu",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysMenuService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_menu_controller)
}
