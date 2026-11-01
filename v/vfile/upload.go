package vfile

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
)

// allowedUploadExts 上传文件扩展名白名单（小写、含点）。
// 排除 html/svg/js 等可被浏览器执行的类型，防止存储型 XSS。
var allowedUploadExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".bmp": true,
	".pdf": true, ".txt": true, ".md": true, ".csv": true,
	".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true,
	".zip": true, ".rar": true, ".7z": true,
	".mp3": true, ".mp4": true, ".wav": true, ".avi": true, ".mov": true,
}

// ValidateUploadExt 校验上传文件的扩展名是否在白名单内，返回规范化的小写扩展名（含点）。
func ValidateUploadExt(fileName string) (string, error) {
	ext := strings.ToLower(gfile.Ext(fileName))
	if !allowedUploadExts[ext] {
		return "", gerror.Newf("不支持的文件类型: %q", ext)
	}
	return ext, nil
}
