package dict

import (
	"embed"
	_ "github.com/vera-byte/vgo/modules/dict/cmd"
	_ "github.com/vera-byte/vgo/modules/dict/packed"

	_ "github.com/vera-byte/vgo/modules/dict/controller"
	"github.com/vera-byte/vgo/v"
)

// func init() {
// 	var (
// 		ctx = gctx.GetInitCtx()
// 	)
// 	g.Log().Debug(ctx, "module dict init start ...")
// 	// v.FillInitData(ctx, "dict", &model.DictInfo{})
// 	// v.FillInitData(ctx, "dict", &model.DictType{})
// 	g.Log().Debug(ctx, "module dict init finished ...")
// }

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("dict", migrationFS)
}
