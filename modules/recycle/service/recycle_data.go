package service

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	baseModel "github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/modules/recycle/model"
	"github.com/vera-byte/vgo/v"
)

// RecycleDataService 数据回收服务（参照 cool-admin-midway recycle 模块）
type RecycleDataService struct {
	*v.Service
}

// NewRecycleDataService 创建数据回收服务实例
func NewRecycleDataService() *RecycleDataService {
	return &RecycleDataService{
		Service: &v.Service{
			Model: model.NewRecycleData(),
			PageQueryOp: &v.QueryOp{
				Select: `recycle_data.*,"u"."name" AS "userName"`,
				Join: []*v.JoinOp{
					{
						Model:     baseModel.NewBaseSysUser(),
						Alias:     "u",
						Type:      "LeftJoin",
						Condition: `"recycle_data"."userId" = "u"."id"`,
					},
				},
				KeyWordField: []string{"url"},
				AddOrderby:   g.MapStrStr{"createTime": "DESC"},
			},
		},
	}
}

// entityInfo 来源信息（对应参考实现 entityInfo 的 dataSourceName/entity）
type entityInfo struct {
	Group string `json:"group"`
	Table string `json:"table"`
}

// Record 删除前记录数据快照
func (s *RecycleDataService) Record(ctx context.Context, group, table string, rows gdb.Result) (err error) {
	dataJSON, err := json.Marshal(rows)
	if err != nil {
		return gerror.Wrap(err, "序列化回收数据失败")
	}
	entityJSON, err := json.Marshal(entityInfo{Group: group, Table: table})
	if err != nil {
		return gerror.Wrap(err, "序列化来源信息失败")
	}
	params := "{}"
	if r := g.RequestFromCtx(ctx); r != nil {
		if b, e := json.Marshal(r.GetMap()); e == nil {
			params = string(b)
		}
	}
	_, err = v.DBM(model.NewRecycleData()).Data(g.Map{
		"entityInfo": string(entityJSON),
		"userId":     v.GetAdmin(ctx).UserId,
		"data":       string(dataJSON),
		"url":        requestURL(ctx),
		"params":     params,
		"count":      len(rows),
	}).Insert()
	return
}

// Restore 恢复数据：把快照写回原表后删除回收记录
func (s *RecycleDataService) Restore(ctx context.Context, ids []int) (err error) {
	for _, id := range ids {
		rec, err := v.DBM(model.NewRecycleData()).Where("id", id).One()
		if err != nil {
			return err
		}
		if rec.IsEmpty() {
			continue
		}
		var info entityInfo
		if err = json.Unmarshal([]byte(rec["entityInfo"].String()), &info); err != nil {
			return gerror.Wrapf(err, "解析来源信息失败 id=%d", id)
		}
		var rows []map[string]interface{}
		if err = json.Unmarshal([]byte(rec["data"].String()), &rows); err != nil {
			return gerror.Wrapf(err, "解析回收数据失败 id=%d", id)
		}
		if len(rows) > 0 {
			if _, err = g.DB(info.Group).Model(info.Table).Data(rows).Insert(); err != nil {
				return gerror.Wrapf(err, "恢复数据失败 id=%d", id)
			}
		}
		if _, err = v.DBM(model.NewRecycleData()).Unscoped().Where("id", id).Delete(); err != nil {
			return err
		}
	}
	return
}

// Clear 清理回收站：按 recycleKeep 配置（base_sys_conf）保留天数清理，未配置则清空
func (s *RecycleDataService) Clear(ctx context.Context) (err error) {
	keep := s.recycleKeep(ctx)
	if keep <= 0 {
		_, err = v.DBM(model.NewRecycleData()).Unscoped().Where("1=1").Delete()
		return
	}
	before := gtime.Now().AddDate(0, 0, -keep).Format("Y-m-d H:i:s")
	_, err = v.DBM(model.NewRecycleData()).Unscoped().Where(`"createTime" < ?`, before).Delete()
	return
}

// recycleKeep 读取回收站保留天数配置
func (s *RecycleDataService) recycleKeep(ctx context.Context) int {
	record, err := v.DBM(baseModel.NewBaseSysConf()).Where("cKey", "recycleKeep").One()
	if err != nil || record.IsEmpty() {
		return 0
	}
	return gconv.Int(record["cValue"].String())
}

func requestURL(ctx context.Context) string {
	if r := g.RequestFromCtx(ctx); r != nil {
		return r.URL.Path
	}
	return ""
}

// init 注册删除快照钩子 + 每日定时清理回收站（对应参考实现的事件与 Schedule）
func init() {
	svc := NewRecycleDataService()
	v.RegisterDeleteRecorder(svc.Record)

	gcron.AddSingleton(gctx.New(), "0 0 * * *", func(ctx context.Context) {
		g.Log().Info(ctx, "回收站数据定时清理开始")
		if err := svc.Clear(ctx); err != nil {
			g.Log().Error(ctx, "回收站数据清理失败:", err)
			return
		}
		g.Log().Info(ctx, "回收站数据定时清理完成")
	}, "recycle_clear")
}
