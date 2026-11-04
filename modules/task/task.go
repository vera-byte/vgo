package task

import (
	"context"
	"embed"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/modules/task/model"
	_ "github.com/vera-byte/vgo/modules/task/packed"

	_ "github.com/vera-byte/vgo/modules/task/cmd"
	_ "github.com/vera-byte/vgo/modules/task/controller"
	_ "github.com/vera-byte/vgo/modules/task/funcs"
	_ "github.com/vera-byte/vgo/modules/task/middleware"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

var taskInfo = model.NewTaskInfo()

func init() {
	v.RegisterMigrationFS("task", migrationFS)
	// 任务调度启动（查 task_info 拉起活动任务）延后到 server 启动，避免任意子命令在 init 阶段触库
	v.RegisterStartupHook("task-scheduler", loadActiveTasks)
}

// loadActiveTasks 加载并启动全部启用中的任务（server 启动钩子）
func loadActiveTasks(ctx context.Context) {
	result, err := v.DBM(taskInfo).Where("status = ?", 1).All()
	if err != nil {
		g.Log().Warning(ctx, "task_info table may not exist, skipping task initialization:", err)
		return
	}
	for _, item := range result {
		id := item["id"].String()
		v.RunFunc(ctx, "TaskAddTask("+id+")")
	}
}
