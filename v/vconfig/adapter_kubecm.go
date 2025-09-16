package vconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
)

// KubecmAdapterConfig Kubecm适配器配置
type KubecmAdapterConfig struct {
	// Kubernetes API配置
	APIServer    string        // API服务器地址
	Token        string        // 访问令牌
	TokenFile    string        // 令牌文件路径，默认 /var/run/secrets/kubernetes.io/serviceaccount/token
	CertFile     string        // 证书文件路径
	Namespace    string        // 命名空间，默认 default
	ConfigMap    string        // ConfigMap名称
	
	// 监听配置
	Watch        bool          // 是否监听配置变化
	Interval     time.Duration // 监听间隔，默认 30s
	Timeout      time.Duration // 请求超时，默认 10s
	
	// 数据格式
	DataKey      string        // ConfigMap中的数据键，默认 config.yaml
	Format       string        // 数据格式，支持 yaml, json, properties
}

// KubecmConfigMap Kubernetes ConfigMap结构
type KubecmConfigMap struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   KubecmMetadata    `json:"metadata"`
	Data       map[string]string `json:"data"`
}

// KubecmMetadata ConfigMap元数据
type KubecmMetadata struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	ResourceVersion string            `json:"resourceVersion"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
}

// KubecmAdapter Kubecm配置适配器
// 基于Kubernetes API实现ConfigMap配置管理
type KubecmAdapter struct {
	config         *KubecmAdapterConfig
	client         *http.Client
	watchers       map[string][]func(*ConfigEvent)
	mutex          sync.RWMutex
	lastData       map[string]interface{}
	lastVersion    string
	stopCh         chan struct{}
}

// NewKubecmAdapter 创建Kubecm配置适配器
// config: 适配器配置
func NewKubecmAdapter(config *KubecmAdapterConfig) (*KubecmAdapter, error) {
	if config == nil {
		return nil, gerror.New("config cannot be nil")
	}
	
	// 设置默认值
	if config.Namespace == "" {
		config.Namespace = "default"
	}
	if config.ConfigMap == "" {
		return nil, gerror.New("configmap name is required")
	}
	if config.DataKey == "" {
		config.DataKey = "config.yaml"
	}
	if config.Format == "" {
		config.Format = "yaml"
	}
	if config.Interval == 0 {
		config.Interval = time.Second * 30
	}
	if config.Timeout == 0 {
		config.Timeout = time.Second * 10
	}
	
	// 自动检测集群内配置
	if config.APIServer == "" {
		config.APIServer = "https://kubernetes.default.svc"
	}
	if config.TokenFile == "" {
		config.TokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}
	
	adapter := &KubecmAdapter{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		watchers:    make(map[string][]func(*ConfigEvent)),
		lastData:    make(map[string]interface{}),
		stopCh:      make(chan struct{}),
	}
	
	// 启动配置监听
	if config.Watch {
		go adapter.startWatch()
	}
	
	return adapter, nil
}

// Name 返回适配器名称
func (k *KubecmAdapter) Name() string {
	return "kubecm"
}

// Available 检查配置源是否可用
// ctx: 上下文
// files: 可选的文件名参数（兼容GoFrame接口）
func (k *KubecmAdapter) Available(ctx context.Context, files ...string) bool {
	// 检查是否能访问Kubernetes API
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/configmaps/%s", 
		k.config.APIServer, k.config.Namespace, k.config.ConfigMap)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}
	
	// 设置认证头
	if err := k.setAuthHeader(req); err != nil {
		return false
	}
	
	resp, err := k.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound
}

// Get 获取指定键的配置值
// ctx: 上下文
// pattern: 配置键模式
func (k *KubecmAdapter) Get(ctx context.Context, pattern string) (*gvar.Var, error) {
	data, err := k.Data(ctx)
	if err != nil {
		return nil, err
	}
	
	// 从数据中查找匹配的键
	value := k.getValueByPattern(data, pattern)
	return g.NewVar(value), nil
}

// Data 获取所有配置数据
// ctx: 上下文
func (k *KubecmAdapter) Data(ctx context.Context) (map[string]interface{}, error) {
	// 获取ConfigMap
	configMap, err := k.getConfigMap(ctx)
	if err != nil {
		return nil, err
	}
	
	if configMap == nil || configMap.Data == nil {
		return make(map[string]interface{}), nil
	}
	
	// 获取配置数据
	configData, exists := configMap.Data[k.config.DataKey]
	if !exists {
		return make(map[string]interface{}), nil
	}
	
	// 解析配置数据
	data, err := k.parseConfigData(configData)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to parse config data")
	}
	
	return data, nil
}

// Set 设置配置值
// ctx: 上下文
// pattern: 配置键模式
// value: 配置值
func (k *KubecmAdapter) Set(ctx context.Context, pattern string, value interface{}) error {
	// 获取当前ConfigMap
	configMap, err := k.getConfigMap(ctx)
	if err != nil {
		return err
	}
	
	if configMap == nil {
		return gerror.New("configmap not found")
	}
	
	// 获取当前配置数据
	currentData := make(map[string]interface{})
	if configData, exists := configMap.Data[k.config.DataKey]; exists {
		currentData, err = k.parseConfigData(configData)
		if err != nil {
			return gerror.Wrapf(err, "failed to parse current config data")
		}
	}
	
	// 设置新值
	k.setValueByPattern(currentData, pattern, value)
	
	// 序列化配置数据
	newConfigData, err := k.serializeConfigData(currentData)
	if err != nil {
		return gerror.Wrapf(err, "failed to serialize config data")
	}
	
	// 更新ConfigMap
	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}
	configMap.Data[k.config.DataKey] = newConfigData
	
	// 提交更新
	return k.updateConfigMap(ctx, configMap)
}

// Watch 监听配置变化
// ctx: 上下文
// pattern: 配置键模式
// callback: 回调函数
func (k *KubecmAdapter) Watch(ctx context.Context, pattern string, callback func(*ConfigEvent)) error {
	if !k.config.Watch {
		return gerror.New("kubecm watching is disabled")
	}
	
	k.mutex.Lock()
	defer k.mutex.Unlock()
	
	k.watchers[pattern] = append(k.watchers[pattern], callback)
	return nil
}

// Close 关闭适配器
// ctx: 上下文
func (k *KubecmAdapter) Close(ctx context.Context) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	
	// 停止监听
	close(k.stopCh)
	
	// 清理资源
	k.watchers = make(map[string][]func(*ConfigEvent))
	k.lastData = make(map[string]interface{})
	k.lastVersion = ""
	
	return nil
}

// getConfigMap 获取ConfigMap
func (k *KubecmAdapter) getConfigMap(ctx context.Context) (*KubecmConfigMap, error) {
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/configmaps/%s", 
		k.config.APIServer, k.config.Namespace, k.config.ConfigMap)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to create request")
	}
	
	// 设置认证头
	if err := k.setAuthHeader(req); err != nil {
		return nil, err
	}
	
	resp, err := k.client.Do(req)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to get configmap")
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, gerror.Newf("kubernetes api returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to read response body")
	}
	
	var configMap KubecmConfigMap
	if err := json.Unmarshal(body, &configMap); err != nil {
		return nil, gerror.Wrapf(err, "failed to unmarshal configmap")
	}
	
	return &configMap, nil
}

// updateConfigMap 更新ConfigMap
func (k *KubecmAdapter) updateConfigMap(ctx context.Context, configMap *KubecmConfigMap) error {
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/configmaps/%s", 
		k.config.APIServer, k.config.Namespace, k.config.ConfigMap)
	
	data, err := json.Marshal(configMap)
	if err != nil {
		return gerror.Wrapf(err, "failed to marshal configmap")
	}
	
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(string(data)))
	if err != nil {
		return gerror.Wrapf(err, "failed to create request")
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	// 设置认证头
	if err := k.setAuthHeader(req); err != nil {
		return err
	}
	
	resp, err := k.client.Do(req)
	if err != nil {
		return gerror.Wrapf(err, "failed to update configmap")
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return gerror.Newf("kubernetes api returned status %d", resp.StatusCode)
	}
	
	return nil
}

// setAuthHeader 设置认证头
func (k *KubecmAdapter) setAuthHeader(req *http.Request) error {
	if k.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+k.config.Token)
		return nil
	}
	
	if k.config.TokenFile != "" && gfile.Exists(k.config.TokenFile) {
		token := gfile.GetContents(k.config.TokenFile)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
			return nil
		}
	}
	
	return gerror.New("no valid authentication token found")
}

// parseValue 解析字符串值为合适的类型
func (k *KubecmAdapter) parseValue(value string) interface{} {
	// 去除引号
	value = strings.Trim(value, `"'`)
	
	// 尝试解析为布尔值
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	
	// 尝试解析为整数
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	
	// 尝试解析为浮点数
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	
	// 返回原始字符串
	return value
}

// parseConfigData 解析配置数据
func (k *KubecmAdapter) parseConfigData(data string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	switch strings.ToLower(k.config.Format) {
	case "json":
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return nil, gerror.Wrapf(err, "failed to parse json data")
		}
	case "yaml", "yml":
		// 简单的YAML解析，支持基本的键值对
		lines := strings.Split(data, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				result[key] = k.parseValue(value)
			}
		}
	case "properties":
		// 解析properties格式
		lines := strings.Split(data, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				result[key] = k.parseValue(value)
			}
		}
	default:
		return nil, gerror.Newf("unsupported format: %s", k.config.Format)
	}
	
	return result, nil
}

// serializeConfigData 序列化配置数据
func (k *KubecmAdapter) serializeConfigData(data map[string]interface{}) (string, error) {
	switch strings.ToLower(k.config.Format) {
	case "json":
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", gerror.Wrapf(err, "failed to marshal json data")
		}
		return string(bytes), nil
	case "yaml", "yml":
		// 简单的YAML序列化
		var lines []string
		for key, value := range data {
			lines = append(lines, fmt.Sprintf("%s: %v", key, value))
		}
		return strings.Join(lines, "\n"), nil
	case "properties":
		// 序列化为properties格式
		var lines []string
		for key, value := range data {
			lines = append(lines, fmt.Sprintf("%s=%v", key, value))
		}
		return strings.Join(lines, "\n"), nil
	default:
		return "", gerror.Newf("unsupported format: %s", k.config.Format)
	}
}

// getValueByPattern 根据模式获取值
func (k *KubecmAdapter) getValueByPattern(data map[string]interface{}, pattern string) interface{} {
	// 支持点号分隔的嵌套键
	parts := strings.Split(pattern, ".")
	current := data
	
	for i, part := range parts {
		if i == len(parts)-1 {
			// 最后一个部分，返回值
			return current[part]
		} else {
			// 中间部分，继续嵌套
			if nested, ok := current[part].(map[string]interface{}); ok {
				current = nested
			} else {
				return nil
			}
		}
	}
	
	return nil
}

// setValueByPattern 根据模式设置值
func (k *KubecmAdapter) setValueByPattern(data map[string]interface{}, pattern string, value interface{}) {
	parts := strings.Split(pattern, ".")
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
func (k *KubecmAdapter) startWatch() {
	ticker := time.NewTicker(k.config.Interval)
	defer ticker.Stop()
	
	ctx := context.Background()
	
	for {
		select {
		case <-k.stopCh:
			return
		case <-ticker.C:
			k.checkConfigChanges(ctx)
		}
	}
}

// checkConfigChanges 检查配置变化
func (k *KubecmAdapter) checkConfigChanges(ctx context.Context) {
	// 获取ConfigMap
	configMap, err := k.getConfigMap(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "failed to get configmap: %v", err)
		return
	}
	
	if configMap == nil {
		return
	}
	
	// 检查资源版本是否变化
	if k.lastVersion == configMap.Metadata.ResourceVersion {
		return
	}
	
	// 获取新的配置数据
	newData, err := k.Data(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "failed to load config data: %v", err)
		return
	}
	
	k.mutex.Lock()
	oldData := k.lastData
	k.lastData = make(map[string]interface{})
	for key, value := range newData {
		k.lastData[key] = value
	}
	k.lastVersion = configMap.Metadata.ResourceVersion
	k.mutex.Unlock()
	
	// 比较配置变化并通知监听器
	k.compareAndNotify(oldData, newData)
}

// compareAndNotify 比较配置变化并通知监听器
func (k *KubecmAdapter) compareAndNotify(oldData, newData map[string]interface{}) {
	k.mutex.RLock()
	watchers := make(map[string][]func(*ConfigEvent))
	for key, callbacks := range k.watchers {
		watchers[key] = make([]func(*ConfigEvent), len(callbacks))
		copy(watchers[key], callbacks)
	}
	k.mutex.RUnlock()
	
	// 扁平化数据进行比较
	oldFlat := k.flattenData(oldData, "")
	newFlat := k.flattenData(newData, "")
	
	// 检查所有配置键的变化
	allKeys := make(map[string]bool)
	for key := range oldFlat {
		allKeys[key] = true
	}
	for key := range newFlat {
		allKeys[key] = true
	}
	
	for key := range allKeys {
		oldValue, oldExists := oldFlat[key]
		newValue, newExists := newFlat[key]
		
		var eventType EventType
		if !oldExists && newExists {
			eventType = EventTypeAdd
		} else if oldExists && !newExists {
			eventType = EventTypeDelete
		} else if oldExists && newExists && !k.isEqual(oldValue, newValue) {
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
			if k.matchPattern(key, pattern) {
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
func (k *KubecmAdapter) flattenData(data map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		if nested, ok := value.(map[string]interface{}); ok {
			// 递归处理嵌套数据
			for nestedKey, nestedValue := range k.flattenData(nested, fullKey) {
				result[nestedKey] = nestedValue
			}
		} else {
			result[fullKey] = value
		}
	}
	
	return result
}

// matchPattern 检查键是否匹配模式
func (k *KubecmAdapter) matchPattern(key, pattern string) bool {
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
func (k *KubecmAdapter) isEqual(a, b interface{}) bool {
	return gconv.String(a) == gconv.String(b)
}