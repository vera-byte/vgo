package v1

import "github.com/gogf/gf/v2/frame/g"

type PersonReq struct {
	g.Meta `path:"person" method:"get" sm:"个人用户信息" tags:"鉴权"`
}

type PersonRes struct {
}

type PermmenuReq struct {
	g.Meta `path:"permmenu" method:"get" sm:"权限菜单" tags:"鉴权"`
}
type Menu struct {
	CreateTime string `json:"createTime" dc:"创建时间"`
	UpdateTime string `json:"updateTime" dc:"更新时间"`
	Icon       string `json:"icon" dc:"图标"`
	Id         int64  `json:"id" dc:"ID"`
	IsShow     bool   `json:"isShow" dc:"是否显示"`
	KeepAlive  bool   `json:"keepAlive" dc:"是否缓存"`
	Name       string `json:"name" dc:"名称"`
	OrderNum   int64  `json:"orderNum" dc:"排序"`
	ParentId   int64  `json:"parentId" dc:"父ID"`
	Perms      string `json:"perms" dc:"权限"`
	Router     string `json:"router" dc:"路由"`
	TenantId   int64  `json:"tenantId" dc:"租户ID"`
	Type       int    `json:"type" dc:"类型"`
	ViewPath   string `json:"viewPath" dc:"路径"`
}
type PermmenuRes struct {
	Menus []Menu   `json:"menus" dc:"菜单"`
	Perms []string `json:"perms"  dc:"权限"`
}
