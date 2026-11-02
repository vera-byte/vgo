package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/v"
)

// rpcMessage 与 pluginsdk 一致的 JSON-RPC 风格消息（单行 JSON、双向）
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

// pluginProcess 插件可执行进程（stdio JSON-RPC）
type pluginProcess struct {
	key     string
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mu      sync.Mutex
	pending map[int64]chan *rpcMessage
	nextID  atomic.Int64
	info    *PluginInfoData
	done    chan struct{}
}

// startPluginProcess 启动插件可执行程序并完成 ready 握手
func startPluginProcess(key, binPath, baseDir string, info *PluginInfoData) (*pluginProcess, error) {
	// 可执行文件用绝对路径，避免相对路径与 cmd.Dir 组合导致 fork/exec 解析失败
	if abs, err := filepath.Abs(binPath); err == nil {
		binPath = abs
	}
	if abs, err := filepath.Abs(baseDir); err == nil {
		baseDir = abs
	}
	cmd := exec.Command(binPath)
	infoJSON, _ := json.Marshal(info)
	cmd.Env = append(os.Environ(),
		"VGO_PLUGIN_INFO="+string(infoJSON),
		"VGO_PLUGIN_BASE_DIR="+baseDir,
	)
	cmd.Dir = baseDir
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return nil, gerror.Wrapf(err, "插件进程启动失败: %s", key)
	}
	p := &pluginProcess{
		key:     key,
		cmd:     cmd,
		stdin:   stdin,
		pending: map[int64]chan *rpcMessage{},
		info:    info,
		done:    make(chan struct{}),
	}
	go p.readLoop(stdout)
	go func() {
		_ = cmd.Wait()
		close(p.done)
	}()
	// ready 握手（单例插件仅启动时一次，非单例每次调用前均执行）
	if _, err = p.call("ready", map[string]interface{}{}, 10*time.Second); err != nil {
		p.Stop()
		return nil, gerror.Wrapf(err, "插件就绪握手失败: %s", key)
	}
	return p, nil
}

// Stop 终止插件进程
func (p *pluginProcess) Stop() {
	_ = p.stdin.Close()
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
}

// call 发起 RPC 并等待结果
func (p *pluginProcess) call(method string, params interface{}, timeout time.Duration) (json.RawMessage, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	id := p.nextID.Add(1)
	ch := make(chan *rpcMessage, 1)
	p.mu.Lock()
	p.pending[id] = ch
	b, _ := json.Marshal(&rpcMessage{ID: id, Method: method, Params: raw})
	_, err = p.stdin.Write(append(b, '\n'))
	p.mu.Unlock()
	if err != nil {
		return nil, err
	}
	select {
	case msg := <-ch:
		if msg.Error != nil {
			return nil, gerror.New(msg.Error.Message)
		}
		return msg.Result, nil
	case <-time.After(timeout):
		p.mu.Lock()
		delete(p.pending, id)
		p.mu.Unlock()
		return nil, gerror.Newf("插件调用超时: %s.%s", p.key, method)
	case <-p.done:
		return nil, gerror.Newf("插件进程已退出: %s", p.key)
	}
}

// readLoop 读取插件消息：响应匹配 + 反向请求（cache/pluginInvoke）
func (p *pluginProcess) readLoop(r io.Reader) {
	scanner := bufio.NewScanner(r)
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
			p.mu.Lock()
			ch := p.pending[msg.ID]
			delete(p.pending, msg.ID)
			p.mu.Unlock()
			if ch != nil {
				ch <- &msg
			}
			continue
		}
		go p.handleHostCall(&msg)
	}
}

// handleHostCall 处理插件发起的反向调用
func (p *pluginProcess) handleHostCall(msg *rpcMessage) {
	var (
		result interface{}
		err    error
	)
	switch msg.Method {
	case "cacheGet":
		var params struct {
			Key string `json:"key"`
		}
		_ = json.Unmarshal(msg.Params, &params)
		val, e := v.CacheManager.Get(context.Background(), params.Key)
		if e != nil {
			err = e
		} else {
			result = val.Val()
		}
	case "cacheSet":
		var params struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}
		_ = json.Unmarshal(msg.Params, &params)
		err = v.CacheManager.Set(context.Background(), params.Key, params.Value, 0)
	case "pluginInvoke":
		var params struct {
			Key    string        `json:"key"`
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}
		_ = json.Unmarshal(msg.Params, &params)
		result, err = Invoke(context.Background(), params.Key, params.Method, params.Params...)
	case "log":
		var params struct {
			Level   string `json:"level"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(msg.Params, &params)
		g.Log().Infof(context.Background(), "[插件:%s] %s", p.key, params.Message)
	default:
		err = fmt.Errorf("未知宿主方法: %s", msg.Method)
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
	b, _ := json.Marshal(resp)
	p.mu.Lock()
	_, _ = p.stdin.Write(append(b, '\n'))
	p.mu.Unlock()
}

// Invoke 实例方法调用（实现 PluginInstance 接口）
func (p *pluginProcess) Invoke(method string, args ...interface{}) (interface{}, error) {
	raw, err := p.call("invoke", map[string]interface{}{
		"method": method,
		"params": args,
	}, 30*time.Second)
	if err != nil {
		return nil, err
	}
	var out interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out, nil
}
