package base

import (
	"context"

	"github.com/vera-byte/vgo/app/admin/internal/dao"
	"github.com/vera-byte/vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterBaseSysLogLogic(NewBaseSysLogLogic())
}

type sBaseSysLogLogic struct {
}

func NewBaseSysLogLogic() *sBaseSysLogLogic {
	return &sBaseSysLogLogic{}
}

// 记录日志
func (c *sBaseSysLogLogic) RecordLog(ctx context.Context, userId *int64, action string, ip string, params string, tenantId *int64, traceId string) (err error) {

	_, err = dao.BaseSysLog.Ctx(ctx).Data(g.Map{
		dao.BaseSysLog.Columns().UserId:   userId,
		dao.BaseSysLog.Columns().Action:   action,
		dao.BaseSysLog.Columns().Ip:       ip,
		dao.BaseSysLog.Columns().Params:   params,
		dao.BaseSysLog.Columns().TenantId: tenantId,
		dao.BaseSysLog.Columns().TraceId:  traceId,
	}).Insert()

	if err != nil {
		g.Log().Error(ctx, "记录日志失败", err)
	}

	return nil
}
