package v1

import "github.com/gogf/gf/v2/frame/g"

// TaskInfoStopReq 停止任务请求结构
type TaskInfoStopReq struct {
	g.Meta `path:"/stop" method:"GET" summary:"停止任务" tags:"任务管理"`
	ID     int64 `json:"id" v:"required#请输入id"`
}

// TaskInfoStartReq 启动任务请求结构
type TaskInfoStartReq struct {
	g.Meta `path:"/start" method:"GET" summary:"启动任务" tags:"任务管理"`
	ID     int64 `json:"id" v:"required#请输入id"`
}

// TaskInfoOnceReq 执行一次任务请求结构
type TaskInfoOnceReq struct {
	g.Meta `path:"/once" method:"POST" summary:"执行一次任务" tags:"任务管理"`
	ID     int64 `json:"id" v:"required#请输入id"`
}

// TaskInfoLogReq 任务日志请求结构
type TaskInfoLogReq struct {
	g.Meta `path:"/log" method:"GET" summary:"获取任务日志" tags:"任务管理"`
	ID     int64 `json:"id"`
	Status int   `json:"status"`
}