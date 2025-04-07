package gengrpccurd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/vera-byte/vgo/vgo-tools/utility/mlog"
)

const (
	CGenGrpcCrudConfig = `vgocli.gen.grpc.crud`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenGrpcCrudConfig`: CGenGrpcCrudConfig,
	})
}

type (
	CGenGrpcCrud      struct{}
	CGenGrpcCrudInput struct {
		g.Meta    `name:"grpc-crud" config:"{CGenGrpcCrudConfig}"`
		Version   string `name:"version" short:"v" d:"v1"`
		Path      string `name:"path"              short:"p"  brief:"{CGenPbEntityBriefPath}" d:"manifest/protobuf/pbentity"`
		Out       string `name:"out"              short:"o"  brief:"{CGenPbEntityBriefPath}" d:"manifest/protobuf"`
		Tables    string `name:"tables"            short:"t"  brief:"{CGenPbEntityBriefTables}"`
		Package   string `name:"package"           short:"k"  brief:"{CGenPbEntityBriefPackage}"`
		GoPackage string `name:"goPackage"           short:"g"  brief:"{CGenPbEntityBriefGoPackage}"`
	}
	CGenGrpcCrudOutput struct{}
)

func (c CGenGrpcCrud) GenGrpcCrud(ctx context.Context, in CGenGrpcCrudInput) (out *CGenGrpcCrudOutput, err error) {
	var (
		config = g.Cfg()
	)
	if config.Available(ctx) {
		v := config.MustGet(ctx, CGenGrpcCrudConfig)
		g.Dump(CGenGrpcCrudConfig)
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				err = genGrpcCrud(ctx, i, in)
			}
		} else {
			err = genGrpcCrud(ctx, -1, in)
		}
	} else {
		err = genGrpcCrud(ctx, -1, in)
	}
	if err != nil {
		mlog.Fatalf("❌生成失败:%s", err.Error())
	}
	mlog.Print("✅完成!")
	return
}
func genGrpcCrud(ctx context.Context, index int, in CGenGrpcCrudInput) (err error) {

	if index >= 0 {
		err = g.Cfg().MustGet(
			ctx,
			fmt.Sprintf(`%s.%d`, CGenGrpcCrudConfig, index),
		).Scan(&in)
		if err != nil {
			mlog.Fatalf(`无效配置 "%s": %+v`, CGenGrpcCrudConfig, err)
		}
	}
	_, err = gproc.ShellExec(ctx, "gf gen pbentity")
	if err != nil {
		return err
	}
	g.Dump(in)
	return

}
