package base

import (
	"context"
	"vgo/app/admin/internal/dao"
	"vgo/app/admin/internal/model/entity"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/container/gset"
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

func (s *sBaseSysMenuLogic) GetPerms(ctx context.Context, roleIds []string) []string {
	var perms []string

	// 转换 roleIds 为整数集合
	roleIdSet := gset.NewIntSetFrom(gconv.Ints(roleIds))

	// 定义存储权限的结构体
	var menus []entity.BaseSysMenu

	// 获取表名和字段
	menuTable := dao.BaseSysMenu.Table()
	roleMenuTable := dao.BaseSysRoleMenu.Table()
	roleTable := dao.BaseSysRole.Table()

	menuIdField := dao.BaseSysMenu.Columns().Id
	menuPermsField := dao.BaseSysMenu.Columns().Perms
	roleIdField := dao.BaseSysRole.Columns().Id
	roleMenuRoleIdField := dao.BaseSysRoleMenu.Columns().RoleId
	roleMenuMenuIdField := dao.BaseSysRoleMenu.Columns().MenuId

	// 查询数据库
	if roleIdSet.Contains(1) {
		// 超级管理员获取所有权限
		err := dao.BaseSysMenu.Ctx(ctx).
			Fields(menuPermsField).
			Scan(&menus)
		if err != nil {
			return nil
		}
	} else {
		// 普通角色权限查询
		err := dao.BaseSysMenu.Ctx(ctx).
			InnerJoin(roleMenuTable, menuTable+"."+menuIdField+"="+roleMenuTable+"."+roleMenuMenuIdField).
			InnerJoin(roleTable, roleMenuTable+"."+roleMenuRoleIdField+"="+roleTable+"."+roleIdField).
			Where(roleTable+"."+roleIdField+" IN(?)", roleIds).
			Fields(menuPermsField).
			Scan(&menus)
		if err != nil {
			return nil
		}
	}

	// 解析权限
	for _, menu := range menus {
		if menu.Perms != "" {
			perms = append(perms, gstr.Split(menu.Perms, ",")...)
		}
	}

	return perms
}

// GetMenus 获取指定角色的菜单列表
// 如果 isAdmin 为 true，则返回所有菜单（超级管理员权限）
// 如果 isAdmin 为 false，则根据角色 ID 查询对应的菜单
func (s *sBaseSysMenuLogic) GetMenus(ctx context.Context, roleIds []string, isAdmin bool) []entity.BaseSysMenu {
	// 定义返回的菜单列表
	var menus []entity.BaseSysMenu

	// 获取数据库表名
	menuTable := dao.BaseSysMenu.Table()
	roleMenuTable := dao.BaseSysRoleMenu.Table()

	// 获取字段名
	menuIdField := dao.BaseSysMenu.Columns().Id
	menuOrderNumField := dao.BaseSysMenu.Columns().OrderNum
	roleMenuRoleIdField := dao.BaseSysRoleMenu.Columns().RoleId
	roleMenuMenuIdField := dao.BaseSysRoleMenu.Columns().MenuId

	// 查询构造器
	m := dao.BaseSysMenu.Ctx(ctx).
		Fields(menuTable + ".*") // 查询所有菜单字段

	if isAdmin {
		// 如果是超级管理员，获取所有菜单
		_ = m.Group(menuIdField).
			Order(menuOrderNumField + " ASC").
			Scan(&menus)
	} else {
		// 普通用户通过角色 ID 关联查询菜单
		_ = m.InnerJoin(roleMenuTable, menuTable+"."+menuIdField+"="+roleMenuTable+"."+roleMenuMenuIdField).
			Where(roleMenuTable+"."+roleMenuRoleIdField+" IN(?)", roleIds).
			Group(menuIdField).
			Order(menuOrderNumField + " ASC").
			Scan(&menus)
	}

	return menus
}
