package vcmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
)

// 定义utils命令
var Utils = &gcmd.Command{
	Name:  "utils",
	Usage: "utils COMMAND [ARGUMENT]",
	Brief: "通用工具命令集合",
	Description: `
提供一系列通用工具命令，用于简化开发和运维工作。
包括文件操作、环境检查、资源打包等功能。
`,
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		g.Log().Info(ctx, "请指定要执行的子命令，使用 -h 查看帮助")
		return nil
	},
}

// 定义pack子命令
var UtilsPack = &gcmd.Command{
	Name:  "pack",
	Usage: "utils pack SOURCE TARGET",
	Brief: "打包资源文件",
	Description: `
将指定的源目录或文件打包到目标Go文件中。
示例:
  vgo utils pack ./public ./internal/packed/public.go
  vgo utils pack ./modules/base/resource/migrates ./modules/base/packed/migrates.go
`,
	Arguments: []gcmd.Argument{
		{
			Name:  "source",
			Short: "s",
			Brief: "源目录或文件路径",
		},
		{
			Name:  "target",
			Short: "t",
			Brief: "目标Go文件路径",
		},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		source := parser.GetArg(1).String()
		target := parser.GetArg(2).String()
		
		if source == "" || target == "" {
			g.Log().Error(ctx, "源路径和目标路径不能为空")
			return fmt.Errorf("源路径和目标路径不能为空")
		}
		
		// 确保源路径存在
		if !gfile.Exists(source) {
			g.Log().Errorf(ctx, "源路径 %s 不存在", source)
			return fmt.Errorf("源路径 %s 不存在", source)
		}
		
		// 确保目标目录存在
		targetDir := filepath.Dir(target)
		if !gfile.Exists(targetDir) {
			if err := gfile.Mkdir(targetDir); err != nil {
				g.Log().Errorf(ctx, "创建目标目录失败: %v", err)
				return err
			}
		}
		
		// 执行打包
		g.Log().Infof(ctx, "正在打包 %s 到 %s", source, target)
		
		// 构建gf pack命令
		packCmd := fmt.Sprintf("gf pack %s %s -y", source, target)
		result, err := gproc.ShellExec(ctx, packCmd)
		if err != nil {
			g.Log().Errorf(ctx, "打包失败: %v", err)
			return err
		}
		
		g.Log().Info(ctx, result)
		g.Log().Info(ctx, "打包完成")
		return nil
	},
}

// 定义clean子命令
var UtilsClean = &gcmd.Command{
	Name:  "clean",
	Usage: "utils clean [TYPE]",
	Brief: "清理临时文件",
	Description: `
清理项目中的临时文件和缓存。
支持的类型:
  all    - 清理所有临时文件和缓存
  temp   - 只清理临时文件
  cache  - 只清理缓存文件
  build  - 清理构建产物
`,
	Arguments: []gcmd.Argument{
		{
			Name:  "type",
			Short: "t",
			Brief: "清理类型",
		},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		cleanType := parser.GetArg(1).String()
		if cleanType == "" {
			cleanType = "all"
		}
		
		g.Log().Infof(ctx, "开始清理 %s 类型的文件", cleanType)
		
		// 获取项目根目录
		rootDir, err := os.Getwd()
		if err != nil {
			g.Log().Error(ctx, "获取当前目录失败:", err)
			return err
		}
		
		// 定义需要清理的目录
		var dirsToClean []string
		
		switch strings.ToLower(cleanType) {
		case "all":
			dirsToClean = []string{
				filepath.Join(rootDir, "temp"),
				filepath.Join(rootDir, "data", "cache"),
				filepath.Join(rootDir, "build"),
			}
		case "temp":
			dirsToClean = []string{
				filepath.Join(rootDir, "temp"),
			}
		case "cache":
			dirsToClean = []string{
				filepath.Join(rootDir, "data", "cache"),
			}
		case "build":
			dirsToClean = []string{
				filepath.Join(rootDir, "build"),
			}
		default:
			g.Log().Error(ctx, "不支持的清理类型:", cleanType)
			return fmt.Errorf("不支持的清理类型: %s", cleanType)
		}
		
		// 执行清理
		for _, dir := range dirsToClean {
			if gfile.Exists(dir) {
				g.Log().Infof(ctx, "正在清理目录: %s", dir)
				if err := gfile.Remove(dir); err != nil {
					g.Log().Warningf(ctx, "清理目录 %s 失败: %v", dir, err)
				}
			}
		}
		
		g.Log().Info(ctx, "清理完成")
		return nil
	},
}

// 定义scan子命令
var UtilsScan = &gcmd.Command{
	Name:  "scan",
	Usage: "utils scan [PATH]",
	Brief: "扫描项目结构",
	Description: `
扫描并显示项目结构，可用于检查模块组织和依赖关系。
示例:
  vgo utils scan            # 扫描整个项目
  vgo utils scan ./modules  # 只扫描模块目录
`,
	Arguments: []gcmd.Argument{
		{
			Name:  "path",
			Short: "p",
			Brief: "要扫描的路径",
		},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		scanPath := parser.GetArg(1).String()
		if scanPath == "" {
			scanPath = "."
		}
		
		g.Log().Infof(ctx, "开始扫描路径: %s", scanPath)
		
		// 获取绝对路径
		absPath, err := filepath.Abs(scanPath)
		if err != nil {
			g.Log().Error(ctx, "获取绝对路径失败:", err)
			return err
		}
		
		// 检查路径是否存在
		if !gfile.Exists(absPath) {
			g.Log().Errorf(ctx, "路径不存在: %s", absPath)
			return fmt.Errorf("路径不存在: %s", absPath)
		}
		
		// 扫描目录结构
		result := ""
		err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// 跳过隐藏文件和目录
			if strings.HasPrefix(filepath.Base(path), ".") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			
			// 计算相对路径
			relPath, err := filepath.Rel(absPath, path)
			if err != nil {
				return err
			}
			
			// 跳过根目录
			if relPath == "." {
				return nil
			}
			
			// 计算缩进
			indent := strings.Repeat("  ", len(strings.Split(relPath, string(filepath.Separator)))-1)
			
			// 添加到结果
			if info.IsDir() {
				result += fmt.Sprintf("%s📁 %s\n", indent, filepath.Base(path))
			} else {
				result += fmt.Sprintf("%s📄 %s\n", indent, filepath.Base(path))
			}
			
			return nil
		})
		
		if err != nil {
			g.Log().Error(ctx, "扫描失败:", err)
			return err
		}
		
		// 输出结果
		fmt.Println(result)
		g.Log().Info(ctx, "扫描完成")
		return nil
	},
}

// 注册所有utils子命令
func init() {
	// 注册utils命令到根命令
	Root.AddCommand(Utils)
	
	// 注册子命令到utils命令
	Utils.AddCommand(UtilsPack)
	Utils.AddCommand(UtilsClean)
	Utils.AddCommand(UtilsScan)
}