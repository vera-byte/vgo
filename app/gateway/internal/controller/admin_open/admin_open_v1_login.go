package admin_open

import (
	"context"

	protobuf "github.com/vera-byte/vgo/app/admin/api/open/v1"

	v1 "github.com/vera-byte/vgo/app/gateway/api/admin_open/v1"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	// 调用微服务
	auth, err := c.AdminBaseOpenClient.Login(ctx, &protobuf.LoginRpcInvoke{
		Username:   req.Username,
		Password:   req.Password,
		CaptchaId:  req.CaptchaId,
		VerifyCode: req.VerifyCode,
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		Expire:        auth.Expire,
		RefreshExpire: auth.RefreshExpire,
		Token:         auth.Token,
		RefreshToken:  auth.RefreshToken,
	}, nil
}
