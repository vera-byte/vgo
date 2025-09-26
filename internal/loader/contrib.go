package loader

import (
	_ "github.com/vera-byte/vgo/internal/packed"

	// 导入Redis适配器
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	_ "github.com/vera-byte/vgo/contrib/files/local"

	// Minio，按需启用
	// _ "github.com/vera-byte/vgo/contrib/files/minio"

	// 阿里云OSS，按需启用
	// _ "github.com/vera-byte/vgo/contrib/files/oss"

	// _ "github.com/vera-byte/vgo/contrib/drivers/sqlite"

	// _ "github.com/vera-byte/vgo/contrib/drivers/mysql"

	_ "github.com/vera-byte/vgo/contrib/drivers/pgsql"

	// 导入各个模块 - 只在server命令时加载
	_ "github.com/vera-byte/vgo/modules"
)
