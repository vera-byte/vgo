package v1

import "github.com/gogf/gf/v2/frame/g"

// RecycleDataRestoreReq 恢复数据请求
type RecycleDataRestoreReq struct {
	g.Meta `path:"/restore" method:"POST" summary:"恢复数据" tags:"数据回收"`
	Ids    []int `json:"ids" v:"required#请选择要恢复的记录"`
}

// RecycleDataClearReq 清理回收站请求
type RecycleDataClearReq struct {
	g.Meta `path:"/clear" method:"POST" summary:"清理回收站" tags:"数据回收"`
}
