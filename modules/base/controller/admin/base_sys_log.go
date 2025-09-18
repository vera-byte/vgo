package admin

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseSysLogController struct {
	*v.Controller
}

func init() {
	var base_sys_log_controller = &BaseSysLogController{
		&v.Controller{
			Perfix:  "/admin/base/sys/log",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewBaseSysLogService(),
		},
	}
	// 注册路由
	v.RegisterController(base_sys_log_controller)
}

// SetKeep 设置保留天数
// 功能: 设置系统日志的保留天数
// 参数: ctx - 上下文, req - 设置保留天数请求
// 返回值: res - 响应结果, err - 错误信息
func (c *BaseSysLogController) SetKeep(ctx g.Ctx, req *v1.SetKeepReq) (res *v.BaseRes, err error) {
	var (
		BaseSysConfService = service.NewBaseSysConfService()
	)
	err = BaseSysConfService.UpdateValue("logKeep", gconv.String(req.Value))
	return
}

// GetKeep 获取保留天数
// 功能: 获取系统日志的保留天数配置
// 参数: ctx - 上下文, req - 获取保留天数请求
// 返回值: res - 响应结果包含保留天数, err - 错误信息
func (c *BaseSysLogController) GetKeep(ctx g.Ctx, req *v1.GetKeepReq) (res *v.BaseRes, err error) {
	var (
		BaseSysConfService = service.NewBaseSysConfService()
	)
	// res.Data = BaseSysConfService.GetValue("logKeep")
	res = v.Ok(BaseSysConfService.GetValue("logKeep"))
	return
}

// Clear 清空日志
// 功能: 清空所有系统日志记录
// 参数: ctx - 上下文, req - 清空日志请求
// 返回值: res - 响应结果, err - 错误信息
func (c *BaseSysLogController) Clear(ctx g.Ctx, req *v1.ClearReq) (res *v.BaseRes, err error) {
	var (
		BaseSysLogService = service.NewBaseSysLogService()
	)
	err = BaseSysLogService.Clear(true)
	return
}
