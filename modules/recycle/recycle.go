// Package recycle 数据回收模块：收集被删除的数据快照，支持查询与恢复。
// 参照 cool-admin-midway 的 recycle 模块实现。
package recycle

import (
	"embed"
	_ "github.com/vera-byte/vgo/modules/recycle/controller"

	_ "github.com/vera-byte/vgo/modules/recycle/service"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("recycle", migrationFS)
}
