package admin

import (
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/plugin/api/v1"
	"github.com/vera-byte/vgo/modules/plugin/service"
	"github.com/vera-byte/vgo/v"
)

// PluginInfoController 插件管理控制器
type PluginInfoController struct {
	*v.Controller
}

func init() {
	svc := service.NewPluginService()
	c := &PluginInfoController{
		&v.Controller{
			Prefix:  "/admin/plugin/info",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page", "Install", "Reload", "Invoke"},
			Service: svc,
		},
	}
	v.RegisterController(c)
	// 启动时加载全部启用插件（延后到 server 启动，避免任意子命令在 init 阶段触库）
	v.RegisterStartupHook("plugin-init-all", func(ctx g.Ctx) {
		svc.InitAll()
	})
}

// Install 安装插件（ZIP 包）
func (c *PluginInfoController) Install(ctx g.Ctx, req *v1.PluginInfoInstallReq) (res *v.BaseRes, err error) {
	file := g.RequestFromCtx(ctx).GetUploadFile("file")
	if file == nil {
		return nil, gerror.New("请上传插件包（file 字段）")
	}
	f, e := file.Open()
	if e != nil {
		return nil, e
	}
	defer f.Close()
	raw, e := io.ReadAll(f)
	if e != nil {
		return nil, e
	}
	force := req.Force || g.RequestFromCtx(ctx).Get("force").Bool()
	result, e := service.NewPluginService().Install(ctx, raw, force)
	if e != nil {
		return nil, e
	}
	res = v.Ok(result)
	return
}

// Invoke 调用插件方法（演示/联调用；生产调用建议在服务代码中经 service.Invoke）
func (c *PluginInfoController) Invoke(ctx g.Ctx, req *v1.PluginInfoInvokeReq) (res *v.BaseRes, err error) {
	data, err := service.Invoke(ctx, req.Key, req.Method, req.Params...)
	if err != nil {
		return
	}
	res = v.Ok(data)
	return
}

// Reload 重载插件（热更新）
func (c *PluginInfoController) Reload(ctx g.Ctx, req *v1.PluginInfoReloadReq) (res *v.BaseRes, err error) {
	service.NewPluginService().ReloadKey(req.Key)
	res = v.Ok(nil)
	return
}
