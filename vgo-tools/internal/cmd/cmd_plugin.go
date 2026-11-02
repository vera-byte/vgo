package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// Plugin 插件脚手架命令：一键生成 Go 插件项目模板
// （plugin.json、main.go、go.mod、README、打包脚本）
var Plugin = &gcmd.Command{
	Name:  "plugin",
	Usage: "vgo-tools plugin [dst] [--key 插件key] [--name 名称]",
	Brief: "创建 Go 插件项目脚手架",
	Description: `生成可执行插件项目模板（基于 pluginsdk）：
  plugin.json / main.go / go.mod / README.md

生成后：
  1) 编写插件方法（main.go 中的导出方法即插件 API）
  2) vgo-tools plugin-pack . 交叉编译多平台并生成安装 ZIP
  3) 管理端安装 ZIP 即可使用`,
	Arguments: []gcmd.Argument{
		{Name: "key", Short: "k", Brief: "插件唯一标识（默认取目录名）"},
		{Name: "name", Brief: "插件名称（默认同 key）"},
		{Name: "dst", Brief: "目标目录（默认当前目录）", IsArg: true},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		dst := parser.GetArg(2).String()
		if dst == "" {
			dst = "."
		}
		key := parser.GetOpt("key").String()
		if key == "" {
			key = gfile.Basename(gfile.RealPath(dst))
		}
		key = strings.ToLower(gstr.ReplaceByMap(key, map[string]string{"-": "_", " ": "_"}))
		if key == "" || key == "plugin" {
			return fmt.Errorf("插件key无效: %q（不能为 plugin 或空，可用 --key 指定）", key)
		}
		name := parser.GetOpt("name").String()
		if name == "" {
			name = key
		}

		// 目标目录检查（非空需确认）
		if !gfile.IsDir(dst) {
			if err = gfile.Mkdir(dst); err != nil {
				return
			}
		} else if !gfile.IsEmpty(dst) {
			s := gcmd.Scanf(`the folder "%s" is not empty, files might be overwrote, continue? [y/n]: `, dst)
			if !strings.EqualFold(s, "y") {
				return nil
			}
		}

		files := pluginTemplateFiles(key, name)
		for path, content := range files {
			full := filepath.Join(dst, path)
			if err = gfile.Mkdir(filepath.Dir(full)); err != nil {
				return
			}
			if err = gfile.PutContents(full, content); err != nil {
				return
			}
		}
		fmt.Printf(`插件项目已生成: %s
  key: %s
下一步:
  1) cd %s && go mod tidy
  2) 编写插件方法（main.go）
  3) vgo-tools plugin-pack . 生成安装 ZIP
`, dst, key, dst)
		return nil
	},
}

// pluginTemplateFiles 插件项目模板
func pluginTemplateFiles(key, name string) map[string]string {
	return map[string]string{
		"plugin.json": fmt.Sprintf(`{
    "key": "%s",
    "name": "%s",
    "hook": "",
    "singleton": true,
    "version": "1.0.0",
    "description": "插件描述",
    "author": "作者",
    "logo": "assets/logo.svg",
    "readme": "README.md",
    "config": {
        "domain": "https://api.example.com"
    }
}
`, key, name),

		"main.go": `package main

import (
	"fmt"

	"github.com/vera-byte/vgo/v/pluginsdk"
)

// MyPlugin 插件：嵌入 pluginsdk.Base 获得 pluginInfo/cache/pluginService
type MyPlugin struct {
	pluginsdk.Base
}

// Ready 插件就绪回调（可选；singleton 插件仅启动时执行一次）
func (p *MyPlugin) Ready() {
	fmt.Println("plugin ready")
}

// Show 示例业务方法：导出方法即插件 API（宿主按方法名调用，大小写不敏感）
// 返回值支持 (T, error) / (error) / (T) / ()
func (p *MyPlugin) Show(a, b int) (interface{}, error) {
	return map[string]interface{}{
		"plugin": p.PluginInfo.Name,
		"sum":    a + b,
	}, nil
}

// Config 读取插件配置（@baseDir/ 前缀会解析为插件安装目录）
func (p *MyPlugin) Config(key string) (interface{}, error) {
	return p.PluginInfo.Config[key], nil
}

// Upload upload Hook 方法示例：plugin.json 中声明 "hook": "upload" 后，
// 文件上传接口（/admin/base/comm/upload）将由本插件接管。
// ctx 形如 {"files":[{"data":"临时文件路径","filename":"原始名"}],"fields":{...}}
//func (p *MyPlugin) Upload(ctx map[string]interface{}) (interface{}, error) {
//	return "https://cdn.example.com/file", nil
//}

func main() {
	pluginsdk.Run(&MyPlugin{})
}
`,

		"go.mod": fmt.Sprintf(`module %s

go 1.27.3

require github.com/vera-byte/vgo/v v1.1.0
`, key),

		"README.md": fmt.Sprintf(`# %s

## 介绍

插件描述（在此说明插件功能）。

## 标识

调用插件时使用的唯一标识，对应 plugin.json 的 key 字段。

- 标识：**%s**

## 配置

`+"```json"+`
{
    "domain": "https://api.example.com"
}
`+"```"+`

## 方法

- show

`+"```"+`
// 展示插件信息
func (p *MyPlugin) Show(a, b int) (interface{}, error)
`+"```"+`

## 依赖

`+"```bash"+`
go mod tidy    # 获取 pluginsdk
`+"```"+`

若 pluginsdk 尚未发版或需要本地调试，在 go.mod 中追加：

`+"```"+`
replace github.com/vera-byte/vgo/v => ../vgo/v
`+"```"+`

## 打包

`+"```bash"+`
vgo-tools plugin-pack .                    # 全平台
vgo-tools plugin-pack . -p "linux-amd64"    # 指定平台
`+"```"+`

## 安装

`+"```bash"+`
curl -X POST http://<host>:8001/admin/plugin/info/install \
  -H "Authorization: <token>" -F "file=@%s-1.0.0.zip"
`+"```"+`

## 调用示例

`+"```go"+`
// 其他插件或宿主内调用
result, err := service.Invoke(ctx, "%s", "Show", 1, 2)
`+"```"+`
`, name, key, key, key),

		"assets/logo.svg": `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64"><rect width="64" height="64" rx="12" fill="#0D9373"/><text x="32" y="40" font-size="28" text-anchor="middle" fill="#fff">P</text></svg>
`,

	}
}


func init() {
	Main.AddCommand(Plugin)
}
