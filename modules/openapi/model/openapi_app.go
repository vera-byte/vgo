package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameOpenapiApp = "openapi_app"

// OpenapiApp 开放平台应用模型，用于存储应用信息和RSA密钥对
type OpenapiApp struct {
	*v.Model
	AppId       string `json:"appId"`
	AppName     string `json:"appName"`
	AppSecret   string `json:"appSecret"`
	PublicKey   string `json:"publicKey"`
	PrivateKey  string `json:"privateKey"`
	Status      *int32 `json:"status"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

// TableName OpenapiApp's table name
// 功能: 返回OpenapiApp模型对应的数据库表名
// 返回值: string - 数据库表名
func (*OpenapiApp) TableName() string {
	return TableNameOpenapiApp
}

// GroupName OpenapiApp's table group
// 功能: 返回OpenapiApp模型所属的数据库组名
// 返回值: string - 数据库组名
func (*OpenapiApp) GroupName() string {
	return "default"
}

// NewOpenapiApp 创建一个新的OpenapiApp实例
// 功能: 创建并初始化一个新的开放平台应用模型实例
// 返回值: *OpenapiApp - 新创建的OpenapiApp实例
func NewOpenapiApp() *OpenapiApp {
	return &OpenapiApp{
		Model: v.NewModel(),
	}
}
