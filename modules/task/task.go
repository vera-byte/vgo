package demo

import (
	_ "github.com/vera-byte/vgo/modules/task/packed"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/vera-byte/vgo/cool"
	_ "github.com/vera-byte/vgo/modules/task/controller"
	_ "github.com/vera-byte/vgo/modules/task/funcs"
	_ "github.com/vera-byte/vgo/modules/task/middleware"
	"github.com/vera-byte/vgo/modules/task/model"
)

func init() {
	var (
		taskInfo = model.NewTaskInfo()
		ctx      = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module task init start ...")
	cool.FillInitData(ctx, "task", taskInfo)

	result, err := cool.DBM(taskInfo).Where("status = ?", 1).All()
	if err != nil {
		panic(err)
	}
	for _, v := range result {
		id := v["id"].String()
		cool.RunFunc(ctx, "TaskAddTask("+id+")")
	}
	g.Log().Debug(ctx, "module task init finished ...")

}
