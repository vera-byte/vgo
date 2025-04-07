package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	Gen = cGen{}
)

type cGen struct {
	g.Meta `name:"gen" eg:"{cGenEg}" brief:"vgo 快速生成"`
	cGenGrpcCrud
	cGenGatewayCrud
}

const (
	cGenEg = `
		vgo gen grpc-crud
		vgo gen gateway-crud
	`
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cGenEg": cGenEg,
	})
}
