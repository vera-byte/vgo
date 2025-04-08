package cmd

import (
	"context"

	"github.com/vera-byte/vgo/app/admin/internal/controller/comm"
	"github.com/vera-byte/vgo/app/admin/internal/controller/open"
	"github.com/vera-byte/vgo/app/admin/internal/controller/system"

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
			system.Register(s)

			s.Run()
			return nil
		},
	}
)
