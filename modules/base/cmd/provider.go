package cmd

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// BaseCommandProvider base模块命令提供者
// 实现CommandProvider接口，提供base模块的所有命令
type BaseCommandProvider struct{}

// GetCommands 获取base模块提供的命令列表
// 返回: 命令列表
func (p *BaseCommandProvider) GetCommands() []*gcmd.Command {
	return []*gcmd.Command{
		{
			Name:  "base",
			Brief: "基础模块命令",
			Arguments: []gcmd.Argument{
				{Name: "action", Short: "a", Brief: "操作类型", Default: "info"},
			},
			Func: func(ctx g.Ctx, parser *gcmd.Parser) error {
				action := parser.GetOpt("action").String()
				fmt.Printf("Base module action: %s\n", action)
				return nil
			},
		},
		{
			Name:  "base-init",
			Brief: "初始化基础模块",
			Func: func(ctx g.Ctx, parser *gcmd.Parser) error {
				fmt.Println("Initializing base module...")
				return nil
			},
		},
	}
}

// GetModuleName 获取模块名称
// 返回: 模块名称
func (p *BaseCommandProvider) GetModuleName() string {
	return "base"
}

// init 初始化base模块命令提供者
// 自动注册到命令注册中心
func init() {
	provider := &BaseCommandProvider{}
	registry := cmd.GetRegistry()
	registry.RegisterProvider(provider)
}