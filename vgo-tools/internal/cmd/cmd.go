package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gtag"
)

var Vgo = cVgo{}

type cVgo struct {
	g.Meta `name:"vgo-tools" ad:"{cGFAd}"`
}

const (
	cGFAd = `
ADDITIONAL
    Use "gf COMMAND -h" for details about a command.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGFAd`: cGFAd,
	})
}

type cVgoInput struct {
	g.Meta  `name:"vgo-tools"`
	Yes     bool `short:"y" name:"yes"     brief:"all yes for all command without prompt ask"   orphan:"true"`
	Version bool `short:"v" name:"version" brief:"show version information of current binary"   orphan:"true"`
	Debug   bool `short:"d" name:"debug"   brief:"show internal detailed debugging information" orphan:"true"`
}
type cGFOutput struct{}

func (c cVgo) Index(ctx context.Context, in cVgoInput) (out *cGFOutput, err error) {
	gcmd.CommandFromCtx(ctx).Print()
	return
}
