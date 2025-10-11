package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// cBase base模块主命令结构体
type cBase struct {
	g.Meta `name:"base" brief:"基础模块命令" dc:"base命令用于管理基础模块的各种操作，包括初始化、信息查看等功能"`
}

// cBaseInput base命令的输入参数
type cBaseInput struct {
	g.Meta `name:"base" brief:"基础模块命令" dc:"base命令用于管理基础模块的各种操作，包括初始化、信息查看等功能"`
	Action string `short:"a" name:"action" brief:"操作类型" dc:"指定要执行的操作类型，如：info、status等" default:"info"`
}

// cBaseOutput base命令的输出
type cBaseOutput struct{}

// Index base命令的执行方法
// 功能：执行基础模块的相关操作
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cBase) Index(ctx context.Context, in cBaseInput) (out *cBaseOutput, err error) {
	fmt.Printf("Base module action: %s\n", in.Action)
	return &cBaseOutput{}, nil
}

// cBaseInit 初始化基础模块的命令结构体
type cBaseInit struct {
	g.Meta `name:"base-init" brief:"初始化基础模块" dc:"初始化基础模块的数据库表、默认数据等"`
}

// cBaseInitInput base init命令的输入参数
type cBaseInitInput struct {
	g.Meta `name:"base-init" brief:"初始化基础模块" dc:"初始化基础模块的数据库表、默认数据等"`
	Force  bool `short:"f" name:"force" brief:"强制初始化" dc:"是否强制重新初始化，会覆盖现有数据"`
}

// cBaseInitOutput base init命令的输出
type cBaseInitOutput struct{}

// Index base init命令的执行方法
// 功能：初始化基础模块
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cBaseInit) Index(ctx context.Context, in cBaseInitInput) (out *cBaseInitOutput, err error) {
	if in.Force {
		fmt.Println("Force initializing base module...")
	} else {
		fmt.Println("Initializing base module...")
	}
	return &cBaseInitOutput{}, nil
}

// BaseCommandProvider base模块命令提供者
type BaseCommandProvider struct{}

// GetCommands 获取base模块的命令列表
// 功能：返回base模块提供的所有命令
// 返回值：命令列表
func (p *BaseCommandProvider) GetCommands() []*gcmd.Command {
	baseCmd, _ := gcmd.NewFromObject(&cBase{})
	baseInitCmd, _ := gcmd.NewFromObject(&cBaseInit{})
	return []*gcmd.Command{
		baseCmd,
		baseInitCmd,
	}
}

// GetModuleName 获取模块名称
// 功能：返回模块名称
// 返回值：模块名称字符串
func (p *BaseCommandProvider) GetModuleName() string {
	return "base"
}

func init() {
	// 注册base模块的命令提供者
	g.Log().Debugf(context.Background(), "Registering base module commands...")
	cmd.GetRegistry().RegisterProvider(&BaseCommandProvider{})
}
