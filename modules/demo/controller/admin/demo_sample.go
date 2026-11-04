package admin

import (
	v1 "vgo-simple/modules/demo/api/v1"
	"vgo-simple/modules/demo/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/v"
)

// DemoSampleController 示例模块控制器：通用 CRUD + 1 个自定义接口
type DemoSampleController struct {
	*v.Controller
}

func init() {
	v.RegisterController(&DemoSampleController{
		&v.Controller{
			Prefix:  "/admin/demo/sample",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewDemoSampleService(),
		},
	})
}

// Welcome 示例自定义接口
func (c *DemoSampleController) Welcome(ctx g.Ctx, req *v1.DemoSampleWelcomeReq) (res *v.BaseRes, err error) {
	res = v.Ok("Welcome to Vgo")
	return
}
