package cmd

import (
	"context"
	"vgo/app/admin/internal/controller/comm"
	"vgo/app/admin/internal/controller/open"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/os/gcmd"
	_ "github.com/vera-byte/vgo/vgo_core_kit"
	"google.golang.org/grpc"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "rpc admin server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			c := grpcx.Server.NewConfig()

			c.Options = append(c.Options, []grpc.ServerOption{
				grpcx.Server.ChainUnary(
					grpcx.Server.UnaryValidate,
				)}...,
			)
			s := grpcx.Server.New(c)
			open.Register(s)
			comm.Register(s)

			s.Run()
			return nil
		},
	}
)
