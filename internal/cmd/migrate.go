package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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
	path := parser.GetOpt("path", "resource/migrates").String()

	database := parser.GetOpt("database", "postgres://postgres:123456@localhost:5432/vgo?sslmode=disable").String()
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

// runMigrateForModules 为多个模块运行迁移
// 参数: ctx - 上下文, pathPattern - 路径模式, database - 数据库连接, operation - 操作类型, steps - 步数
// 返回值: error - 错误信息
func runMigrateForModules(ctx context.Context, pathPattern, database, operation string, steps int) error {
	// 获取模块目录
	basePath := strings.Replace(pathPattern, "/*/resource/migrates", "", 1)
	modules, err := getModules(basePath)
	if err != nil {
		return fmt.Errorf("获取模块列表失败: %v", err)
	}

	g.Log().Infof(ctx, "找到 %d 个模块", len(modules))

	// 为每个模块运行迁移
	for _, module := range modules {
		migratePath := filepath.Join(basePath, module, "resource", "migrates")
		if !gfile.Exists(migratePath) {
			g.Log().Warningf(ctx, "模块 %s 的迁移目录不存在: %s", module, migratePath)
			continue
		}

		g.Log().Infof(ctx, "正在迁移模块: %s", module)
		if err := runMigrate(ctx, migratePath, database, operation, steps); err != nil {
			g.Log().Errorf(ctx, "模块 %s 迁移失败: %v", module, err)
			return err
		}
		g.Log().Infof(ctx, "模块 %s 迁移完成", module)
	}

	return nil
}

// runMigrate 执行数据库迁移
// 参数: ctx - 上下文, migratePath - 迁移文件路径, database - 数据库连接, operation - 操作类型, steps - 步数
// 返回值: error - 错误信息
func runMigrate(ctx context.Context, migratePath, database, operation string, steps int) error {
	// 检查迁移目录是否存在
	if !gfile.Exists(migratePath) {
		return fmt.Errorf("迁移目录不存在: %s", migratePath)
	}

	// 获取迁移文件
	files, err := gfile.ScanDir(migratePath, "*.sql", false)
	if err != nil {
		return fmt.Errorf("扫描迁移文件失败: %v", err)
	}

	if len(files) == 0 {
		g.Log().Warningf(ctx, "迁移目录中没有找到SQL文件: %s", migratePath)
		return nil
	}

	g.Log().Infof(ctx, "找到 %d 个迁移文件", len(files))

	// 创建数据库连接
	db, err := sql.Open("postgres", database)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	if pingErr := db.Ping(); pingErr != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 创建postgres驱动实例
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("创建数据库驱动失败: %v", err)
	}

	// 创建migrate实例
	sourceURL := fmt.Sprintf("file://%s", migratePath)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("创建迁移实例失败: %v", err)
	}
	defer m.Close()

	// 执行迁移操作
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
		g.Log().Infof(ctx, "当前版本: %d, 状态: %v", version, dirty)
		return nil
	case "force":
		if steps == 0 {
			return fmt.Errorf("force操作需要指定版本号")
		}
		err = m.Force(steps)
	default:
		return fmt.Errorf("不支持的操作类型: %s", operation)
	}

	// 处理迁移结果
	if err != nil {
		if err == migrate.ErrNoChange {
			g.Log().Infof(ctx, "没有需要执行的迁移")
			return nil
		}
		return fmt.Errorf("执行迁移失败: %v", err)
	}

	g.Log().Infof(ctx, "迁移执行成功")
	return nil
}

// getModules 获取模块列表
// 参数: basePath - 基础路径
// 返回值: []string - 模块名列表, error - 错误信息
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
		if gfile.IsDir(filepath.Join(basePath, dir)) {
			modules = append(modules, dir)
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
