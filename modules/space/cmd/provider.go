package cmd

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/internal/cmd"
)

// SpaceCommandProvider dict模块命令提供者
// 实现CommandProvider接口，提供dict模块的所有命令
type SpaceCommandProvider struct{}

// GetCommands 获取dict模块提供的命令列表
// 返回: 命令列表
func (p *SpaceCommandProvider) GetCommands() []*gcmd.Command {
	return []*gcmd.Command{}
}

// GetModuleName 获取模块名称
// 返回: 模块名称
func (p *SpaceCommandProvider) GetModuleName() string {
	return "space"
}

// init 初始化dict模块命令提供者
// 自动注册到命令注册中心
func init() {
	provider := &SpaceCommandProvider{}
	registry := cmd.GetRegistry()
	registry.RegisterProvider(provider)
}
