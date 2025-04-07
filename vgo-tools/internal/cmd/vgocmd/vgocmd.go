package vgocmd

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/vera-byte/vgo/vgo-tools/internal/cmd"
	"github.com/vera-byte/vgo/vgo-tools/utility/allyes"
	"github.com/vera-byte/vgo/vgo-tools/utility/mlog"
)

const cliFolderName = `hack`

type Command struct {
	*gcmd.Command
}

// Run starts running the command according the command line arguments and options.
func (c *Command) Run(ctx context.Context) {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				mlog.Print(err.Error())
			} else {
				panic(gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception))
			}
		}
	}()

	// CLI configuration, using the `hack/vgo.yaml` in priority.
	if path, _ := gfile.Search(cliFolderName); path != "" {
		if adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
			if err := adapter.SetPath(path); err != nil {
				mlog.Fatal(err)
			}
			adapter.SetFileName("vgo.yaml")
		}
	}

	// -y option checks.
	allyes.Init()

	// just run.
	if err := c.RunWithError(ctx); err != nil {
		// Exit with error message and exit code 1.
		// It is very important to exit the command process with code 1.
		mlog.Fatalf(`%+v`, err)
	}
}
func GetCommand(ctx context.Context) (*Command, error) {
	root, err := gcmd.NewFromObject(cmd.Vgo)
	if err != nil {
		return nil, err
	}
	err = root.AddObject(
		// cmd.Up,
		// cmd.Env,
		// cmd.Fix,
		// cmd.Run,
		// cmd.Gen,
		// cmd.Tpl,
		cmd.Init,
		// cmd.Pack,
		// cmd.Build,
		// cmd.Docker,
		// cmd.Install,
		cmd.Version,
		cmd.Gen,
	// cmd.Doc,
	)
	if err != nil {
		return nil, err
	}
	command := &Command{
		root,
	}
	return command, nil
}
