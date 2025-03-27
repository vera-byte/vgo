package open

import (
	"context"
	v1 "vgo/app/admin/api/open/v1"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedBaseOpenServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterBaseOpenServer(s.Server, &Controller{})
}

func (*Controller) Captcha(ctx context.Context, req *v1.CaptchaRpcInvoke) (res *v1.CaptchaRpcRes, err error) {
	id, b64s, _, err := service.BaseSysLoginLogic().GenerateCaptcha(ctx, int(req.Widget), int(req.Height))
	if err != nil {
		return nil, gerror.New("生成验证码失败,服务端异常")
	}
	return &v1.CaptchaRpcRes{
		CaptchaId: id,
		Data:      b64s,
	}, nil
}

func (*Controller) Login(ctx context.Context, req *v1.LoginRpcInvoke) (res *v1.LoginRpcRes, err error) {
	var (
		baseSysLoginService = service.BaseSysLoginLogic()
	)
	// 核验用户信息
	result, err := baseSysLoginService.Login(ctx, req.CaptchaId, req.Password, req.Username, req.VerifyCode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginRpcRes{
		Expire:        result.Expire,
		RefreshExpire: result.RefreshExpire,
		Token:         result.Token,
		RefreshToken:  result.RefreshToken,
	}, nil
}

func (*Controller) RefreshToken(ctx context.Context, req *v1.RefreshTokenInvoke) (res *v1.LoginRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
