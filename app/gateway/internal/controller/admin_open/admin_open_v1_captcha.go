package admin_open

import (
	"context"

	protobuf "vgo/app/admin/api/open/v1"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "vgo/app/gateway/api/admin_open/v1"
)

func (c *ControllerV1) Captcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error) {
	_res, _err := c.AdminBaseOpenClient.Captcha(ctx, &protobuf.CaptchaRpcInvoke{
		Height: req.Height,
		Widget: req.Width,
	})
	if _err != nil {
		g.Log().Error(ctx, _err)
		return nil, gerror.New("生成验证码失败")
	}
	return &v1.CaptchaRes{
		CaptchaId: _res.CaptchaId,
		Data:      _res.Data,
	}, nil
}
