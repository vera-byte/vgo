package v1

import "github.com/gogf/gf/v2/frame/g"

// CreateAppReq 创建应用请求结构
type CreateAppReq struct {
	g.Meta      `path:"/createApp" method:"POST" summary:"创建应用" tags:"开放平台管理"`
	AppName     string `json:"appName" v:"required#应用名称不能为空"`
	Description string `json:"description"`
	KeyBits     int    `json:"keyBits" v:"min:2048#密钥位数不能小于2048"`
}

// GetAppInfoReq 获取应用信息请求结构
type GetAppInfoReq struct {
	g.Meta `path:"/getAppInfo" method:"GET" summary:"获取应用信息" tags:"开放平台管理"`
	AppId  string `json:"appId" v:"required#应用ID不能为空"`
}

// TestSignReq 测试签名请求结构
type TestSignReq struct {
	g.Meta    `path:"/testSign" method:"POST" summary:"测试签名" tags:"开放平台管理"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
}

// VerifySignReq 验证签名请求结构
type VerifySignReq struct {
	g.Meta    `path:"/verifySign" method:"POST" summary:"验证签名" tags:"开放平台管理"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature" v:"required#签名不能为空"`
}