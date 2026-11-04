// Package plugin 插件模块：插件包安装/启停/卸载，运行时执行插件脚本，
// 支持 Hook 插件接管内建行为（如 upload）。参照 cool-admin-midway plugin 模块。
package plugin

import (
	"embed"
	_ "github.com/vera-byte/vgo/modules/plugin/controller"
	_ "github.com/vera-byte/vgo/modules/plugin/service"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("plugin", migrationFS)
}
