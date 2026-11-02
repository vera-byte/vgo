package service

import (
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

// maxLogParamsLen 请求参数日志长度上限
const maxLogParamsLen = 2048

// sensitiveLogFieldPattern 敏感字段脱敏（password/secret/token 等值替换为 ***）
var sensitiveLogFieldPattern = regexp.MustCompile(`(?i)("?(?:password|passwd|secret|token|accessToken|refreshToken|appSecret|appSecretKey|privateKey)"?\s*[:=]\s*)("[^"]*"|'[^']*'|[^&"',}\s]+)`)

// requestLogParams 提取请求参数日志：multipart 不落二进制、敏感字段脱敏、UTF-8 清洗、截断长度
func requestLogParams(r *ghttp.Request) string {
	// multipart 上传只记录字段名与文件名，避免二进制内容破坏编码
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		var parts []string
		if mf := r.GetMultipartForm(); mf != nil {
			for field, files := range mf.File {
				for _, f := range files {
					parts = append(parts, "file:"+field+"="+f.Filename)
				}
			}
			for k := range mf.Value {
				parts = append(parts, k+"=<form>")
			}
		}
		return "[multipart] " + strings.Join(parts, "&")
	}
	params := r.GetBodyString()
	params = sensitiveLogFieldPattern.ReplaceAllString(params, "${1}***")
	params = strings.ToValidUTF8(params, "")
	if len(params) > maxLogParamsLen {
		params = params[:maxLogParamsLen] + "...(truncated)"
	}
	return params
}

type BaseSysLogService struct {
	*v.Service
}

func NewBaseSysLogService() *BaseSysLogService {
	return &BaseSysLogService{
		&v.Service{
			Model: model.NewBaseSysLog(),
			PageQueryOp: &v.QueryOp{
				KeyWordField: []string{"name", "params", "ipAddr"},
				Select:       `base_sys_log.*,"u"."name" `,
				Join: []*v.JoinOp{
					{
						Model:     model.NewBaseSysUser(),
						Alias:     "u",
						Type:      "LeftJoin",
						Condition: `"u"."id" = "base_sys_log"."userId"`,
					},
				},
			},
		},
	}
}

// Record 记录日志（请求参数经脱敏/清洗后入库，失败仅告警不影响主流程）
func (s *BaseSysLogService) Record(ctx g.Ctx) {
	var (
		admin = v.GetAdmin(ctx)
		r     = g.RequestFromCtx(ctx)
	)
	if _, err := v.DBM(s.Model).Insert(g.Map{
		"userId": admin.UserId,
		"action": r.Method + ":" + r.URL.Path,
		"ip":     r.GetClientIp(),
		"ipAddr": r.GetClientIp(),
		"params": requestLogParams(r),
	}); err != nil {
		g.Log().Warningf(ctx, "记录操作日志失败: %v", err)
	}
}

// Clear 清除日志
func (s *BaseSysLogService) Clear(isAll bool) (err error) {
	BaseSysConfService := NewBaseSysConfService()
	m := v.DBM(s.Model)
	if isAll {
		_, err = m.Delete("1=1")
	} else {
		keepDays := gconv.Int(BaseSysConfService.GetValue("logKeep"))
		_, err = m.Delete("createTime < ?", gtime.Now().AddDate(0, 0, -keepDays).String())
	}
	return
}
