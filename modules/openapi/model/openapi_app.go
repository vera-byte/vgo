package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameOpenapiApp = "openapi_app"

// OpenapiApp 开放平台应用模型，用于存储应用信息和RSA密钥对
type OpenapiApp struct {
	*v.Model
	AppId       string `gorm:"column:app_id;type:varchar(64);not null;uniqueIndex;comment:应用ID" json:"appId"`
	AppName     string `gorm:"column:app_name;type:varchar(255);not null;comment:应用名称" json:"appName"`
	AppSecret   string `gorm:"column:app_secret;type:varchar(255);not null;comment:应用密钥" json:"appSecret"`
	PublicKey   string `gorm:"column:public_key;type:text;not null;comment:RSA公钥" json:"publicKey"`
	PrivateKey  string `gorm:"column:private_key;type:text;not null;comment:RSA私钥" json:"privateKey"`
	Status      *int32 `gorm:"column:status;type:int;not null;default:1;comment:状态 0:禁用 1:启用" json:"status"`
	Description string `gorm:"column:description;type:varchar(500);comment:应用描述" json:"description"`
	Remark      string `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
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

// init 创建表
// 功能: 模块初始化时自动创建OpenapiApp数据库表
func init() {
	v.CreateTable(&OpenapiApp{})
}