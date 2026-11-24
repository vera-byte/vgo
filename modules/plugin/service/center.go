package service

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/vera-byte/vgo/v"
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

// pluginCmd 插件进程启动命令：生产为已安装可执行文件，开发模式为 go run 源码
type pluginCmd struct {
	Name string   // 可执行文件路径或工具名（如 go）
	Args []string // 附加参数
}

// binCmd 已安装插件的可执行文件（绝对路径，避免相对路径与 cmd.Dir 组合导致 fork/exec 解析失败）
func binCmd(binPath string) pluginCmd {
	if abs, err := filepath.Abs(binPath); err == nil {
		binPath = abs
	}
	return pluginCmd{Name: binPath}
}

// devCmd 开发模式源码启动：在源码目录内 go run .（依赖该模块 go.mod 的 replace 指向工作区 SDK）
func devCmd() pluginCmd {
	return pluginCmd{Name: "go", Args: []string{"run", "."}}
}

// LoadProcessPlugin 加载并注册插件进程
// singleton=true 常驻进程；false 每次调用拉起进程（对齐参考实现的实例语义）
func LoadProcessPlugin(keyName, hook string, cmd pluginCmd, info *PluginInfoData) error {
	pc := GetPluginCenter()
	var (
		instance PluginInstance
		err      error
	)
	if info.Singleton {
		instance, err = startPluginProcess(keyName, cmd, info.BaseDir, info)
		if err != nil {
			return err
		}
	} else {
		instance = &lazyProcessFactory{key: keyName, cmd: cmd, info: info}
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
	key  string
	cmd  pluginCmd
	info *PluginInfoData
	mu   sync.Mutex
	proc *pluginProcess // 当前调用占用的进程
}

func (l *lazyProcessFactory) Invoke(method string, args ...interface{}) (interface{}, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	proc, err := startPluginProcess(l.key, l.cmd, l.info.BaseDir, l.info)
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

// vPluginInvoker 将调用能力注册到 v.Plugin：业务模块经 v.Plugin.Invoke 调用插件，无需依赖 plugin 模块
type vPluginInvoker struct{}

func (vPluginInvoker) Invoke(ctx context.Context, key, method string, params ...interface{}) (interface{}, error) {
	return Invoke(ctx, key, method, params...)
}

func init() {
	v.RegisterPluginInvoker(vPluginInvoker{})
}

// HasHook 是否有已启用的 Hook 插件
func HasHook(hook string) bool {
	return GetPluginCenter().GetInstance(hook) != nil
}
