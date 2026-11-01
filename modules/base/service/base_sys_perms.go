package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/vera-byte/vgo/v"
)

type BaseSysPermsService struct {
}

func NewBaseSysPermsService() *BaseSysPermsService {
	return &BaseSysPermsService{}
}

// permmenu 方法
func (c *BaseSysPermsService) Permmenu(ctx context.Context, roleIds []string) (res interface{}, err error) {
	type permmenu struct {
		Perms []string   `json:"perms"`
		Menus gdb.Result `json:"menus"`
	}
	var (
		baseSysMenuService = NewBaseSysMenuService()
		admin              = v.GetAdmin(ctx)
	)
	perms, err := baseSysMenuService.GetPerms(roleIds, admin.IsAdmin)
	if err != nil {
		return nil, err
	}
	menus, err := baseSysMenuService.GetMenus(admin.RoleIds, admin.IsAdmin)
	if err != nil {
		return nil, err
	}
	res = &permmenu{
		Perms: perms,
		Menus: menus,
	}

	return

}

// refreshPerms(userId) 刷新指定用户的 sa-token 角色与权限
func (c *BaseSysPermsService) RefreshPerms(ctx context.Context, userId uint) (err error) {
	roleIds, err := NewBaseSysRoleService().GetByUser(userId)
	if err != nil {
		return err
	}
	_, err = satokenSetupUserAuth(userId, roleIds)
	return
}
