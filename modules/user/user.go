// Package user 用户模块：APP/小程序/公众号 C端用户体系
// （登录认证、资料管理、收货地址；参照 cool-admin-midway user 模块）
package user

import (
	"embed"
	_ "github.com/vera-byte/vgo/modules/user/controller"

	_ "github.com/vera-byte/vgo/modules/user/middleware"

	_ "github.com/vera-byte/vgo/modules/user/service"
	"github.com/vera-byte/vgo/v"
)

//go:embed resource/migrates
var migrationFS embed.FS

func init() {
	v.RegisterMigrationFS("user", migrationFS)
}
