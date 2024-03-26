package admin

import (
	"github.com/vera-byte/vgo/modules/space/service"
	"github.com/vera-byte/vgo/v"
)

type SpaceTypeController struct {
	*v.Controller
}

func init() {
	var space_type_controller = &SpaceTypeController{
		&v.Controller{
			Perfix:  "/admin/space/type",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewSpaceTypeService(),
		},
	}
	// 注册路由
	v.RegisterController(space_type_controller)
}
