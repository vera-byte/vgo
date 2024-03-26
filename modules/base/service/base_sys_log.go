package service

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseSysLogService struct {
	*v.Service
}

func NewBaseSysLogService() *BaseSysLogService {
	return &BaseSysLogService{
		&v.Service{
			Model: model.NewBaseSysLog(),
			PageQueryOp: &v.QueryOp{
				KeyWordField: []string{"name", "params", "ipAddr"},
				Select:       "base_sys_log.*,user.name ",
				Join: []*v.JoinOp{
					{
						Model:     model.NewBaseSysUser(),
						Alias:     "user",
						Type:      "LeftJoin",
						Condition: "user.id = base_sys_log.userID",
					},
				},
			},
		},
	}
}

// Record 记录日志
func (s *BaseSysLogService) Record(ctx g.Ctx) {
	var (
		admin = v.GetAdmin(ctx)
		r     = g.RequestFromCtx(ctx)
	)
	baseSysLog := model.NewBaseSysLog()
	baseSysLog.UserID = admin.UserId
	baseSysLog.Action = r.Method + ":" + r.URL.Path
	baseSysLog.IP = r.GetClientIp()
	baseSysLog.IPAddr = r.GetClientIp()
	baseSysLog.Params = r.GetBodyString()
	m := v.DBM(s.Model)
	m.Insert(g.Map{
		"userId": baseSysLog.UserID,
		"action": baseSysLog.Action,
		"ip":     baseSysLog.IP,
		"ipAddr": baseSysLog.IPAddr,
		"params": baseSysLog.Params,
	})
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
