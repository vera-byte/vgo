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

// TestKubecmAdapter_Basic 测试Kubecm适配器基本功能
func TestKubecmAdapter_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟的Kubernetes API服务器
		server := createMockKubernetesServer(t)
		defer server.Close()

		// 创建适配器配置
		config := &KubecmAdapterConfig{
			APIServer: server.URL,
			Token:     "test-token",
			Namespace: "default",
			ConfigMap: "test-config",
			DataKey:   "config.yaml",
			Format:    "yaml",
			Timeout:   time.Second * 5,
		}

		// 创建适配器
		adapter, err := NewKubecmAdapter(config)
		t.AssertNil(err)
		t.AssertNE(adapter, nil)

		// 测试适配器名称
		t.Assert(adapter.Name(), "kubecm")

		// 测试可用性检查
		ctx := context.Background()
		available := adapter.Available(ctx)
		t.Assert(available, true)

		// 测试获取配置数据
		data, err := adapter.Data(ctx)
		t.AssertNil(err)
		t.AssertNE(data, nil)
		t.Assert(data["app"].(map[string]interface{})["name"], "test-app")

		// 测试获取特定配置
		value, err := adapter.Get(ctx, "app.name")
		t.AssertNil(err)
		t.Assert(value.String(), "test-app")

		// 关闭适配器
		err = adapter.Close(ctx)
		t.AssertNil(err)
	})
}

// TestKubecmAdapter_WithFileBackup 测试Kubecm适配器与文件备用方案
func TestKubecmAdapter_WithFileBackup(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "kubecm_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")
		configContent := `
app:
  name: backup-app
  version: 1.0.0
database:
  host: localhost
  port: 3306
`
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		t.AssertNil(err)

		// 创建文件适配器作为备用
		fileAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
		})
		t.AssertNil(err)

		// 创建不可用的Kubecm适配器（使用错误的服务器地址）
		kubecmConfig := &KubecmAdapterConfig{
			APIServer: "https://invalid-server:8443",
			Token:     "test-token",
			Namespace: "default",
			ConfigMap: "test-config",
			Timeout:   time.Millisecond * 100, // 短超时以快速失败
		}

		kubecmAdapter, err := NewKubecmAdapter(kubecmConfig)
		t.AssertNil(err)

		ctx := context.Background()

		// 测试Kubecm适配器不可用
		available := kubecmAdapter.Available(ctx)
		t.Assert(available, false)

		// 使用文件适配器作为备用
		if !available {
			// 从文件适配器获取配置
			data, err := fileAdapter.Data(ctx)
			t.AssertNil(err)
			t.AssertNE(data, nil)
			t.Assert(data["app"].(map[string]interface{})["name"], "backup-app")

			value, err := fileAdapter.Get(ctx, "app.name")
			t.AssertNil(err)
			t.Assert(fmt.Sprintf("%v", value), "backup-app")
		}

		// 关闭适配器
		err = kubecmAdapter.Close(ctx)
		t.AssertNil(err)
		err = fileAdapter.Close(ctx)
		t.AssertNil(err)
	})
}

// TestKubecmAdapter_Watch 测试Kubecm适配器监听功能
func TestKubecmAdapter_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟的Kubernetes API服务器
		server := createMockKubernetesServerWithWatch(t)
		defer server.Close()

		// 创建适配器配置
		config := &KubecmAdapterConfig{
			APIServer: server.URL,
			Token:     "test-token",
			Namespace: "default",
			ConfigMap: "test-config",
			DataKey:   "config.yaml",
			Format:    "yaml",
			Watch:     true,
			Interval:  time.Millisecond * 100,
			Timeout:   time.Second * 5,
		}

		// 创建适配器
		adapter, err := NewKubecmAdapter(config)
		t.AssertNil(err)

		ctx := context.Background()
		eventReceived := make(chan *ConfigEvent, 1)

		// 设置监听回调
		err = adapter.Watch(ctx, "app.name", func(event *ConfigEvent) {
			eventReceived <- event
		})
		t.AssertNil(err)

		// 等待配置变化事件
		select {
		case event := <-eventReceived:
			t.AssertNE(event, nil)
			t.Assert(event.Key, "app.name")
			t.AssertNE(event.Value, nil)
		case <-time.After(time.Second * 2):
			t.Error("未收到配置变化事件")
		}

		// 关闭适配器
		err = adapter.Close(ctx)
		t.AssertNil(err)
	})
}

// TestKubecmAdapter_InvalidConfig 测试Kubecm适配器无效配置处理
func TestKubecmAdapter_InvalidConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试nil配置
		adapter, err := NewKubecmAdapter(nil)
		t.AssertNE(err, nil)
		t.Assert(adapter, nil)

		// 测试缺少ConfigMap名称
		config := &KubecmAdapterConfig{
			APIServer: "https://kubernetes.default.svc",
			Token:     "test-token",
			Namespace: "default",
		}
		adapter, err = NewKubecmAdapter(config)
		t.AssertNE(err, nil)
		t.Assert(adapter, nil)

		// 测试有效配置的默认值设置
		config = &KubecmAdapterConfig{
			APIServer: "https://kubernetes.default.svc",
			Token:     "test-token",
			ConfigMap: "test-config",
		}
		adapter, err = NewKubecmAdapter(config)
		t.AssertNil(err)
		t.AssertNE(adapter, nil)
		t.Assert(adapter.config.Namespace, "default")
		t.Assert(adapter.config.DataKey, "config.yaml")
		t.Assert(adapter.config.Format, "yaml")
		t.Assert(adapter.config.Interval, time.Second*30)
		t.Assert(adapter.config.Timeout, time.Second*10)

		err = adapter.Close(context.Background())
		t.AssertNil(err)
	})
}

// TestKubecmAdapter_Set 测试Kubecm适配器设置操作
func TestKubecmAdapter_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟的Kubernetes API服务器
		server := createMockKubernetesServerWithUpdate(t)
		defer server.Close()

		// 创建适配器配置
		config := &KubecmAdapterConfig{
			APIServer: server.URL,
			Token:     "test-token",
			Namespace: "default",
			ConfigMap: "test-config",
			DataKey:   "config.yaml",
			Format:    "yaml",
			Timeout:   time.Second * 5,
		}

		// 创建适配器
		adapter, err := NewKubecmAdapter(config)
		t.AssertNil(err)

		ctx := context.Background()

		// 测试设置配置值
		err = adapter.Set(ctx, "app.version", "2.0.0")
		t.AssertNil(err)

		// 验证设置的值
		value, err := adapter.Get(ctx, "app.version")
		t.AssertNil(err)
		t.Assert(value.String(), "2.0.0")

		// 关闭适配器
		err = adapter.Close(ctx)
		t.AssertNil(err)
	})
}

// TestKubecmAdapter_ConcurrentAccess 测试Kubecm适配器并发访问
func TestKubecmAdapter_ConcurrentAccess(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建模拟的Kubernetes API服务器
		server := createMockKubernetesServer(t)
		defer server.Close()

		// 创建适配器配置
		config := &KubecmAdapterConfig{
			APIServer: server.URL,
			Token:     "test-token",
			Namespace: "default",
			ConfigMap: "test-config",
			DataKey:   "config.yaml",
			Format:    "yaml",
			Timeout:   time.Second * 5,
		}

		// 创建适配器
		adapter, err := NewKubecmAdapter(config)
		t.AssertNil(err)

		ctx := context.Background()
		done := make(chan bool, 10)

		// 并发读取配置
		for i := 0; i < 10; i++ {
			go func(index int) {
				defer func() { done <- true }()
				
				value, err := adapter.Get(ctx, "app.name")
				t.AssertNil(err)
				t.Assert(value.String(), "test-app")
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < 10; i++ {
			<-done
		}

		// 关闭适配器
		err = adapter.Close(ctx)
		t.AssertNil(err)
	})
}

// TestKubecmAdapter_FallbackToFile 测试Kubecm适配器故障回退到文件
func TestKubecmAdapter_FallbackToFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建临时目录和配置文件
		tempDir, err := os.MkdirTemp("", "kubecm_fallback_test")
		t.AssertNil(err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")
		configContent := `
app:
  name: fallback-app
  version: 1.0.0
  environment: production
database:
  host: prod-db
  port: 5432
`
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		t.AssertNil(err)

		// 创建文件适配器
		fileAdapter, err := NewFileAdapter(&FileAdapterConfig{
			Path:     tempDir,
			FileName: "config.yaml",
		})
		t.AssertNil(err)

		ctx := context.Background()

		// 模拟Kubecm不可用的情况，直接使用文件适配器
		// 测试从文件获取配置
		data, err := fileAdapter.Data(ctx)
		t.AssertNil(err)
		t.AssertNE(data, nil)
		t.Assert(data["app"].(map[string]interface{})["name"], "fallback-app")
		t.Assert(data["app"].(map[string]interface{})["environment"], "production")

		// 测试获取特定配置
		value, err := fileAdapter.Get(ctx, "database.host")
		t.AssertNil(err)
		t.Assert(fmt.Sprintf("%v", value), "prod-db")

		// 关闭适配器
		err = fileAdapter.Close(ctx)
		t.AssertNil(err)
	})
}

// createMockKubernetesServer 创建模拟的Kubernetes API服务器
func createMockKubernetesServer(t *gtest.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证认证头
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 模拟ConfigMap响应
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/api/v1/namespaces/default/configmaps/test-config") {
			configMap := KubecmConfigMap{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Metadata: KubecmMetadata{
					Name:            "test-config",
					Namespace:       "default",
					ResourceVersion: "12345",
				},
				Data: map[string]string{
					"config.yaml": `
app:
  name: test-app
  version: 1.0.0
database:
  host: localhost
  port: 3306
`,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(configMap)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

// createMockKubernetesServerWithWatch 创建支持监听的模拟Kubernetes API服务器
func createMockKubernetesServerWithWatch(t *gtest.T) *httptest.Server {
	version := "12345"
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证认证头
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 模拟ConfigMap响应
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/api/v1/namespaces/default/configmaps/test-config") {
			// 模拟配置变化
			if version == "12345" {
				version = "12346"
			}

			configMap := KubecmConfigMap{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Metadata: KubecmMetadata{
					Name:            "test-config",
					Namespace:       "default",
					ResourceVersion: version,
				},
				Data: map[string]string{
					"config.yaml": fmt.Sprintf(`
app:
  name: test-app-updated-%s
  version: 1.0.1
database:
  host: localhost
  port: 3306
`, version),
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(configMap)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

// createMockKubernetesServerWithUpdate 创建支持更新的模拟Kubernetes API服务器
func createMockKubernetesServerWithUpdate(t *gtest.T) *httptest.Server {
	configData := map[string]string{
		"config.yaml": `
app:
  name: test-app
  version: 1.0.0
database:
  host: localhost
  port: 3306
`,
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证认证头
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 处理GET请求
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/api/v1/namespaces/default/configmaps/test-config") {
			configMap := KubecmConfigMap{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Metadata: KubecmMetadata{
					Name:            "test-config",
					Namespace:       "default",
					ResourceVersion: "12345",
				},
				Data: configData,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(configMap)
			return
		}

		// 处理PUT请求（更新ConfigMap）
		if r.Method == "PUT" && strings.Contains(r.URL.Path, "/api/v1/namespaces/default/configmaps/test-config") {
			var updateConfigMap KubecmConfigMap
			err := json.NewDecoder(r.Body).Decode(&updateConfigMap)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// 更新配置数据
			configData = updateConfigMap.Data

			// 返回更新后的ConfigMap
			updateConfigMap.Metadata.ResourceVersion = "12346"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateConfigMap)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}