package service

import (
	"context"
	"image/color"
	"strings"
	"time"

	"github.com/golang-module/base64Captcha"
	"github.com/golang-module/base64Captcha/driver"
	"github.com/sa-tokens/sa-token-go/stputil"
	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type BaseSysLoginService struct {
	*v.Service
}

type TokenResult struct {
	Expire        uint   `json:"expire"`
	Token         string `json:"token"`
	RefreshExpire uint   `json:"refreshExpire"`
	RefreshToken  string `json:"refreshToken"`
}

// Login 登录
func (s *BaseSysLoginService) Login(ctx context.Context, req *v1.BaseOpenLoginReq) (result *TokenResult, err error) {
	var (
		captchaId   = req.CaptchaId
		verifyCode  = req.VerifyCode
		password    = req.Password
		username    = req.Username
		baseSysUser = model.NewBaseSysUser()
	)

	if !s.VerifyCaptcha(ctx, captchaId, verifyCode) {
		err = gerror.New("验证码错误")
		return
	}
	var user *model.BaseSysUser
	v.DBM(baseSysUser).Where("username=?", username).Scan(&user)
	if user == nil {
		err = gerror.New("账户或密码不正确~")
		return
	}
	ok, legacy := checkPassword(password, user.Password)
	if !ok {
		err = gerror.New("账户或密码不正确~")
		return
	}
	// 账号状态白名单：仅"正常"可登录
	if err = checkLoginStatus(ctx, user.Status); err != nil {
		return
	}
	if legacy {
		// 旧 MD5 哈希登录成功后升级为 bcrypt 存储
		hashed, e := hashPassword(password)
		if e != nil {
			g.Log().Warning(ctx, "密码哈希升级失败:", e)
		} else if _, e = v.DBM(baseSysUser).Where("id", user.ID).Data(g.Map{"password": hashed}).Update(); e != nil {
			g.Log().Warning(ctx, "密码哈希升级落库失败:", e)
		}
	}

	result, err = s.generateTokenByUser(ctx, user)
	if err != nil {
		return
	}

	return
}

// Captcha 图形验证码
func (*BaseSysLoginService) Captcha(ctx context.Context, req *v1.BaseOpenCaptchaReq) (interface{}, error) {
	type capchaInfo struct {
		CaptchaId string `json:"captchaId"`
		Data      string `json:"data"`
	}
	result := &capchaInfo{}

	// 题目文字由库随机生成深色，必须用浅色背景保证对比度
	bgColor := color.RGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}
	base64Driver := driver.NewDriverMath(driver.DriverMath{
		Width:      req.Width,
		Height:     req.Height,
		NoiseCount: 0,
		// ShowLineOptions: driver.OptionShowSineLine | driver.OptionShowSlimeLine,
		BgColor: &bgColor,
	})

	captcha := base64Captcha.NewCaptcha(base64Driver, newCacheStore(ctx))
	// answer 由 cacheStore 持久化；验证码明文严禁进入日志
	id, src, _, err := captcha.Generate()
	if err != nil {
		return nil, gerror.Wrap(err, "获取验证码失败")
	}

	result.CaptchaId = id
	result.Data = src
	return result, nil
}

// VerifyCaptcha 一次性校验验证码：无论成败验证码均失效，防止重放
func (*BaseSysLoginService) VerifyCaptcha(ctx context.Context, captchaId, code string) bool {
	return newCacheStore(ctx).Verify(captchaId, code, true)
}

const captchaTTL = 5 * time.Minute

// cacheStore 基于 v.CacheManager 存储验证码（可切 Redis，多实例共享），
// 同时满足 base64Captcha 的 store.Store 接口
type cacheStore struct {
	ctx context.Context
}

func newCacheStore(ctx context.Context) *cacheStore {
	return &cacheStore{ctx: ctx}
}

func (s *cacheStore) Set(id string, value string) error {
	return v.CacheManager.Set(s.ctx, id, value, captchaTTL)
}

func (s *cacheStore) Get(id string, clear bool) string {
	val, err := v.CacheManager.Get(s.ctx, id)
	if err != nil {
		g.Log().Warning(s.ctx, "验证码缓存读取失败:", err)
		return ""
	}
	if val == nil || val.IsNil() {
		return ""
	}
	if clear {
		v.CacheManager.Remove(s.ctx, id)
	}
	return val.String()
}

func (s *cacheStore) Verify(id, answer string, clear bool) bool {
	stored := s.Get(id, clear)
	if stored == "" {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(stored), strings.TrimSpace(answer))
}

// Logout 退出登录
func (*BaseSysLoginService) Logout(ctx context.Context) (err error) {
	logoutSatoken(v.GetAdmin(ctx).UserId)
	return
}

// RefreshToken 刷新token
// 校验旧 refresh token 后撤销作废，并重新走登录流程签发新的 token 对（refresh token 轮换）
func (s *BaseSysLoginService) RefreshToken(ctx context.Context, token string) (result *TokenResult, err error) {
	pair, err := stputil.RefreshAccessToken(token)
	if err != nil {
		err = gerror.New("刷新Token失败")
		return
	}
	userId := gconv.Uint(pair.LoginID)
	// 用户必须仍然启用
	var user *model.BaseSysUser
	v.DBM(model.NewBaseSysUser()).Where("id", userId).Where("status", 1).Scan(&user)
	if user == nil {
		_ = stputil.RevokeRefreshToken(pair.RefreshToken)
		err = gerror.New("用户不存在")
		return
	}
	// 旧 refresh token 作废，重新签发
	_ = stputil.RevokeRefreshToken(token)
	roleIds, err := NewBaseSysRoleService().GetByUser(user.ID)
	if err != nil {
		return
	}
	if len(roleIds) == 0 {
		err = gerror.New("该用户未设置任何角色，无法登录~")
		return
	}
	result, err = loginSatoken(user, roleIds)
	if err != nil {
		return
	}
	if _, err = satokenSetupUserAuth(user.ID, roleIds); err != nil {
		logoutSatoken(user.ID)
		return nil, err
	}
	return
}

// 根据用户生成前端需要的Token信息
func (s *BaseSysLoginService) generateTokenByUser(_ context.Context, user *model.BaseSysUser) (result *TokenResult, err error) {
	// 获取用户角色
	roleIds, err := NewBaseSysRoleService().GetByUser(user.ID)
	if err != nil {
		return
	}
	// 如果没有角色，则报错
	if len(roleIds) == 0 {
		err = gerror.New("该用户未设置任何角色，无法登录~")
		return
	}
	// 先登录再写入角色/权限：首次登录会重建会话，先写入的权限将丢失
	result, err = loginSatoken(user, roleIds)
	if err != nil {
		return
	}
	if _, err = satokenSetupUserAuth(user.ID, roleIds); err != nil {
		logoutSatoken(user.ID)
		return nil, err
	}
	return
}

// checkLoginStatus 账号状态白名单校验：仅"正常"放行，其余（封禁/未知/缺省）返回拒绝错误
func checkLoginStatus(ctx context.Context, status *int32) error {
	if status == nil || *status != model.UserStatusNormal {
		return gerror.New(statusRefusedMessage(ctx, status))
	}
	return nil
}

// statusRefusedMessage 状态拒绝登录的提示文案：
// 优先取该状态字典条目的备注（value → 显示名 → 条目），缺省返回默认文案
func statusRefusedMessage(ctx context.Context, status *int32) string {
	const defaultMsg = "当前账号状态不支持登录~"
	if status == nil {
		return defaultMsg
	}
	label := v.Dict.GetLabel(ctx, model.DictKeyUserStatus, gconv.String(*status))
	if label == "" {
		return defaultMsg
	}
	if e := v.Dict.GetInfo(ctx, model.DictKeyUserStatus, label); e != nil && e.Remark != "" {
		return e.Remark
	}
	return defaultMsg
}

// NewBaseSysLoginService 创建一个新的BaseSysLoginService
func NewBaseSysLoginService() *BaseSysLoginService {
	return &BaseSysLoginService{}
}
