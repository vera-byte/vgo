// Package pluginsdk 是 vgo 插件开发 SDK：插件编译为独立可执行程序，
// 通过 stdio 上的 JSON-RPC 与宿主（vgo modules/plugin）通信。
// 契约对齐 cool-admin-midway-plugin 的 BasePlugin（pluginInfo/cache/pluginService/ready）。
//
// 插件作者用法：
//
//	type MyPlugin struct{ pluginsdk.Base }
//	func (p *MyPlugin) Show(a, b int) (interface{}, error) { return p.PluginInfo, nil }
//	func main() { pluginsdk.Run(&MyPlugin{}) }
package pluginsdk

import (
	"encoding/json"
	"fmt"
)

// PluginInfo 插件信息（对应 plugin.json + 安装元数据）
type PluginInfo struct {
	Name        string                 `json:"name"`
	Key         string                 `json:"key"`
	Hook        string                 `json:"hook"`
	Singleton   bool                   `json:"singleton"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Logo        string                 `json:"logo"`
	Readme      string                 `json:"readme"`
	BaseDir     string                 `json:"baseDir"` // 插件安装目录（@baseDir 解析基准）
	Config      map[string]interface{} `json:"config"`
}

// ConfigString 读取字符串配置（支持 @baseDir 相对路径解析为绝对路径）
func (i *PluginInfo) ConfigString(key string) string {
	if i == nil || i.Config == nil {
		return ""
	}
	return fmt.Sprintf("%v", i.Config[key])
}

// Cache 宿主缓存（实现由 SDK 通过反向 RPC 调宿主，开发环境也可本地实现）
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

// PluginService 跨插件调用（反向 RPC 到宿主，由宿主路由到其他插件）
type PluginService interface {
	// Invoke 调用其他插件：key 插件标识，method 方法名
	Invoke(key, method string, params ...interface{}) (interface{}, error)
}

// Base 插件基类：嵌入后自动获得 pluginInfo/cache/pluginService（对应 BasePlugin）
type Base struct {
	PluginInfo    *PluginInfo
	Cache         Cache
	PluginService PluginService
}

// Ready 插件就绪回调（可选实现；单例插件仅启动时执行一次）
// 使用方可在自己的插件结构体上定义 Ready() 覆写。

// readyer 可选就绪接口
type readyer interface {
	Ready()
}

var _ = json.Marshal
