package base

import (
	"embed"

	_ "github.com/vera-byte/vgo/modules/base/cmd"
	_ "github.com/vera-byte/vgo/modules/base/controller"
	_ "github.com/vera-byte/vgo/modules/base/funcs"
	_ "github.com/vera-byte/vgo/modules/base/middleware"
	_ "github.com/vera-byte/vgo/modules/base/packed"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("base", migrationFS)
}
