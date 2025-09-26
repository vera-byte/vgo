package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

func init() {
	packModuleCmd := &gcmd.Command{
		Name:  "packmodule",
		Usage: "packmodule",
		Brief: "打包模块资源文件",
		Description: `
自动打包指定模块或所有模块的资源文件，包括SQL迁移文件。
按照项目规范，会自动查找并打包以下资源：
1. resource/migrates 目录下的SQL文件 -> packed/migrates.go
2. resource/public 目录下的静态文件 -> packed/public.go

示例:
  vgo packmodule           # 打包所有模块的资源
  vgo packmodule -m base   # 只打包base模块的资源
`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			moduleName := parser.GetOpt("m").String()

			// 获取项目根目录
			rootDir, err := os.Getwd()
			if err != nil {
				g.Log().Error(ctx, "获取当前目录失败:", err)
				return err
			}

			// 模块目录
			modulesDir := filepath.Join(rootDir, "modules")

			// 检查模块目录是否存在
			if !gfile.Exists(modulesDir) {
				g.Log().Error(ctx, "模块目录不存在:", modulesDir)
				return fmt.Errorf("模块目录不存在: %s", modulesDir)
			}

			// 如果指定了模块名，只打包该模块
			if moduleName != "" {
				moduleDir := filepath.Join(modulesDir, moduleName)
				if !gfile.Exists(moduleDir) {
					g.Log().Errorf(ctx, "模块 %s 不存在", moduleName)
					return fmt.Errorf("模块 %s 不存在", moduleName)
				}

				return packSingleModule(ctx, moduleDir, moduleName)
			}

			// 否则打包所有模块
			modules, err := gfile.ScanDir(modulesDir, "*", false)
			if err != nil {
				g.Log().Error(ctx, "扫描模块目录失败:", err)
				return err
			}

			for _, moduleDir := range modules {
				// 跳过非目录文件
				if !gfile.IsDir(moduleDir) {
					continue
				}

				moduleName := filepath.Base(moduleDir)
				g.Log().Infof(ctx, "正在处理模块: %s", moduleName)

				if err := packSingleModule(ctx, moduleDir, moduleName); err != nil {
					g.Log().Warningf(ctx, "打包模块 %s 时出错: %v", moduleName, err)
					// 继续处理其他模块
					continue
				}
			}

			g.Log().Info(ctx, "所有模块资源打包完成")
			return nil
		},
	}

	// 不需要显式添加命令行参数，直接在Func中通过parser.GetOpt获取

	// 注册packmodule命令到根命令
	Root.AddCommand(packModuleCmd)

	// 递归查找指定目录下的所有resource目录
	findResourceDirs := func(rootPath string) ([]string, error) {
		var resourceDirs []string

		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 跳过.git、.github等隐藏目录
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			// 跳过node_modules目录
			if info.IsDir() && info.Name() == "node_modules" {
				return filepath.SkipDir
			}

			// 找到resource目录
			if info.IsDir() && info.Name() == "resource" {
				resourceDirs = append(resourceDirs, path)
			}

			return nil
		})

		return resourceDirs, err
	}

	// 创建打包根目录资源的命令
	packRootCmd := &gcmd.Command{
		Name:  "packroot",
		Usage: "packroot",
		Brief: "打包项目所有resource目录资源文件",
		Description: `
自动查找并打包项目中所有的resource目录。
使用gres.PackWithOption方法进行打包，按照实际路径生成对应的Go文件。

示例:
  vgo packroot           # 打包所有resource目录
`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 获取项目根目录
			rootDir, err := os.Getwd()
			if err != nil {
				g.Log().Error(ctx, "获取当前目录失败:", err)
				return err
			}

			// 查找所有resource目录
			g.Log().Infof(ctx, "正在查找项目中所有的resource目录")
			resourceDirs, err := findResourceDirs(rootDir)
			if err != nil {
				g.Log().Errorf(ctx, "查找resource目录失败: %v", err)
				return err
			}

			g.Log().Infof(ctx, "找到 %d 个resource目录", len(resourceDirs))

			// 遍历所有resource目录进行打包
			for _, resourceDir := range resourceDirs {
				// 计算相对路径，用于生成包名
				relPath, err := filepath.Rel(rootDir, resourceDir)
				if err != nil {
					g.Log().Errorf(ctx, "计算相对路径失败: %v", err)
					continue
				}

				// 确定打包后的Go文件存放位置
				// 如果是模块内的resource，则放在模块的packed目录下
				// 如果是根目录的resource，则放在根目录的packed目录下
				var packedDir, goFileName string

				if strings.HasPrefix(relPath, "modules/") {
					// 模块内的resource
					parts := strings.Split(relPath, "/")
					if len(parts) >= 2 {
						moduleName := parts[1]
						moduleDir := filepath.Join(rootDir, "modules", moduleName)
						packedDir = filepath.Join(moduleDir, "packed")
						// 将生成的文件名从 resource.go 改为 packed.go
						goFileName = "packed.go"
					}
				} else {
					// 根目录的resource
					packedDir = filepath.Join(rootDir, "packed")

					// 根据resource目录的相对路径生成文件名
					// 例如：i18n/resource -> i18n_resource.go
					goFileName = strings.ReplaceAll(relPath, "/", "_") + ".go"
					if goFileName == "resource.go" {
						// 如果是根目录下的resource目录，使用packed.go
						goFileName = "packed.go"
					}
				}

				// 确保packed目录存在
				if !gfile.Exists(packedDir) {
					if err := gfile.Mkdir(packedDir); err != nil {
						g.Log().Errorf(ctx, "创建packed目录失败: %v", err)
						continue
					}
				}

				goFilePath := filepath.Join(packedDir, goFileName)

				g.Log().Infof(ctx, "正在打包 %s", relPath)

				// 使用gres.PackToGoFile方法打包资源到Go文件
				// 获取目录名作为包名
				pkgName := filepath.Base(packedDir)
				if pkgName == "packed" {
					// 如果是packed目录，使用上一级目录名作为包名
					parentDir := filepath.Dir(packedDir)
					pkgName = filepath.Base(parentDir)
				}

				if err := gres.PackToGoFile(resourceDir, goFilePath, pkgName, "resource"); err != nil {
					g.Log().Errorf(ctx, "打包 %s 失败: %v", relPath, err)
					continue
				}

				g.Log().Infof(ctx, "成功打包 %s 到 %s", relPath, goFilePath)
			}

			g.Log().Infof(ctx, "所有resource目录打包完成")

			return nil
		},
	}

	// 注册packroot命令到根命令
	Root.AddCommand(packRootCmd)
}

// packSingleModule 打包单个模块的资源
func packSingleModule(ctx context.Context, moduleDir, moduleName string) error {
	// 确保packed目录存在
	packedDir := filepath.Join(moduleDir, "packed")
	if !gfile.Exists(packedDir) {
		if err := gfile.Mkdir(packedDir); err != nil {
			g.Log().Errorf(ctx, "创建模块 %s 的packed目录失败: %v", moduleName, err)
			return err
		}
	}

	// 打包整个resource目录
	resourceDir := filepath.Join(moduleDir, "resource")
	if gfile.Exists(resourceDir) {
		// 将生成的文件名从 resource.go 改为 packed.go
		resourceGoFile := filepath.Join(packedDir, "packed.go")
		g.Log().Infof(ctx, "正在打包模块 %s 的resource目录", moduleName)

		// 使用gres.PackToGoFile方法打包资源到Go文件
		if err := gres.PackToGoFile(resourceDir, resourceGoFile, "packed", "resource"); err != nil {
			g.Log().Errorf(ctx, "打包模块 %s 的resource目录失败: %v", moduleName, err)
			return err
		}

		g.Log().Infof(ctx, "模块 %s 的resource目录打包完成", moduleName)
	} else {
		g.Log().Infof(ctx, "模块 %s 没有resource目录，跳过", moduleName)
	}

	return nil
}
