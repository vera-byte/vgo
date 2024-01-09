package base

import (
	_ "github.com/vera-byte/vgo/modules/base/packed"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/vera-byte/vgo/cool"
	_ "github.com/vera-byte/vgo/modules/base/controller"
	_ "github.com/vera-byte/vgo/modules/base/funcs"
	_ "github.com/vera-byte/vgo/modules/base/middleware"
	"github.com/vera-byte/vgo/modules/base/model"
)

func init() {
	var (
		ctx = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module base init start ...")

	cool.FillInitData(ctx, "base", &model.BaseSysMenu{})
	cool.FillInitData(ctx, "base", &model.BaseSysUser{})
	cool.FillInitData(ctx, "base", &model.BaseSysUserRole{})
	cool.FillInitData(ctx, "base", &model.BaseSysRole{})
	cool.FillInitData(ctx, "base", &model.BaseSysRoleMenu{})
	cool.FillInitData(ctx, "base", &model.BaseSysDepartment{})
	cool.FillInitData(ctx, "base", &model.BaseSysRoleDepartment{})
	cool.FillInitData(ctx, "base", &model.BaseSysParam{})

	g.Log().Debug(ctx, "module base init finished ...")

}
