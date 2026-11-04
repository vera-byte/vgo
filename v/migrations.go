package v

import (
	"io/fs"
	"sort"
	"sync"
)

// MigrationSource 模块内嵌迁移源
type MigrationSource struct {
	Name string
	FS   fs.FS
}

var migrationSources sync.Map // module name -> fs.FS

// RegisterMigrationFS 注册模块内嵌迁移目录（//go:embed resource/migrates）。
// migrate 命令据此在任意项目布局（含模块以依赖形式引入的生成项目）下发现并执行模块迁移。
func RegisterMigrationFS(module string, fsys fs.FS) {
	migrationSources.Store(module, fsys)
}

// MigrationFS 返回按模块名排序的内嵌迁移集合（base < demo < dict < ... 保证依赖次序）。
func MigrationFS() []MigrationSource {
	out := []MigrationSource{}
	migrationSources.Range(func(k, val any) bool {
		out = append(out, MigrationSource{Name: k.(string), FS: val.(fs.FS)})
		return true
	})
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
