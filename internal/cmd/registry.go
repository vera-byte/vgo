package cmd

import (
	"context"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
)

// CommandProvider 命令提供者接口
// 模块需要实现此接口来提供命令
type CommandProvider interface {
	// GetCommands 获取模块提供的命令列表
	// 返回: 命令列表
	GetCommands() []*gcmd.Command
	
	// GetModuleName 获取模块名称
	// 返回: 模块名称
	GetModuleName() string
}

// CommandRegistry 命令注册中心
// 负责管理和注册所有模块的命令
type CommandRegistry struct {
	providers []CommandProvider
	commands  map[string]*gcmd.Command
	mutex     sync.RWMutex
}

var (
	// registry 全局命令注册中心实例
	registry *CommandRegistry
	// registryOnce 确保注册中心只初始化一次
	registryOnce sync.Once
)

// GetRegistry 获取全局命令注册中心实例
// 使用单例模式确保全局唯一
// 返回: 命令注册中心实例
func GetRegistry() *CommandRegistry {
	registryOnce.Do(func() {
		registry = &CommandRegistry{
			providers: make([]CommandProvider, 0),
			commands:  make(map[string]*gcmd.Command),
		}
	})
	return registry
}

// RegisterProvider 注册命令提供者
// provider: 命令提供者实例
func (r *CommandRegistry) RegisterProvider(provider CommandProvider) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.providers = append(r.providers, provider)
	g.Log().Infof(context.Background(), "Registered command provider: %s", provider.GetModuleName())
}

// RegisterCommands 注册所有提供者的命令到根命令
// 遍历所有已注册的命令提供者，将其命令注册到根命令中
func (r *CommandRegistry) RegisterCommands() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for _, provider := range r.providers {
		commands := provider.GetCommands()
		for _, cmd := range commands {
			if cmd != nil {
				// 检查命令名称是否已存在
				if existingCmd, exists := r.commands[cmd.Name]; exists {
					g.Log().Warningf(context.Background(), 
						"Command '%s' already exists (from %s), skipping registration from %s", 
						cmd.Name, r.getCommandSource(existingCmd), provider.GetModuleName())
					continue
				}
				
				// 注册命令到根命令
				Root.AddCommand(cmd)
				r.commands[cmd.Name] = cmd
				
				g.Log().Infof(context.Background(), 
					"Registered command '%s' from module '%s'", 
					cmd.Name, provider.GetModuleName())
			}
		}
	}
}

// getCommandSource 获取命令来源信息
// cmd: 命令实例
// 返回: 命令来源描述
func (r *CommandRegistry) getCommandSource(cmd *gcmd.Command) string {
	if cmd.Func != nil {
		funcPtr := reflect.ValueOf(cmd.Func).Pointer()
		funcInfo := runtime.FuncForPC(funcPtr)
		if funcInfo != nil {
			file, _ := funcInfo.FileLine(funcPtr)
			return filepath.Base(file)
		}
	}
	return "unknown"
}

// GetRegisteredCommands 获取已注册的命令列表
// 返回: 命令名称到命令实例的映射
func (r *CommandRegistry) GetRegisteredCommands() map[string]*gcmd.Command {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	result := make(map[string]*gcmd.Command)
	for name, cmd := range r.commands {
		result[name] = cmd
	}
	return result
}

// AutoDiscoverModules 自动发现模块命令
// 扫描modules目录下的所有模块，尝试自动注册命令
// modulesPath: 模块目录路径
func (r *CommandRegistry) AutoDiscoverModules(modulesPath string) {
	if !gfile.Exists(modulesPath) {
		g.Log().Warningf(context.Background(), "Modules path does not exist: %s", modulesPath)
		return
	}
	
	// 获取所有模块目录
	dirs, err := gfile.ScanDir(modulesPath, "*", false)
	if err != nil {
		g.Log().Errorf(context.Background(), "Failed to scan modules directory: %v", err)
		return
	}
	
	for _, dir := range dirs {
		if gfile.IsDir(dir) {
			moduleName := filepath.Base(dir)
			// 跳过隐藏目录和特殊目录
			if strings.HasPrefix(moduleName, ".") || moduleName == "cmd.go" {
				continue
			}
			
			g.Log().Infof(context.Background(), "Discovered module: %s", moduleName)
		}
	}
}

// ListProviders 列出所有已注册的命令提供者
// 返回: 提供者列表
func (r *CommandRegistry) ListProviders() []CommandProvider {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	result := make([]CommandProvider, len(r.providers))
	copy(result, r.providers)
	return result
}

// GetProviderByName 根据名称获取命令提供者
// name: 提供者名称
// 返回: 命令提供者实例，如果未找到则返回nil
func (r *CommandRegistry) GetProviderByName(name string) CommandProvider {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, provider := range r.providers {
		if provider.GetModuleName() == name {
			return provider
		}
	}
	return nil
}

// PrintRegistryStatus 打印注册中心状态信息
// 用于调试和监控命令注册情况
func (r *CommandRegistry) PrintRegistryStatus() {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	ctx := context.Background()
	g.Log().Infof(ctx, "=== Command Registry Status ===")
	g.Log().Infof(ctx, "Total providers: %d", len(r.providers))
	g.Log().Infof(ctx, "Total commands: %d", len(r.commands))
	
	g.Log().Infof(ctx, "Registered providers:")
	for i, provider := range r.providers {
		commands := provider.GetCommands()
		g.Log().Infof(ctx, "  %d. %s (%d commands)", i+1, provider.GetModuleName(), len(commands))
	}
	
	g.Log().Infof(ctx, "Registered commands:")
	for name := range r.commands {
		g.Log().Infof(ctx, "  - %s", name)
	}
	g.Log().Infof(ctx, "===============================")
}