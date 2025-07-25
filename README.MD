# Vgo

[![Build Release](https://github.com/vera-byte/vgo/actions/workflows/release.yml/badge.svg)](https://github.com/vera-byte/vgo/actions/workflows/release.yml)

vAdmin 的 Go 版本后端

- [docs](https://v-team-official.github.io/vgo/) - 开发文档
- [changelog](docs/changelog.md) - 更新日志
- [github](https://github.com/vera-byte/vgo) - vgo 仓库
- [release](https://github.com/vera-byte/vgo/releases) - vgo-tools 编译版本下载
- [issues](https://github.com/vera-byte/vgo/issues) - 问答区
- [vAdmin](https://v-js.com) - vAdmin 官网
- [GoFrame](https://goframe.org) - GoFrame 官网

## TODO

- [x] 基础 CURD 结构搭建
- [x] model 层快速生成
- [x] service 层快速生成
- [x] controller 层快速生成
- [x] 中间件错误状态及错误码与 node 版本统一
- [x] 全局错误码统一
- [x] controller 层接口输入与 node 版本统一
- [x] 优化开发工具安装体验
- [x] service 层逻辑补全
- [x] base 模块逻辑补全 感谢 @imhuso
- [ ] 开发文档编制
- [ ] 数据库驱动补充
- [ ] 文件上传驱动补充
- [ ] 其他项目需求

## 使用说明

首先下载安装 Go 语言环境，如果没有安装，请自行安装，如果已经安装，请自行配置环境变量

### 1.安装`vgo-tools`开发工具

从 github 下载

```bash
wget -O vgo-tools https://github.com/vera-byte/vgo/releases/latest/download/vgo-tools_$(go env GOOS)_$(go env GOARCH) && chmod +x vgo-tools && ./vgo-tools install  && rm ./vgo-tools
```

如果到直接下载速度过慢，可使用`pgit`加速安装

```bash
pgit wget -O vgo-tools https://github.com/vera-byte/vgo/releases/latest/download/vgo-tools_$(go env GOOS)_$(go env GOARCH) && chmod +x vgo-tools && ./vgo-tools install  && rm ./vgo-tools
```

通过`go install`安装

```bash
go install github.com/vera-byte/vgo/vgo-tools@latest
```

验证

```bash
vgo-tools version
```

### 2.创建新项目

```bash
vgo-tools init hello-world
```

### 3.安装依赖

```bash
cd hellow-world
go mod tidy
```

### 5.启动开发环境

```bash
vgo-tools run main.go
```


enjoy it!
