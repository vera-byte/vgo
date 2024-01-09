package admin

import (
	"github.com/vera-byte/vgo/cool"
	"github.com/vera-byte/vgo/modules/dict/service"
)

type DictTypeController struct {
	*cool.Controller
}

func init() {
	var dict_type_controller = &DictTypeController{
		&cool.Controller{
			Perfix:  "/admin/dict/type",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewDictTypeService(),
		},
	}
	// 注册路由
	cool.RegisterController(dict_type_controller)
}
