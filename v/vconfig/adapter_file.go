package vconfig

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
)

// FileAdapterConfig 文件适配器配置
type FileAdapterConfig struct {
	Path     string        // 配置文件路径
	FileName string        // 配置文件名
	Watch    bool          // 是否监听文件变化
	Interval time.Duration // 监听间隔
}

// FileAdapter 文件配置适配器
// 基于GoFrame的gcfg组件实现，支持多种文件格式
type FileAdapter struct {
	*AdapterWrapper
	config   *FileAdapterConfig
	gcfg     *gcfg.Config
	watchers map[string][]func(*ConfigEvent)
	mutex    sync.RWMutex
	lastData map[string]interface{}
}

// NewFileAdapter 创建文件配置适配器
// config: 适配器配置
func NewFileAdapter(config *FileAdapterConfig) (*FileAdapter, error) {
	if config == nil {
		return nil, gerror.New("config cannot be nil")
	}
	
	// 设置默认值
	if config.FileName == "" {
		config.FileName = "config.yaml"
	}
	if config.Interval == 0 {
		config.Interval = time.Second * 3
	}
	
	// 创建gcfg实例
	cfg, err := gcfg.New()
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to create gcfg instance")
	}
	
	adapter := &FileAdapter{
		AdapterWrapper: NewAdapterWrapper(cfg.GetAdapter(), "file").(*AdapterWrapper),
		config:         config,
		gcfg:           cfg,
		watchers:       make(map[string][]func(*ConfigEvent)),
		lastData:       make(map[string]interface{}),
	}
	
	// 设置文件路径和名称
	if config.Path != "" {
		cfg.GetAdapter().(*gcfg.AdapterFile).SetPath(config.Path)
	}
	cfg.GetAdapter().(*gcfg.AdapterFile).SetFileName(config.FileName)
	
	// 启动文件监听
	if config.Watch {
		if err := adapter.startWatch(); err != nil {
			g.Log().Warningf(context.Background(), "failed to start file watch: %v", err)
		}
	}
	
	return adapter, nil
}

// Name 返回适配器名称
func (f *FileAdapter) Name() string {
	return "file"
}

// Available 检查配置源是否可用
// ctx: 上下文
// files: 可选的文件名参数（兼容GoFrame接口）
func (f *FileAdapter) Available(ctx context.Context, files ...string) bool {
	return f.gcfg.Available(ctx, files...)
}

// Get 获取指定键的配置值
// ctx: 上下文
// pattern: 配置键模式
func (f *FileAdapter) Get(ctx context.Context, pattern string) (any, error) {
	v, err := f.gcfg.Get(ctx, pattern)
	if err != nil {
		return nil, err
	}
	return v.Val(), nil
}

// Data 获取所有配置数据
// ctx: 上下文
func (f *FileAdapter) Data(ctx context.Context) (map[string]interface{}, error) {
	return f.gcfg.Data(ctx)
}

// Set 设置配置值（文件适配器不支持动态设置）
// ctx: 上下文
// pattern: 配置键模式
// value: 配置值
func (f *FileAdapter) Set(ctx context.Context, pattern string, value interface{}) error {
	return gerror.New("file adapter does not support dynamic configuration setting")
}

// Watch 监听配置变化
// ctx: 上下文
// pattern: 配置键模式
// callback: 回调函数
func (f *FileAdapter) Watch(ctx context.Context, pattern string, callback func(*ConfigEvent)) error {
	if !f.config.Watch {
		return gerror.New("file watching is disabled")
	}
	
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	f.watchers[pattern] = append(f.watchers[pattern], callback)
	return nil
}

// Close 关闭适配器
// ctx: 上下文
func (f *FileAdapter) Close(ctx context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	// 清理监听器
	f.watchers = make(map[string][]func(*ConfigEvent))
	f.lastData = make(map[string]interface{})
	
	return nil
}

// startWatch 启动文件监听
func (f *FileAdapter) startWatch() error {
	// 获取配置文件完整路径
	var filePath string
	if f.config.Path != "" {
		filePath = filepath.Join(f.config.Path, f.config.FileName)
	} else {
		// 使用gcfg的默认搜索路径
		searchPaths := []string{
			".",
			"config",
			"manifest/config",
		}
		
		for _, path := range searchPaths {
			testPath := filepath.Join(path, f.config.FileName)
			if gfile.Exists(testPath) {
				filePath = testPath
				break
			}
		}
	}
	
	if filePath == "" || !gfile.Exists(filePath) {
		return gerror.Newf("config file not found: %s", f.config.FileName)
	}
	
	// 监听文件变化
	_, err := gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
		if event.IsWrite() || event.IsRename() {
			f.handleFileChange()
		}
	})
	
	return err
}

// handleFileChange 处理文件变化
func (f *FileAdapter) handleFileChange() {
	ctx := context.Background()
	
	// 等待一小段时间，确保文件写入完成
	time.Sleep(100 * time.Millisecond)
	
	// 获取新的配置数据
	newData, err := f.gcfg.Data(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "failed to load config data after file change: %v", err)
		return
	}
	
	f.mutex.Lock()
	oldData := f.lastData
	f.lastData = make(map[string]interface{})
	for k, v := range newData {
		f.lastData[k] = v
	}
	f.mutex.Unlock()
	
	// 比较配置变化并通知监听器
	f.compareAndNotify(oldData, newData)
}

// compareAndNotify 比较配置变化并通知监听器
func (f *FileAdapter) compareAndNotify(oldData, newData map[string]interface{}) {
	f.mutex.RLock()
	watchers := make(map[string][]func(*ConfigEvent))
	for k, v := range f.watchers {
		watchers[k] = make([]func(*ConfigEvent), len(v))
		copy(watchers[k], v)
	}
	f.mutex.RUnlock()
	
	// 检查所有配置键的变化
	allKeys := make(map[string]bool)
	for k := range oldData {
		allKeys[k] = true
	}
	for k := range newData {
		allKeys[k] = true
	}
	
	for key := range allKeys {
		oldValue, oldExists := oldData[key]
		newValue, newExists := newData[key]
		
		var eventType EventType
		if !oldExists && newExists {
			eventType = EventTypeAdd
		} else if oldExists && !newExists {
			eventType = EventTypeDelete
		} else if oldExists && newExists && !f.isEqual(oldValue, newValue) {
			eventType = EventTypeUpdate
		} else {
			continue // 没有变化
		}
		
		event := &ConfigEvent{
			Key:      key,
			Value:    newValue,
			OldValue: oldValue,
			Type:     eventType,
		}
		
		// 通知匹配的监听器
		for pattern, callbacks := range watchers {
			if f.matchPattern(key, pattern) {
				for _, callback := range callbacks {
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
		}
	}
}

// matchPattern 检查键是否匹配模式
func (f *FileAdapter) matchPattern(key, pattern string) bool {
	// 简单的模式匹配，支持通配符 *
	if pattern == "*" {
		return true
	}
	
	if strings.Contains(pattern, "*") {
		// 简单的通配符匹配
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			prefix, suffix := parts[0], parts[1]
			return strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix)
		}
	}
	
	return key == pattern || strings.HasPrefix(key, pattern+".")
}

// isEqual 比较两个值是否相等
func (f *FileAdapter) isEqual(a, b interface{}) bool {
	return g.NewVar(a).String() == g.NewVar(b).String()
}