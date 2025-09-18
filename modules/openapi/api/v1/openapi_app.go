package v1

import "github.com/gogf/gf/v2/frame/g"

// GenerateSignReq 生成签名请求结构
type GenerateSignReq struct {
	g.Meta    `path:"/generate" method:"POST" summary:"生成签名" tags:"开放平台应用"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
}

// AppVerifySignReq 验证签名请求结构
type AppVerifySignReq struct {
	g.Meta    `path:"/verify" method:"POST" summary:"验证签名" tags:"开放平台应用"`
	AppId     string      `json:"appId" v:"required#应用ID不能为空"`
	Timestamp int64       `json:"timestamp" v:"required#时间戳不能为空"`
	Nonce     string      `json:"nonce" v:"required#随机数不能为空"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature" v:"required#签名不能为空"`
}

// GetPublicKeyReq 获取公钥请求结构
type GetPublicKeyReq struct {
	g.Meta `path:"/publicKey" method:"GET" summary:"获取公钥" tags:"开放平台应用"`
	AppId  string `json:"appId" v:"required#应用ID不能为空"`
}