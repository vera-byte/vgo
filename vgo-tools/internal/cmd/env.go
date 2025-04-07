package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gproc"
)

var Env = cEnv{}

type cEnv struct {
	g.Meta `name:"env" ad:"{cGFAd}"`
}
type cEnvInput struct {
	g.Meta `name:"env"`
}

func (c cEnv) Index(ctx context.Context, in cEnvInput) (out *cGFOutput, err error) {
	getEnv()
	gcmd.CommandFromCtx(ctx).Print()
	return
}
func getEnv() (str string, ok bool) {
	env, err := gproc.ShellExec(context.Background(), "gf env")
	if err != nil {
		return "", false
	}
	return env, true
}
