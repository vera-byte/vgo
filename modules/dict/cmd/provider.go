package cmd

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// DictCommandProvider dict模块命令提供者
// 实现CommandProvider接口，提供dict模块的所有命令
type DictCommandProvider struct{}

// GetCommands 获取dict模块提供的命令列表
// 返回: 命令列表
func (p *DictCommandProvider) GetCommands() []*gcmd.Command {
	return []*gcmd.Command{
		{
			Name:  "dict",
			Brief: "字典管理命令",
			Arguments: []gcmd.Argument{
				{Name: "name", Short: "n", Brief: "名字", Default: "Dict"},
			},
			Func: func(ctx g.Ctx, parser *gcmd.Parser) error {
				name := parser.GetOpt("name").String()
				fmt.Printf("Hello from Dict module, %s 👋\n", name)
				return nil
			},
		},
		{
			Name:  "dict-list",
			Brief: "列出所有字典",
			Func: func(ctx g.Ctx, parser *gcmd.Parser) error {
				fmt.Println("Listing all dictionaries...")
				return nil
			},
		},
	}
}

// GetModuleName 获取模块名称
// 返回: 模块名称
func (p *DictCommandProvider) GetModuleName() string {
	return "dict"
}

// init 初始化dict模块命令提供者
// 自动注册到命令注册中心
func init() {
	provider := &DictCommandProvider{}
	registry := cmd.GetRegistry()
	registry.RegisterProvider(provider)
}