package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// cSpace space模块主命令结构体
type cSpace struct {
	g.Meta `name:"space" brief:"空间模块命令" dc:"space命令用于管理空间模块的各种操作，包括空间管理、查看等功能"`
}

// cSpaceInput space命令的输入参数
type cSpaceInput struct {
	g.Meta `name:"space" brief:"空间模块命令" dc:"space命令用于管理空间模块的各种操作，包括空间管理、查看等功能"`
	Action string `short:"a" name:"action" brief:"操作类型" dc:"指定要执行的操作类型，如：info、status等" default:"info"`
}

// cSpaceOutput space命令的输出
type cSpaceOutput struct{}

// Index space命令的执行方法
// 功能：执行空间模块的相关操作
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cSpace) Index(ctx context.Context, in cSpaceInput) (out *cSpaceOutput, err error) {
	fmt.Printf("Space module action: %s\n", in.Action)
	return &cSpaceOutput{}, nil
}

// cSpaceList 空间列表命令结构体
type cSpaceList struct {
	g.Meta `name:"space-list" brief:"查看空间列表" dc:"查看系统中所有空间的列表信息"`
}

// cSpaceListInput space list命令的输入参数
type cSpaceListInput struct {
	g.Meta `name:"space-list" brief:"查看空间列表" dc:"查看系统中所有空间的列表信息"`
	Type   string `short:"t" name:"type" brief:"空间类型" dc:"指定要查看的空间类型，为空则显示所有"`
}

// cSpaceListOutput space list命令的输出
type cSpaceListOutput struct{}

// Index space list命令的执行方法
// 功能：查看空间列表
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cSpaceList) Index(ctx context.Context, in cSpaceListInput) (out *cSpaceListOutput, err error) {
	if in.Type != "" {
		fmt.Printf("Listing spaces with type: %s\n", in.Type)
	} else {
		fmt.Println("Listing all spaces...")
	}
	return &cSpaceListOutput{}, nil
}

// cSpaceCreate 创建空间命令结构体
type cSpaceCreate struct {
	g.Meta `name:"space-create" brief:"创建空间" dc:"创建一个新的空间"`
}

// cSpaceCreateInput space create命令的输入参数
type cSpaceCreateInput struct {
	g.Meta `name:"space-create" brief:"创建空间" dc:"创建一个新的空间"`
	Name   string `short:"n" name:"name" brief:"空间名称" dc:"指定要创建的空间名称" v:"required"`
	Type   string `short:"t" name:"type" brief:"空间类型" dc:"指定空间类型" default:"default"`
}

// cSpaceCreateOutput space create命令的输出
type cSpaceCreateOutput struct{}

// Index space create命令的执行方法
// 功能：创建新空间
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cSpaceCreate) Index(ctx context.Context, in cSpaceCreateInput) (out *cSpaceCreateOutput, err error) {
	fmt.Printf("Creating space '%s' with type: %s\n", in.Name, in.Type)
	return &cSpaceCreateOutput{}, nil
}

// SpaceCommandProvider space模块命令提供者
type SpaceCommandProvider struct{}

// GetCommands 获取space模块的命令列表
// 功能：返回space模块提供的所有命令
// 返回值：命令列表
func (p *SpaceCommandProvider) GetCommands() []*gcmd.Command {
	spaceCmd, _ := gcmd.NewFromObject(&cSpace{})
	spaceListCmd, _ := gcmd.NewFromObject(&cSpaceList{})
	spaceCreateCmd, _ := gcmd.NewFromObject(&cSpaceCreate{})
	return []*gcmd.Command{
		spaceCmd,
		spaceListCmd,
		spaceCreateCmd,
	}
}

// GetModuleName 获取模块名称
// 功能：返回模块名称
// 返回值：模块名称字符串
func (p *SpaceCommandProvider) GetModuleName() string {
	return "space"
}

func init() {
	// 注册space模块的命令提供者
	cmd.GetRegistry().RegisterProvider(&SpaceCommandProvider{})
}
