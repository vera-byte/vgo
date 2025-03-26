package utility

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	r := g.RequestFromCtx(ctx)
	gatewayAdminToken := r.GetHeader("Authorization")
	if len(gatewayAdminToken) > 0 {
		newMD := metadata.Pairs("AdminToken", gatewayAdminToken)

		// ✅ 传递 metadata 到新的 context
		newCtx := metadata.NewOutgoingContext(ctx, newMD)

		// 调用 gRPC 方法
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
	return invoker(ctx, method, req, reply, cc, opts...)

}
