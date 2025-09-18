package v1

import "github.com/gogf/gf/v2/frame/g"

// DictInfoDataReq 获取字典数据请求参数
type DictInfoDataReq struct {
	g.Meta `path:"/data" method:"POST" summary:"获取字典数据" tags:"字典管理"`
	Types  []string `json:"types"`
}