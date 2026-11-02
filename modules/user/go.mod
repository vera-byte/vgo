module github.com/vera-byte/vgo/modules/user

go 1.27.3

require (
	github.com/gogf/gf/v2 v2.10.4
	github.com/redis/go-redis/v9 v9.19.0
	github.com/sa-tokens/sa-token-go/core v0.2.3
	github.com/sa-tokens/sa-token-go/integrations/gf v0.2.3
	github.com/sa-tokens/sa-token-go/storage/memory v0.2.3
	github.com/sa-tokens/sa-token-go/storage/redis v0.2.3
	github.com/vera-byte/vgo/modules/base v1.2.0
	github.com/vera-byte/vgo/modules/plugin v1.2.0
	github.com/vera-byte/vgo/v v1.2.0
	golang.org/x/crypto v0.58.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/emirpasic/gods/v2 v2.0.0-alpha // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grokify/html-strip-tags-go v0.1.0 // indirect
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
	github.com/sa-tokens/sa-token-go/stputil v0.2.3 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/gorm v1.31.0 // indirect
)

// 本地开发统一指向工作区源码（发布走 tag）
replace (
	github.com/vera-byte/vgo/modules/base => ../base
	github.com/vera-byte/vgo/modules/plugin => ../plugin
	github.com/vera-byte/vgo/v => ../../v
)
