package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Vgo = cVgo{}

type cVgo struct {
	g.Meta `name:"vgo"`
}

type cVgoInput struct {
	g.Meta  `name:"vgo"`
	Version bool `short:"v" name:"version" brief:"显示当前二进制文件的版本信息"   orphan:"true"`
}
type cGFOutput struct{}

func (c cVgo) Index(ctx context.Context, in cVgoInput) (out *cGFOutput, err error) {
	gcmd.CommandFromCtx(ctx).Print()
	return
}
