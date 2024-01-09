package admin

import (
	"context"

	"github.com/vera-byte/vgo/modules/base/service"
	"github.com/vera-byte/vgo/v"

	"github.com/gogf/gf/v2/frame/g"
)

type BaseCommController struct {
	*v.ControllerSimple
}

func init() {
	var base_comm_controller = &BaseCommController{
		&v.ControllerSimple{
			Perfix: "/app/base/comm",
			//    Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			//    Service: service.NewBaseCommService(),
		},
	}
	// 注册路由
	v.RegisterControllerSimple(base_comm_controller)
}

// eps 接口请求
type BaseCommControllerEpsReq struct {
	g.Meta `path:"/eps" method:"GET"`
}

// eps 接口
func (c *BaseCommController) Eps(ctx context.Context, req *BaseCommControllerEpsReq) (res *v.BaseRes, err error) {
	if !v.Config.Eps {
		g.Log().Error(ctx, "eps is not open")
		res = v.Ok(nil)
		return
	}
	baseOpenService := service.NewBaseOpenService()
	data, err := baseOpenService.AppEPS(ctx)
	if err != nil {
		g.Log().Error(ctx, "eps error", err)
		return v.Fail(err.Error()), err
	}
	res = v.Ok(data)
	return
}
