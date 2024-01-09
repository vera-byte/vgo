package admin

import (
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysRoleController struct {
	*v.Controller
}

func init() {
	var base_sys_role_controller = &BaseSysRoleController{
		&v.Controller{
			Perfix:  "/admin/base/sys/role",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysRoleService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_role_controller)
}
