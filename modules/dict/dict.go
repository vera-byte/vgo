package dict

import (
	"embed"

	_ "github.com/vera-byte/vgo/modules/dict/cmd"
	_ "github.com/vera-byte/vgo/modules/dict/packed"

	_ "github.com/vera-byte/vgo/modules/dict/controller"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("dict", migrationFS)
}
