package admin

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "github.com/vera-byte/vgo/modules/task/api/v1"
	"github.com/vera-byte/vgo/modules/task/service"
	"github.com/vera-byte/vgo/v"
)

type TaskInfoController struct {
	*v.Controller
}

func init() {
	var task_info_controller = &TaskInfoController{
		&v.Controller{
			Perfix:  "/admin/task/info",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page", "Start", "Stop"},
			Service: service.NewTaskInfoService(),
		},
	}
	// 注册路由
	v.RegisterController(task_info_controller)
}

// Stop 停止任务
// 功能: 停止指定的定时任务
// 参数: ctx - 上下文, req - 停止任务请求
// 返回值: res - 响应结果, err - 错误信息
func (c *TaskInfoController) Stop(ctx g.Ctx, req *v1.TaskInfoStopReq) (res *v.BaseRes, err error) {

	err = v.ClusterRunFunc(ctx, "TaskStopFunc("+gconv.String(req.ID)+")")
	if err != nil {
		return v.Fail(err.Error()), err
	}
	res = v.Ok("停止成功")
	return
}

// Start 启动任务
// 功能: 启动指定的定时任务
// 参数: ctx - 上下文, req - 启动任务请求
// 返回值: res - 响应结果, err - 错误信息
func (c *TaskInfoController) Start(ctx g.Ctx, req *v1.TaskInfoStartReq) (res *v.BaseRes, err error) {

	err = v.ClusterRunFunc(ctx, "TaskStartFunc("+gconv.String(req.ID)+")")
	if err != nil {
		return v.Fail(err.Error()), err
	}
	res = v.Ok("启动成功")
	return
}

// Once 执行一次
// 功能: 立即执行一次指定的定时任务
// 参数: ctx - 上下文, req - 执行一次任务请求
// 返回值: res - 响应结果, err - 错误信息
func (c *TaskInfoController) Once(ctx g.Ctx, req *v1.TaskInfoOnceReq) (res *v.BaseRes, err error) {
	err = c.Service.(*service.TaskInfoService).Once(ctx, req.ID)
	if err != nil {
		return v.Fail(err.Error()), err
	}
	res = v.Ok("执行成功")
	return
}

// Log 任务日志
// 功能: 获取指定任务的执行日志
// 参数: ctx - 上下文, req - 任务日志请求
// 返回值: res - 响应结果包含日志数据, err - 错误信息
func (c *TaskInfoController) Log(ctx g.Ctx, req *v1.TaskInfoLogReq) (res *v.BaseRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	param := r.GetQueryMapStrStr()
	data, err := c.Service.(*service.TaskInfoService).Log(ctx, param)
	if err != nil {
		return v.Fail(err.Error()), err
	}
	res = v.Ok(data)
	return
}
