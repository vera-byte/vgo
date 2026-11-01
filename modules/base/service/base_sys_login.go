package service

import (
	"context"
	"time"

	"github.com/sa-tokens/sa-token-go/stputil"

	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gogf/gf/v2/util/guid"
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

	vcode, err := v.CacheManager.Get(ctx, captchaId)
	if err != nil {
		err = gerror.Wrap(err, "系统错误")
		return
	}
	if vcode.String() != verifyCode {
		err = gerror.New("验证码错误")
		return
	}
	var user *model.BaseSysUser
	v.DBM(baseSysUser).Where("username=?", username).Where("status=?", 1).Scan(&user)
	if user == nil {
		err = gerror.New("账户或密码不正确~")
		return
	}
	ok, legacy := checkPassword(password, user.Password)
	if !ok {
		err = gerror.New("账户或密码不正确~")
		return
	}
	if legacy {
		// 旧 MD5 哈希登录成功后升级为 bcrypt 存储
		if hashed, e := hashPassword(password); e == nil {
			_, _ = v.DBM(baseSysUser).Where("id=?", user.ID).Update(g.Map{"password": hashed})
		} else {
			g.Log().Warning(ctx, "密码哈希升级失败:", e)
		}
	}

	result, err = s.generateTokenByUser(ctx, user)
	if err != nil {
		return
	}

	return
}

// Captcha 图形验证码
func (*BaseSysLoginService) Captcha(req *v1.BaseOpenCaptchaReq) (interface{}, error) {
	type capchaInfo struct {
		CaptchaId string `json:"captchaId"`
		Data      string `json:"data"`
	}
	var (
		ctx g.Ctx
		err error

		result = &capchaInfo{}
	)
	captchaText := grand.Digits(4)
	svg := `<svg width="` + gconv.String(req.Width) + `" height="` + gconv.String(req.Height) + `" xmlns="http://www.w3.org/2000/svg"><text x="75" y="25" text-anchor="middle" font-size="25" fill="#` + gconv.String(req.Color) + `
	">` + captchaText + `</text></svg>`
	svgbase64 := gbase64.EncodeString(svg)

	result.Data = `data:image/svg+xml;base64,` + svgbase64
	result.CaptchaId = guid.S()
	v.CacheManager.Set(ctx, result.CaptchaId, captchaText, 1800*time.Second)
	// 注意：验证码明文严禁进入日志
	g.Log().Debug(ctx, "生成验证码", result.CaptchaId)
	return result, err
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
func (s *BaseSysLoginService) generateTokenByUser(ctx context.Context, user *model.BaseSysUser) (result *TokenResult, err error) {
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

// NewBaseSysLoginService 创建一个新的BaseSysLoginService
func NewBaseSysLoginService() *BaseSysLoginService {
	return &BaseSysLoginService{}
}
