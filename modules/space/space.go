package space

import (
	"embed"

	_ "github.com/vera-byte/vgo/modules/space/cmd"
	_ "github.com/vera-byte/vgo/modules/space/controller"
	_ "github.com/vera-byte/vgo/modules/space/middleware"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("space", migrationFS)
}
