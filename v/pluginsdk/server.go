package pluginsdk

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

// rpcMessage JSON-RPC 风格消息（单行 JSON，双向）
type rpcMessage struct {
	ID     int64           `json:"id,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// hostBridge 反向 RPC：cache / pluginService 调宿主
type hostBridge struct {
	mu      sync.Mutex
	writer  *bufio.Writer
	pending map[int64]chan *rpcMessage
	nextID  atomic.Int64
	plugin  any
	info    *PluginInfo
}

func (h *hostBridge) call(method string, params interface{}) (json.RawMessage, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	id := h.nextID.Add(1)
	ch := make(chan *rpcMessage, 1)
	h.mu.Lock()
	h.pending[id] = ch
	err = h.write(&rpcMessage{ID: id, Method: method, Params: raw})
	h.mu.Unlock()
	if err != nil {
		return nil, err
	}
	msg := <-ch
	if msg.Error != nil {
		return nil, errors.New(msg.Error.Message)
	}
	return msg.Result, nil
}

func (h *hostBridge) write(msg *rpcMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err = h.writer.Write(append(b, '\n')); err != nil {
		return err
	}
	return h.writer.Flush()
}

// ---- Cache / PluginService 反向 RPC 实现 ----

type bridgeCache struct{ h *hostBridge }

func (c bridgeCache) Get(key string) (interface{}, error) {
	raw, err := c.h.call("cacheGet", map[string]interface{}{"key": key})
	if err != nil {
		return nil, err
	}
	var out interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out, nil
}

func (c bridgeCache) Set(key string, value interface{}) error {
	_, err := c.h.call("cacheSet", map[string]interface{}{"key": key, "value": value})
	return err
}

type bridgePluginService struct{ h *hostBridge }

func (s bridgePluginService) Invoke(key, method string, params ...interface{}) (interface{}, error) {
	raw, err := s.h.call("pluginInvoke", map[string]interface{}{
		"key": key, "method": method, "params": params,
	})
	if err != nil {
		return nil, err
	}
	var out interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out, nil
}

// Run 启动插件服务（在插件 main 中调用）：从环境读取插件信息，
// 注入 pluginInfo/cache/pluginService，进入 stdio JSON-RPC 服务循环。
func Run(plugin any) {
	info := loadPluginInfo()
	bridge := &hostBridge{
		writer:  bufio.NewWriter(os.Stdout),
		pending: map[int64]chan *rpcMessage{},
		info:    info,
		plugin:  plugin,
	}
	injectBase(plugin, info, bridgeCache{bridge}, bridgePluginService{bridge})

	go readLoop(bridge)

	// 阻塞主进程
	select {}
}

// loadPluginInfo 从环境变量 VGO_PLUGIN_INFO 读取插件信息并解析 @baseDir
func loadPluginInfo() *PluginInfo {
	info := &PluginInfo{}
	if raw := os.Getenv("VGO_PLUGIN_INFO"); raw != "" {
		_ = json.Unmarshal([]byte(raw), info)
	}
	if info.BaseDir == "" {
		info.BaseDir = os.Getenv("VGO_PLUGIN_BASE_DIR")
	}
	// @baseDir/xxx → 绝对路径
	for k, v := range info.Config {
		if s, ok := v.(string); ok && strings.HasPrefix(s, "@baseDir/") {
			info.Config[k] = filepath.Join(info.BaseDir, strings.TrimPrefix(s, "@baseDir/"))
		}
	}
	return info
}

// injectBase 注入 Base 字段（插件结构体嵌入 pluginsdk.Base）
func injectBase(plugin any, info *PluginInfo, cache Cache, svc PluginService) {
	v := reflect.ValueOf(plugin)
	if v.Kind() != reflect.Ptr {
		return
	}
	f := v.Elem().FieldByName("Base")
	if f.IsValid() && f.CanSet() {
		if b, ok := f.Interface().(Base); ok {
			b.PluginInfo = info
			b.Cache = cache
			b.PluginService = svc
			f.Set(reflect.ValueOf(b))
			return
		}
	}
	if bp := v.Elem().FieldByName("Base"); bp.IsValid() && bp.Kind() == reflect.Ptr {
		bp.Set(reflect.ValueOf(&Base{PluginInfo: info, Cache: cache, PluginService: svc}))
	}
}

// readLoop 读取宿主消息：请求派发 / 响应匹配
func readLoop(h *hostBridge) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var msg rpcMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}
		if msg.Method == "" {
			// 响应
			h.mu.Lock()
			ch := h.pending[msg.ID]
			delete(h.pending, msg.ID)
			h.mu.Unlock()
			if ch != nil {
				ch <- &msg
			}
			continue
		}
		go dispatch(h, &msg)
	}
	os.Exit(0)
}

// dispatch 派发宿主请求（ready / invoke）
func dispatch(h *hostBridge, msg *rpcMessage) {
	var (
		result interface{}
		err    error
	)
	switch msg.Method {
	case "ready":
		if r, ok := h.plugin.(readyer); ok {
			r.Ready()
		}
	case "invoke":
		var params struct {
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}
		if e := json.Unmarshal(msg.Params, &params); e != nil {
			err = e
			break
		}
		result, err = invokeMethod(h.plugin, params.Method, params.Params)
	default:
		err = fmt.Errorf("未知方法: %s", msg.Method)
	}

	resp := &rpcMessage{ID: msg.ID}
	if err != nil {
		resp.Error = &rpcError{Code: -1, Message: err.Error()}
	} else {
		raw, e := json.Marshal(result)
		if e != nil {
			resp.Error = &rpcError{Code: -1, Message: e.Error()}
		} else {
			resp.Result = raw
		}
	}
	h.mu.Lock()
	_ = h.write(resp)
	h.mu.Unlock()
}

// invokeMethod 反射调用插件方法（方法名大小写不敏感；参数按 JSON 解码）
func invokeMethod(plugin any, method string, params []interface{}) (interface{}, error) {
	v := reflect.ValueOf(plugin)
	m := v.MethodByName(method)
	if !m.IsValid() {
		// 大小写不敏感匹配（JS 契约方法名小写开头，Go 导出方法大写开头）
		t := v.Type()
		for i := 0; i < t.NumMethod(); i++ {
			if strings.EqualFold(t.Method(i).Name, method) {
				m = v.MethodByName(t.Method(i).Name)
				break
			}
		}
	}
	if !m.IsValid() {
		return nil, fmt.Errorf("插件方法不存在: %s", method)
	}
	mt := m.Type()
	args := make([]reflect.Value, 0, len(params))
	for i := 0; i < mt.NumIn(); i++ {
		var raw interface{}
		if i < len(params) {
			raw = params[i]
		}
		b, _ := json.Marshal(raw)
		arg := reflect.New(mt.In(i))
		_ = json.Unmarshal(b, arg.Interface())
		args = append(args, arg.Elem())
	}
	out := m.Call(args)

	// 返回值归一：支持 (T, error) / (error) / (T) / ()
	if len(out) == 0 {
		return nil, nil
	}
	last := out[len(out)-1]
	if last.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !last.IsNil() {
			return nil, last.Interface().(error)
		}
		if len(out) == 1 {
			return nil, nil
		}
		return out[0].Interface(), nil
	}
	return out[0].Interface(), nil
}
