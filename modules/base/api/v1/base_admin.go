package v1

import "github.com/gogf/gf/v2/frame/g"

// BaseOpenEpsReq EPS接口请求参数
type BaseOpenEpsReq struct {
	g.Meta `path:"/eps" method:"GET" summary:"获取EPS接口信息" tags:"管理员"`
}

// RefreshTokenReq 刷新token请求参数
type RefreshTokenReq struct {
	g.Meta       `path:"/refreshToken" method:"GET" summary:"刷新访问令牌" tags:"管理员"`
	RefreshToken string `json:"refreshToken" v:"required#refreshToken不能为空"`
}

// UserMoveReq 用户移动请求参数
type UserMoveReq struct {
	g.Meta        `path:"/move" method:"POST" summary:"移动用户" tags:"管理员"`
	Authorization string `json:"Authorization" in:"header"`
}

// OrderReq 排序请求参数
type OrderReq struct {
	g.Meta        `path:"/order" method:"POST" summary:"排序操作" tags:"管理员"`
	Authorization string `json:"Authorization" in:"header"`
}