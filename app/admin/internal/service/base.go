// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"github.com/vera-byte/vgo/app/admin/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	vck_config "github.com/vera-byte/vgo/vgo_core_kit/config"
	vck_request "github.com/vera-byte/vgo/vgo_core_kit/request"
)

type (
	IBaseSysDepartmentLogic interface {
		// GetByRoleIds 获取部门
		// GetByRoleIds 根据角色 ID 获取部门 ID 列表
		// 如果 isAdmin 为 true，则返回所有部门 ID（超级管理员权限）
		// 如果 isAdmin 为 false，则返回该角色所属的部门 ID
		GetByRoleIds(ctx context.Context, roleIds []string, isAdmin bool) (res []uint)
		// Order 排序部门
		Order(ctx g.Ctx) (err error)
		// 部门列表
		List(ctx context.Context) (res []entity.BaseSysDepartment, err error)
	}
	IBaseSysLogLogic interface {
		// 记录日志
		RecordLog(ctx context.Context, userId *int64, action string, ip string, params string, tenantId *int64, traceId string) (err error)
	}
	IBaseSysLoginLogic interface {
		// 生成验证码
		GenerateCaptcha(ctx context.Context, width int, height int) (id string, b64s string, answer string, err error)
		// 验证验证码
		VerifyCaptcha(id string, answer string) bool
		// 退出登录
		LoginOut(ctx context.Context) (err error)
		// 密码登录 此处只验证密码和验证码 Token由其他函数生成
		Login(ctx context.Context, captchaId string, password string, userName string, code string) (result *vck_config.TokenResult, err error)
		RefreshToken(ctx context.Context, _refreshToken string) (result *vck_config.TokenResult, err error)
	}
	IBaseSysMenuLogic interface {
		GetPerms(ctx context.Context, roleIds []string) []string
		// GetMenus 获取指定角色的菜单列表
		// 如果 isAdmin 为 true，则返回所有菜单（超级管理员权限）
		// 如果 isAdmin 为 false，则根据角色 ID 查询对应的菜单
		GetMenus(ctx context.Context, roleIds []string, isAdmin bool) []entity.BaseSysMenu
	}
	IBaseSysRoleLogic interface {
		// 通过用户ID获取角色集合
		GetByUser(ctx context.Context, userId int64) (roles []string, err error)
	}
	IBaseSysUserLogic interface {
		// 用户信息
		Person(ctx context.Context, userId int64) (user *entity.BaseSysUser, err error)
		// 更新用户信息
		PersonUpdate(ctx context.Context, userId int64, k string, v string, o string) error
		// page
		Page(ctx context.Context, pageReq *vck_request.PageReq, departmentIds []int64) (res []*entity.BaseSysUser, pagination *vck_request.Pagination, err error)
	}
)

var (
	localBaseSysDepartmentLogic IBaseSysDepartmentLogic
	localBaseSysLogLogic        IBaseSysLogLogic
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

func BaseSysLogLogic() IBaseSysLogLogic {
	if localBaseSysLogLogic == nil {
		panic("implement not found for interface IBaseSysLogLogic, forgot register?")
	}
	return localBaseSysLogLogic
}

func RegisterBaseSysLogLogic(i IBaseSysLogLogic) {
	localBaseSysLogLogic = i
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
