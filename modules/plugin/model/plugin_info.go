package model

import "github.com/vera-byte/vgo/v"

const TableNamePluginInfo = "plugin_info"

// PluginInfo 插件信息（参照 cool-admin-midway plugin 模块 PluginInfoEntity）
type PluginInfo struct {
	*v.Model
	Name        string `json:"name"`        // 名称
	Description string `json:"description"` // 简介
	KeyName     string `json:"keyName"`     // Key名（唯一标识）
	Hook        string `json:"hook"`        // Hook 名（如 upload），同 Hook 同时只能启用一个
	Readme      string `json:"readme"`      // 使用说明（markdown）
	Version     string `json:"version"`     // 版本
	Logo        string `json:"logo"`        // Logo（base64）
	Author      string `json:"author"`      // 作者
	Status      int    `json:"status"`      // 状态 0-禁用 1-启用
	Content     string `json:"content"`     // 插件内容（JSON：{type,data}，data 为可执行脚本）
	TsContent   string `json:"tsContent"`   // TS 源码（JSON：{type:'ts',data}，留档参考）
	PluginJson  string `json:"pluginJson"`  // 插件的 plugin.json（JSON）
	Config      string `json:"config"`      // 配置（JSON，支持 @env 多环境）
}

// TableName PluginInfo 的表名
func (*PluginInfo) TableName() string {
	return TableNamePluginInfo
}

// NewPluginInfo 创建一个 PluginInfo 模型实例
func NewPluginInfo() *PluginInfo {
	return &PluginInfo{
		Model: v.NewModel(),
	}
}
