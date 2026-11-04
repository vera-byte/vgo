package v

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ErrNoUploadHook 未注册上传 Hook 插件（调用方回退默认上传实现）
var ErrNoUploadHook = gerror.New("未注册上传Hook插件")

// startupHook 服务启动钩子：init 期只注册不执行，
// 由 server 命令启动时统一运行，避免 migrate/pack 等子命令在 init 阶段触库。
type startupHook struct {
	name string
	fn   func(ctx context.Context)
}

var startupHooks []startupHook

// RegisterStartupHook 注册服务启动钩子（同名覆盖）。
func RegisterStartupHook(name string, fn func(ctx context.Context)) {
	for i := range startupHooks {
		if startupHooks[i].name == name {
			startupHooks[i].fn = fn
			return
		}
	}
	startupHooks = append(startupHooks, startupHook{name: name, fn: fn})
}

// RunStartupHooks 运行全部启动钩子（server 启动时调用；单个失败不阻断其余）。
func RunStartupHooks(ctx context.Context) {
	for _, h := range startupHooks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					g.Log().Errorf(ctx, "启动钩子[%s]执行失败: %v", h.name, r)
				}
			}()
			h.fn(ctx)
		}()
	}
}

// DeleteRecorderFunc 删除前的数据快照回调。
// 由 modules/recycle 注册：记录被删数据到回收站，支持后续恢复。
type DeleteRecorderFunc func(ctx context.Context, group, table string, rows gdb.Result) error

var deleteRecorder DeleteRecorderFunc

// RegisterDeleteRecorder 注册删除快照回调（全局唯一，通常在回收站模块初始化时调用）
func RegisterDeleteRecorder(f DeleteRecorderFunc) {
	deleteRecorder = f
}

// UploadHookFunc 上传 Hook：由 modules/plugin 注册。
// 存在启用中的 upload Hook 插件时接管文件上传；否则返回 ErrNoUploadHook 由默认实现处理。
type UploadHookFunc func(ctx context.Context) (interface{}, error)

// UploadHook 上传钩子（插件模块注册）
var UploadHook UploadHookFunc

// RegisterUploadHook 注册上传钩子
func RegisterUploadHook(f UploadHookFunc) {
	UploadHook = f
}
