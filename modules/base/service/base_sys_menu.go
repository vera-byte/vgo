package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseSysMenuService struct {
	*v.Service
}

// GetPerms 获取菜单的权限（isAdmin 为 true 时返回全部权限）
func (s *BaseSysMenuService) GetPerms(roleIds []string, isAdmin bool) (perms []string, err error) {
	var (
		result gdb.Result
	)
	m := v.DBM(s.Model).As("a")
	group := s.Model.GroupName()
	if isAdmin {
		result, err = m.Fields("a.perms").All()
	} else {
		result, err = m.InnerJoin("base_sys_role_menu b", v.QuoteIdent(group, "a.id")+"="+v.QuoteIdent(group, "b.menuId")).
			InnerJoin("base_sys_role c", v.QuoteIdent(group, "b.roleId")+"="+v.QuoteIdent(group, "c.id")).
			WhereIn(v.QuoteIdent(group, "c.id"), roleIds).Fields("a.perms").All()
	}
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		vmap := v.Map()
		if vmap["perms"] != nil {
			p := gstr.Split(vmap["perms"].(string), ",")
			perms = append(perms, p...)
		}
	}
	return perms, nil
}

// GetMenus 获取菜单
func (s *BaseSysMenuService) GetMenus(roleIds []string, isAdmin bool) (result gdb.Result, err error) {
	// 屏蔽 base_sys_role_menu.id 防止部分权限的用户登录时菜单渲染错误
	m := v.DBM(s.Model).As("a").Fields("a.*")
	group := s.Model.GroupName()
	if isAdmin {
		result, err = m.Group("a.id").Order("a.orderNum asc").All()
	} else {
		result, err = m.InnerJoin("base_sys_role_menu b", v.QuoteIdent(group, "a.id")+"="+v.QuoteIdent(group, "b.menuId")).
			WhereIn(v.QuoteIdent(group, "b.roleId"), roleIds).Group("a.id").Order("a.orderNum asc").All()
	}
	return
}

// ModifyAfter 修改后
func (s *BaseSysMenuService) ModifyAfter(ctx context.Context, method string, param g.MapStrAny) (err error) {
	if method == "Delete" {
		ids := gconv.Ints(param["ids"])
		if len(ids) > 0 {
			_, err = v.DBM(s.Model).WhereIn("parentId", ids).Delete()
		}
		return
	}
	return
}

// ServiceAdd 添加
func (s *BaseSysMenuService) ServiceAdd(ctx context.Context, req *v.AddReq) (data interface{}, err error) {
	r := g.RequestFromCtx(ctx)
	rjson, err := r.GetJson()
	if err != nil {
		return
	}
	m := v.DBM(s.Model)
	lastInsertId, err := m.Data(rjson).InsertAndGetId()
	if err != nil {
		return
	}
	data = g.Map{"id": lastInsertId}
	return
}

// NewBaseSysMenuService 创建一个BaseSysMenuService实例
func NewBaseSysMenuService() *BaseSysMenuService {
	return &BaseSysMenuService{
		&v.Service{
			Model: model.NewBaseSysMenu(),
		},
	}
}
