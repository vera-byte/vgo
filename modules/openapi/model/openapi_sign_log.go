package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameOpenapiSignLog = "openapi_sign_log"

// OpenapiSignLog 开放平台签名日志模型，用于记录签名请求的详细信息
type OpenapiSignLog struct {
	*v.Model
	AppId       string `gorm:"column:app_id;type:varchar(64);not null;index;comment:应用ID" json:"appId"`
	RequestId   string `gorm:"column:request_id;type:varchar(64);not null;uniqueIndex;comment:请求ID" json:"requestId"`
	Timestamp   int64  `gorm:"column:timestamp;type:bigint;not null;comment:时间戳" json:"timestamp"`
	Nonce       string `gorm:"column:nonce;type:varchar(64);not null;comment:随机数" json:"nonce"`
	RequestBody string `gorm:"column:request_body;type:text;comment:请求体" json:"requestBody"`
	Signature   string `gorm:"column:signature;type:text;not null;comment:生成的签名" json:"signature"`
	ClientIp    string `gorm:"column:client_ip;type:varchar(45);comment:客户端IP" json:"clientIp"`
	UserAgent   string `gorm:"column:user_agent;type:varchar(500);comment:用户代理" json:"userAgent"`
	Status      *int32 `gorm:"column:status;type:int;not null;default:1;comment:状态 0:失败 1:成功" json:"status"`
	ErrorMsg    string `gorm:"column:error_msg;type:varchar(500);comment:错误信息" json:"errorMsg"`
}

// TableName OpenapiSignLog's table name
// 功能: 返回OpenapiSignLog模型对应的数据库表名
// 返回值: string - 数据库表名
func (*OpenapiSignLog) TableName() string {
	return TableNameOpenapiSignLog
}

// GroupName OpenapiSignLog's table group
// 功能: 返回OpenapiSignLog模型所属的数据库组名
// 返回值: string - 数据库组名
func (*OpenapiSignLog) GroupName() string {
	return "default"
}

// NewOpenapiSignLog 创建一个新的OpenapiSignLog实例
// 功能: 创建并初始化一个新的签名日志模型实例
// 返回值: *OpenapiSignLog - 新创建的OpenapiSignLog实例
func NewOpenapiSignLog() *OpenapiSignLog {
	return &OpenapiSignLog{
		Model: v.NewModel(),
	}
}
