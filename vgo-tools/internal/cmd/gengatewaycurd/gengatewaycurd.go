package gengatewaycurd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	CGenGatewayCrud      struct{}
	CGenGatewayCrudInput struct {
		g.Meta `name:"gateway-crud"`
	}
	CGenGatewayCrudOutput struct{}
)

func (c CGenGatewayCrud) GenGatewayCrud(ctx context.Context, in CGenGatewayCrudInput) (*CGenGatewayCrudOutput, error) {
	g.Log().Debug(ctx, "生成网关 CRUD")
	g.Dump(in)
	return nil, nil
}
