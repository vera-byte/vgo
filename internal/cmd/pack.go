package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

// cPack pack命令的主结构体
type cPack struct {
	g.Meta `name:"pack" brief:"打包资源文件" dc:"pack命令用于管理项目资源文件的打包操作，支持模块级和项目级的资源打包"`
}

// cPackInput pack命令的输入参数
type cPackInput struct {
	g.Meta `name:"pack" brief:"打包资源文件" dc:"pack命令用于管理项目资源文件的打包操作，支持模块级和项目级的资源打包"`
}

// cPackOutput pack命令的输出
type cPackOutput struct{}

// Index pack命令的执行方法
// 功能：显示pack命令的帮助信息
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cPack) Index(ctx context.Context, in cPackInput) (out *cPackOutput, err error) {
	// 当直接执行 pack 命令时，显示帮助信息
	// 帮助信息现在通过 g.Meta 标签的结构化定义自动生成
	return &cPackOutput{}, nil
}

// cPackModule 打包模块资源文件的命令结构体
type cPackModule struct {
	g.Meta `name:"module" brief:"打包指定模块的资源文件" dc:"打包指定模块或所有模块的资源文件到packed目录"`
}

// cPackModuleInput pack module命令的输入参数
type cPackModuleInput struct {
	g.Meta `name:"module" brief:"打包指定模块的资源文件" dc:"打包指定模块或所有模块的资源文件到packed目录"`
	Module string `short:"m" name:"module" brief:"指定要打包的模块名称，不指定则打包所有模块" dc:"模块名称，如：base、dict、task等，留空则打包所有模块"`
}

// cPackModuleOutput pack module命令的输出
type cPackModuleOutput struct{}

// Index pack module命令的执行方法
// 功能：打包指定模块或所有模块的资源文件
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cPackModule) Index(ctx context.Context, in cPackModuleInput) (out *cPackModuleOutput, err error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		g.Log().Error(ctx, "获取当前目录失败:", err)
		return nil, err
	}

	// 确定modules目录路径
	modulesDir := filepath.Join(currentDir, "modules")
	if !gfile.Exists(modulesDir) {
		g.Log().Error(ctx, "modules目录不存在")
		return nil, err
	}

	if in.Module != "" {
		// 打包指定模块
		moduleDir := filepath.Join(modulesDir, in.Module)
		if !gfile.Exists(moduleDir) {
			g.Log().Errorf(ctx, "模块 %s 不存在", in.Module)
			return nil, err
		}

		g.Log().Infof(ctx, "正在打包模块: %s", in.Module)
		if err := packSingleModule(ctx, moduleDir, in.Module); err != nil {
			g.Log().Errorf(ctx, "打包模块 %s 失败: %v", in.Module, err)
			return nil, err
		}
		g.Log().Infof(ctx, "模块 %s 打包完成", in.Module)
	} else {
		// 打包所有模块
		g.Log().Info(ctx, "正在打包所有模块...")

		// 遍历modules目录下的所有子目录
		entries, err := gfile.ScanDir(modulesDir, "*", false)
		if err != nil {
			g.Log().Errorf(ctx, "扫描modules目录失败: %v", err)
			return nil, err
		}

		for _, entry := range entries {
			if gfile.IsDir(entry) {
				moduleName := gfile.Basename(entry)
				g.Log().Infof(ctx, "正在打包模块: %s", moduleName)

				if err := packSingleModule(ctx, entry, moduleName); err != nil {
					g.Log().Errorf(ctx, "打包模块 %s 失败: %v", moduleName, err)
					continue
				}

				g.Log().Infof(ctx, "模块 %s 打包完成", moduleName)
			}
		}

		g.Log().Info(ctx, "所有模块打包完成")
	}

	return &cPackModuleOutput{}, nil
}

// cPackRoot 打包项目根目录资源的命令结构体
type cPackRoot struct {
	g.Meta `name:"root" brief:"打包项目所有resource目录资源文件" dc:"打包项目根目录下所有resource目录的资源文件"`
}

// cPackRootInput pack root命令的输入参数
type cPackRootInput struct {
	g.Meta `name:"root" brief:"打包项目所有resource目录资源文件" dc:"打包项目根目录下所有resource目录的资源文件"`
}

// cPackRootOutput pack root命令的输出
type cPackRootOutput struct{}

// Index pack root命令的执行方法
// 功能：打包项目中所有resource目录的资源文件
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cPackRoot) Index(ctx context.Context, in cPackRootInput) (out *cPackRootOutput, err error) {
	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		g.Log().Error(ctx, "获取当前目录失败:", err)
		return nil, err
	}

	// 查找所有resource目录
	g.Log().Infof(ctx, "正在查找项目中所有的resource目录")
	resourceDirs, err := findResourceDirs(rootDir)
	if err != nil {
		g.Log().Errorf(ctx, "查找resource目录失败: %v", err)
		return nil, err
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
				goFileName = "resource.go"
			}
		} else {
			// 根目录的resource
			packedDir = filepath.Join(rootDir, "packed")
			// 使用相对路径作为文件名，将路径分隔符替换为下划线
			goFileName = strings.ReplaceAll(relPath, string(filepath.Separator), "_") + ".go"
		}

		if packedDir == "" || goFileName == "" {
			g.Log().Errorf(ctx, "无法确定打包目标路径: %s", relPath)
			continue
		}

		// 确保packed目录存在
		if !gfile.Exists(packedDir) {
			if err := gfile.Mkdir(packedDir); err != nil {
				g.Log().Errorf(ctx, "创建packed目录失败: %v", err)
				continue
			}
		}

		// 执行打包
		targetFile := filepath.Join(packedDir, goFileName)
		g.Log().Infof(ctx, "正在打包 %s 到 %s", resourceDir, targetFile)

		if err := gres.PackToGoFile(resourceDir, targetFile, "packed"); err != nil {
			g.Log().Errorf(ctx, "打包失败: %v", err)
			continue
		}

		g.Log().Infof(ctx, "打包成功: %s", targetFile)
	}

	g.Log().Info(ctx, "所有resource目录打包完成")
	return &cPackRootOutput{}, nil
}

// cPackList 列出打包状态的命令结构体
type cPackList struct {
	g.Meta `name:"list" brief:"列出模块打包状态" dc:"列出模块打包状态或所有资源文件信息"`
}

// cPackListInput pack list命令的输入参数
type cPackListInput struct {
	g.Meta `name:"list" brief:"列出模块打包状态" dc:"列出模块打包状态或所有资源文件信息"`
	All    bool `short:"a" name:"all" brief:"显示所有资源文件信息" dc:"是否显示所有资源文件的详细信息，包括未打包的文件"`
}

// cPackListOutput pack list命令的输出
type cPackListOutput struct{}

// Index pack list命令的执行方法
// 功能：列出模块打包状态或所有资源文件信息
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cPackList) Index(ctx context.Context, in cPackListInput) (out *cPackListOutput, err error) {
	if in.All {
		// 显示所有资源文件信息
		g.Log().Info(ctx, "显示所有资源文件信息:")
		if err := listAllResources(ctx); err != nil {
			g.Log().Errorf(ctx, "列出所有资源失败: %v", err)
			return nil, err
		}
	} else {
		// 显示模块打包状态
		// 获取当前工作目录
		currentDir, err := os.Getwd()
		if err != nil {
			g.Log().Error(ctx, "获取当前目录失败:", err)
			return nil, err
		}

		// 确定modules目录路径
		modulesDir := filepath.Join(currentDir, "modules")
		if !gfile.Exists(modulesDir) {
			g.Log().Error(ctx, "modules目录不存在")
			return nil, err
		}

		g.Log().Info(ctx, "模块打包状态:")
		if err := listPackedModules(ctx, modulesDir); err != nil {
			g.Log().Errorf(ctx, "列出模块打包状态失败: %v", err)
			return nil, err
		}
	}

	return &cPackListOutput{}, nil
}

// cPackClear 清理打包文件的命令结构体
type cPackClear struct {
	g.Meta `name:"clear" brief:"清理模块打包文件" dc:"清理指定模块或所有模块的打包文件"`
}

// cPackClearInput pack clear命令的输入参数
type cPackClearInput struct {
	g.Meta `name:"clear" brief:"清理模块打包文件" dc:"清理指定模块或所有模块的打包文件"`
	Module string `short:"m" name:"module" brief:"指定要清理的模块名称，不指定则清理所有模块" dc:"模块名称，如：base、dict、task等，留空则清理所有模块的打包文件"`
}

// cPackClearOutput pack clear命令的输出
type cPackClearOutput struct{}

// Index pack clear命令的执行方法
// 功能：清理指定模块或所有模块的打包文件
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cPackClear) Index(ctx context.Context, in cPackClearInput) (out *cPackClearOutput, err error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		g.Log().Error(ctx, "获取当前目录失败:", err)
		return nil, err
	}

	// 确定modules目录路径
	modulesDir := filepath.Join(currentDir, "modules")
	if !gfile.Exists(modulesDir) {
		g.Log().Error(ctx, "modules目录不存在")
		return nil, err
	}

	if in.Module != "" {
		// 清理指定模块
		moduleDir := filepath.Join(modulesDir, in.Module)
		if !gfile.Exists(moduleDir) {
			g.Log().Errorf(ctx, "模块 %s 不存在", in.Module)
			return nil, err
		}

		g.Log().Infof(ctx, "正在清理模块: %s", in.Module)
		if err := clearSingleModule(ctx, moduleDir, in.Module); err != nil {
			g.Log().Errorf(ctx, "清理模块 %s 失败: %v", in.Module, err)
			return nil, err
		}
		g.Log().Infof(ctx, "模块 %s 清理完成", in.Module)
	} else {
		// 清理所有模块
		g.Log().Info(ctx, "正在清理所有模块...")
		if err := clearAllModules(ctx, modulesDir); err != nil {
			g.Log().Errorf(ctx, "清理所有模块失败: %v", err)
			return nil, err
		}
		g.Log().Info(ctx, "所有模块清理完成")
	}

	return &cPackClearOutput{}, nil
}

// cPackAll 打包所有资源的命令结构体
type cPackAll struct {
	g.Meta `name:"all" brief:"打包所有模块和项目资源文件" dc:"打包所有模块和项目资源文件，包括模块资源和项目根目录资源"`
}

// cPackAllInput pack all命令的输入参数
type cPackAllInput struct {
	g.Meta `name:"all" brief:"打包所有模块和项目资源文件" dc:"打包所有模块和项目资源文件，包括模块资源和项目根目录资源"`
}

// cPackAllOutput pack all命令的输出
type cPackAllOutput struct{}

// Index pack all命令的执行方法
// 功能：打包所有模块和项目资源文件
// 参数：ctx - 上下文，in - 输入参数
// 返回值：out - 输出结果，err - 错误信息
func (c *cPackAll) Index(ctx context.Context, in cPackAllInput) (out *cPackAllOutput, err error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		g.Log().Error(ctx, "获取当前目录失败:", err)
		return nil, err
	}

	g.Log().Info(ctx, "开始打包所有资源...")

	// 1. 先打包所有模块
	modulesDir := filepath.Join(currentDir, "modules")
	if gfile.Exists(modulesDir) {
		g.Log().Info(ctx, "正在打包所有模块...")

		// 获取所有模块目录
		dirs, err := gfile.ScanDir(modulesDir, "*", false)
		if err != nil {
			g.Log().Errorf(ctx, "扫描modules目录失败: %v", err)
			return nil, err
		}

		// 打包每个模块
		for _, dir := range dirs {
			if gfile.IsDir(dir) {
				moduleName := filepath.Base(dir)
				g.Log().Infof(ctx, "正在打包模块: %s", moduleName)

				if err := packSingleModule(ctx, dir, moduleName); err != nil {
					g.Log().Errorf(ctx, "打包模块 %s 失败: %v", moduleName, err)
					return nil, err
				}
				g.Log().Infof(ctx, "模块 %s 打包完成", moduleName)
			}
		}
		g.Log().Info(ctx, "所有模块打包完成")
	}

	// 2. 再打包项目根目录的resource
	g.Log().Info(ctx, "正在打包项目根目录资源...")

	// 查找所有resource目录
	resourceDirs, err := findResourceDirs(currentDir)
	if err != nil {
		g.Log().Errorf(ctx, "查找resource目录失败: %v", err)
		return nil, err
	}

	if len(resourceDirs) == 0 {
		g.Log().Warning(ctx, "未找到任何resource目录")
	} else {
		// 打包所有resource目录
		for _, resourceDir := range resourceDirs {
			g.Log().Infof(ctx, "正在打包资源目录: %s", resourceDir)

			// 获取相对于项目根目录的路径作为包名
			relPath, err := filepath.Rel(currentDir, resourceDir)
			if err != nil {
				g.Log().Errorf(ctx, "获取相对路径失败: %v", err)
				continue
			}

			// 将路径分隔符替换为下划线作为包名
			packageName := strings.ReplaceAll(relPath, string(filepath.Separator), "_")

			if err := gres.PackToFile(resourceDir, filepath.Join(currentDir, "packed", packageName+".go")); err != nil {
				g.Log().Errorf(ctx, "打包资源目录 %s 失败: %v", resourceDir, err)
				return nil, err
			}
			g.Log().Infof(ctx, "资源目录 %s 打包完成", resourceDir)
		}
		g.Log().Info(ctx, "项目根目录资源打包完成")
	}

	g.Log().Info(ctx, "所有资源打包完成！")
	return &cPackAllOutput{}, nil
}

func init() {
	// 注册pack主命令
	packCmd, err := gcmd.NewFromObject(&cPack{})
	if err != nil {
		panic(err)
	}

	// 注册pack的子命令
	moduleCmd, err := gcmd.NewFromObject(&cPackModule{})
	if err != nil {
		panic(err)
	}
	packCmd.AddCommand(moduleCmd)

	rootCmd, err := gcmd.NewFromObject(&cPackRoot{})
	if err != nil {
		panic(err)
	}
	packCmd.AddCommand(rootCmd)

	listCmd, err := gcmd.NewFromObject(&cPackList{})
	if err != nil {
		panic(err)
	}
	packCmd.AddCommand(listCmd)

	clearCmd, err := gcmd.NewFromObject(&cPackClear{})
	if err != nil {
		panic(err)
	}
	packCmd.AddCommand(clearCmd)

	allCmd, err := gcmd.NewFromObject(&cPackAll{})
	if err != nil {
		panic(err)
	}
	packCmd.AddCommand(allCmd)

	// 将pack命令添加到根命令
	Root.AddCommand(packCmd)
}
