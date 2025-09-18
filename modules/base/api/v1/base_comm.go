package v1

import "github.com/gogf/gf/v2/frame/g"

// BaseCommPersonReq 获取个人信息请求参数
type BaseCommPersonReq struct {
	g.Meta        `path:"/person" method:"GET" summary:"获取个人信息" tags:"通用接口"`
	Authorization string `json:"Authorization" in:"header"`
}

// BaseCommPermmenuReq 获取权限菜单请求参数
type BaseCommPermmenuReq struct {
	g.Meta        `path:"/permmenu" method:"GET" summary:"获取权限菜单" tags:"通用接口"`
	Authorization string `json:"Authorization" in:"header"`
}

// BaseCommLogoutReq 退出登录请求参数
type BaseCommLogoutReq struct {
	g.Meta        `path:"/logout" method:"POST" summary:"退出登录" tags:"通用接口"`
	Authorization string `json:"Authorization" in:"header"`
}

// BaseCommUploadModeReq 获取上传模式请求参数
type BaseCommUploadModeReq struct {
	g.Meta        `path:"/uploadMode" method:"GET" summary:"获取上传模式" tags:"通用接口"`
	Authorization string `json:"Authorization" in:"header"`
}

// BaseCommUploadReq 文件上传请求参数
type BaseCommUploadReq struct {
	g.Meta        `path:"/upload" method:"POST" summary:"文件上传" tags:"通用接口"`
	Authorization string `json:"Authorization" in:"header"`
}

// PersonUpdateReq 更新个人信息请求参数
type PersonUpdateReq struct {
	g.Meta        `path:"/personUpdate" method:"POST" summary:"更新个人信息" tags:"通用接口"`
	Authorization string `json:"Authorization" in:"header"`
}

// BaseCommControllerEpsReq EPS接口请求参数
type BaseCommControllerEpsReq struct {
	g.Meta `path:"/eps" method:"GET" summary:"获取控制器EPS信息" tags:"通用接口"`
}