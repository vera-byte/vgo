package demo

import (
	_ "github.com/vera-byte/vgo/modules/task/packed"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/vera-byte/vgo/modules/task/cmd"
	_ "github.com/vera-byte/vgo/modules/task/controller"
	_ "github.com/vera-byte/vgo/modules/task/funcs"
	_ "github.com/vera-byte/vgo/modules/task/middleware"
	"github.com/vera-byte/vgo/modules/task/model"
	"github.com/vera-byte/vgo/v"
)

func init() {
	var (
		taskInfo = model.NewTaskInfo()
		ctx      = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module task init start ...")
	// v.FillInitData(ctx, "task", taskInfo)

	result, err := v.DBM(taskInfo).Where("status = ?", 1).All()
	if err != nil {
		panic(err)
	}
	for _, item := range result {
		id := item["id"].String()
		v.RunFunc(ctx, "TaskAddTask("+id+")")
	}
	g.Log().Debug(ctx, "module task init finished ...")

}
