package openapi

import (
	"embed"

	_ "github.com/vera-byte/vgo/modules/openapi/cmd"
	_ "github.com/vera-byte/vgo/modules/openapi/controller/admin"
	_ "github.com/vera-byte/vgo/modules/openapi/controller/open"
	_ "github.com/vera-byte/vgo/modules/openapi/middleware"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("openapi", migrationFS)
}
