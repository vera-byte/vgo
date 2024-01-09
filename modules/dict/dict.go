package dict

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/vera-byte/vgo/cool"
	_ "github.com/vera-byte/vgo/modules/dict/packed"

	_ "github.com/vera-byte/vgo/modules/dict/controller"
	"github.com/vera-byte/vgo/modules/dict/model"
)

func init() {
	var (
		ctx = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module dict init start ...")
	cool.FillInitData(ctx, "dict", &model.DictInfo{})
	cool.FillInitData(ctx, "dict", &model.DictType{})
	g.Log().Debug(ctx, "module dict init finished ...")
}
