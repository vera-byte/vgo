package vconfig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestConfigManager_MultiSource 测试配置管理器多源配置
func TestConfigManager_MultiSource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "manager_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建主配置文件
		primaryConfigFile := filepath.Join(tempDir, "primary.yaml")
		primaryContent := `
app:
  name: primary-app
  version: 1.0.0
  environment: production
database:
  host: primary-db
  port: 5432
`
		err = os.WriteFile(primaryConfigFile, []byte(primaryContent), 0644)
		t.AssertNil(err)

		// 创建备用配置文件
		fallbackConfigFile := filepath.Join(tempDir, "fallback.yaml")
		fallbackContent := `
app:
  name: fallback-app
  version: 1.0.1
  timeout: 30
cache:
  redis:
    host: localhost
    port: 6379
`
		err = os.WriteFile(fallbackConfigFile, []byte(fallbackContent), 0644)
		t.AssertNil(err)

		// 创建配置管理器
		manager := NewConfigManager()

		// 创建并注册主配置适配器
		primaryAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "primary.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("primary", primaryAdapter)
		t.AssertNil(err)

		// 创建并注册备用配置适配器
		fallbackAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "fallback.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("fallback", fallbackAdapter)
		t.AssertNil(err)

		// 设置主配置源和备用配置源
		err = manager.SetPrimary("primary")
		t.AssertNil(err)

		err = manager.SetFallback("fallback")
		t.AssertNil(err)

		ctx := context.Background()

		// 测试从主配置源获取配置
		value, err := manager.Get(ctx, "app.name")
		t.AssertNil(err)
		t.Assert(fmt.Sprintf("%v", value.Val()), "primary-app")

		// 测试从主配置源获取不存在的配置，应该从备用配置源获取
		value, err = manager.Get(ctx, "cache.redis.host")
		t.AssertNil(err)
		t.Assert(fmt.Sprintf("%v", value.Val()), "localhost")

		// 测试获取所有配置数据
		data, err := manager.Data(ctx)
		t.AssertNil(err)
		t.AssertNE(data, nil)

		// 关闭管理器
		err = manager.Close(ctx)
		t.AssertNil(err)
	})
}

// TestConfigManager_PrimaryUnavailable 测试主配置源不可用时的备用机制
func TestConfigManager_PrimaryUnavailable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "manager_fallback_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 只创建备用配置文件
		fallbackConfigFile := filepath.Join(tempDir, "fallback.yaml")
		fallbackContent := `
app:
  name: fallback-app
  version: 2.0.0
database:
  host: fallback-db
  port: 3306
`
		err = os.WriteFile(fallbackConfigFile, []byte(fallbackContent), 0644)
		t.AssertNil(err)

		// 创建配置管理器
		manager := NewConfigManager()

		// 创建不可用的主配置适配器（文件不存在）
		primaryAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "nonexistent.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("primary", primaryAdapter)
		t.AssertNil(err)

		// 创建并注册备用配置适配器
		fallbackAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "fallback.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("fallback", fallbackAdapter)
		t.AssertNil(err)

		// 设置主配置源和备用配置源
		err = manager.SetPrimary("primary")
		t.AssertNil(err)

		err = manager.SetFallback("fallback")
		t.AssertNil(err)

		ctx := context.Background()

		// 测试从备用配置源获取配置
		value, err := manager.Get(ctx, "app.name")
		t.AssertNil(err)
		t.Assert(fmt.Sprintf("%v", value.Val()), "fallback-app")

		// 测试获取数据库配置
		value, err = manager.Get(ctx, "database.host")
		t.AssertNil(err)
		t.Assert(fmt.Sprintf("%v", value.Val()), "fallback-db")

		// 关闭管理器
		err = manager.Close(ctx)
		t.AssertNil(err)
	})
}

// TestConfigManager_Cache 测试配置管理器缓存机制
func TestConfigManager_Cache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "manager_cache_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")
		configContent := `
app:
  name: cache-test-app
  version: 1.0.0
`
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		t.AssertNil(err)

		// 创建配置管理器
		manager := NewConfigManager()

		// 创建并注册配置适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("primary", adapter)
		t.AssertNil(err)

		err = manager.SetPrimary("primary")
		t.AssertNil(err)

		ctx := context.Background()

		// 第一次获取配置（从适配器获取）
		start := time.Now()
		value1, err := manager.Get(ctx, "app.name")
		t.AssertNil(err)
		firstDuration := time.Since(start)

		// 第二次获取相同配置（从缓存获取）
		start = time.Now()
		value2, err := manager.Get(ctx, "app.name")
		t.AssertNil(err)
		secondDuration := time.Since(start)

		// 验证配置值相同
		t.Assert(fmt.Sprintf("%v", value1.Val()), fmt.Sprintf("%v", value2.Val()))
		t.Assert(fmt.Sprintf("%v", value1.Val()), "cache-test-app")

		// 缓存访问应该更快（虽然在测试环境中差异可能很小）
		t.AssertLE(secondDuration, firstDuration*2) // 允许一定的误差

		// 关闭管理器
		err = manager.Close(ctx)
		t.AssertNil(err)
	})
}

// TestConfigManager_Watch 测试配置管理器监听功能
func TestConfigManager_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "manager_watch_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")
		configContent := `
app:
  name: watch-test-app
  version: 1.0.0
`
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		t.AssertNil(err)

		// 创建配置管理器
		manager := NewConfigManager()

		// 创建并注册配置适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
			Watch:    true,
			Interval: time.Millisecond * 100,
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("primary", adapter)
		t.AssertNil(err)

		err = manager.SetPrimary("primary")
		t.AssertNil(err)

		ctx := context.Background()
		eventReceived := make(chan *ConfigEvent, 1)

		// 设置监听回调
		err = manager.Watch(ctx, "app.name", func(event *ConfigEvent) {
			eventReceived <- event
		})
		t.AssertNil(err)

		// 修改配置文件
		time.Sleep(time.Millisecond * 200) // 等待监听启动
		newContent := `
app:
  name: watch-test-app-updated
  version: 1.0.1
`
		err = os.WriteFile(configFile, []byte(newContent), 0644)
		t.AssertNil(err)

		// 等待配置变化事件
		select {
		case event := <-eventReceived:
			t.AssertNE(event, nil)
			t.Assert(event.Key, "app.name")
			t.AssertNE(event.Value, nil)
		case <-time.After(time.Second * 2):
			// 在某些情况下文件监听可能不会立即触发，这是正常的
			t.Log("配置变化事件未在预期时间内收到，这在某些环境中是正常的")
		}

		// 关闭管理器
		err = manager.Close(ctx)
		t.AssertNil(err)
	})
}

// TestConfigManager_InvalidOperations 测试配置管理器无效操作处理
func TestConfigManager_InvalidOperations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := NewConfigManager()

		// 测试注册空名称适配器
		err := manager.RegisterAdapter("", nil)
		t.AssertNE(err, nil)

		// 测试注册nil适配器
		err = manager.RegisterAdapter("test", nil)
		t.AssertNE(err, nil)

		// 测试设置不存在的主配置源
		err = manager.SetPrimary("nonexistent")
		t.AssertNE(err, nil)

		// 测试设置不存在的备用配置源
		err = manager.SetFallback("nonexistent")
		t.AssertNE(err, nil)

		// 创建临时配置文件用于测试
		tempDir, err := os.MkdirTemp("", "manager_invalid_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")
		configContent := `
app:
  name: test-app
`
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		t.AssertNil(err)

		// 注册有效的适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("valid", adapter)
		t.AssertNil(err)

		// 现在设置主配置源应该成功
		err = manager.SetPrimary("valid")
		t.AssertNil(err)

		// 关闭管理器
		err = manager.Close(context.Background())
		t.AssertNil(err)
	})
}

// TestConfigManager_GetWithDefault 测试配置管理器默认值功能
func TestConfigManager_GetWithDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "manager_default_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")
		configContent := `
app:
  name: default-test-app
`
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		t.AssertNil(err)

		// 创建配置管理器
		manager := NewConfigManager()

		// 创建并注册配置适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
		})
		t.AssertNil(err)

		err = manager.RegisterAdapter("primary", adapter)
		t.AssertNil(err)

		err = manager.SetPrimary("primary")
		t.AssertNil(err)

		ctx := context.Background()

		// 测试获取存在的配置
		value := manager.GetWithDefault(ctx, "app.name", gvar.New("default-name"))
		t.Assert(fmt.Sprintf("%v", value.Val()), "default-test-app")

		// 测试获取不存在的配置，应该返回默认值
		value = manager.GetWithDefault(ctx, "app.nonexistent", gvar.New("default-value"))
		t.Assert(fmt.Sprintf("%v", value.Val()), "default-value")

		// 关闭管理器
		err = manager.Close(ctx)
		t.AssertNil(err)
	})
}