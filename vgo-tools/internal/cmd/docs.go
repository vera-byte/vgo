package cmd

import (
	"context"
	"net"
	"os"
	"os/exec"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/vgo-tools/internal/utility/mlog"
)

var (
	Docs = gcmd.Command{
		Name:  "docs",
		Usage: "vgo-tools docs [--dev] [--dir docs]",
		Brief: "查看帮助文档",
		Description: "内置文档服务器：托管随 vgo-tools 发布的 Mintlify 静态文档（无需安装任何依赖）。\n" +
			"--dev 则在本地 docs 目录启动 Mintlify 开发预览（需要 Mintlify CLI 与 docs.json，文档作者用）。",
		Arguments: []gcmd.Argument{
			{
				Name:  "dev",
				Brief: "以 Mintlify 开发模式预览本地 docs 目录（需要 npm i -g mint）",
			},
			{
				Name:  "dir",
				Brief: "开发模式的文档目录（默认 docs）",
			},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) error {
			if parser.GetOpt("dev", false).Bool() {
				return runMintDev(parser.GetOpt("dir", "docs").String())
			}
			if !gres.Contains("docs/index.html") {
				mlog.Fatalf("内置文档缺失：请使用完整发布的 vgo-tools 二进制，或用 --dev 预览本地 docs 目录")
				return nil
			}
			port := freePort()
			s := g.Server("docs")
			s.SetServerRoot("docs")
			s.SetPort(port)
			mlog.Printf("Vgo docs server is running at %s", "http://127.0.0.1:"+gconv.String(port))
			s.Run()
			return nil
		},
	}
)

// runMintDev 以 Mintlify 开发模式预览本地文档（可选路径，需要 Mintlify CLI）
func runMintDev(dir string) error {
	if !gfile.Exists(gfile.Join(dir, "docs.json")) {
		mlog.Fatalf("未找到 %s/docs.json：请在包含 Mintlify 文档的项目根目录执行（或用 --dir 指定）", dir)
		return nil
	}
	bin, args := mintCommand()
	if bin == "" {
		mlog.Fatalf("未找到 Mintlify CLI（npm i -g mint）；不安装也可以直接运行 vgo-tools docs 查看内置文档")
		return nil
	}
	mlog.Printf("启动 Mintlify 开发预览（%s）...", dir)
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// mintCommand 查找 Mintlify CLI（mint 优先，回退 npx mint）
func mintCommand() (string, []string) {
	if path, err := exec.LookPath("mint"); err == nil {
		return path, []string{"dev"}
	}
	if path, err := exec.LookPath("npx"); err == nil {
		return path, []string{"--yes", "mint", "dev"}
	}
	return "", nil
}

// freePort 获取本机未占用的端口
func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 3210
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 3210
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return port
}

func init() {
	Main.AddCommand(&Docs)
}
