package main

import (
	_ "github.com/vera-byte/vgo/internal/packed"
	"github.com/vera-byte/vgo/v/cmd"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/vera-byte/vgo/contrib/files/local"

	// Minio，按需启用
	// _ "github.com/vera-byte/vgo/contrib/files/minio"

	// 阿里云OSS，按需启用
	// _ "github.com/vera-byte/vgo/contrib/files/oss"

	// _ "github.com/vera-byte/vgo/contrib/drivers/sqlite"

	_ "github.com/vera-byte/vgo/contrib/drivers/mysql"

	_ "github.com/vera-byte/vgo/contrib/drivers/pgsql"

	_ "github.com/vera-byte/vgo/modules"
)

func main() {
	cmd.Root.Run(gctx.New())
}
