package model

import "github.com/vera-byte/vgo/v"

const TableNameBaseSysLog = "base_sys_log"

// BaseSysLog mapped from table <base_sys_log>
type BaseSysLog struct {
	*v.Model
	UserID uint   `json:"userId"` // 用户ID
	Action string `json:"action"` // 行为
	IP     string `json:"ip"`     // ip
	IPAddr string `json:"ipAddr"` // ip地址
	Params string `json:"params"` // 参数
}

// TableName BaseSysLog's table name
func (*BaseSysLog) TableName() string {
	return TableNameBaseSysLog
}

// init 创建表
func init() {
	v.CreateTable(&BaseSysLog{})
}

// NewBaseSysLog 创建实例
func NewBaseSysLog() *BaseSysLog {
	return &BaseSysLog{
		Model: v.NewModel(),
	}
}
