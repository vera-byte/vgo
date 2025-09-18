package base

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/vera-byte/vgo/modules/base/cmd"
	_ "github.com/vera-byte/vgo/modules/base/controller"
	_ "github.com/vera-byte/vgo/modules/base/funcs"
	_ "github.com/vera-byte/vgo/modules/base/middleware"
	_ "github.com/vera-byte/vgo/modules/base/packed"
)

func init() {
	var (
		ctx = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module base init start ...")

	// v.FillInitData(ctx, "base", &model.BaseSysMenu{})
	// v.FillInitData(ctx, "base", &model.BaseSysUser{})
	// v.FillInitData(ctx, "base", &model.BaseSysUserRole{})
	// v.FillInitData(ctx, "base", &model.BaseSysRole{})
	// v.FillInitData(ctx, "base", &model.BaseSysRoleMenu{})
	// v.FillInitData(ctx, "base", &model.BaseSysDepartment{})
	// v.FillInitData(ctx, "base", &model.BaseSysRoleDepartment{})
	// v.FillInitData(ctx, "base", &model.BaseSysParam{})

	g.Log().Debug(ctx, "module base init finished ...")

}
