package admin_comm

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"vgo/app/gateway/api/admin_comm/v1"
)

func (c *ControllerV1) Person(ctx context.Context, req *v1.PersonReq) (res *v1.PersonRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
