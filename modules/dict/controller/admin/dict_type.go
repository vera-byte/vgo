package admin

import (
	"github.com/vera-byte/vgo/modules/dict/service"
	"github.com/vera-byte/vgo/v"
)

type DictTypeController struct {
	*v.Controller
}

func init() {
	var dict_type_controller = &DictTypeController{
		&v.Controller{
			Perfix:  "/admin/dict/type",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewDictTypeService(),
		},
	}
	// 注册路由
	v.RegisterController(dict_type_controller)
}
