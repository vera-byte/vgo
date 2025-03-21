package utility

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func GrpcClientTimeout(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
