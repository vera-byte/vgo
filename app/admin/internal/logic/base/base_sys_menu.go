package base

import (
	"context"
	"vgo/app/admin/internal/dao"
	"vgo/app/admin/internal/model/entity"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterBaseSysMenuLogic(NewBaseSysMenuLogic())

}

type sBaseSysMenuLogic struct {
}

func NewBaseSysMenuLogic() *sBaseSysMenuLogic {
	return &sBaseSysMenuLogic{}
}

// GetPerms 获取菜单的权限
func (s *sBaseSysMenuLogic) GetPerms(ctx context.Context, roleIds []string) []string {
	var (
		perms  []string
		result gdb.Result
	)
	m := dao.BaseSysMenu.Ctx(ctx).As("a")
	// 如果roldIds 包含 1 则表示是超级管理员，则返回所有权限
	if garray.NewIntArrayFrom(gconv.Ints(roleIds)).Contains(1) {
		result, _ = m.Fields("a.perms").All()
	} else {
		result, _ = m.InnerJoin("base_sys_role_menu b", "a.id=b.menu_id").InnerJoin("base_sys_role c", "b.role_id=c.id").Where("c.id IN (?)", roleIds).Fields("a.perms").All()
	}
	for _, v := range result {
		vmap := v.Map()
		if vmap["perms"] != nil {
			p := gstr.Split(vmap["perms"].(string), ",")
			perms = append(perms, p...)
		}
	}
	return perms
}

// GetMenus 获取菜单
func (s *sBaseSysMenuLogic) GetMenus(ctx context.Context, roleIds []string, isAdmin bool) []entity.BaseSysMenu {
	var (
		menus = &[]entity.BaseSysMenu{}
	)
	// 屏蔽 base_sys_role_menu.id 防止部分权限的用户登录时菜单渲染错误
	m := dao.BaseSysMenu.Ctx(ctx).As("a").Fields("a.*")
	if isAdmin {
		_ = m.Group("a.id").Order("a.order_num asc").Scan(menus)
	} else {
		_ = m.InnerJoin("base_sys_role_menu b", "a.id=b.menu_id").Where("b.role_id IN (?)", roleIds).Group("a.id").Order("a.order_num asc").Scan(menus)
	}
	return *menus

}
