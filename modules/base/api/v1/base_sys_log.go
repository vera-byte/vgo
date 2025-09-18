package v1

import "github.com/gogf/gf/v2/frame/g"

// SetKeepReq 设置保留天数请求结构
type SetKeepReq struct {
	g.Meta `method:"POST" path:"/setKeep" summary:"设置保留天数" tags:"系统日志"`
	Value  int `json:"value" v:"required#请输入保留天数"`
}

// GetKeepReq 获取日志保留天数请求参数
type GetKeepReq struct {
	g.Meta        `path:"/getKeep" method:"GET" summary:"获取日志保留天数" tags:"系统日志"`
	Authorization string `json:"Authorization" in:"header"`
}

// ClearReq 清理日志请求参数
type ClearReq struct {
	g.Meta        `path:"/clear" method:"POST" summary:"清理系统日志" tags:"系统日志"`
	Authorization string `json:"Authorization" in:"header"`
}