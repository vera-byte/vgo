package admin

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/base/service"
)

type BaseSysRoleController struct {
	*cool.Controller
}

func init() {
	var base_sys_role_controller = &BaseSysRoleController{
		&cool.Controller{
			Perfix:  "/admin/base/sys/role",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysRoleService(),
		},
	}
	// 注册路由
	cool.RegisterController(base_sys_role_controller)
}
