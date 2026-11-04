// Package demo 模块脚手架示例：`vgo-tools module <name>` 复制本目录生成新模块。
package demo

import (
	"embed"

	_ "vgo-simple/modules/demo/controller"
	_ "vgo-simple/modules/demo/service"

	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("demo", migrationFS)
}
