package demo

import (
	_ "github.com/vera-byte/vgo/modules/task/packed"

	_ "github.com/vera-byte/vgo/modules/task/cmd"
	_ "github.com/vera-byte/vgo/modules/task/controller"
	_ "github.com/vera-byte/vgo/modules/task/funcs"
	_ "github.com/vera-byte/vgo/modules/task/middleware"
)

// func init() {
// 	var (
// 		taskInfo = model.NewTaskInfo()
// 		ctx      = gctx.GetInitCtx()
// 	)
// 	g.Log().Debug(ctx, "module task init start ...")
// 	v.FillInitData(ctx, "task", taskInfo)

// 	// 尝试查询，如果表不存在则跳过
// 	result, err := v.DBM(taskInfo).Where("status = ?", 1).All()
// 	if err != nil {
// 		g.Log().Warning(ctx, "task_info table may not exist, skipping task initialization:", err)
// 		g.Log().Debug(ctx, "module task init finished ...")
// 		return
// 	}
// 	for _, item := range result {
// 		id := item["id"].String()
// 		v.RunFunc(ctx, "TaskAddTask("+id+")")
// 	}
// 	g.Log().Debug(ctx, "module task init finished ...")

// }
