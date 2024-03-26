package admin

import (
	"github.com/vera-byte/vgo/modules/space/service"
	"github.com/vera-byte/vgo/v"
)

type SpaceInfoController struct {
	*v.Controller
}

func init() {
	var space_info_controller = &SpaceInfoController{
		&v.Controller{
			Perfix:  "/admin/space/info",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewSpaceInfoService(),
		},
	}
	// 注册路由
	v.RegisterController(space_info_controller)
}
