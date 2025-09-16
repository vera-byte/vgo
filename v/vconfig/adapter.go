package vconfig

import (
	"context"
	"github.com/gogf/gf/v2/os/gcfg"
)

// VConfigAdapter 扩展GoFrame官方gcfg.Adapter接口
// 在官方接口基础上增加配置源管理功能
type VConfigAdapter interface {
	gcfg.Adapter // 继承GoFrame官方接口
	
	// Name 返回适配器名称
	Name() string
	
	// Watch 监听配置变化（扩展功能）
	Watch(ctx context.Context, pattern string, callback func(event *ConfigEvent)) error
	
	// Close 关闭适配器（扩展功能）
	Close(ctx context.Context) error
}

// ConfigEvent 配置变更事件
type ConfigEvent struct {
	Key      string      // 配置键
	Value    interface{} // 新值
	OldValue interface{} // 旧值
	Type     EventType   // 事件类型
}

// EventType 事件类型
type EventType int

const (
	EventTypeAdd    EventType = iota // 添加
	EventTypeUpdate                  // 更新
	EventTypeDelete                  // 删除
)

// String 返回事件类型字符串
func (e EventType) String() string {
	switch e {
	case EventTypeAdd:
		return "ADD"
	case EventTypeUpdate:
		return "UPDATE"
	case EventTypeDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// AdapterWrapper 适配器包装器
// 将GoFrame官方适配器包装为VConfigAdapter
type AdapterWrapper struct {
	gcfg.Adapter
	name string
}

// NewAdapterWrapper 创建适配器包装器
// adapter: GoFrame官方适配器
// name: 适配器名称
func NewAdapterWrapper(adapter gcfg.Adapter, name string) VConfigAdapter {
	return &AdapterWrapper{
		Adapter: adapter,
		name:    name,
	}
}

// Name 返回适配器名称
func (w *AdapterWrapper) Name() string {
	return w.name
}

// Watch 监听配置变化（默认实现）
func (w *AdapterWrapper) Watch(ctx context.Context, pattern string, callback func(event *ConfigEvent)) error {
	// 默认不支持监听，子类可以重写
	return nil
}

// Close 关闭适配器（默认实现）
func (w *AdapterWrapper) Close(ctx context.Context) error {
	// 默认不需要关闭操作，子类可以重写
	return nil
}