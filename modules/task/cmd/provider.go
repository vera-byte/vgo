package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// cTask task模块主命令结构体
type cTask struct {
	g.Meta `name:"task" brief:"任务模块命令" dc:"task命令用于管理任务模块的各种操作，包括任务管理、查看等功能"`
}

// cTaskInput task命令的输入参数
type cTaskInput struct {
	g.Meta `name:"task" brief:"任务模块命令" dc:"task命令用于管理任务模块的各种操作，包括任务管理、查看等功能"`
	Action string `short:"a" name:"action" brief:"操作类型" dc:"指定要执行的操作类型，如：info、status等" default:"info"`
}

// cTaskOutput task命令的输出
type cTaskOutput struct{}

// Index task命令的执行方法
// 功能：执行任务模块的相关操作
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cTask) Index(ctx context.Context, in cTaskInput) (out *cTaskOutput, err error) {
	fmt.Printf("Task module action: %s\n", in.Action)
	return &cTaskOutput{}, nil
}

// cTaskList 任务列表命令结构体
type cTaskList struct {
	g.Meta `name:"task-list" brief:"查看任务列表" dc:"查看系统中所有任务的列表信息"`
}

// cTaskListInput task list命令的输入参数
type cTaskListInput struct {
	g.Meta `name:"task-list" brief:"查看任务列表" dc:"查看系统中所有任务的列表信息"`
	Status string `short:"s" name:"status" brief:"任务状态" dc:"指定要查看的任务状态，为空则显示所有"`
}

// cTaskListOutput task list命令的输出
type cTaskListOutput struct{}

// Index task list命令的执行方法
// 功能：查看任务列表
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cTaskList) Index(ctx context.Context, in cTaskListInput) (out *cTaskListOutput, err error) {
	if in.Status != "" {
		fmt.Printf("Listing tasks with status: %s\n", in.Status)
	} else {
		fmt.Println("Listing all tasks...")
	}
	return &cTaskListOutput{}, nil
}

// cTaskStart 启动任务命令结构体
type cTaskStart struct {
	g.Meta `name:"task-start" brief:"启动任务" dc:"启动指定的任务"`
}

// cTaskStartInput task start命令的输入参数
type cTaskStartInput struct {
	g.Meta `name:"task-start" brief:"启动任务" dc:"启动指定的任务"`
	ID     int64 `short:"i" name:"id" brief:"任务ID" dc:"指定要启动的任务ID" v:"required"`
}

// cTaskStartOutput task start命令的输出
type cTaskStartOutput struct{}

// Index task start命令的执行方法
// 功能：启动指定任务
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cTaskStart) Index(ctx context.Context, in cTaskStartInput) (out *cTaskStartOutput, err error) {
	fmt.Printf("Starting task with ID: %d\n", in.ID)
	return &cTaskStartOutput{}, nil
}

// TaskCommandProvider task模块命令提供者
type TaskCommandProvider struct{}

// GetCommands 获取task模块的命令列表
// 功能：返回task模块提供的所有命令
// 返回值：命令列表
func (p *TaskCommandProvider) GetCommands() []*gcmd.Command {
	taskCmd, _ := gcmd.NewFromObject(&cTask{})
	taskListCmd, _ := gcmd.NewFromObject(&cTaskList{})
	taskStartCmd, _ := gcmd.NewFromObject(&cTaskStart{})
	return []*gcmd.Command{
		taskCmd,
		taskListCmd,
		taskStartCmd,
	}
}

// GetModuleName 获取模块名称
// 功能：返回模块名称
// 返回值：模块名称字符串
func (p *TaskCommandProvider) GetModuleName() string {
	return "task"
}

func init() {
	// 注册task模块的命令提供者
	cmd.GetRegistry().RegisterProvider(&TaskCommandProvider{})
}
