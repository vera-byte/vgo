package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	vgo_request "github.com/vera-byte/vgo/vgo_core_kit/request"
)

type LogPageReq struct {
	g.Meta `path:"log/page" method:"post" sm:"日志分页" tags:"系统"`
}
type LogPageRes struct {
}

type DepartmentListReq struct {
	g.Meta `path:"department/list" method:"post" sm:"部门列表" tags:"系统"`
}
type DepartmentListRes []DepartmentResultItem
type DepartmentResultItem struct {
	Id         int64  `json:"id" sm:"ID" dc:"部门ID"`
	Name       string `json:"name" sm:"部门名称" dc:"部门名称"`
	OrderNum   int64  `json:"orderNum" sm:"排序" dc:"排序"`
	ParentId   int64  `json:"parentId" sm:"父ID" dc:"父ID"`
	TenantId   int64  `json:"tenantId" sm:"租户ID" dc:"租户ID"`
	UpdateTime int64  `json:"updateTime" sm:"更新时间" dc:"更新时间"`
	CreateTime int64  `json:"createTime" sm:"创建时间" dc:"创建时间"`
}

type UserPageReq struct {
	g.Meta        `path:"user/page" method:"post" sm:"用户分页" tags:"用户"`
	DepartmentIds []int64 `json:"departmentIds" sm:"部门ID集合" dc:"部门ID"`
	vgo_request.PageReq
}
type UserPageRes struct {
	Pagination vgo_request.Pagination `json:"pagination" sm:"分页信息" dc:"分页信息"`
	List       []UserItem             `json:"list" sm:"用户列表" dc:"用户列表"`
}

type UserItem struct {
	Id             int64   `json:"id" sm:"用户ID" dc:"用户ID"`
	UpdateTime     int64   `json:"updateTime" sm:"更新时间" dc:"更新时间"`
	CreateTime     int64   `json:"createTime" sm:"创建时间" dc:"创建时间"`
	DepartmentId   int64   `json:"departmentId" sm:"部门ID" dc:"部门ID"`
	DepartmentName string  `json:"departmentName" sm:"部门名称" dc:"部门名称"`
	Email          string  `json:"email" sm:"邮箱" dc:"邮箱"`
	HeadImg        string  `json:"headImg" sm:"头像" dc:"头像"`
	Name           string  `json:"name" sm:"名称" dc:"名称"`
	NickName       string  `json:"nickName" sm:"昵称" dc:"昵称"`
	Phone          string  `json:"phone" sm:"手机号" dc:"手机号"`
	Remark         string  `json:"remark" sm:"备注" dc:"备注"`
	RoleIds        []int64 `json:"roleIds" sm:"身份" dc:"身份"`
	RoleName       string  `json:"roleName" sm:"身份名称" dc:"身份名称"`
	Status         int     `json:"status" sm:"状态" dc:"状态"`
	Username       string  `json:"username" sm:"用户名" dc:"用户名"`
}
