package funcs

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseFuncClearLog struct {
}

// Func
func (f *BaseFuncClearLog) Func(ctx g.Ctx, param string) (err error) {
	g.Log().Info(ctx, "清理日志 BaseFuncClearLog.Func", "param", param)
	baseSysLogService := service.NewBaseSysLogService()
	if param == "true" {
		err = baseSysLogService.Clear(true)
	} else {
		err = baseSysLogService.Clear(false)
	}
	return
}

// IsSingleton
func (f *BaseFuncClearLog) IsSingleton() bool {
	return true
}

// IsAllWorker
func (f *BaseFuncClearLog) IsAllWorker() bool {
	return false
}

// init
func init() {
	v.RegisterFunc("BaseFuncClearLog", &BaseFuncClearLog{})
}
