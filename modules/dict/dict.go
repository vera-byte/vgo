package dict

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/vera-byte/vgo/modules/dict/cmd"
	_ "github.com/vera-byte/vgo/modules/dict/packed"

	_ "github.com/vera-byte/vgo/modules/dict/controller"
)

func init() {
	var (
		ctx = gctx.GetInitCtx()
	)
	g.Log().Debug(ctx, "module dict init start ...")
	// v.FillInitData(ctx, "dict", &model.DictInfo{})
	// v.FillInitData(ctx, "dict", &model.DictType{})
	g.Log().Debug(ctx, "module dict init finished ...")
}
