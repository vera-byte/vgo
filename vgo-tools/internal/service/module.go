package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

const demoModulePrefix = "vgo-simple/modules/demo"

var pkgRe = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// CreatModule 在项目 modules/ 下创建新模块（复制内嵌脚手架 modules/demo）。
func CreatModule(ctx context.Context, moduleName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = gerror.Newf("创建模块失败: %v", r)
		}
	}()
	// 检测当前目录是否存在go.mod文件
	if !gfile.Exists("go.mod") {
		return gerror.New("当前目录不存在go.mod文件,请在项目根目录下执行")
	}
	module := modulePathFromGoMod()
	if module == "" {
		return gerror.New("go.mod文件中不存在module行")
	}
	// 内嵌脚手架缺失时给出明确错误（而非空目录/panic）
	if !gres.Contains(demoModulePrefix) {
		return gerror.New("内嵌模板缺少 modules/demo 脚手架，vgo-tools 资源包可能已损坏")
	}
	moduleDir := gfile.Join(gfile.Pwd(), "modules", moduleName)
	if gfile.Exists(moduleDir) {
		return gerror.New("模块已经存在,请先删除原有模块")
	}
	if err = gfile.Mkdir(moduleDir); err != nil {
		return
	}
	// 创建模块目录结构
	if err = gres.Export(demoModulePrefix, moduleDir, gres.ExportOption{
		RemovePrefix: demoModulePrefix,
	}); err != nil {
		return gerror.Wrap(err, "导出模块脚手架失败")
	}
	// 替换import路径：vgo-simple/modules/demo -> <module>/modules/<moduleName>
	if err = gfile.ReplaceDir(demoModulePrefix, module+"/modules/"+moduleName, moduleDir, "*", true); err != nil {
		return
	}
	// 包名同步为合法标识符（demo -> <pkg>）
	pkg := strings.ToLower(strings.ReplaceAll(moduleName, "-", ""))
	if !pkgRe.MatchString(pkg) {
		return gerror.Newf("模块名 %q 无法派生合法 Go 包名（需小写字母/数字，可含 -）", moduleName)
	}
	if err = gfile.ReplaceDir("package demo", "package "+pkg, moduleDir, "*.go", true); err != nil {
		return
	}
	// 重命名 demo.go 为 moduleName.go
	if err = gfile.Rename(gfile.Join(moduleDir, "demo.go"), gfile.Join(moduleDir, moduleName+".go")); err != nil {
		return
	}
	// 接线 modules/modules.go
	wireModulesFile(module, moduleName)
	println("创建模块成功:", moduleDir)
	return nil
}

// modulePathFromGoMod 读取 go.mod 的 module 行（只认行首 module 前缀，避免误匹配 require 路径）。
func modulePathFromGoMod() string {
	module := ""
	_ = gfile.ReadLines("go.mod", func(text string) error {
		line := strings.TrimSpace(text)
		if strings.HasPrefix(line, "module ") && module == "" {
			module = strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
		return nil
	})
	return module
}

// wireModulesFile 把新模块的 blank import 追加到 modules/modules.go（存在且未引入时）。
func wireModulesFile(module, moduleName string) {
	f := gfile.Join("modules", "modules.go")
	if !gfile.Exists(f) {
		return
	}
	importPath := fmt.Sprintf("%s/modules/%s", module, moduleName)
	content := gfile.GetContents(f)
	if strings.Contains(content, importPath) {
		return
	}
	old := "import ("
	idx := strings.Index(content, old)
	if idx < 0 {
		return
	}
	insert := fmt.Sprintf("%s\n\t_ %q", old, importPath)
	content = strings.Replace(content, old, insert, 1)
	_ = gfile.PutContents(f, content)
}
