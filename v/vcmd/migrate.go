package vcmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/vera-byte/vgo/v"
)

func init() {

	migrateCmd := &gcmd.Command{
		Name:  "migrate",
		Usage: "migrate",
		Brief: "数据库迁移",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return MigrateFunc(ctx, parser)
		},
	}
	Root.AddCommand(migrateCmd)

}

// MigrateFunc 数据库迁移命令处理函数
// 参数: ctx - 上下文对象, parser - 命令行解析器
// 返回值: error - 错误信息
func MigrateFunc(ctx context.Context, parser *gcmd.Parser) error {
	// 获取命令行参数
	path := parser.GetOpt("path", "modules/*/resource/migrates").String()

	database := parser.GetOpt("database", defaultMigrateDSN(ctx)).String()
	operation := parser.GetOpt("operation", "up").String()
	steps := parser.GetOpt("steps", "0").Int()

	g.Log().Infof(ctx, "开始数据库迁移...")
	g.Log().Infof(ctx, "迁移路径: %s", path)
	g.Log().Infof(ctx, "数据库连接: %s", maskPassword(database))
	g.Log().Infof(ctx, "操作类型: %s", operation)

	// 如果路径在gres资源中，则从gres中读取
	// 检查更通用的路径模式
	resourcePath := path
	if !gfile.Exists(path) && (gres.Contains(path) || gres.Contains("resource/"+path) || gres.Contains("resources/"+path)) {
		// 尝试不同的路径格式
		if gres.Contains(path) {
			resourcePath = path
		} else if gres.Contains("resource/" + path) {
			resourcePath = "resource/" + path
		} else if gres.Contains("resources/" + path) {
			resourcePath = "resources/" + path
		}

		tempPath, err := extractMigrationsFromGres(resourcePath)
		if err != nil {
			return fmt.Errorf("从gres中提取迁移文件失败: %v", err)
		}
		defer gfile.Remove(tempPath)
		path = tempPath
	}

	// 如果路径包含通配符，处理多个模块
	if strings.Contains(path, "*") {
		return runMigrateForModules(ctx, path, database, operation, steps)
	}

	// 单个路径迁移
	return runMigrate(ctx, path, database, operation, steps)
}

// extractMigrationsFromGres 从gres资源中提取迁移文件到临时目录
// 参数: path - gres中的路径
// 返回值: string - 临时目录路径, error - 错误信息
func extractMigrationsFromGres(path string) (string, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 从gres中获取文件列表
	files := gres.ScanDirFile(path, "*", false)
	if len(files) == 0 {
		// 尝试添加前缀再次查找
		pathsToTry := []string{path, "resource/" + path, "resources/" + path}
		for _, tryPath := range pathsToTry {
			files = gres.ScanDirFile(tryPath, "*", false)
			if len(files) > 0 {
				break
			}
		}

		if len(files) == 0 {
			return "", fmt.Errorf("在gres中未找到迁移文件: %s", path)
		}
	}

	// 提取文件到临时目录
	for _, file := range files {
		// 获取文件名（去掉路径前缀）
		fileName := filepath.Base(file.Name())
		filePath := filepath.Join(tempDir, fileName)
		if err := gfile.PutBytes(filePath, file.Content()); err != nil {
			return "", fmt.Errorf("写入文件失败 %s: %v", filePath, err)
		}
	}

	g.Log().Infof(context.Background(), "从gres中提取了 %d 个文件到临时目录: %s", len(files), tempDir)

	// 记录提取的文件名
	fileList, _ := gfile.ScanDir(tempDir, "*", false)
	g.Log().Infof(context.Background(), "临时目录中的文件: %v", fileList)

	return tempDir, nil
}

// runMigrateForModules 为多个模块运行迁移。
// 迁移源 = 本地 modules/*/resource/migrates 目录 ∪ 各模块内嵌迁移（v.RegisterMigrationFS），
// 按模块名排序执行（base < demo < dict < ...，保证 update_updated_time_column 等公共对象先建）。
func runMigrateForModules(ctx context.Context, pathPattern, database, operation string, steps int) error {
	type migSrc struct {
		name string
		dir  string // 本地目录，空则用 fs
		fs   fs.FS
	}
	srcs := map[string]migSrc{}

	// 本地模块目录
	basePath := strings.Replace(pathPattern, "/*/resource/migrates", "", 1)
	if modules, err := getModules(basePath); err == nil {
		for _, module := range modules {
			migratePath := filepath.Join(basePath, module, "resource", "migrates")
			if gfile.Exists(migratePath) {
				srcs[module] = migSrc{name: module, dir: migratePath}
			}
		}
	}

	// 模块内嵌迁移（生成项目中框架模块以依赖形式引入，本地无目录）
	for _, ms := range v.MigrationFS() {
		if _, ok := srcs[ms.Name]; !ok {
			srcs[ms.Name] = migSrc{name: ms.Name, fs: ms.FS}
		}
	}

	names := make([]string, 0, len(srcs))
	for name := range srcs {
		names = append(names, name)
	}
	sort.Strings(names)

	g.Log().Infof(ctx, "找到 %d 个模块: %v", len(names), names)

	for _, name := range names {
		src := srcs[name]
		g.Log().Infof(ctx, "正在迁移模块: %s", name)
		var err error
		if src.dir != "" {
			err = runMigrate(ctx, src.dir, database, operation, steps)
		} else {
			err = runMigrateEmbedded(ctx, src.fs, name, database, operation, steps)
		}
		if err != nil {
			g.Log().Errorf(ctx, "模块 %s 迁移失败: %v", name, err)
			return err
		}
		g.Log().Infof(ctx, "模块 %s 迁移完成", name)
	}

	return nil
}


// defaultMigrateDSN 从 config.yaml 的 database.default 生成默认 PostgreSQL DSN，
// 使 `go run main.go migrate` 无需显式传参即可运行。
func defaultMigrateDSN(ctx context.Context) string {
	node, err := g.Cfg().Get(ctx, "database.default")
	if err != nil || node.IsEmpty() {
		return ""
	}
	cfg := node.Map()
	host := gconv.String(cfg["host"])
	port := gconv.String(cfg["port"])
	user := gconv.String(cfg["user"])
	pass := gconv.String(cfg["pass"])
	name := gconv.String(cfg["name"])
	if host == "" || name == "" {
		return ""
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(user), url.QueryEscape(pass), host, port, name)
}

// moduleNameFromMigratePath 从迁移目录路径（如 modules/base/resource/migrates）中提取模块名（base）。
func moduleNameFromMigratePath(migratePath string) string {
	parts := strings.Split(filepath.ToSlash(migratePath), "/")
	for i, part := range parts {
		if part == "resource" && i > 0 {
			return parts[i-1]
		}
	}
	return ""
}

// runMigrate 执行数据库迁移（本地文件目录源）
// 参数: ctx - 上下文, migratePath - 迁移文件路径, database - 数据库连接, operation - 操作类型, steps - 步数
// 返回值: error - 错误信息
func runMigrate(ctx context.Context, migratePath, database, operation string, steps int) error {
	if !gfile.Exists(migratePath) {
		return fmt.Errorf("迁移目录不存在: %s", migratePath)
	}
	files, err := gfile.ScanDir(migratePath, "*.sql", false)
	if err != nil {
		return fmt.Errorf("扫描迁移文件失败: %v", err)
	}
	if len(files) == 0 {
		g.Log().Warningf(ctx, "迁移目录中没有找到SQL文件: %s", migratePath)
		return nil
	}
	g.Log().Infof(ctx, "找到 %d 个迁移文件", len(files))

	m, err := newMigrateFromSourceURL(fmt.Sprintf("file://%s", migratePath), database, moduleNameFromMigratePath(migratePath))
	if err != nil {
		return err
	}
	defer m.Close()
	return execMigrateOp(m, operation, steps)
}

// runMigrateEmbedded 执行模块内嵌迁移（go:embed 资源源）
func runMigrateEmbedded(ctx context.Context, fsys fs.FS, moduleName, database, operation string, steps int) error {
	src, err := iofs.New(fsys, "resource/migrates")
	if err != nil {
		return fmt.Errorf("创建内嵌迁移源失败: %v", err)
	}
	g.Log().Infof(ctx, "使用模块 %s 的内嵌迁移", moduleName)
	m, err := newMigrateFromSourceDriver(src, "iofs", database, moduleName)
	if err != nil {
		return err
	}
	defer m.Close()
	return execMigrateOp(m, operation, steps)
}

// newMigrateFromSourceDriver 创建 migrate 实例（每模块独立版本表，避免版本号冲突）
func newMigrateFromSourceDriver(src source.Driver, sourceName, database, moduleName string) (*migrate.Migrate, error) {
	driver, err := newPGDriver(database, moduleName)
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithInstance(sourceName, src, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("创建迁移实例失败: %v", err)
	}
	return m, nil
}

// newMigrateFromSourceURL 经 source registry（file:// 等）创建 migrate 实例
func newMigrateFromSourceURL(sourceURL, database, moduleName string) (*migrate.Migrate, error) {
	driver, err := newPGDriver(database, moduleName)
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("创建迁移实例失败: %v", err)
	}
	return m, nil
}

// newPGDriver 创建 PostgreSQL 迁移驱动（每模块独立版本表）
func newPGDriver(dsn, moduleName string) (database.Driver, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}
	if pingErr := db.Ping(); pingErr != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %v", pingErr)
	}
	migrationsTable := "schema_migrations"
	if moduleName != "" {
		migrationsTable = "schema_migrations_" + moduleName
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: migrationsTable})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("创建数据库驱动失败: %v", err)
	}
	return driver, nil
}

// execMigrateOp 执行迁移操作
func execMigrateOp(m *migrate.Migrate, operation string, steps int) error {
	var err error
	switch operation {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	case "version":
		version, dirty, verr := m.Version()
		if verr != nil {
			return fmt.Errorf("获取版本信息失败: %v", verr)
		}
		g.Log().Infof(context.Background(), "当前版本: %d, 状态: %v", version, dirty)
		return nil
	case "force":
		if steps == 0 {
			return fmt.Errorf("force操作需要指定版本号")
		}
		err = m.Force(steps)
	default:
		return fmt.Errorf("不支持的操作类型: %s", operation)
	}
	if err != nil {
		if err == migrate.ErrNoChange {
			g.Log().Infof(context.Background(), "没有需要执行的迁移")
			return nil
		}
		return fmt.Errorf("执行迁移失败: %v", err)
	}
	g.Log().Infof(context.Background(), "迁移执行成功")
	return nil
}

func getModules(basePath string) ([]string, error) {
	if !gfile.Exists(basePath) {
		return nil, fmt.Errorf("模块目录不存在: %s", basePath)
	}

	dirs, err := gfile.ScanDir(basePath, "*", false)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, dir := range dirs {
		if gfile.IsDir(dir) {
			modules = append(modules, filepath.Base(dir))
		}
	}

	return modules, nil
}

// maskPassword 隐藏数据库连接字符串中的密码
// 参数: dsn - 数据库连接字符串
// 返回值: string - 隐藏密码后的连接字符串
func maskPassword(dsn string) string {
	// 简单的密码隐藏逻辑
	if strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		if len(parts) >= 2 {
			userPart := parts[0]
			if strings.Contains(userPart, ":") {
				userParts := strings.Split(userPart, ":")
				if len(userParts) >= 2 {
					userParts[len(userParts)-1] = "****"
					parts[0] = strings.Join(userParts, ":")
				}
			}
			return strings.Join(parts, "@")
		}
	}
	return dsn
}
