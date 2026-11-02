package cmd

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
)

// defaultPackPlatforms 默认交叉编译平台
const defaultPackPlatforms = "linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64"

// PluginPack 插件打包命令：交叉编译多平台二进制并生成可安装的插件 ZIP
// （vgo-tools 插件打包命令）
var PluginPack = &gcmd.Command{
	Name:  "plugin-pack",
	Usage: "vgo-tools plugin-pack [dir] [-o 输出.zip] [-p 平台列表] [-m 主包路径]",
	Brief: "打包 Go 插件（交叉编译多平台并生成安装 ZIP）",
	Description: `在插件项目目录中执行交叉编译，生成符合宿主安装约定的 ZIP：

  plugin.json / README.md / assets/... / bin/<goos>-<goarch>/plugin[.exe]

平台格式为 goos-goarch，空格分隔；Windows 产物自动使用 plugin.exe。
构建参数：CGO_ENABLED=0 -trimpath -ldflags "-s -w"（静态、精简）。
打包自动排除 go build 产物与源码文件。`,
	Arguments: []gcmd.Argument{
		{Name: "output", Short: "o", Brief: "输出 ZIP 路径（默认 <key>-<version>.zip）"},
		{Name: "platforms", Short: "p", Brief: "平台列表，空格分隔（默认全平台）"},
		{Name: "main", Short: "m", Brief: "main 包相对路径（默认 .）"},
		{Name: "dir", Brief: "插件项目目录（默认当前目录）", IsArg: true},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		dir := parser.GetArg(2).String()
		if dir == "" {
			dir = "."
		}
		return packPlugin(dir,
			parser.GetOpt("output", parser.GetOpt("o").String()).String(),
			parser.GetOpt("platforms", parser.GetOpt("p").String()).String(),
			parser.GetOpt("main", parser.GetOpt("m").String()).String(),
		)
	},
}

// packPlugin 打包插件（dir 插件项目目录；output/platforms/mainPkg 为空取默认）
func packPlugin(dir, output, platforms, mainPkg string) (err error) {
	if !gfile.Exists(filepath.Join(dir, "plugin.json")) {
		return fmt.Errorf("目录缺少 plugin.json: %s", dir)
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	// plugin.json → key / version
	j, err := gjson.Load(filepath.Join(dir, "plugin.json"))
	if err != nil {
		return fmt.Errorf("plugin.json 解析失败: %w", err)
	}
	key := j.Get("key").String()
	version := j.Get("version").String()
	if key == "" {
		return fmt.Errorf("plugin.json 缺少 key")
	}
	if version == "" {
		version = "0.0.0"
	}

	if output == "" {
		output = fmt.Sprintf("%s-%s.zip", key, version)
	}
	if output, err = filepath.Abs(output); err != nil {
		return err
	}
	if platforms == "" {
		platforms = defaultPackPlatforms
	}
	if mainPkg == "" {
		mainPkg = "."
	}

	fmt.Printf("== 打包插件: %s v%s\n   源码: %s\n   主包: %s\n   平台: %s\n", key, version, absDir, mainPkg, platforms)

	// 临时打包目录
	stage, err := os.MkdirTemp("", "plugin-pack-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)

	// go build 产物（模块名二进制）不入包
	binNames := map[string]bool{"plugin": true, "plugin.exe": true}
	if modName := moduleNameFromGoMod(filepath.Join(dir, "go.mod")); modName != "" {
		binNames[modName] = true
		binNames[modName+".exe"] = true
	}

	// 静态资源：plugin.json + 源码/构建产物以外的全部文件
	if err = copyFile(filepath.Join(dir, "plugin.json"), filepath.Join(stage, "plugin.json")); err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if name == "plugin.json" || shouldSkipOnPack(name) || binNames[name] {
			continue
		}
		src := filepath.Join(dir, name)
		dst := filepath.Join(stage, name)
		if e.IsDir() {
			err = copyDir(src, dst)
		} else {
			err = copyFile(src, dst)
		}
		if err != nil {
			return err
		}
	}

	// 逐平台交叉编译
	for _, platform := range strings.Fields(platforms) {
		parts := strings.SplitN(platform, "-", 2)
		if len(parts) != 2 {
			return fmt.Errorf("平台格式错误: %q（应为 goos-goarch）", platform)
		}
		goos, goarch := parts[0], parts[1]
		outName := "plugin"
		if goos == "windows" {
			outName = "plugin.exe"
		}
		binDir := filepath.Join(stage, "bin", platform)
		if err = os.MkdirAll(binDir, 0o755); err != nil {
			return err
		}
		outBin := filepath.Join(binDir, outName)

		fmt.Printf("   构建 %s ... ", platform)
		cmdBuild := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", outBin, mainPkg)
		cmdBuild.Dir = absDir
		// 插件项目按自身 go.mod 独立构建，不受宿主仓库 go.work 影响
		cmdBuild.Env = append(os.Environ(), "CGO_ENABLED=0", "GOWORK=off", "GOOS="+goos, "GOARCH="+goarch)
		cmdBuild.Stdout = os.Stdout
		cmdBuild.Stderr = os.Stderr
		if err = cmdBuild.Run(); err != nil {
			fmt.Println("失败")
			return fmt.Errorf("构建 %s 失败: %w", platform, err)
		}
		info, _ := os.Stat(outBin)
		fmt.Printf("ok (%s)\n", humanSize(info.Size()))
	}

	// 打包 ZIP
	if err = zipDir(stage, output); err != nil {
		return err
	}
	info, _ := os.Stat(output)
	fmt.Printf("== 完成: %s (%s)\n", output, humanSize(info.Size()))
	fmt.Printf("   安装: curl -X POST http://<host>:8001/admin/plugin/info/install -H \"Authorization: <token>\" -F \"file=@%s\"\n", output)
	return nil
}

// moduleNameFromGoMod 读取 go.mod 的模块名（用于排除 go build 产物）
func moduleNameFromGoMod(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			mod := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if i := strings.LastIndex(mod, "/"); i >= 0 {
				mod = mod[i+1:]
			}
			return mod
		}
	}
	return ""
}

// shouldSkipOnPack 打包时跳过的源码/构建产物
func shouldSkipOnPack(name string) bool {
	switch name {
	case "go.mod", "go.sum", "build", "bin":
		return true
	}
	if strings.HasSuffix(name, ".go") || strings.HasSuffix(name, ".zip") {
		return true
	}
	return false
}

// copyFile 单文件拷贝
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// copyDir 目录递归拷贝
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

// zipDir 以目录为根打包 ZIP（保证 bin/... 相对路径正确）
func zipDir(root, output string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	zw := zip.NewWriter(out)
	defer zw.Close()

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)
		if info.IsDir() {
			header.Name += "/"
			_, err = zw.CreateHeader(header)
			return err
		}
		// 压缩存储；可执行文件保留执行权限
		header.Method = zip.Deflate
		if info.Mode()&0o111 != 0 {
			header.SetMode(info.Mode())
		}
		w, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		return err
	})
}

// humanSize 人类可读大小
func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(n)/float64(div), "KMGT"[exp])
}

func init() {
	Main.AddCommand(PluginPack)
}
