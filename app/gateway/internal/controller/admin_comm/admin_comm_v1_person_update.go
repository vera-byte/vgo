package admin_comm

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	protobuf "vgo/app/admin/api/comm/v1"
	v1 "vgo/app/gateway/api/admin_comm/v1"
)

func (c *ControllerV1) PersonUpdate(ctx context.Context, req *v1.PersonUpdateReq) (res *v1.PersonUpdateRes, err error) {
	updateFields := map[string]string{
		"headImg":  req.HeadImg,
		"nickName": req.NickName,
		"password": req.Password,
	}

	// 遍历需要更新的字段
	for key, value := range updateFields {
		if value != "" {
			InvokeArg := &protobuf.PersonUpdateInvoke{
				Key:   key,
				Value: value,
			}

			// 处理密码特殊逻辑
			if key == "password" {
				InvokeArg.Other = req.OldPassword
			}

			// 记录日志
			g.Log().Infof(ctx, "用户更改 %s", key)

			// 调用更新方法
			if _, err := c.AdminBaseCommClient.PersonUpdate(ctx, InvokeArg); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}
