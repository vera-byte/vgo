package cmd

import (
	"github.com/gogf/gf/v2/os/gcmd"
)

// Root 根命令定义
// 该变量定义了vgo命令行工具的根命令，所有子命令都会注册到这个根命令下
var Root = &gcmd.Command{
	Name:  "vgo",
	Brief: "vgo 命令行工具",
}

// Main 主命令，包含所有子命令
// 在程序启动时会自动注册所有模块的命令
var Main = &gcmd.Command{
	Name:  "vgo",
	Brief: "vgo 命令行工具",
}

// init 初始化命令系统
// 设置Main命令并启动命令注册流程
func init() {
	// 将Root设置为Main的别名，保持向后兼容
	Main = Root
	
	// 启动命令注册流程
	initializeCommands()
}

// initializeCommands 初始化命令注册
// 自动发现并注册所有模块的命令
func initializeCommands() {
	registry := GetRegistry()
	
	// 自动发现模块（可选，主要用于调试）
	registry.AutoDiscoverModules("modules")
	
	// 注册所有已发现的命令
	registry.RegisterCommands()
}
