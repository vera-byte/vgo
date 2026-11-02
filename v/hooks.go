package v

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ErrNoUploadHook 未注册上传 Hook 插件（调用方回退默认上传实现）
var ErrNoUploadHook = gerror.New("未注册上传Hook插件")

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
