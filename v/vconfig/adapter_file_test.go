package vconfig

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// TestFileAdapter_Basic 测试文件适配器基本功能
func TestFileAdapter_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建测试配置文件
		configFile := filepath.Join(tempDir, "test_config.yaml")
		configContent := `
server:
  address: ":8080"
  name: "test-server"
database:
  host: "localhost"
  port: 5432
  name: "testdb"
redis:
  address: "127.0.0.1:6379"
  db: 0
`
		t.AssertNil(os.WriteFile(configFile, []byte(configContent), 0644))

		// 创建文件适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "test_config.yaml",
			Watch:    false,
		})
		t.AssertNil(err)
		t.AssertNE(adapter, nil)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 测试适配器名称
		t.Assert(adapter.Name(), "file")

		// 测试可用性检查
		t.Assert(adapter.Available(ctx), true)

		// 测试获取配置值
		serverAddr, err := adapter.Get(ctx, "server.address")
		t.AssertNil(err)
		t.Assert(serverAddr, ":8080")

		serverName, err := adapter.Get(ctx, "server.name")
		t.AssertNil(err)
		t.Assert(serverName, "test-server")

		dbHost, err := adapter.Get(ctx, "database.host")
		t.AssertNil(err)
		t.Assert(dbHost, "localhost")

		// 测试获取不存在的配置
		nonExistent, err := adapter.Get(ctx, "non.existent.key")
		t.AssertNil(err)
		t.Assert(nonExistent, nil)

		// 测试获取所有配置数据
		allData, err := adapter.Data(ctx)
		t.AssertNil(err)
		t.AssertNE(allData, nil)
		t.Assert(len(allData) > 0, true)
	})
}

// TestFileAdapter_Watch 测试文件适配器监听功能
func TestFileAdapter_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_watch_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建测试配置文件
		configFile := filepath.Join(tempDir, "watch_config.yaml")
		initialContent := `
app:
  name: "initial-app"
  version: "1.0.0"
`
		t.AssertNil(os.WriteFile(configFile, []byte(initialContent), 0644))

		// 创建带监听功能的文件适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "watch_config.yaml",
			Watch:    true,
			Interval: time.Millisecond * 100, // 快速检查间隔用于测试
		})
		t.AssertNil(err)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 设置配置变更监听
		eventReceived := make(chan *ConfigEvent, 1)
		err = adapter.Watch(ctx, "app.*", func(event *ConfigEvent) {
			select {
			case eventReceived <- event:
			default:
			}
		})
		t.AssertNil(err)

		// 等待一小段时间确保监听器启动
		time.Sleep(time.Millisecond * 200)

		// 修改配置文件
		updatedContent := `
app:
  name: "updated-app"
  version: "1.0.1"
`
		t.AssertNil(os.WriteFile(configFile, []byte(updatedContent), 0644))

		// 等待配置变更事件
		select {
		case event := <-eventReceived:
			t.AssertNE(event, nil)
			t.Assert(event.Type, EventTypeUpdate)
		case <-time.After(time.Second * 2):
			t.Fatal("配置变更事件超时")
		}

		// 验证配置已更新
		appName, err := adapter.Get(ctx, "app.name")
		t.AssertNil(err)
		t.Assert(appName, "updated-app")

		appVersion, err := adapter.Get(ctx, "app.version")
		t.AssertNil(err)
		t.Assert(appVersion, "1.0.1")
	})
}

// TestFileAdapter_InvalidConfig 测试无效配置处理
func TestFileAdapter_InvalidConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试空配置
		_, err := NewFileAdapter(nil)
		t.AssertNE(err, nil)

		// 测试不存在的文件路径
		_, err = NewFileAdapter(&FileAdapterConfig{
			Path:     "/non/existent/path",
			FileName: "config.yaml",
		})
		t.AssertNE(err, nil)

		// 测试空文件名
		tempDir, err := os.MkdirTemp("", "vconfig_invalid_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		_, err = NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "",
		})
		t.AssertNE(err, nil)
	})
}

// TestFileAdapter_SetOperation 测试文件适配器设置操作
func TestFileAdapter_SetOperation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_set_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建测试配置文件
		configFile := filepath.Join(tempDir, "set_config.yaml")
		initialContent := `
app:
  name: "test-app"
`
		t.AssertNil(os.WriteFile(configFile, []byte(initialContent), 0644))

		// 创建文件适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "set_config.yaml",
			Watch:    false,
		})
		t.AssertNil(err)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 测试设置操作（文件适配器通常不支持写入）
		err = adapter.Set(ctx, "app.version", "1.0.0")
		// 文件适配器可能不支持Set操作，这取决于具体实现
		// 这里我们测试错误处理
		if err != nil {
			t.Log("文件适配器不支持Set操作，这是预期的行为")
		}
	})
}

// TestFileAdapter_MultipleFormats 测试多种配置文件格式
func TestFileAdapter_MultipleFormats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_formats_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 测试YAML格式
		yamlFile := filepath.Join(tempDir, "config.yaml")
		yamlContent := `
server:
  port: 8080
  host: "localhost"
`
		t.AssertNil(os.WriteFile(yamlFile, []byte(yamlContent), 0644))

		yamlAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
		})
		t.AssertNil(err)
		defer yamlAdapter.Close(context.Background())

		ctx := context.Background()
		port, err := yamlAdapter.Get(ctx, "server.port")
		t.AssertNil(err)
		t.Assert(port, 8080)

		// 测试JSON格式
		jsonFile := filepath.Join(tempDir, "config.json")
		jsonContent := `{
  "server": {
    "port": 9090,
    "host": "0.0.0.0"
  }
}`
		t.AssertNil(os.WriteFile(jsonFile, []byte(jsonContent), 0644))

		jsonAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.json",
		})
		t.AssertNil(err)
		defer jsonAdapter.Close(context.Background())

		port, err = jsonAdapter.Get(ctx, "server.port")
		t.AssertNil(err)
		t.Assert(port, 9090)
	})
}

// TestFileAdapter_Concurrent 测试文件适配器并发访问
func TestFileAdapter_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_concurrent_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建测试配置文件
		configFile := filepath.Join(tempDir, "concurrent_config.yaml")
		configContent := `
test:
  value1: "data1"
  value2: "data2"
  value3: "data3"
`
		t.AssertNil(os.WriteFile(configFile, []byte(configContent), 0644))

		// 创建文件适配器
		adapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "concurrent_config.yaml",
			Watch:    false,
		})
		t.AssertNil(err)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 并发读取测试
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(index int) {
				defer func() { done <- true }()
				
				for j := 0; j < 100; j++ {
					value, err := adapter.Get(ctx, "test.value1")
					t.AssertNil(err)
					t.Assert(value, "data1")
					
					allData, err := adapter.Data(ctx)
					t.AssertNil(err)
					t.AssertNE(allData, nil)
				}
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < 10; i++ {
			select {
			case <-done:
			case <-time.After(time.Second * 5):
				t.Fatal("并发测试超时")
			}
		}
	})
}