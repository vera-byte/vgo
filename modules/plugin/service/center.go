package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
)

// PluginInstance 插件实例抽象（可执行进程 / 内置实现统一派发）
type PluginInstance interface {
	// Invoke 调用插件方法
	Invoke(method string, args ...interface{}) (interface{}, error)
}

// ErrNoHookPlugin 未注册对应 Hook 插件（调用方回退默认实现）
var ErrNoHookPlugin = gerror.New("未注册对应的Hook插件")

// PluginInfoData 插件信息（对应 pluginsdk.PluginInfo / 参考实现 PluginInfo）
type PluginInfoData struct {
	Name        string                 `json:"name"`
	Key         string                 `json:"key"`
	Hook        string                 `json:"hook"`
	Singleton   bool                   `json:"singleton"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Logo        string                 `json:"logo"`
	Readme      string                 `json:"readme"`
	BaseDir     string                 `json:"baseDir"`
	Config      map[string]interface{} `json:"config"`
}

// PluginCenter 插件中心（单例注册表，对应参考实现 PluginCenterService）
type PluginCenter struct {
	mu          sync.RWMutex
	plugins     map[string]PluginInstance
	pluginInfos map[string]*PluginInfoData
}

var pluginCenter = &PluginCenter{
	plugins:     map[string]PluginInstance{},
	pluginInfos: map[string]*PluginInfoData{},
}

// GetPluginCenter 获取插件中心实例
func GetPluginCenter() *PluginCenter {
	return pluginCenter
}

// Register 注册插件实例（key：普通插件用 keyName，Hook 插件用 hook 名）
func (c *PluginCenter) Register(key string, instance PluginInstance, info *PluginInfoData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.plugins[key] = instance
	c.pluginInfos[key] = info
}

// Remove 移除插件并停止其进程
func (c *PluginCenter) Remove(key string, isHook bool) {
	c.mu.Lock()
	inst := c.plugins[key]
	delete(c.plugins, key)
	delete(c.pluginInfos, key)
	c.mu.Unlock()
	if sp, ok := inst.(*pluginProcess); ok {
		sp.Stop()
	}
	if lf, ok := inst.(*lazyProcessFactory); ok {
		lf.Stop()
	}
}

// GetPluginInfo 获取插件信息
func (c *PluginCenter) GetPluginInfo(key string) *PluginInfoData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pluginInfos[key]
}

// GetInstance 获取插件实例
func (c *PluginCenter) GetInstance(key string) PluginInstance {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.plugins[key]
}

// LoadProcessPlugin 加载并注册插件可执行程序
// singleton=true 常驻进程；false 每次调用拉起进程（对齐参考实现的实例语义）
func LoadProcessPlugin(keyName, hook, binPath string, info *PluginInfoData) error {
	pc := GetPluginCenter()
	var (
		instance PluginInstance
		err      error
	)
	if info.Singleton {
		instance, err = startPluginProcess(keyName, binPath, info.BaseDir, info)
		if err != nil {
			return err
		}
	} else {
		instance = &lazyProcessFactory{key: keyName, binPath: binPath, info: info}
	}
	regKey := keyName
	if hook != "" {
		regKey = hook
	}
	pc.Register(regKey, instance, info)
	return nil
}

// lazyProcessFactory 非单例插件：每次 Invoke 拉起独立进程，调用后回收
type lazyProcessFactory struct {
	key     string
	binPath string
	info    *PluginInfoData
	mu      sync.Mutex
	proc    *pluginProcess // 当前调用占用的进程
}

func (l *lazyProcessFactory) Invoke(method string, args ...interface{}) (interface{}, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	proc, err := startPluginProcess(l.key, l.binPath, l.info.BaseDir, l.info)
	if err != nil {
		return nil, err
	}
	l.proc = proc
	defer func() {
		proc.Stop()
		l.proc = nil
	}()
	return proc.Invoke(method, args...)
}

// Stop 停止工厂（如有进行中的进程）
func (l *lazyProcessFactory) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.proc != nil {
		l.proc.Stop()
		l.proc = nil
	}
}

// InvokeHook 调用 Hook 插件方法；未注册返回 ErrNoHookPlugin（调用方回退默认实现）
func InvokeHook(ctx context.Context, hook, method string, args ...interface{}) (interface{}, error) {
	inst := GetPluginCenter().GetInstance(hook)
	if inst == nil {
		return nil, ErrNoHookPlugin
	}
	return inst.Invoke(method, args...)
}

// Invoke 调用指定插件方法（跨插件调用入口，插件 SDK pluginService.Invoke 走此路径）
func Invoke(ctx context.Context, key, method string, params ...interface{}) (interface{}, error) {
	inst := GetPluginCenter().GetInstance(key)
	if inst == nil {
		return nil, fmt.Errorf("插件[%s]不存在或已禁用", key)
	}
	return inst.Invoke(method, params...)
}

// HasHook 是否有已启用的 Hook 插件
func HasHook(hook string) bool {
	return GetPluginCenter().GetInstance(hook) != nil
}
