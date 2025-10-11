package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

// findResourceDirs 查找指定根目录下所有名为"resource"的目录
// 功能：递归遍历目录，查找所有resource目录
// 参数：rootPath - 根目录路径
// 返回值：resourceDirs - resource目录列表，err - 错误信息
func findResourceDirs(rootPath string) ([]string, error) {
	var resourceDirs []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过隐藏目录和node_modules
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" {
				return filepath.SkipDir
			}

			// 如果目录名为resource，则添加到列表
			if name == "resource" {
				resourceDirs = append(resourceDirs, path)
			}
		}

		return nil
	})

	return resourceDirs, err
}

// packSingleModule 打包单个模块的资源文件
// 功能：检查模块是否有resource目录，如果有则打包到packed目录
// 参数：ctx - 上下文，moduleDir - 模块目录路径，moduleName - 模块名称
// 返回值：err - 错误信息
func packSingleModule(ctx context.Context, moduleDir, moduleName string) error {
	resourceDir := filepath.Join(moduleDir, "resource")
	
	// 检查resource目录是否存在
	if !gfile.Exists(resourceDir) {
		g.Log().Infof(ctx, "模块 %s 没有resource目录，跳过", moduleName)
		return nil
	}

	// 确保packed目录存在
	packedDir := filepath.Join(moduleDir, "packed")
	if !gfile.Exists(packedDir) {
		if err := gfile.Mkdir(packedDir); err != nil {
			return err
		}
	}

	// 打包resource目录到resource.go文件
	targetFile := filepath.Join(packedDir, "resource.go")
	g.Log().Infof(ctx, "正在打包 %s 到 %s", resourceDir, targetFile)
	
	if err := gres.PackToGoFile(resourceDir, targetFile, "packed"); err != nil {
		return err
	}

	g.Log().Infof(ctx, "打包成功: %s", targetFile)
	return nil
}

// clearSingleModule 清理单个模块的打包文件
// 功能：删除模块packed目录下的resource.go文件
// 参数：ctx - 上下文，moduleDir - 模块目录路径，moduleName - 模块名称
// 返回值：err - 错误信息
func clearSingleModule(ctx context.Context, moduleDir, moduleName string) error {
	packedDir := filepath.Join(moduleDir, "packed")
	resourceFile := filepath.Join(packedDir, "resource.go")
	
	if gfile.Exists(resourceFile) {
		if err := gfile.Remove(resourceFile); err != nil {
			return err
		}
		g.Log().Infof(ctx, "已删除: %s", resourceFile)
	} else {
		g.Log().Infof(ctx, "模块 %s 没有打包文件", moduleName)
	}
	
	return nil
}

// clearAllModules 清理所有模块的打包文件
// 功能：遍历modules目录，清理所有模块的打包文件
// 参数：ctx - 上下文，modulesDir - modules目录路径
// 返回值：err - 错误信息
func clearAllModules(ctx context.Context, modulesDir string) error {
	// 遍历modules目录下的所有子目录
	entries, err := gfile.ScanDir(modulesDir, "*", false)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if gfile.IsDir(entry) {
			moduleName := gfile.Basename(entry)
			g.Log().Infof(ctx, "正在清理模块: %s", moduleName)
			
			if err := clearSingleModule(ctx, entry, moduleName); err != nil {
				g.Log().Errorf(ctx, "清理模块 %s 失败: %v", moduleName, err)
				continue
			}
		}
	}
	
	return nil
}

// listPackedModules 列出所有模块的打包状态
// 功能：遍历modules目录，显示每个模块的打包状态
// 参数：ctx - 上下文，modulesDir - modules目录路径
// 返回值：err - 错误信息
func listPackedModules(ctx context.Context, modulesDir string) error {
	// 遍历modules目录下的所有子目录
	entries, err := gfile.ScanDir(modulesDir, "*", false)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if gfile.IsDir(entry) {
			moduleName := gfile.Basename(entry)
			resourceDir := filepath.Join(entry, "resource")
			packedFile := filepath.Join(entry, "packed", "resource.go")
			
			hasResource := gfile.Exists(resourceDir)
			hasPacked := gfile.Exists(packedFile)
			
			status := "无资源"
			if hasResource && hasPacked {
				status = "已打包"
			} else if hasResource && !hasPacked {
				status = "未打包"
			}
			
			g.Log().Infof(ctx, "模块 %s: %s", moduleName, status)
		}
	}
	
	return nil
}

// listAllResources 列出所有资源文件信息
// 功能：显示项目中所有resource目录的详细信息
// 参数：ctx - 上下文
// 返回值：err - 错误信息
func listAllResources(ctx context.Context) error {
	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 查找所有resource目录
	resourceDirs, err := findResourceDirs(rootDir)
	if err != nil {
		return err
	}

	g.Log().Infof(ctx, "找到 %d 个resource目录:", len(resourceDirs))
	
	for _, resourceDir := range resourceDirs {
		// 计算相对路径
		relPath, err := filepath.Rel(rootDir, resourceDir)
		if err != nil {
			relPath = resourceDir
		}
		
		// 统计文件数量
		files, err := gfile.ScanDirFile(resourceDir, "*", true)
		if err != nil {
			g.Log().Errorf(ctx, "扫描目录 %s 失败: %v", relPath, err)
			continue
		}
		
		g.Log().Infof(ctx, "  %s (%d 个文件)", relPath, len(files))
	}
	
	return nil
}