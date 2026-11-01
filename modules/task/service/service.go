package service

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/vera-byte/vgo/v"
)

// EnableTask 启用任务
func EnableTask(ctx g.Ctx, cronId string, funcstring string, cron string, startDate string) (err error) {
	funcName := gstr.SubStr(funcstring, 0, gstr.Pos(funcstring, "("))
	if _, ok := v.FuncMap[funcName]; !ok {
		err = gerror.New("函数不存在" + funcName)
		return
	}
	taskInfoService := NewTaskInfoService()

	// startDate 为启用时间（Y-m-d H:i:s），未到时间不执行；空值/不可解析表示立即生效
	var startTime *gtime.Time
	if startDate != "" {
		startTime = gtime.NewFromStr(startDate)
	}
	reached := func() bool {
		return startTime == nil || !gtime.Now().Before(startTime)
	}

	if v.FuncMap[funcName].IsSingleton() {
		gcron.Remove(cronId)
		_, err = gcron.AddSingleton(ctx, cron, func(ctx g.Ctx) {
			if !reached() {
				g.Log().Debug(ctx, "当前时间小于启用时间, 不执行单例函数", funcName)
				return
			}
			err := v.RunFunc(ctx, funcstring)
			if err != nil {
				g.Log().Error(ctx, err)
				if e := taskInfoService.Record(ctx, cronId, 0, err.Error()); e != nil {
					g.Log().Error(ctx, "记录任务执行结果失败:", e)
				}
			} else {
				if e := taskInfoService.Record(ctx, cronId, 1, "任务执行成功"); e != nil {
					g.Log().Error(ctx, "记录任务执行结果失败:", e)
				}
			}
		}, cronId)
	} else {
		gcron.Remove(cronId)
		_, err = gcron.Add(ctx, cron, func(ctx g.Ctx) {
			if !reached() {
				g.Log().Debug(ctx, "当前时间小于启用时间, 不执行函数", funcName)
				return
			}
			err := v.RunFunc(ctx, funcstring)
			if err != nil {
				g.Log().Error(ctx, err)
				if e := taskInfoService.Record(ctx, cronId, 0, gstr.AddSlashes(err.Error())); e != nil {
					g.Log().Error(ctx, "记录任务执行结果失败:", e)
				}
			} else {
				if e := taskInfoService.Record(ctx, cronId, 1, gstr.AddSlashes("任务执行成功")); e != nil {
					g.Log().Error(ctx, "记录任务执行结果失败:", e)
				}
			}
			if e := taskInfoService.SetNextRunTime(ctx, cronId, cron); e != nil {
				g.Log().Error(ctx, "更新下次执行时间失败:", e)
			}
		}, cronId)
	}
	err = taskInfoService.SetNextRunTime(ctx, cronId, cron)
	return
}

// DisableTask 禁用任务
func DisableTask(ctx g.Ctx, cronId string) (err error) {
	gcron.Remove(cronId)
	return
}
