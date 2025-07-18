package v

import (
	"context"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gbuild"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"gorm.io/gorm"
)

var (
	GormDBS      = make(map[string]*gorm.DB) // 定义全局gorm.DB对象集合 仅供内部使用
	CacheEPS     = gcache.New()              // 定义全局缓存对象	供EPS使用
	CacheManager = gcache.New()              // 定义全局缓存对象	供其他业务使用
	ProcessFlag  = guid.S()                  // 定义全局进程标识
	RunMode      = "dev"                     // 定义全局运行模式
	IsRedisMode  = false                     // 定义全局是否为redis模式
	I18n         = gi18n.New()               // 定义全局国际化对象
)

func init() {
	var (
		ctx         = gctx.GetInitCtx()
		redisConfig = &gredis.Config{}
	)
	g.Log().Debug(ctx, "module v init start ...")
	buildData := gbuild.Data()
	if _, ok := buildData["mode"]; ok {
		RunMode = buildData["mode"].(string)
	}
	if RunMode == "vgo-tools" {
		return
	}
	redisVar, err := g.Cfg().Get(ctx, "redis.v")
	if err != nil {
		g.Log().Error(ctx, "初始化缓存失败,请检查配置文件")
		panic(err)
	}
	if !redisVar.IsEmpty() {
		redisVar.Struct(redisConfig)
		redis, err := gredis.New(redisConfig)
		if err != nil {
			panic(err)
		}
		CacheManager.SetAdapter(gcache.NewAdapterRedis(redis))
		IsRedisMode = true
	}
	g.Log().Debug(ctx, "当前运行模式", RunMode)
	g.Log().Debug(ctx, "当前实例ID:", ProcessFlag)
	g.Log().Debug(ctx, "是否缓存模式:", IsRedisMode)
	g.Log().Debug(ctx, "module v init finished ...")

}

// v.OK 正常返回
type BaseRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 返回正常结果
func Ok(data interface{}) *BaseRes {

	return &BaseRes{
		Code:    1000,
		Message: I18n.Translate(context.TODO(), "BaseResMessage"),
		Data:    data,
	}
}

// 失败返回结果
func Fail(message string) *BaseRes {
	return &BaseRes{
		Code:    1001,
		Message: message,
	}
}

// 分布式函数
// func DistributedFunc(ctx g.Ctx, f func(ctx g.Ctx) (interface{}, error)) (interface{}, error) {
// 	if ProcessFlag == ctx.Request.Header.Get("processFlag") {
// 		return f(ctx)
// 	}
// 	return nil, nil
// }
