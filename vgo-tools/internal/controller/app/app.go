package app

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	v1 "github.com/vera-byte/vgo/vgo-tools/api/app/v1"
)

type Controller struct {
	v1.UnimplementedBaseSysConfCrudServiceServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterBaseSysConfCrudServiceServer(s.Server, &Controller{})
}

func (*Controller) Add(ctx context.Context, req *v1.BaseSysConfAddRpcInvoke) (res *v1.BaseSysConfAddRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Update(ctx context.Context, req *v1.BaseSysConfUpdateRpcInvoke) (res *v1.BaseSysConfUpdateRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Info(ctx context.Context, req *v1.BaseSysConfInfoRpcInvoke) (res *v1.BaseSysConfInfoRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Delete(ctx context.Context, req *v1.BaseSysConfDeleteRpcInvoke) (res *v1.BaseSysConfDeleteRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Page(ctx context.Context, req *v1.BaseSysConfPageRpcInvoke) (res *v1.BaseSysConfPageRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) List(ctx context.Context, req *v1.BaseSysConfListRpcInvoke) (res *v1.BaseSysConfListRpcRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
