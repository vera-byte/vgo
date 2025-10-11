package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// cDict dict模块主命令结构体
type cDict struct {
	g.Meta `name:"dict" brief:"字典模块命令" dc:"dict命令用于管理字典模块的各种操作，包括字典管理、查看等功能"`
}

// cDictInput dict命令的输入参数
type cDictInput struct {
	g.Meta `name:"dict" brief:"字典模块命令" dc:"dict命令用于管理字典模块的各种操作，包括字典管理、查看等功能"`
	Action string `short:"a" name:"action" brief:"操作类型" dc:"指定要执行的操作类型，如：info、status等" default:"info"`
}

// cDictOutput dict命令的输出
type cDictOutput struct{}

// Index dict命令的执行方法
// 功能：执行字典模块的相关操作
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cDict) Index(ctx context.Context, in cDictInput) (out *cDictOutput, err error) {
	fmt.Printf("Dict module action: %s\n", in.Action)
	return &cDictOutput{}, nil
}

// cDictList 字典列表命令结构体
type cDictList struct {
	g.Meta `name:"dict-list" brief:"查看字典列表" dc:"查看系统中所有字典类型和字典项的列表"`
}

// cDictListInput dict list命令的输入参数
type cDictListInput struct {
	g.Meta `name:"dict-list" brief:"查看字典列表" dc:"查看系统中所有字典类型和字典项的列表"`
	Type   string `short:"t" name:"type" brief:"字典类型" dc:"指定要查看的字典类型，为空则显示所有"`
}

// cDictListOutput dict list命令的输出
type cDictListOutput struct{}

// Index dict list命令的执行方法
// 功能：查看字典列表
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cDictList) Index(ctx context.Context, in cDictListInput) (out *cDictListOutput, err error) {
	if in.Type != "" {
		fmt.Printf("Listing dict items for type: %s\n", in.Type)
	} else {
		fmt.Println("Listing all dict types and items...")
	}
	return &cDictListOutput{}, nil
}

// DictCommandProvider dict模块命令提供者
type DictCommandProvider struct{}

// GetCommands 获取dict模块的命令列表
// 功能：返回dict模块提供的所有命令
// 返回值：命令列表
func (p *DictCommandProvider) GetCommands() []*gcmd.Command {
	dictCmd, _ := gcmd.NewFromObject(&cDict{})
	dictListCmd, _ := gcmd.NewFromObject(&cDictList{})
	return []*gcmd.Command{
		dictCmd,
		dictListCmd,
	}
}

// GetModuleName 获取模块名称
// 功能：返回模块名称
// 返回值：模块名称字符串
func (p *DictCommandProvider) GetModuleName() string {
	return "dict"
}

func init() {
	// 注册dict模块的命令提供者
	cmd.GetRegistry().RegisterProvider(&DictCommandProvider{})
}