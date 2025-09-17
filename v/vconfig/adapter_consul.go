package vconfig

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// ConsulAdapterConfig Consul适配器配置
type ConsulAdapterConfig struct {
	Address    string        // Consul地址，默认 127.0.0.1:8500
	Scheme     string        // 协议，默认 http
	Datacenter string        // 数据中心
	Token      string        // 访问令牌
	Prefix     string        // 键前缀
	Watch      bool          // 是否监听配置变化
	Interval   time.Duration // 监听间隔，默认 30s
	Timeout    time.Duration // 请求超时，默认 10s
}

// ConsulKVPair Consul KV键值对
type ConsulKVPair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// ConsulAdapter Consul配置适配器
// 基于HTTP API实现配置管理，不依赖consul/api包
type ConsulAdapter struct {
	config   *ConsulAdapterConfig
	baseURL  string
	client   *http.Client
	watchers map[string][]func(*ConfigEvent)
	mutex    sync.RWMutex
	lastData map[string]interface{}
	stopCh   chan struct{}
}

// NewConsulAdapter 创建Consul配置适配器
// config: 适配器配置
func NewConsulAdapter(config *ConsulAdapterConfig) (*ConsulAdapter, error) {
	if config == nil {
		return nil, gerror.New("config cannot be nil")
	}
	
	// 设置默认值
	if config.Address == "" {
		config.Address = "127.0.0.1:8500"
	}
	if config.Scheme == "" {
		config.Scheme = "http"
	}
	if config.Prefix == "" {
		config.Prefix = "config/"
	}
	if config.Interval == 0 {
		config.Interval = time.Second * 30
	}
	if config.Timeout == 0 {
		config.Timeout = time.Second * 10
	}
	
	baseURL := fmt.Sprintf("%s://%s", config.Scheme, config.Address)
	
	adapter := &ConsulAdapter{
		config:  config,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		watchers: make(map[string][]func(*ConfigEvent)),
		lastData: make(map[string]interface{}),
		stopCh:   make(chan struct{}),
	}
	
	// 启动配置监听
	if config.Watch {
		go adapter.startWatch()
	}
	
	return adapter, nil
}

// Name 返回适配器名称
func (c *ConsulAdapter) Name() string {
	return "consul"
}

// Available 检查配置源是否可用
// ctx: 上下文
// files: 可选的文件名参数（兼容GoFrame接口）
func (c *ConsulAdapter) Available(ctx context.Context, files ...string) bool {
	// 尝试连接Consul
	url := fmt.Sprintf("%s/v1/kv/health-check", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}
	
	if c.config.Token != "" {
		req.Header.Set("X-Consul-Token", c.config.Token)
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound
}

// Get 获取指定键的配置值
// ctx: 上下文
// pattern: 配置键模式
func (c *ConsulAdapter) Get(ctx context.Context, pattern string) (*gvar.Var, error) {
	key := c.buildKey(pattern)
	
	url := fmt.Sprintf("%s/v1/kv/%s", c.baseURL, url.QueryEscape(key))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to create request for key: %s", key)
	}
	
	if c.config.Token != "" {
		req.Header.Set("X-Consul-Token", c.config.Token)
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to get config from consul: %s", key)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return g.NewVar(nil), nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, gerror.Newf("consul returned status %d for key: %s", resp.StatusCode, key)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to read response body for key: %s", key)
	}
	
	var pairs []ConsulKVPair
	if err := json.Unmarshal(body, &pairs); err != nil {
		return nil, gerror.Wrapf(err, "failed to unmarshal response for key: %s", key)
	}
	
	if len(pairs) == 0 {
		return g.NewVar(nil), nil
	}
	
	// Base64解码值
	valueBytes, err := c.decodeBase64(pairs[0].Value)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to decode value for key: %s", key)
	}
	
	// 尝试解析JSON
	var value interface{}
	if err := json.Unmarshal(valueBytes, &value); err != nil {
		// 如果不是JSON，直接返回字符串
		value = string(valueBytes)
	}
	
	return g.NewVar(value), nil
}

// Data 获取所有配置数据
// ctx: 上下文
func (c *ConsulAdapter) Data(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/kv/%s?recurse=true", c.baseURL, url.QueryEscape(c.config.Prefix))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to create request for prefix: %s", c.config.Prefix)
	}
	
	if c.config.Token != "" {
		req.Header.Set("X-Consul-Token", c.config.Token)
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to list configs from consul")
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return make(map[string]interface{}), nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, gerror.Newf("consul returned status %d for prefix: %s", resp.StatusCode, c.config.Prefix)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to read response body")
	}
	
	var pairs []ConsulKVPair
	if err := json.Unmarshal(body, &pairs); err != nil {
		return nil, gerror.Wrapf(err, "failed to unmarshal response")
	}
	
	result := make(map[string]interface{})
	for _, pair := range pairs {
		key := strings.TrimPrefix(pair.Key, c.config.Prefix)
		if key == "" {
			continue
		}
		
		// Base64解码值
		valueBytes, err := c.decodeBase64(pair.Value)
		if err != nil {
			g.Log().Warningf(ctx, "failed to decode value for key %s: %v", key, err)
			continue
		}
		
		// 尝试解析JSON
		var value interface{}
		if err := json.Unmarshal(valueBytes, &value); err != nil {
			// 如果不是JSON，直接使用字符串
			value = string(valueBytes)
		}
		
		// 构建嵌套结构
		c.setNestedValue(result, key, value)
	}
	
	return result, nil
}

// Set 设置配置值
// ctx: 上下文
// pattern: 配置键模式
// value: 配置值
func (c *ConsulAdapter) Set(ctx context.Context, pattern string, value interface{}) error {
	key := c.buildKey(pattern)
	
	// 序列化值
	var data []byte
	var err error
	
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return gerror.Wrapf(err, "failed to marshal value for key: %s", key)
		}
	}
	
	url := fmt.Sprintf("%s/v1/kv/%s", c.baseURL, url.QueryEscape(key))
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(string(data)))
	if err != nil {
		return gerror.Wrapf(err, "failed to create request for key: %s", key)
	}
	
	if c.config.Token != "" {
		req.Header.Set("X-Consul-Token", c.config.Token)
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return gerror.Wrapf(err, "failed to set config in consul: %s", key)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return gerror.Newf("consul returned status %d for key: %s", resp.StatusCode, key)
	}
	
	return nil
}

// Watch 监听配置变化
// ctx: 上下文
// pattern: 配置键模式
// callback: 回调函数
func (c *ConsulAdapter) Watch(ctx context.Context, pattern string, callback func(*ConfigEvent)) error {
	if !c.config.Watch {
		return gerror.New("consul watching is disabled")
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.watchers[pattern] = append(c.watchers[pattern], callback)
	return nil
}

// Close 关闭适配器
// ctx: 上下文
func (c *ConsulAdapter) Close(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// 停止监听
	close(c.stopCh)
	
	// 清理资源
	c.watchers = make(map[string][]func(*ConfigEvent))
	c.lastData = make(map[string]interface{})
	
	return nil
}

// buildKey 构建完整的Consul键
func (c *ConsulAdapter) buildKey(pattern string) string {
	if strings.HasPrefix(pattern, c.config.Prefix) {
		return pattern
	}
	return c.config.Prefix + pattern
}

// decodeBase64 解码Base64字符串
func (c *ConsulAdapter) decodeBase64(encoded string) ([]byte, error) {
	// Consul API返回的值是Base64编码的
	return base64.StdEncoding.DecodeString(encoded)
}

// setNestedValue 设置嵌套值
func (c *ConsulAdapter) setNestedValue(data map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, "/")
	current := data
	
	for i, part := range parts {
		if i == len(parts)-1 {
			// 最后一个部分，设置值
			current[part] = value
		} else {
			// 中间部分，创建嵌套map
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]interface{})
			}
			if nested, ok := current[part].(map[string]interface{}); ok {
				current = nested
			}
		}
	}
}

// startWatch 启动配置监听
func (c *ConsulAdapter) startWatch() {
	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()
	
	ctx := context.Background()
	
	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.checkConfigChanges(ctx)
		}
	}
}

// checkConfigChanges 检查配置变化
func (c *ConsulAdapter) checkConfigChanges(ctx context.Context) {
	// 获取新的配置数据
	newData, err := c.Data(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "failed to load config data from consul: %v", err)
		return
	}
	
	c.mutex.Lock()
	oldData := c.lastData
	c.lastData = make(map[string]interface{})
	for k, v := range newData {
		c.lastData[k] = v
	}
	c.mutex.Unlock()
	
	// 比较配置变化并通知监听器
	c.compareAndNotify(oldData, newData)
}

// compareAndNotify 比较配置变化并通知监听器
func (c *ConsulAdapter) compareAndNotify(oldData, newData map[string]interface{}) {
	c.mutex.RLock()
	watchers := make(map[string][]func(*ConfigEvent))
	for k, v := range c.watchers {
		watchers[k] = make([]func(*ConfigEvent), len(v))
		copy(watchers[k], v)
	}
	c.mutex.RUnlock()
	
	// 扁平化数据进行比较
	oldFlat := c.flattenData(oldData, "")
	newFlat := c.flattenData(newData, "")
	
	// 检查所有配置键的变化
	allKeys := make(map[string]bool)
	for k := range oldFlat {
		allKeys[k] = true
	}
	for k := range newFlat {
		allKeys[k] = true
	}
	
	for key := range allKeys {
		oldValue, oldExists := oldFlat[key]
		newValue, newExists := newFlat[key]
		
		var eventType EventType
		if !oldExists && newExists {
			eventType = EventTypeAdd
		} else if oldExists && !newExists {
			eventType = EventTypeDelete
		} else if oldExists && newExists && !c.isEqual(oldValue, newValue) {
			eventType = EventTypeUpdate
		} else {
			continue // 没有变化
		}
		
		event := &ConfigEvent{
			Key:      key,
			Value:    newValue,
			OldValue: oldValue,
			Type:     eventType,
		}
		
		// 通知匹配的监听器
		for pattern, callbacks := range watchers {
			if c.matchPattern(key, pattern) {
				for _, callback := range callbacks {
					go func(cb func(*ConfigEvent)) {
						defer func() {
							if r := recover(); r != nil {
								g.Log().Errorf(context.Background(), "config watcher panic: %v", r)
							}
						}()
						cb(event)
					}(callback)
				}
			}
		}
	}
}

// flattenData 扁平化数据
func (c *ConsulAdapter) flattenData(data map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		if nested, ok := value.(map[string]interface{}); ok {
			// 递归处理嵌套数据
			for k, v := range c.flattenData(nested, fullKey) {
				result[k] = v
			}
		} else {
			result[fullKey] = value
		}
	}
	
	return result
}

// matchPattern 检查键是否匹配模式
func (c *ConsulAdapter) matchPattern(key, pattern string) bool {
	// 简单的模式匹配，支持通配符 *
	if pattern == "*" {
		return true
	}
	
	if strings.Contains(pattern, "*") {
		// 简单的通配符匹配
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			prefix, suffix := parts[0], parts[1]
			return strings.HasPrefix(key, prefix) && strings.HasSuffix(key, suffix)
		}
	}
	
	return key == pattern || strings.HasPrefix(key, pattern+".")
}

// isEqual 比较两个值是否相等
func (c *ConsulAdapter) isEqual(a, b interface{}) bool {
	return gconv.String(a) == gconv.String(b)
}