package vgocmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/vgo-tools/internal/cmd"
)

type Command struct {
	*gcmd.Command
}

func GetCommand(ctx context.Context) (*Command, error) {
	root, err := gcmd.NewFromObject(cmd.Vgo)
	if err != nil {
		return nil, err
	}
	err = root.AddObject(
		// cmd.Up,
		cmd.Env,
		// cmd.Fix,
		// cmd.Run,
		// cmd.Gen,
		// cmd.Tpl,
		// cmd.Init,
		// cmd.Pack,
		// cmd.Build,
		// cmd.Docker,
		// cmd.Install,
		cmd.Version,
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
