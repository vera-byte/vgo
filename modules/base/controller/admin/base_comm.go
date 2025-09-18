package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"
)

type BaseCommController struct {
	*v.ControllerSimple
}

func init() {
	var base_comm_controller = &BaseCommController{
		ControllerSimple: &v.ControllerSimple{
			Perfix: "/admin/base/comm",
		},
	}
	// 注册路由
	v.RegisterControllerSimple(base_comm_controller)
}

// BaseCommPersonReq 接口请求参数
// Person 方法
func (c *BaseCommController) Person(ctx context.Context, req *v1.BaseCommPersonReq) (res *v.BaseRes, err error) {
	var (
		baseSysUserService = service.NewBaseSysUserService()
		admin              = v.GetAdmin(ctx)
	)
	data, err := baseSysUserService.Person(admin.UserId)
	res = v.Ok(data)
	return
}

// BaseCommPermmenuReq 接口请求参数
// Permmenu 方法
func (c *BaseCommController) Permmenu(ctx context.Context, req *v1.BaseCommPermmenuReq) (res *v.BaseRes, err error) {

	var (
		baseSysPermsService = service.NewBaseSysPermsService()
		admin               = v.GetAdmin(ctx)
	)
	res = v.Ok(baseSysPermsService.Permmenu(ctx, admin.RoleIds))
	return
}

// Logout BaseCommLogout 方法
func (c *BaseCommController) Logout(ctx context.Context, req *v1.BaseCommLogoutReq) (res *v.BaseRes, err error) {
	var (
		BaseSysLoginService = service.NewBaseSysLoginService()
	)
	err = BaseSysLoginService.Logout(ctx)
	res = v.Ok(nil)
	return
}

// UploadMode 方法
func (c *BaseCommController) UploadMode(ctx context.Context, req *v1.BaseCommUploadModeReq) (res *v.BaseRes, err error) {
	data, err := v.File().GetMode()
	res = v.Ok(data)
	return
}

// Upload 方法
func (c *BaseCommController) Upload(ctx context.Context, req *v1.BaseCommUploadReq) (res *v.BaseRes, err error) {
	data, err := v.File().Upload(ctx)
	res = v.Ok(data)
	return
}

// PersonUpdate 方法
func (c *BaseCommController) PersonUpdate(ctx g.Ctx, req *v1.PersonUpdateReq) (res *v.BaseRes, err error) {
	var (
		baseSysUserService = service.NewBaseSysUserService()
	)

	_, err = baseSysUserService.ServiceUpdate(ctx, &v.UpdateReq{})
	if err != nil {
		return
	}

	res = v.Ok(nil)
	return
}
