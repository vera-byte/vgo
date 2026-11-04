package service

import (
	"context"
	"github.com/vera-byte/vgo/v"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sa-tokens/sa-token-go/core/adapter"
	sagf "github.com/sa-tokens/sa-token-go/integrations/gf"
	samemory "github.com/sa-tokens/sa-token-go/storage/memory"
	saredis "github.com/sa-tokens/sa-token-go/storage/redis"
	"github.com/vera-byte/vgo/modules/user/config"
)

// userDevice C端会话设备标识
const userDevice = "app"

// init 注册 C端用户会话初始化（含 Redis Ping 探测，延后到 server 启动执行）
func init() {
	v.RegisterStartupHook("satoken-user", func(ctx context.Context) {
		initUserSatoken()
	})
}

// initUserSatoken 初始化 C端用户会话（sa-token gf 集成，独立 KeyPrefix 与管理端会话隔离）
func initUserSatoken() {
	cfg := sagf.DefaultConfig()
	cfg.TokenName = "Authorization"
	cfg.Timeout = int64(config.Config.Token.Expire)
	// 单会话：新登录挤掉旧登录（IsShare 必须为 false，否则并发登录共享 token）
	cfg.IsConcurrent = false
	cfg.IsShare = false
	cfg.TokenStyle = sagf.TokenStyleRandom64
	cfg.KeyPrefix = "satoken:user:"
	cfg.IsPrintBanner = false
	sagf.SetManager(sagf.NewManager(newUserStorage(), cfg))
}

// newUserStorage 会话存储：Redis 优先（生产），不可用降级内存（开发）
func newUserStorage() adapter.Storage {
	ctx := context.Background()
	redisVar, err := g.Cfg().Get(ctx, "redis.v")
	if err == nil && !redisVar.IsEmpty() {
		cfg := redisVar.Map()
		client := goredis.NewClient(&goredis.Options{
			Addr:     gconv.String(cfg["address"]),
			Password: gconv.String(cfg["pass"]),
			DB:       gconv.Int(cfg["db"]),
		})
		if err = client.Ping(ctx).Err(); err == nil {
			return saredis.NewStorageFromClient(client)
		}
		g.Log().Warningf(ctx, "用户会话 Redis 存储不可用，降级为内存存储: %v", err)
	}
	return samemory.NewStorage()
}

// userSessionTenantKey 会话中的租户键
const userSessionTenantKey = "tenantId"

// setUserSession 登录后写入会话数据（租户等）
func setUserSession(id, tenantID uint) error {
	sess, err := sagf.GetSession(id)
	if err != nil {
		return err
	}
	sess.Set(userSessionTenantKey, tenantID)
	return nil
}

// UserTenantFromSession 从会话读取租户ID（供中间件使用）
func UserTenantFromSession(token string) uint {
	sess, err := sagf.GetSessionByToken(token)
	if err != nil {
		return 0
	}
	val, _ := sess.Get(userSessionTenantKey)
	return gconv.Uint(val)
}
