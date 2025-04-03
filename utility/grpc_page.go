package utility

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func PageRequestInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	r := g.RequestFromCtx(ctx)
	m := gmap.New(true)
	m.Set("page", r.Get("page", 1).Int())
	m.Set("size", r.Get("size", 20).Int())
	order := r.Get("order").String()
	if len(order) > 0 {
		m.Set("order", order)
	}
	sort := r.Get("sort").String()
	if len(sort) > 0 {
		m.Set("sort", sort)
	}
	newMD := metadata.Pairs("PageReq", m.String())

	// ✅ 传递 metadata 到新的 context
	newCtx := metadata.NewOutgoingContext(ctx, newMD)

	// 调用 gRPC 方法
	return invoker(newCtx, method, req, reply, cc, opts...)

}
