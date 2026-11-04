package space

import (
	"embed"
	_ "github.com/vera-byte/vgo/modules/space/cmd"
	_ "github.com/vera-byte/vgo/modules/space/controller"
	_ "github.com/vera-byte/vgo/modules/space/middleware"
	"github.com/vera-byte/vgo/v"
)

// func init() {
// 	var (
// 		ctx = gctx.GetInitCtx()
// 	)
// 	g.Log().Debug(ctx, "module space init start ...")
// 	g.Log().Debug(ctx, "module space init finished ...")
// }

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("space", migrationFS)
}
