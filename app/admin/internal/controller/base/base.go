package base

import (
	"context"
	v1 "vgo/app/admin/api/base/v1"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedBaseServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterBaseServer(s.Server, &Controller{})
}

func (*Controller) Captcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error) {
	id, b64s, _, err := service.Captcha().GenerateCaptcha(ctx, int(req.Widget), int(req.Height), req.Color)
	if err != nil {
		return nil, gerror.New("生成验证码失败")
	}
	return &v1.CaptchaRes{
		CaptchaId: id,
		Data:      b64s,
	}, nil
}
