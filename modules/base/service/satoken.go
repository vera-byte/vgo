package service

import (
	"context"
	"encoding"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sa-tokens/sa-token-go/core"
	"github.com/sa-tokens/sa-token-go/core/adapter"
	"github.com/sa-tokens/sa-token-go/core/listener"
	samemory "github.com/sa-tokens/sa-token-go/storage/memory"
	saredis "github.com/sa-tokens/sa-token-go/storage/redis"
	"github.com/sa-tokens/sa-token-go/stputil"
	"github.com/vera-byte/vgo/modules/base/config"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

// adminRoleLabel 超级管理员角色标识（数据驱动，替代 UserId==1 硬编码）
const adminRoleLabel = "admin"

// init 初始化 sa-token 会话认证：
// Token 风格为 Random64 服务端会话（支持踢人/封禁/续期），
// 会话存储优先 Redis（生产），不可用时降级为内存（开发）。
func init() {
	// 会话管理器构建含 Redis Ping 探测，延后到 server 启动执行
	v.RegisterStartupHook("satoken-admin", func(ctx context.Context) {
		initAdminSatoken()
	})
}

func initAdminSatoken() {
	mgr := core.NewBuilder().
		Storage(newSatokenStorage()).
		TokenName("Authorization").
		Timeout(int64(config.Config.Jwt.Token.Expire)).
		// 单会话模式：同一账号新登录挤掉旧登录（IsShare 必须为 false，
		// 否则并发登录共享同一 token，顶号语义失效）
		IsConcurrent(false).
		IsShare(false).
		TokenStyle(core.TokenStyleRandom64).
		IsPrintBanner(false).
		Build()
	stputil.SetManager(mgr)

	// 认证事件审计日志（配合登录安全与问题追溯）
	logEvent := func(event listener.Event, action string) {
		mgr.RegisterFunc(event, func(data *listener.EventData) {
			g.Log().Infof(context.Background(), "认证事件[%s] userId:%s device:%s", action, data.LoginID, data.Device)
		})
	}
	logEvent(listener.EventLogin, "登录")
	logEvent(listener.EventLogout, "登出")
	logEvent(listener.EventKickout, "踢出")
	logEvent(listener.EventReplaced, "顶号")
}

func newSatokenStorage() adapter.Storage {
	ctx := context.Background()
	// 注意不能依赖 v.IsRedisMode：该标志在 NewVgo() 运行时才置位，包 init 阶段恒为 false
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
		g.Log().Warningf(ctx, "sa-token Redis 存储不可用，降级为内存存储: %v", err)
	}
	return binaryStorage{samemory.NewStorage()}
}

// binaryStorage 包装存储：将实现了 encoding.BinaryMarshaler 的值序列化为 []byte 再存储，
// 保证与 Redis 存储行为一致（memory 存储直接存结构体时 sa-token 读回转换会失败）。
type binaryStorage struct {
	adapter.Storage
}

func (s binaryStorage) Set(key string, value any, ttl time.Duration) error {
	if b, ok := value.(encoding.BinaryMarshaler); ok {
		data, err := b.MarshalBinary()
		if err != nil {
			return err
		}
		value = data
	}
	return s.Storage.Set(key, value, ttl)
}

// satokenIsAdmin 判断角色集合中是否包含超级管理员角色（按角色 label 识别）
func satokenIsAdmin(roleIds []string) bool {
	if len(roleIds) == 0 {
		return false
	}
	count, err := v.DBM(model.NewBaseSysRole()).
		WhereIn("id", roleIds).
		Where("label", adminRoleLabel).
		Count()
	if err != nil {
		g.Log().Errorf(context.Background(), "查询管理员角色失败: %v", err)
		return false
	}
	return count > 0
}

// satokenSetupUserAuth 在登录/权限刷新时写入 sa-token 的角色与权限。
// 超管权限为 "*"（sa-token 通配，匹配一切）；普通用户为菜单 perms 列表。
func satokenSetupUserAuth(userId uint, roleIds []string) (isAdmin bool, err error) {
	isAdmin = satokenIsAdmin(roleIds)
	labels, err := roleLabels(roleIds)
	if err != nil {
		return
	}
	if err = stputil.SetRoles(userId, labels); err != nil {
		return
	}
	var perms []string
	if isAdmin {
		perms = []string{"*"}
	} else {
		perms, err = NewBaseSysMenuService().GetPerms(roleIds, false)
		if err != nil {
			return
		}
	}
	err = stputil.SetPermissions(userId, perms)
	return
}

// roleLabels 查询角色 label 列表（用于 sa-token 角色校验）
func roleLabels(roleIds []string) (labels []string, err error) {
	res, err := v.DBM(model.NewBaseSysRole()).WhereIn("id", roleIds).Fields("label").Array()
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		labels = append(labels, gconv.String(item))
	}
	return labels, nil
}

// AdminFromToken 从请求 token 构建上下文用户信息（供中间件填充 ctx）
func AdminFromToken(token string) (*v.Admin, error) {
	loginId, err := stputil.GetLoginID(token)
	if err != nil {
		return nil, err
	}
	userId := gconv.Uint(loginId)
	roleIds, err := NewBaseSysRoleService().GetByUser(userId)
	if err != nil {
		return nil, err
	}
	sess, err := stputil.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}
	return &v.Admin{
		UserId:   userId,
		Username: sess.GetString("username"),
		RoleIds:  roleIds,
		IsAdmin:  stputil.HasRole(userId, adminRoleLabel),
	}, nil
}

// satokenDevice Web 端设备标识
const satokenDevice = "web"

// loginSatoken 为用户签发 access/refresh token 并写入会话数据
func loginSatoken(user *model.BaseSysUser, roleIds []string) (*TokenResult, error) {
	pair, err := stputil.LoginWithRefreshToken(user.ID, satokenDevice)
	if err != nil {
		return nil, err
	}
	sess, err := stputil.GetSession(user.ID)
	if err != nil {
		return nil, err
	}
	sess.Set("username", user.Username)
	sess.Set("roleIds", roleIds)
	sess.Set("refreshToken", pair.RefreshToken)
	return &TokenResult{
		Expire:        config.Config.Jwt.Token.Expire,
		Token:         pair.AccessToken,
		RefreshExpire: config.Config.Jwt.Token.RefreshExpire,
		RefreshToken:  pair.RefreshToken,
	}, nil
}

// logoutSatoken 登出：撤销 refresh token 并清除该账号全部会话
// （密码变更等需要强制下线的场景同样调用此函数）
// 注意必须携带登录时的 device，否则默认设备会话不匹配、无法真正登出。
func logoutSatoken(userId uint) {
	if sess, err := stputil.GetSession(userId); err == nil {
		if rt := sess.GetString("refreshToken"); rt != "" {
			_ = stputil.RevokeRefreshToken(rt)
		}
	}
	_ = stputil.Logout(userId, satokenDevice)
}
