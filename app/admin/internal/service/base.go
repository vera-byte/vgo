// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"vgo/app/admin/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IBaseSysDepartmentLogic interface {
		// GetByRoleIds 获取部门
		GetByRoleIds(ctx context.Context, roleIds []string, isAdmin bool) (res []uint)
		// Order 排序部门
		Order(ctx g.Ctx) (err error)
	}
	IBaseSysLoginLogic interface {
		// 生成验证码
		GenerateCaptcha(ctx context.Context, width int, height int) (id string, b64s string, answer string, err error)
		// 验证验证码
		VerifyCaptcha(id string, answer string) bool
		// 密码登录 此处只验证密码和验证码 Token由其他函数生成
		Login(ctx context.Context, captchaId string, password string, userName string, code string) (expire *int64, refreshExpire *int64, token *string, refreshToken *string, err error)
	}
	IBaseSysMenuLogic interface {
		// GetPerms 获取菜单的权限
		GetPerms(ctx context.Context, roleIds []string) []string
		// GetMenus 获取菜单
		GetMenus(ctx context.Context, roleIds []string, isAdmin bool) []entity.BaseSysMenu
	}
	IBaseSysRoleLogic interface {
		// 通过用户ID获取角色集合
		GetByUser(ctx context.Context, userId int64) (roles []string, err error)
	}
	IBaseSysUserLogic interface {
		Person(ctx context.Context, userId int64) (user interface{}, err error)
	}
)

var (
	localBaseSysDepartmentLogic IBaseSysDepartmentLogic
	localBaseSysLoginLogic      IBaseSysLoginLogic
	localBaseSysMenuLogic       IBaseSysMenuLogic
	localBaseSysRoleLogic       IBaseSysRoleLogic
	localBaseSysUserLogic       IBaseSysUserLogic
)

func BaseSysDepartmentLogic() IBaseSysDepartmentLogic {
	if localBaseSysDepartmentLogic == nil {
		panic("implement not found for interface IBaseSysDepartmentLogic, forgot register?")
	}
	return localBaseSysDepartmentLogic
}

func RegisterBaseSysDepartmentLogic(i IBaseSysDepartmentLogic) {
	localBaseSysDepartmentLogic = i
}

func BaseSysLoginLogic() IBaseSysLoginLogic {
	if localBaseSysLoginLogic == nil {
		panic("implement not found for interface IBaseSysLoginLogic, forgot register?")
	}
	return localBaseSysLoginLogic
}

func RegisterBaseSysLoginLogic(i IBaseSysLoginLogic) {
	localBaseSysLoginLogic = i
}

func BaseSysMenuLogic() IBaseSysMenuLogic {
	if localBaseSysMenuLogic == nil {
		panic("implement not found for interface IBaseSysMenuLogic, forgot register?")
	}
	return localBaseSysMenuLogic
}

func RegisterBaseSysMenuLogic(i IBaseSysMenuLogic) {
	localBaseSysMenuLogic = i
}

func BaseSysRoleLogic() IBaseSysRoleLogic {
	if localBaseSysRoleLogic == nil {
		panic("implement not found for interface IBaseSysRoleLogic, forgot register?")
	}
	return localBaseSysRoleLogic
}

func RegisterBaseSysRoleLogic(i IBaseSysRoleLogic) {
	localBaseSysRoleLogic = i
}

func BaseSysUserLogic() IBaseSysUserLogic {
	if localBaseSysUserLogic == nil {
		panic("implement not found for interface IBaseSysUserLogic, forgot register?")
	}
	return localBaseSysUserLogic
}

func RegisterBaseSysUserLogic(i IBaseSysUserLogic) {
	localBaseSysUserLogic = i
}
