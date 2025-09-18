package admin

import (
	"context"

	v1 "github.com/vera-byte/vgo/modules/dict/api/v1"
	"github.com/vera-byte/vgo/modules/dict/service"
	"github.com/vera-byte/vgo/v"
)

type DictInfoController struct {
	*v.Controller
}

func init() {
	var dict_info_controller = &DictInfoController{
		&v.Controller{
			Perfix:  "/admin/dict/info",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewDictInfoService(),
		},
	}
	// 注册路由
	v.RegisterController(dict_info_controller)
}

// Data 方法 获得字典数据
func (c *DictInfoController) Data(ctx context.Context, req *v1.DictInfoDataReq) (res *v.BaseRes, err error) {
	service := service.NewDictInfoService()
	data, err := service.Data(ctx, req.Types)
	res = v.Ok(data)
	return
}
