module github.com/vera-byte/vgo

go 1.27.3

require (
	github.com/gogf/gf/contrib/nosql/redis/v2 v2.10.4
	github.com/gogf/gf/v2 v2.10.4
	github.com/golang-migrate/migrate/v4 v4.19.0
	github.com/lib/pq v1.10.9
	github.com/vera-byte/vgo/contrib/drivers/pgsql v1.1.0
	github.com/vera-byte/vgo/contrib/files/local v1.1.0
	github.com/vera-byte/vgo/modules/base v1.1.0
	github.com/vera-byte/vgo/modules/dict v1.1.0
	github.com/vera-byte/vgo/modules/openapi v1.1.0
	github.com/vera-byte/vgo/modules/plugin v1.1.0
	github.com/vera-byte/vgo/pluginsdk v1.1.0
	github.com/vera-byte/vgo/modules/recycle v1.1.0
	github.com/vera-byte/vgo/modules/space v1.1.0
	github.com/vera-byte/vgo/modules/task v1.1.0
	github.com/vera-byte/vgo/v v1.1.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/dlclark/regexp2/v2 v2.11.0 // indirect
	github.com/dop251/goja v0.0.0-20261102140139-622d10077b6c // indirect
	github.com/emirpasic/gods/v2 v2.0.0-alpha // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/gogf/gf/contrib/drivers/pgsql/v2 v2.10.4 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grokify/html-strip-tags-go v0.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mattn/go-runewidth v0.0.23 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.3.0 // indirect
	github.com/olekukonko/ll v0.1.8 // indirect
	github.com/olekukonko/tablewriter v1.1.4 // indirect
	github.com/panjf2000/ants/v2 v2.12.0 // indirect
	github.com/redis/go-redis/v9 v9.19.0 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/sa-tokens/sa-token-go/core v0.2.3 // indirect
	github.com/sa-tokens/sa-token-go/storage/memory v0.2.3 // indirect
	github.com/sa-tokens/sa-token-go/storage/redis v0.2.3 // indirect
	github.com/sa-tokens/sa-token-go/stputil v0.2.3 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
	gorm.io/gorm v1.31.0 // indirect
)

// 本地开发统一指向工作区源码（发布走 tag；详见 README「模块与依赖」）
replace (
	github.com/vera-byte/vgo/contrib/drivers/pgsql => ./contrib/drivers/pgsql
	github.com/vera-byte/vgo/contrib/files/local => ./contrib/files/local
	github.com/vera-byte/vgo/modules/base => ./modules/base
	github.com/vera-byte/vgo/modules/dict => ./modules/dict
	github.com/vera-byte/vgo/modules/openapi => ./modules/openapi
	github.com/vera-byte/vgo/modules/plugin => ./modules/plugin
	github.com/vera-byte/vgo/pluginsdk => ./pluginsdk
	github.com/vera-byte/vgo/modules/recycle => ./modules/recycle
	github.com/vera-byte/vgo/modules/space => ./modules/space
	github.com/vera-byte/vgo/modules/task => ./modules/task
	github.com/vera-byte/vgo/v => ./v
)
