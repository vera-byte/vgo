package vconfig

import (
	"context"
	"fmt"
	"sync"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ConfigManager 配置管理器
// 支持多种配置源的统一管理和配置热更新
type ConfigManager struct {
	adapters map[string]VConfigAdapter // 配置适配器映射
	primary  string                    // 主配置源名称
	fallback []string                  // 备用配置源列表
	cache    map[string]*gvar.Var      // 配置缓存
	mutex    sync.RWMutex              // 读写锁
	watchers map[string][]func(*ConfigEvent) // 配置监听器
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		adapters: make(map[string]VConfigAdapter),
		cache:    make(map[string]*gvar.Var),
		watchers: make(map[string][]func(*ConfigEvent)),
	}
}

// RegisterAdapter 注册配置适配器
// name: 适配器名称
// adapter: 配置适配器实例
func (m *ConfigManager) RegisterAdapter(name string, adapter VConfigAdapter) error {
	if name == "" {
		return gerror.New("adapter name cannot be empty")
	}
	if adapter == nil {
		return gerror.New("adapter cannot be nil")
	}
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.adapters[name] = adapter
	return nil
}

// SetPrimary 设置主配置源
// name: 配置源名称
func (m *ConfigManager) SetPrimary(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.adapters[name]; !exists {
		return gerror.Newf("adapter '%s' not found", name)
	}
	
	m.primary = name
	return nil
}

// SetFallback 设置备用配置源列表
// names: 备用配置源名称列表
func (m *ConfigManager) SetFallback(names ...string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	for _, name := range names {
		if _, exists := m.adapters[name]; !exists {
			return gerror.Newf("adapter '%s' not found", name)
		}
	}
	
	m.fallback = names
	return nil
}

// Get 获取配置值
// ctx: 上下文
// pattern: 配置键模式
func (m *ConfigManager) Get(ctx context.Context, pattern string) (*gvar.Var, error) {
	// 先从缓存获取
	m.mutex.RLock()
	if cached, exists := m.cache[pattern]; exists {
		m.mutex.RUnlock()
		return cached, nil
	}
	m.mutex.RUnlock()
	
	// 从主配置源获取
	if m.primary != "" {
		if adapter, exists := m.adapters[m.primary]; exists && adapter.Available(ctx) {
			if value, err := adapter.Get(ctx, pattern); err == nil && value != nil {
				gvarValue := gvar.New(value)
				if !gvarValue.IsEmpty() {
					m.setCache(pattern, gvarValue)
					return gvarValue, nil
				}
			}
		}
	}
	
	// 从备用配置源获取
	for _, name := range m.fallback {
		if adapter, exists := m.adapters[name]; exists && adapter.Available(ctx) {
			if value, err := adapter.Get(ctx, pattern); err == nil && value != nil {
				gvarValue := gvar.New(value)
				if !gvarValue.IsEmpty() {
					m.setCache(pattern, gvarValue)
					return gvarValue, nil
				}
			}
		}
	}
	
	return g.NewVar(nil), gerror.Newf("config key '%s' not found", pattern)
}

// GetWithDefault 获取配置值，如果不存在则返回默认值
// ctx: 上下文
// pattern: 配置键模式
// defaultValue: 默认值
func (m *ConfigManager) GetWithDefault(ctx context.Context, pattern string, defaultValue *gvar.Var) *gvar.Var {
	value, err := m.Get(ctx, pattern)
	if err != nil || value.IsEmpty() || value.IsNil() {
		return defaultValue
	}
	return value
}

// Data 获取所有配置数据
// ctx: 上下文
func (m *ConfigManager) Data(ctx context.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	// 从主配置源获取
	if m.primary != "" {
		if adapter, exists := m.adapters[m.primary]; exists && adapter.Available(ctx) {
			if data, err := adapter.Data(ctx); err == nil {
				for k, v := range data {
					result[k] = v
				}
			}
		}
	}
	
	// 从备用配置源获取（不覆盖已存在的键）
	for _, name := range m.fallback {
		if adapter, exists := m.adapters[name]; exists && adapter.Available(ctx) {
			if data, err := adapter.Data(ctx); err == nil {
				for k, v := range data {
					if _, exists := result[k]; !exists {
						result[k] = v
					}
				}
			}
		}
	}
	
	return result, nil
}

// Watch 监听配置变化
// ctx: 上下文
// pattern: 配置键模式
// callback: 回调函数
func (m *ConfigManager) Watch(ctx context.Context, pattern string, callback func(*ConfigEvent)) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// 添加到监听器列表
	m.watchers[pattern] = append(m.watchers[pattern], callback)
	
	// 在所有适配器上设置监听
	for _, adapter := range m.adapters {
		if err := adapter.Watch(ctx, pattern, func(event *ConfigEvent) {
			// 清除缓存
			m.clearCache(event.Key)
			// 通知所有监听器
			m.notifyWatchers(pattern, event)
		}); err != nil {
			g.Log().Warningf(ctx, "failed to watch config on adapter %s: %v", adapter.Name(), err)
		}
	}
	
	return nil
}

// setCache 设置缓存
func (m *ConfigManager) setCache(key string, value *gvar.Var) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.cache[key] = value
}

// clearCache 清除缓存
func (m *ConfigManager) clearCache(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.cache, key)
}

// notifyWatchers 通知监听器
func (m *ConfigManager) notifyWatchers(pattern string, event *ConfigEvent) {
	m.mutex.RLock()
	watchers := m.watchers[pattern]
	m.mutex.RUnlock()
	
	for _, callback := range watchers {
		go func(cb func(*ConfigEvent)) {
			defer func() {
				if r := recover(); r != nil {
					g.Log().Errorf(context.Background(), "config watcher panic: %v", r)
				}
			}()
			cb(event)
		}(callback)
	}
}

// Close 关闭配置管理器
// ctx: 上下文
func (m *ConfigManager) Close(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	var errs []error
	for name, adapter := range m.adapters {
		if err := adapter.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close adapter %s: %w", name, err))
		}
	}
	
	// 清理资源
	m.adapters = make(map[string]VConfigAdapter)
	m.cache = make(map[string]*gvar.Var)
	m.watchers = make(map[string][]func(*ConfigEvent))
	
	if len(errs) > 0 {
		return gerror.Newf("errors occurred while closing adapters: %v", errs)
	}
	
	return nil
}

// 全局配置管理器实例
var defaultManager = NewConfigManager()

// RegisterAdapter 注册配置适配器到默认管理器
func RegisterAdapter(name string, adapter VConfigAdapter) error {
	return defaultManager.RegisterAdapter(name, adapter)
}

// SetPrimary 设置默认管理器的主配置源
func SetPrimary(name string) error {
	return defaultManager.SetPrimary(name)
}

// SetFallback 设置默认管理器的备用配置源
func SetFallback(names ...string) error {
	return defaultManager.SetFallback(names...)
}

// GetConfig 从默认管理器获取配置值
func GetConfig(ctx context.Context, pattern string) (*gvar.Var, error) {
	return defaultManager.Get(ctx, pattern)
}

// GetConfigWithDefault 从默认管理器获取配置值，带默认值
func GetConfigWithDefault(ctx context.Context, pattern string, defaultValue *gvar.Var) *gvar.Var {
	return defaultManager.GetWithDefault(ctx, pattern, defaultValue)
}

// WatchConfig 监听默认管理器的配置变化
func WatchConfig(ctx context.Context, pattern string, callback func(*ConfigEvent)) error {
	return defaultManager.Watch(ctx, pattern, callback)
}