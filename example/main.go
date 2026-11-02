package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/vera-byte/vgo/v/pluginsdk"
)

// DemoPlugin 插件演示：嵌入 pluginsdk.Base 获得 pluginInfo / cache / pluginService
// 对齐 cool-admin BasePlugin 契约（ready / pluginInfo / cache / pluginService）
type DemoPlugin struct {
	pluginsdk.Base
}

// Ready 插件就绪回调（可选实现；singleton 插件仅启动时执行一次）
func (p *DemoPlugin) Ready() {
	fmt.Fprintln(os.Stderr, "[demo-plugin] ready, key =", p.PluginInfo.Key)
}

// Show 示例业务方法：导出方法即插件 API（宿主按方法名调用，大小写不敏感）
// 返回值支持 (T, error) / (error) / (T) / ()
func (p *DemoPlugin) Show(a, b int) (interface{}, error) {
	return map[string]any{
		"plugin": p.PluginInfo.Name,
		"key":    p.PluginInfo.Key,
		"sum":    a + b,
	}, nil
}

// Info 返回插件信息与配置（@baseDir/ 占位已解析为插件安装目录）
func (p *DemoPlugin) Info() (interface{}, error) {
	return map[string]any{
		"name":    p.PluginInfo.Name,
		"version": p.PluginInfo.Version,
		"baseDir": p.PluginInfo.BaseDir,
		"config":  p.PluginInfo.Config,
	}, nil
}

// Config 读取指定配置项
func (p *DemoPlugin) Config(key string) (interface{}, error) {
	return p.PluginInfo.Config[key], nil
}

// UseCache 读写宿主缓存（跨插件/持久共享，落宿主 Redis）
func (p *DemoPlugin) UseCache(key, value string) (interface{}, error) {
	if err := p.Cache.Set(key, value); err != nil {
		return nil, err
	}
	return p.Cache.Get(key)
}

// CallOther 调用其他插件的方法（跨插件协作；目标插件需已安装并启用）
func (p *DemoPlugin) CallOther(key, method string) (interface{}, error) {
	return p.PluginService.Invoke(key, method, 1, 2)
}

// Hash 纯函数示例：SHA256 摘要
func (p *DemoPlugin) Hash(text string) (string, error) {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:]), nil
}

// Greet 按配置打招呼（演示 @baseDir 路径解析与配置读取）
func (p *DemoPlugin) Greet(name string) (string, error) {
	greeting := strings.TrimSpace(p.PluginInfo.ConfigString("greeting"))
	if greeting == "" {
		greeting = "hello"
	}
	return fmt.Sprintf("%s, %s", greeting, name), nil
}

// Upload upload Hook 方法示例：
// plugin.json 中声明 "hook": "upload" 后，文件上传接口（/admin/base/comm/upload）
// 将由本插件接管；ctx 形如 {"files":[{"data":"临时文件路径","filename":"原始名"}],"fields":{...}}
// 当前 hook 为空（不接管），保留方法供演示/启用。
//func (p *DemoPlugin) Upload(ctx map[string]any) (interface{}, error) {
//	files, _ := ctx["files"].([]any)
//	if len(files) > 0 {
//		if f, ok := files[0].(map[string]any); ok {
//			return "https://cdn.example.com/" + fmt.Sprint(f["filename"]), nil
//		}
//	}
//	return "", nil
//}

func main() {
	// 启动 stdio JSON-RPC 服务循环（宿主 spawn 后调用 ready / invoke）
	pluginsdk.Run(&DemoPlugin{})
}
