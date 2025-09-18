package admin

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysParamController struct {
	*v.Controller
}

func init() {
	var base_sys_param_controller = &BaseSysParamController{
		&v.Controller{
			Perfix:  "/admin/base/sys/param",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysParamService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_param_controller)
}

// Html 根据配置参数key获取网页内容(富文本)
func (c *BaseSysParamController) Html(ctx g.Ctx, req *v1.BaseSysParamHtmlReq) (res *v.BaseRes, err error) {
	var (
		BaseSysParamService = service.NewBaseSysParamService()
		r                   = ghttp.RequestFromCtx(ctx)
	)
	r.Response.WriteExit(BaseSysParamService.HtmlByKey(req.Key))
	return
}
