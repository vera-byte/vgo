package v1

import "github.com/gogf/gf/v2/frame/g"

// PluginInfoInstallReq 安装插件请求（ZIP 包上传）
type PluginInfoInstallReq struct {
	g.Meta `path:"/install" method:"POST" summary:"安装插件" tags:"插件管理"`
	Force  bool `json:"force"` // 强制安装（跳过检查）
}

// PluginInfoReloadReq 重载插件请求
type PluginInfoReloadReq struct {
	g.Meta `path:"/reload" method:"POST" summary:"重载插件" tags:"插件管理"`
	Key    string `json:"key" v:"required#插件key不能为空"`
}
