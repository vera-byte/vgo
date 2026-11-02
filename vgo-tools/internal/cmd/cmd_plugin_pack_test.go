package cmd

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTestFile 写入测试文件
func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// zipEntries 读取 ZIP 内条目名（含目录）与文件大小
func zipEntries(t *testing.T, path string) map[string]uint64 {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("打开 zip 失败: %v", err)
	}
	defer zr.Close()
	out := map[string]uint64{}
	for _, f := range zr.File {
		out[f.Name] = f.UncompressedSize64
	}
	return out
}

func TestModuleNameFromGoMod(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		file    string
		content string
		want    string
	}{
		{"simple.mod", "module demo-plugin\n\ngo 1.27.3\n", "demo-plugin"},
		{"path.mod", "module github.com/foo/bar/baz\n\ngo 1.27.3\n", "baz"},
		{"nomodule.mod", "go 1.27.3\n", ""},
		{"missing.mod", "", ""},
	}
	for _, c := range cases {
		p := filepath.Join(dir, c.file)
		if c.content != "" {
			writeTestFile(t, p, c.content)
		}
		if got := moduleNameFromGoMod(p); got != c.want {
			t.Fatalf("%s: got %q want %q", c.file, got, c.want)
		}
	}
}

func TestShouldSkipOnPack(t *testing.T) {
	skip := []string{"go.mod", "go.sum", "build", "bin", "main.go", "helper.go", "plugin.zip"}
	keep := []string{"plugin.json", "README.md", "assets", "logo.png"}
	for _, n := range skip {
		if !shouldSkipOnPack(n) {
			t.Fatalf("应跳过: %s", n)
		}
	}
	for _, n := range keep {
		if shouldSkipOnPack(n) {
			t.Fatalf("不应跳过: %s", n)
		}
	}
}

func TestHumanSize(t *testing.T) {
	cases := map[int64]string{
		512:        "512B",
		2048:       "2.0K",
		5 * 1048576: "5.0M",
	}
	for n, want := range cases {
		if got := humanSize(n); got != want {
			t.Fatalf("humanSize(%d)=%q want %q", n, got, want)
		}
	}
}

func TestZipDir(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "a.txt"), "hello hello hello hello hello")
	writeTestFile(t, filepath.Join(dir, "sub", "b.txt"), "world")
	bin := filepath.Join(dir, "bin", "run")
	writeTestFile(t, bin, "#!/bin/sh\necho hi\n")
	if err := os.Chmod(bin, 0o755); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(t.TempDir(), "out.zip")
	if err := zipDir(dir, out); err != nil {
		t.Fatalf("zipDir: %v", err)
	}
	entries := zipEntries(t, out)
	for _, name := range []string{"a.txt", "sub/b.txt", "bin/run"} {
		if _, ok := entries[name]; !ok {
			t.Fatalf("zip 缺少条目: %s (实际: %v)", name, entries)
		}
	}
	if entries["a.txt"] != 29 {
		t.Fatalf("a.txt 大小异常: %d", entries["a.txt"])
	}

	// 可执行权限保留
	zr, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name == "bin/run" && f.Mode()&0o111 == 0 {
			t.Fatal("可执行文件权限丢失")
		}
	}
}

// newTestPlugin 生成最小插件项目（无第三方依赖，go build 无需联网）
func newTestPlugin(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "plugin.json"),
		`{"key": "demo", "name": "Demo", "version": "1.2.3"}`)
	writeTestFile(t, filepath.Join(dir, "go.mod"), "module demo-plugin\n\ngo 1.27.3\n")
	writeTestFile(t, filepath.Join(dir, "main.go"), "package main\n\nfunc main() {}\n")
	writeTestFile(t, filepath.Join(dir, "README.md"), "# demo")
	writeTestFile(t, filepath.Join(dir, "assets", "logo.txt"), "logo")
	// 伪造 go build 产物与源码文件：均不应入包
	writeTestFile(t, filepath.Join(dir, "demo-plugin"), "fake-binary")
	writeTestFile(t, filepath.Join(dir, "helper.go"), "package main")
	return dir
}

func TestPackPlugin_EndToEnd(t *testing.T) {
	dir := newTestPlugin(t)
	out := filepath.Join(t.TempDir(), "demo-1.2.3.zip")
	if err := packPlugin(dir, out, "linux-amd64", "."); err != nil {
		t.Fatalf("packPlugin: %v", err)
	}

	entries := zipEntries(t, out)
	for _, name := range []string{"plugin.json", "README.md", "assets/logo.txt", "bin/linux-amd64/plugin"} {
		if _, ok := entries[name]; !ok {
			t.Fatalf("zip 缺少: %s (实际: %v)", name, keys(entries))
		}
	}
	for _, name := range []string{"demo-plugin", "helper.go", "go.mod", "main.go"} {
		if _, ok := entries[name]; ok {
			t.Fatalf("zip 不应包含: %s", name)
		}
	}
	if entries["bin/linux-amd64/plugin"] < 100 {
		t.Fatalf("插件二进制异常: %d 字节", entries["bin/linux-amd64/plugin"])
	}
}

func TestPackPlugin_DefaultOutput(t *testing.T) {
	dir := newTestPlugin(t)
	// 输出为空 → 在当前工作目录生成 <key>-<version>.zip
	t.Chdir(t.TempDir())
	if err := packPlugin(dir, "", "linux-amd64", "."); err != nil {
		t.Fatalf("packPlugin: %v", err)
	}
	if _, err := os.Stat("demo-1.2.3.zip"); err != nil {
		t.Fatalf("默认输出文件未生成: %v", err)
	}
}

func TestPackPlugin_MultiPlatform(t *testing.T) {
	dir := newTestPlugin(t)
	out := filepath.Join(t.TempDir(), "multi.zip")
	if err := packPlugin(dir, out, "linux-amd64 windows-amd64", "."); err != nil {
		t.Fatalf("packPlugin: %v", err)
	}
	entries := zipEntries(t, out)
	if _, ok := entries["bin/linux-amd64/plugin"]; !ok {
		t.Fatal("缺少 linux 二进制")
	}
	if _, ok := entries["bin/windows-amd64/plugin.exe"]; !ok {
		t.Fatal("缺少 windows 二进制（plugin.exe）")
	}
}

func TestPackPlugin_Errors(t *testing.T) {
	// 缺少 plugin.json
	if err := packPlugin(t.TempDir(), "", "linux-amd64", "."); err == nil {
		t.Fatal("缺少 plugin.json 应报错")
	}
	// 缺少 key
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "plugin.json"), `{"name": "x"}`)
	if err := packPlugin(dir, filepath.Join(dir, "o.zip"), "linux-amd64", "."); err == nil {
		t.Fatal("缺少 key 应报错")
	}
	// 平台格式错误
	dir2 := newTestPlugin(t)
	err := packPlugin(dir2, filepath.Join(dir2, "o.zip"), "linux-amd64 badplatform", ".")
	if err == nil || !strings.Contains(err.Error(), "平台格式错误") {
		t.Fatalf("平台格式错误未识别: %v", err)
	}
	// 构建失败（main 包不存在）
	err = packPlugin(dir2, filepath.Join(dir2, "o2.zip"), "linux-amd64", "./not-exist")
	if err == nil {
		t.Fatal("构建失败应报错")
	}
}

func keys(m map[string]uint64) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
