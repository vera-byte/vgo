package vck

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	vck_config "github.com/vera-byte/vgo/vgo_core_kit/config"
)

var (
	CacheManager          = gcache.New()                     // 定义全局缓存对象	供其他业务使用
	ProcessFlag           = guid.S()                         // 定义全局进程标识
	GetAdminConfig        = vck_config.NewAdminConfig()      // 配置中的v节相关配置
	GetAdminAtGateway     = vck_config.GetAdminAtGateway     // 在网关中获取传入ctx 中的 admin 对象
	GetAdminAtGrpcService = vck_config.GetAdminAtGrpcService // 在微服务中获取metadata中的 admin 对象
	GetCfgWithDefault     = vck_config.GetCfgWithDefault     // GetCfgWithDefault 获取配置，如果配置不存在，则使用默认值
	EtcdManager           = vck_config.NewChainableEtcdClient()
)

func init() {
	var (
		ctx         = gctx.GetInitCtx()
		redisConfig = &gredis.Config{}
	)
	g.Log().Debug(ctx, "vck init start ...")

	redisVar, err := g.Cfg().Get(ctx, "redis")
	if err != nil {
		g.Log().Warning(ctx, "初始化缓存失败,请检查配置文件")
		// panic(err)
	}
	if !redisVar.IsEmpty() {
		redisVar.Struct(redisConfig)
		redis, err := gredis.New(redisConfig)
		if err != nil {
			panic(err)
		}
		if redisConfig.Address == "" {
			panic(gerror.New("redis 配置错误"))
		}
		g.Log().Info(ctx, "初始化缓存成功")
		CacheManager.SetAdapter(gcache.NewAdapterRedis(redis))

	}
	go vck_config.WatchVckConfig()
	g.Log().Debug(ctx, "当前实例ID:", ProcessFlag)

}
