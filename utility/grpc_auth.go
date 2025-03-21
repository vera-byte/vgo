package utility

import (
	"context"
	"encoding/json"

	vck "github.com/vera-byte/vgo/vgo_core_kit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	admin := vck.GetAdminAtGateway(ctx)

	// ✅ 将 Admin 结构体转换为 JSON 存入 metadata
	adminJSON, _ := json.Marshal(admin)
	newMD := metadata.Pairs("Admin", string(adminJSON))

	// ✅ 传递 metadata 到新的 context
	newCtx := metadata.NewOutgoingContext(ctx, newMD)

	// 调用 gRPC 方法
	return invoker(newCtx, method, req, reply, cc, opts...)
}
