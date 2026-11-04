# Vgo

[![CI](https://github.com/vera-byte/vgo/actions/workflows/ci.yml/badge.svg)](https://github.com/vera-byte/vgo/actions/workflows/ci.yml)
[![Build Release](https://github.com/vera-byte/vgo/actions/workflows/release.yml/badge.svg)](https://github.com/vera-byte/vgo/actions/workflows/release.yml)

vAdmin 的 Go 版本后端——基于 GoFrame v2 的后台管理系统，开箱自带 RBAC 权限、数据字典、任务调度、数据回收、开放平台签名与可执行插件体系。

- [文档](https://github.com/vera-byte/vgo/tree/master/docs) — 使用指南与参考（Mintlify，`cd docs && mint dev` 本地预览）
- [更新日志](docs/changelog.mdx) · [发布流程](docs/guides/release.mdx) · [已知问题](docs/known-issues.mdx)
- [Releases](https://github.com/vera-byte/vgo/releases) — vgo-tools 多平台二进制下载
- [Issues](https://github.com/vera-byte/vgo/issues) — 问答与反馈

## 技术栈

- **Go 1.27.3** · **GoFrame v2.10.4**（多模块 go.work 工作区）
- **sa-token-go** 会话认证（Random64 服务端会话、顶号/刷新轮换/改密下线）
- **PostgreSQL**（迁移命令 PG-only）· **Redis**（会话/缓存，可降级内存）
- 密钥经 **`VGO_*` 环境变量**注入（`.env` 自动加载，严禁入库）

## 内置模块

| 模块 | 能力 |
|------|------|
| base | 用户/角色/部门/菜单/参数、RBAC、登录认证、文件上传（local/minio/oss）、系统日志 |
| dict | 数据字典 |
| task | 任务调度 |
| recycle | 数据回收站（删除快照/恢复/定期清理） |
| plugin | 可执行插件体系（ZIP 安装、Hook 接管上传、热重载） |
| user | C 端用户体系（手机/密码/微信登录、资料、收货地址） |
| openapi | 开放平台应用与签名（防重放） |
| space | 空间管理 |

插件为 **Go 二进制可执行程序**（`pluginsdk` + stdio JSON-RPC），`vgo-tools plugin` 一键脚手架、`vgo-tools plugin-pack` 多平台打包，详见[插件开发指南](docs/guides/plugin-dev.mdx)。

## 快速开始

### 1. 安装 vgo-tools

```bash
wget -O vgo-tools https://github.com/vera-byte/vgo/releases/latest/download/vgo-tools_$(go env GOOS)_$(go env GOARCH) \
  && chmod +x vgo-tools && ./vgo-tools install && rm ./vgo-tools
```

下载慢可用 `pgit` 加速（`pgit wget ...`），或 `go install github.com/vera-byte/vgo/vgo-tools@latest`。验证：`vgo-tools version`。

### 2. 创建项目并配置

```bash
vgo-tools init hello-world && cd hello-world
cp .env.example .env            # 填数据库/Redis 凭证（VGO_* 变量）
```

### 3. 建库与启动

```bash
go mod tidy
go run main.go migrate          # 迁移全部模块（含内嵌迁移）
go run main.go server           # 启动 Web 服务（默认 :8001）
```

命令树：`server` / `migrate` / `pack` / `utils` + 各模块管理命令（`go run main.go` 查看）。

## 本仓库开发

```bash
git clone git@github.com:vera-byte/vgo.git && cd vgo
cp .env.example .env && cp manifest/config/config.example.yaml manifest/config/config.yaml
make test && make vet           # 门禁
go run . server                 # 或 make dev
```

多模块：根与各子模块 `go.mod` 均 `replace` 指向本地源码，克隆即可构建；`gowork.sh` 可生成 `go.work` 便于跨模块联调。


## License

见 [LICENSE](LICENSE)。
