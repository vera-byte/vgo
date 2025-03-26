package admin_open

import (
	"context"

	protobuf "vgo/app/admin/api/open/v1"
	v1 "vgo/app/gateway/api/admin_open/v1"
)

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	auth, err := c.AdminBaseOpenClient.RefreshToken(ctx, &protobuf.RefreshTokenInvoke{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RefreshTokenRes{
		Expire:        auth.Expire,
		RefreshExpire: auth.RefreshExpire,
		Token:         auth.Token,
		RefreshToken:  auth.RefreshToken,
	}, nil
}
