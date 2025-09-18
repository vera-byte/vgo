package cmd

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/vera-byte/vgo/v/cmd"
)

func init() {
	dictCmd := &gcmd.Command{
		Name:  "dictCmd",
		Brief: "打招呼",
		Arguments: []gcmd.Argument{
			{Name: "name", Short: "n", Brief: "名字", Default: "GoFrame"},
		},
		Func: func(ctx g.Ctx, parser *gcmd.Parser) error {
			name := parser.GetOpt("name").String()
			fmt.Printf("Hello, %s 👋\n", name)
			return nil
		},
	}
	cmd.Root.AddCommand(dictCmd)

}
