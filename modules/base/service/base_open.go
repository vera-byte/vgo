package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

type BaseOpenService struct {
	*v.Service
}

func NewBaseOpenService() *BaseOpenService {
	return &BaseOpenService{
		&v.Service{},
	}
}

// AdminEPS 获取eps
func (s *BaseOpenService) AdminEPS(ctx g.Ctx) (result *g.Var, err error) {
	c := v.CacheEPS
	result, err = c.GetOrSetFunc(ctx, "adminEPS", func(ctx g.Ctx) (interface{}, error) {
		return s.creatAdminEPS(ctx)
	}, 0)

	return
}

// creatAdminEPS 创建eps
func (s *BaseOpenService) creatAdminEPS(ctx g.Ctx) (adminEPS interface{}, err error) {
	var (
		baseEpsAdmin = model.NewBaseEpsAdmin()
	)

	type Api struct {
		Module  string `json:"module"`  // 所属模块名称 例如：base
		Method  string `json:"method"`  // 请求方法 例如：GET
		Path    string `json:"path"`    // 请求路径 例如：/welcome
		Prefix  string `json:"prefix"`  // 路由前缀 例如：/admin/base/open
		Summary string `json:"summary"` // 描述 例如：欢迎页面
		Tag     string `json:"tag"`     // 标签 例如：base  好像暂时不用
		Dts     string `json:"dts"`     // 未知 例如：{} 好像暂时不用
	}
	// type Column struct {
	// }
	type Module struct {
		Api     []*Api          `json:"api"`
		Columns []*v.ColumnInfo `json:"columns"`
		Module  string          `json:"module"`
		Prefix  string          `json:"prefix"`
	}
	admineps := make(map[string][]*Module)
	// 获取所有路由并更新到数据库表 base_eps_admin（事务内删除重建，避免中途失败导致数据缺失）
	// 注意：必须硬删除——软删残留会占用 (module,prefix,path,method) 唯一键导致重建失败
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		if _, err = v.DBM(baseEpsAdmin).TX(tx).Unscoped().Where("1=1").Delete(); err != nil {
			return
		}
		seen := make(map[string]bool)
		routers := g.Server().GetRoutes()
		for _, router := range routers {
			if router.Type == ghttp.HandlerTypeMiddleware || router.Type == ghttp.HandlerTypeHook {
				continue
			}
			if router.Method == "ALL" {
				continue
			}
			routeSplite := gstr.Split(router.Route, "/")
			if len(routeSplite) < 5 {
				continue
			}
			if routeSplite[1] != "admin" {
				continue
			}
			module := routeSplite[2]
			method := router.Method
			// 获取最后一个元素加前缀 / 为 path
			path := "/" + routeSplite[len(routeSplite)-1]
			// 获取前面的元素为prefix
			prefix := gstr.Join(routeSplite[0:len(routeSplite)-1], "/")
			// 获取最后一个元素为summary
			summary := routeSplite[len(routeSplite)-1]
			// 同一路由可能被多次注册（如通配路由），按唯一键去重
			key := module + "|" + prefix + "|" + path + "|" + method
			if seen[key] {
				continue
			}
			seen[key] = true
			if _, err = v.DBM(baseEpsAdmin).TX(tx).Insert(&Api{
				Module:  module,
				Method:  method,
				Path:    path,
				Prefix:  prefix,
				Summary: summary,
				Tag:     "",
				Dts:     "",
			}); err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	// 读取数据库表生成eps
	items, err := v.DBM(baseEpsAdmin).Fields("DISTINCT module,prefix").All()
	if err != nil {
		return
	}
	for _, item := range items {
		module := item["module"].String()
		prefix := item["prefix"].String()
		apis, e := v.DBM(baseEpsAdmin).Where("module=? AND prefix=?", module, prefix).All()
		if e != nil {
			return nil, e
		}
		var apiList []*Api
		for _, api := range apis {
			apiList = append(apiList, &Api{
				Module:  api["module"].String(),
				Method:  api["method"].String(),
				Path:    api["path"].String(),
				Prefix:  api["prefix"].String(),
				Summary: api["summary"].String(),
				Tag:     api["tag"].String(),
				Dts:     api["dts"].String(),
			})
		}
		admineps[module] = append(admineps[module], &Module{
			Api:     apiList,
			Columns: v.ModelInfo[prefix],
			Module:  module,
			Prefix:  prefix,
		})

	}

	adminEPS = gjson.New(admineps)
	return

}

// AdminEPS 获取eps
func (s *BaseOpenService) AppEPS(ctx g.Ctx) (result *g.Var, err error) {
	c := v.CacheEPS
	result, err = c.GetOrSetFunc(ctx, "appEPS", func(ctx g.Ctx) (interface{}, error) {
		return s.creatAppEPS(ctx)
	}, 0)

	return
}

// creatAppEPS 创建app eps
func (s *BaseOpenService) creatAppEPS(ctx g.Ctx) (appEPS interface{}, err error) {
	var (
		baseEpsApp = model.NewBaseEpsApp()
	)

	type Api struct {
		Module  string `json:"module"`  // 所属模块名称 例如：base
		Method  string `json:"method"`  // 请求方法 例如：GET
		Path    string `json:"path"`    // 请求路径 例如：/welcome
		Prefix  string `json:"prefix"`  // 路由前缀 例如：/admin/base/open
		Summary string `json:"summary"` // 描述 例如：欢迎页面
		Tag     string `json:"tag"`     // 标签 例如：base  好像暂时不用
		Dts     string `json:"dts"`     // 未知 例如：{} 好像暂时不用
	}
	// type Column struct {
	// }
	type Module struct {
		Api     []*Api          `json:"api"`
		Columns []*v.ColumnInfo `json:"columns"`
		Module  string          `json:"module"`
		Prefix  string          `json:"prefix"`
	}
	appeps := make(map[string][]*Module)
	// 获取所有路由并更新到数据库表 base_eps_app（事务内硬删除重建，语义同 creatAdminEPS）
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		if _, err = v.DBM(baseEpsApp).TX(tx).Unscoped().Where("1=1").Delete(); err != nil {
			return
		}
		seen := make(map[string]bool)
		routers := g.Server().GetRoutes()
		for _, router := range routers {
			if router.Type == ghttp.HandlerTypeMiddleware || router.Type == ghttp.HandlerTypeHook {
				continue
			}
			if router.Method == "ALL" {
				continue
			}
			routeSplite := gstr.Split(router.Route, "/")
			if len(routeSplite) < 5 {
				continue
			}
			if routeSplite[1] != "app" {
				continue
			}
			module := routeSplite[2]
			method := router.Method
			// 获取最后一个元素加前缀 / 为 path
			path := "/" + routeSplite[len(routeSplite)-1]
			// 获取前面的元素为prefix
			prefix := gstr.Join(routeSplite[0:len(routeSplite)-1], "/")
			// 获取最后一个元素为summary
			summary := routeSplite[len(routeSplite)-1]
			key := module + "|" + prefix + "|" + path + "|" + method
			if seen[key] {
				continue
			}
			seen[key] = true
			if _, err = v.DBM(baseEpsApp).TX(tx).Insert(&Api{
				Module:  module,
				Method:  method,
				Path:    path,
				Prefix:  prefix,
				Summary: summary,
				Tag:     "",
				Dts:     "",
			}); err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	// 读取数据库表生成eps
	items, err := v.DBM(baseEpsApp).Fields("DISTINCT module,prefix").All()
	if err != nil {
		return
	}
	for _, item := range items {
		module := item["module"].String()
		prefix := item["prefix"].String()
		apis, e := v.DBM(baseEpsApp).Where("module=? AND prefix=?", module, prefix).All()
		if e != nil {
			return nil, e
		}
		var apiList []*Api
		for _, api := range apis {
			apiList = append(apiList, &Api{
				Module:  api["module"].String(),
				Method:  api["method"].String(),
				Path:    api["path"].String(),
				Prefix:  api["prefix"].String(),
				Summary: api["summary"].String(),
				Tag:     api["tag"].String(),
				Dts:     api["dts"].String(),
			})
		}
		appeps[module] = append(appeps[module], &Module{
			Api:     apiList,
			Columns: v.ModelInfo[prefix],
			Module:  module,
			Prefix:  prefix,
		})

	}

	appEPS = gjson.New(appeps)
	return
}
