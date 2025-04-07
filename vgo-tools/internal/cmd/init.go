package cmd

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/vera-byte/vgo/vgo-tools/utility/allyes"
	"github.com/vera-byte/vgo/vgo-tools/utility/mlog"
)

var (
	// Init .
	Init = cInitVgoProject{}
)

type (
	cInitVgoProject struct {
		g.Meta `name:"init" brief:"{cInitBrief}" eg:"{cInitEg}"`
	}
	CInitVgoProjectInput struct {
		g.Meta `name:"init"`
		Name   string `name:"NAME" arg:"true" v:"required" brief:"{cInitNameBrief}"`
	}
	CInitVgoProjectOutput struct{}
)

const (
	cInitBrief = `初始化一个带有后台管理系统的VGO项目`
	cInitEg    = `
	vgo init my-project
	`
	cInitNameBrief = `
	项目的名称。它将在当前目录中创建一个带有 NAME 的文件夹。NAME 也将是项目的模块名称。
	`
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cInitBrief":     cInitBrief,
		"cInitEg":        cInitEg,
		"cInitNameBrief": cInitNameBrief,
	})
}

func (c cInitVgoProject) Index(ctx context.Context, in CInitVgoProjectInput) (out *CInitVgoProjectOutput, err error) {
	g.Log().Debug(ctx, "项目初始化")
	var overwrote = false
	if !gfile.IsEmpty(in.Name) && !allyes.Check() {
		s := gcmd.Scanf(`the folder "%s" is not empty, files might be overwrote, continue? [y/n]: `, in.Name)
		if strings.EqualFold(s, "n") {
			return
		}
		overwrote = true
	}
	mlog.Print("initializing...", overwrote)
	return nil, nil
}
