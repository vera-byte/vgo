package v1

import "github.com/gogf/gf/v2/frame/g"

// BaseSysParamHtmlReq 系统参数HTML请求参数
type BaseSysParamHtmlReq struct {
	g.Meta        `path:"/html" method:"GET" summary:"获取系统参数HTML" tags:"系统参数"`
	Authorization string `json:"Authorization" in:"header"`
	Key           string `v:"required#请输入key"`
}