package v

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// PluginInvoker 插件调用实现，由 plugin 模块注册（见 RegisterPluginInvoker）
type PluginInvoker interface {
	// Invoke 调用指定插件方法（key 插件标识，method 方法名）
	Invoke(ctx context.Context, key string, method string, params ...interface{}) (interface{}, error)
}

var pluginInvoker PluginInvoker

// RegisterPluginInvoker 注册插件调用实现（plugin 模块初始化时调用）
func RegisterPluginInvoker(p PluginInvoker) {
	pluginInvoker = p
}

type pluginFacade struct{}

// Plugin 插件调用入口
var Plugin pluginFacade

// Invoke 调用指定插件方法；未注册实现（未引入 plugin 模块）时返回错误
func (pluginFacade) Invoke(ctx context.Context, key string, method string, params ...interface{}) (interface{}, error) {
	if pluginInvoker == nil {
		return nil, gerror.New("插件调用未注册（未引入 plugin 模块）")
	}
	return pluginInvoker.Invoke(ctx, key, method, params...)
}
