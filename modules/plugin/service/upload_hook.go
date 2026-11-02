package service

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/vera-byte/vgo/v"
)

// InitUploadHook 注册上传钩子：存在启用中的 upload Hook 插件时接管上传
//（对应 cool-admin 中 Hook 插件替换默认上传实现的行为）
func InitUploadHook() {
	v.RegisterUploadHook(func(ctx context.Context) (interface{}, error) {
		if !HasHook("upload") {
			return nil, v.ErrNoUploadHook
		}
		r := g.RequestFromCtx(ctx)
		file := r.GetUploadFile("file")
		if file == nil {
			return nil, gerror.New("上传文件为空")
		}
		// 写入临时文件供插件读取（插件契约：ctx.files[0].data 为文件路径）
		tmp, err := os.CreateTemp("", "vgo-plugin-upload-*"+filepath.Ext(file.Filename))
		if err != nil {
			return nil, err
		}
		tmpName := tmp.Name()
		_ = tmp.Close()
		defer os.Remove(tmpName)
		src, e := file.Open()
		if e != nil {
			return nil, e
		}
		defer src.Close()
		dst, e := os.Create(tmpName)
		if e != nil {
			return nil, e
		}
		if _, e = io.Copy(dst, src); e != nil {
			dst.Close()
			return nil, e
		}
		dst.Close()
		fields := map[string]interface{}{}
		for k := range r.GetMap() {
			fields[k] = r.Get(k).String()
		}
		ctxData := map[string]interface{}{
			"files": []map[string]interface{}{
				{"data": tmpName, "filename": file.Filename},
			},
			"fields": fields,
		}
		result, err := InvokeHook(ctx, "upload", "upload", ctxData)
		if err != nil {
			if gerror.Is(err, ErrNoHookPlugin) {
				return nil, v.ErrNoUploadHook
			}
			return nil, gerror.Wrap(err, "上传插件执行失败")
		}
		return result, nil
	})
}

func init() {
	InitUploadHook()
}
