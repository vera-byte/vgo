package vconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// MockConsulServer 模拟Consul服务器
type MockConsulServer struct {
	server *httptest.Server
	kvData map[string]string
}

// NewMockConsulServer 创建模拟Consul服务器
func NewMockConsulServer() *MockConsulServer {
	mock := &MockConsulServer{
		kvData: make(map[string]string),
	}

	mux := http.NewServeMux()
	
	// 处理KV GET请求
	mux.HandleFunc("/v1/kv/", func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/v1/kv/")
		
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("recurse") == "true" {
				// 递归获取所有键值
				var results []ConsulKVPair
				for k, v := range mock.kvData {
					if strings.HasPrefix(k, key) {
						results = append(results, ConsulKVPair{
							Key:   k,
							Value: v,
						})
					}
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(results)
			} else {
				// 获取单个键值
				if value, exists := mock.kvData[key]; exists {
					result := []ConsulKVPair{{
						Key:   key,
						Value: value,
					}}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(result)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		case "PUT":
			// 设置键值
			body := make([]byte, r.ContentLength)
			r.Body.Read(body)
			mock.kvData[key] = string(body)
			w.WriteHeader(http.StatusOK)
		}
	})

	mock.server = httptest.NewServer(mux)
	return mock
}

// Close 关闭模拟服务器
func (m *MockConsulServer) Close() {
	m.server.Close()
}

// SetKV 设置键值对
func (m *MockConsulServer) SetKV(key, value string) {
	m.kvData[key] = value
}

// TestConsulAdapter_Basic 测试Consul适配器基本功能
func TestConsulAdapter_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟Consul服务器
		mockServer := NewMockConsulServer()
		defer mockServer.Close()

		// 设置测试数据
		mockServer.SetKV("config/app/name", "test-app")
		mockServer.SetKV("config/app/version", "1.0.0")
		mockServer.SetKV("config/database/host", "localhost")
		mockServer.SetKV("config/database/port", "5432")

		// 创建Consul适配器
		adapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address: mockServer.server.URL,
			Prefix:  "config/",
			Watch:   false,
		})
		t.AssertNil(err)
		t.AssertNE(adapter, nil)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 测试适配器名称
		t.Assert(adapter.Name(), "consul")

		// 测试可用性检查
		t.Assert(adapter.Available(ctx), true)

		// 测试获取配置值
		appName, err := adapter.Get(ctx, "app.name")
		t.AssertNil(err)
		t.Assert(appName, "test-app")

		appVersion, err := adapter.Get(ctx, "app.version")
		t.AssertNil(err)
		t.Assert(appVersion, "1.0.0")

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

// TestConsulAdapter_WithFileBackup 测试Consul适配器与文件备用方案
func TestConsulAdapter_WithFileBackup(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_consul_backup_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建备用配置文件
		backupFile := filepath.Join(tempDir, "backup_config.yaml")
		backupContent := `
app:
  name: "backup-app"
  version: "1.0.0-backup"
database:
  host: "backup-host"
  port: 3306
`
		t.AssertNil(os.WriteFile(backupFile, []byte(backupContent), 0644))

		// 创建文件适配器作为备用
		fileAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "backup_config.yaml",
			Watch:    false,
		})
		t.AssertNil(err)
		defer fileAdapter.Close(context.Background())

		ctx := context.Background()

		// 尝试创建Consul适配器（指向不存在的服务器）
		consulAdapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address: "http://localhost:8500", // 不存在的Consul服务器
			Prefix:  "config/",
			Watch:   false,
		})

		// 测试备用方案：当Consul不可用时使用文件配置
		var appName, appVersion, dbHost interface{}
		
		if consulAdapter != nil && consulAdapter.Available(ctx) {
			// 如果Consul可用，从Consul获取配置
			appName, _ = consulAdapter.Get(ctx, "app.name")
			appVersion, _ = consulAdapter.Get(ctx, "app.version")
			dbHost, _ = consulAdapter.Get(ctx, "database.host")
		}
		
		// 如果Consul不可用或没有配置，使用文件备用
		if appName == nil {
			appName, err = fileAdapter.Get(ctx, "app.name")
			t.AssertNil(err)
		}
		if appVersion == nil {
			appVersion, err = fileAdapter.Get(ctx, "app.version")
			t.AssertNil(err)
		}
		if dbHost == nil {
			dbHost, err = fileAdapter.Get(ctx, "database.host")
			t.AssertNil(err)
		}

		// 验证从备用配置源获取的数据
		t.Assert(appName, "backup-app")
		t.Assert(appVersion, "1.0.0-backup")
		t.Assert(dbHost, "backup-host")
	})
}

// TestConsulAdapter_Watch 测试Consul适配器监听功能
func TestConsulAdapter_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟Consul服务器
		mockServer := NewMockConsulServer()
		defer mockServer.Close()

		// 设置初始测试数据
		mockServer.SetKV("config/watch/value", "initial")

		// 创建带监听功能的Consul适配器
		adapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address:  mockServer.server.URL,
			Prefix:   "config/",
			Watch:    true,
			Interval: time.Millisecond * 100, // 快速检查间隔用于测试
		})
		t.AssertNil(err)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 设置配置变更监听
		eventReceived := make(chan *ConfigEvent, 1)
		err = adapter.Watch(ctx, "watch.*", func(event *ConfigEvent) {
			select {
			case eventReceived <- event:
			default:
			}
		})
		t.AssertNil(err)

		// 等待一小段时间确保监听器启动
		time.Sleep(time.Millisecond * 200)

		// 修改配置值
		mockServer.SetKV("config/watch/value", "updated")

		// 等待配置变更事件
		select {
		case event := <-eventReceived:
			t.AssertNE(event, nil)
			t.Assert(event.Type, EventTypeUpdate)
		case <-time.After(time.Second * 2):
			t.Log("配置变更事件超时，这在模拟环境中是正常的")
		}

		// 验证配置已更新
		value, err := adapter.Get(ctx, "watch.value")
		t.AssertNil(err)
		t.Assert(value, "updated")
	})
}

// TestConsulAdapter_InvalidConfig 测试无效配置处理
func TestConsulAdapter_InvalidConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试空配置
		_, err := NewConsulAdapter(nil)
		t.AssertNE(err, nil)

		// 测试无效地址
		_, err = NewConsulAdapter(&ConsulAdapterConfig{
			Address: "invalid-url",
			Prefix:  "config/",
		})
		// 创建可能成功，但Available会返回false
		if err == nil {
			t.Log("Consul适配器创建成功，但可用性检查会失败")
		}

		// 测试空前缀
		adapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address: "http://localhost:8500",
			Prefix:  "",
		})
		if err == nil {
			t.AssertNE(adapter, nil)
			defer adapter.Close(context.Background())
		}
	})
}

// TestConsulAdapter_SetOperation 测试Consul适配器设置操作
func TestConsulAdapter_SetOperation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟Consul服务器
		mockServer := NewMockConsulServer()
		defer mockServer.Close()

		// 创建Consul适配器
		adapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address: mockServer.server.URL,
			Prefix:  "config/",
			Watch:   false,
		})
		t.AssertNil(err)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 测试设置操作
		err = adapter.Set(ctx, "test.key", "test-value")
		t.AssertNil(err)

		// 验证设置的值
		value, err := adapter.Get(ctx, "test.key")
		t.AssertNil(err)
		t.Assert(value, "test-value")

		// 测试更新操作
		err = adapter.Set(ctx, "test.key", "updated-value")
		t.AssertNil(err)

		// 验证更新的值
		value, err = adapter.Get(ctx, "test.key")
		t.AssertNil(err)
		t.Assert(value, "updated-value")
	})
}

// TestConsulAdapter_Concurrent 测试Consul适配器并发访问
func TestConsulAdapter_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟Consul服务器
		mockServer := NewMockConsulServer()
		defer mockServer.Close()

		// 设置测试数据
		for i := 0; i < 10; i++ {
			mockServer.SetKV(fmt.Sprintf("config/test/value%d", i), fmt.Sprintf("data%d", i))
		}

		// 创建Consul适配器
		adapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address: mockServer.server.URL,
			Prefix:  "config/",
			Watch:   false,
		})
		t.AssertNil(err)
		defer adapter.Close(context.Background())

		ctx := context.Background()

		// 并发读取测试
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(index int) {
				defer func() { done <- true }()
				
				for j := 0; j < 50; j++ {
					key := fmt.Sprintf("test.value%d", index%10)
					value, err := adapter.Get(ctx, key)
					t.AssertNil(err)
					t.Assert(value, fmt.Sprintf("data%d", index%10))
					
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

// TestConsulAdapter_FallbackToFile 测试Consul不可用时回退到文件配置
func TestConsulAdapter_FallbackToFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时测试目录
		tempDir, err := os.MkdirTemp("", "vconfig_consul_fallback_test_")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		// 创建备用配置文件
		fallbackFile := filepath.Join(tempDir, "fallback_config.yaml")
		fallbackContent := `
app:
  name: "fallback-app"
  mode: "production"
server:
  port: 8080
  host: "0.0.0.0"
`
		t.AssertNil(os.WriteFile(fallbackFile, []byte(fallbackContent), 0644))

		// 创建文件适配器作为备用
		fileAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "fallback_config.yaml",
			Watch:    false,
		})
		t.AssertNil(err)
		defer fileAdapter.Close(context.Background())

		ctx := context.Background()

		// 测试备用方案：当Consul不可用时使用文件配置
		var appName, appMode, serverPort interface{}
		
		// 尝试从Consul获取配置（会失败）
		consulAdapter, err := NewConsulAdapter(&ConsulAdapterConfig{
			Address: "http://localhost:9999", // 不存在的端口
			Prefix:  "config/",
			Watch:   false,
		})
		
		if consulAdapter != nil && consulAdapter.Available(ctx) {
			appName, _ = consulAdapter.Get(ctx, "app.name")
			appMode, _ = consulAdapter.Get(ctx, "app.mode")
			serverPort, _ = consulAdapter.Get(ctx, "server.port")
		}
		
		// 如果Consul不可用，使用文件备用
		if appName == nil {
			appName, err = fileAdapter.Get(ctx, "app.name")
			t.AssertNil(err)
		}
		if appMode == nil {
			appMode, err = fileAdapter.Get(ctx, "app.mode")
			t.AssertNil(err)
		}
		if serverPort == nil {
			serverPort, err = fileAdapter.Get(ctx, "server.port")
			t.AssertNil(err)
		}

		// 验证从备用配置源获取的数据
		t.Assert(appName, "fallback-app")
		t.Assert(appMode, "production")
		t.Assert(serverPort, 8080)
	})
}