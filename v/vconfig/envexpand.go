package vconfig

import (
	"os"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
)

// expandVar 展开配置值中的环境变量占位符；非字符串值原样返回。
func expandVar(value *gvar.Var) *gvar.Var {
	if s, ok := value.Val().(string); ok {
		return gvar.New(ExpandEnv(s))
	}
	return value
}

// envPlaceholderPattern 匹配 ${VAR} 形式的环境变量占位符。
var envPlaceholderPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// ExpandEnv 将字符串中的 ${VAR} 占位符替换为对应环境变量的值。
// 未设置的变量替换为空字符串。
func ExpandEnv(s string) string {
	if !strings.Contains(s, "${") {
		return s
	}
	return envPlaceholderPattern.ReplaceAllStringFunc(s, func(match string) string {
		return os.Getenv(envPlaceholderPattern.FindStringSubmatch(match)[1])
	})
}

// loadDotEnv 加载工作目录下的 .env 文件，为未显式设置的变量补位。
// 这样 `go run . server` / IDE 直接启动也能读到密钥，无需手动 export。
func loadDotEnv() {
	content := gfile.GetContents(".env")
	if content == "" {
		return
	}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

// init 在其他包读取配置前展开全局配置文件中的环境变量占位符。
// gdb/gredis 等组件直接经 gcfg 消费 database/redis 配置节，不经过 GetCfgWithDefault，
// 因此需要在全局配置适配器层面统一完成展开。
func init() {
	loadDotEnv()
	content := gfile.GetContents("manifest/config/config.yaml")
	if content == "" {
		return
	}
	adapter, err := gcfg.NewAdapterContent(ExpandEnv(content))
	if err != nil {
		return
	}
	g.Cfg().SetAdapter(adapter)
}
