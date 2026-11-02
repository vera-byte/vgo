# 演示插件（demo-plugin）

vgo 插件开发示例：可执行插件（Go 二进制）+ `pluginsdk` SDK，展示
`pluginInfo` / `cache` / `pluginService`、就绪回调与 Hook 方法的完整用法。

## 介绍

安装后可通过插件管理接口调用本插件导出的方法（导出方法即插件 API）。
本示例覆盖：

- `Ready()` 就绪回调
- 业务方法（`Show` / `Greet` / `Hash` / `Info` / `Config`）
- 宿主缓存读写（`UseCache`）
- 跨插件调用（`CallOther`）
- upload Hook 方法示例（默认未启用）

## 标识

调用插件时使用的唯一标识，对应 `plugin.json` 的 `key` 字段。

- 标识：**demo**

## 配置

```json
{
    "greeting": "hello from plugin",
    "greetingFile": "@baseDir/README.md",
    "domain": "https://api.example.com"
}
```

`@baseDir/` 前缀会在加载时解析为插件安装目录的绝对路径。

## 方法

- show

  ```go
  // 展示插件信息与两数之和
  func (p *DemoPlugin) Show(a, b int) (interface{}, error)
  ```

- info / config

  ```go
  func (p *DemoPlugin) Info() (interface{}, error)
  func (p *DemoPlugin) Config(key string) (interface{}, error)
  ```

- useCache

  ```go
  // 读写宿主缓存（Redis）
  func (p *DemoPlugin) UseCache(key, value string) (interface{}, error)
  ```

- callOther

  ```go
  // 跨插件调用
  func (p *DemoPlugin) CallOther(key, method string) (interface{}, error)
  ```

- hash / greet

  ```go
  func (p *DemoPlugin) Hash(text string) (string, error)
  func (p *DemoPlugin) Greet(name string) (string, error)
  ```

## 打包

```bash
# 在仓库根目录执行（vgo-tools）
go run ./vgo-tools plugin-pack . -o demo-1.0.0.zip
# 或在本目录（已安装 vgo-tools）
vgo-tools plugin-pack . -p "linux-amd64 darwin-arm64"
```

产物 ZIP 布局：`plugin.json` / `README.md` / `assets/` / `bin/<goos>-<goarch>/plugin`。

## 安装

```bash
curl -X POST http://<host>:8001/admin/plugin/info/install \
  -H "Authorization: <token>" -F "file=@demo-1.0.0.zip"
```

安装即启用（进程拉起 + Ready 握手）。启停/卸载/热重载：
`/admin/plugin/info/update|delete|reload`。

## 启用 upload Hook（可选）

修改 `plugin.json` 为 `"hook": "upload"`，并取消 `main.go` 中 `Upload` 方法注释后重新打包，
文件上传接口（`/admin/base/comm/upload`）将由本插件接管（同 Hook 同时只能启用一个插件）。

## 调用示例

```go
// 其他插件内
result, err := p.PluginService.Invoke("demo", "Show", 1, 2)

// 宿主业务代码内
import pluginService "github.com/vera-byte/vgo/modules/plugin/service"
result, err := pluginService.Invoke(ctx, "demo", "Greet", "vgo")
```

## 本地开发

`go.mod` 中 `replace github.com/vera-byte/vgo/v => ../v` 指向本仓库的 v 模块；
若作为独立项目发布，改为 `go get github.com/vera-byte/vgo/v@latest` 并移除 replace。
