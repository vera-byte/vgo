package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/v"
)

type BaseSysPermsService struct {
}

func NewBaseSysPermsService() *BaseSysPermsService {
	return &BaseSysPermsService{}
}

// permmenu 方法
func (c *BaseSysPermsService) Permmenu(ctx context.Context, roleIds []string) (res interface{}) {
	type permmenu struct {
		Perms []string   `json:"perms"`
		Menus gdb.Result `json:"menus"`
	}
	var (
		baseSysMenuService = NewBaseSysMenuService()
		admin              = v.GetAdmin(ctx)
	)

	res = &permmenu{
		Perms: baseSysMenuService.GetPerms(roleIds),
		Menus: baseSysMenuService.GetMenus(admin.RoleIds, admin.UserId == 1),
	}

	return

}

// refreshPerms(userId)
func (c *BaseSysPermsService) RefreshPerms(ctx context.Context, userId uint) (err error) {
	var (
		baseSysUserRoleService   = NewBaseSysRoleService()
		baseSysMenuService       = NewBaseSysMenuService()
		baseSysDepartmentService = NewBaseSysDepartmentService()
		roleIds                  = baseSysUserRoleService.GetByUser(userId)
		perms                    = baseSysMenuService.GetPerms(roleIds)
	)
	v.CacheManager.Set(ctx, "admin:perms:"+gconv.String(userId), perms, 0)
	// 更新部门权限
	departments := baseSysDepartmentService.GetByRoleIds(roleIds, userId == 1)
	v.CacheManager.Set(ctx, "admin:department:"+gconv.String(userId), departments, 0)

	return
}
