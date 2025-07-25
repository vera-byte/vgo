package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	v1 "github.com/vera-byte/vgo/modules/base/api/v1"
	"github.com/vera-byte/vgo/modules/base/config"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"

	"github.com/gogf/gf/v2/crypto/gmd5"
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

	vcode, _ := v.CacheManager.Get(ctx, captchaId)
	if vcode.String() != verifyCode {
		err = gerror.New("验证码错误")
		return
	}
	md5password, _ := gmd5.Encrypt(password)

	var user *model.BaseSysUser
	v.DBM(baseSysUser).Where("username=?", username).Where("password=?", md5password).Where("status=?", 1).Scan(&user)
	if user == nil {
		err = gerror.New("账户或密码不正确~")
		return
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
	g.Log().Debug(ctx, "验证码", result.CaptchaId, captchaText)
	return result, err
}

// Logout 退出登录
func (*BaseSysLoginService) Logout(ctx context.Context) (err error) {
	userId := v.GetAdmin(ctx).UserId
	v.CacheManager.Remove(ctx, "admin:department:"+gconv.String(userId))
	v.CacheManager.Remove(ctx, "admin:perms:"+gconv.String(userId))
	v.CacheManager.Remove(ctx, "admin:token:"+gconv.String(userId))
	v.CacheManager.Remove(ctx, "admin:token:refresh:"+gconv.String(userId))
	return
}

// RefreshToken 刷新token
func (s *BaseSysLoginService) RefreshToken(ctx context.Context, token string) (result *TokenResult, err error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &v.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Jwt.Secret), nil
	})
	if err != nil {
		return
	}
	claims, ok := tokenClaims.Claims.(*v.Claims)
	if !ok {
		err = gerror.New("tokenClaims.Claims.(*Claims) error")
		return
	}
	if !tokenClaims.Valid {
		err = gerror.New("tokenClaims.Valid error")
		return
	}
	if !claims.IsRefresh {
		err = gerror.New("claims.IsRefresh error")
		return
	}

	if !(claims.UserId > 0) {
		err = gerror.New("claims.UserId error")
		return
	}

	var (
		user        *model.BaseSysUser
		baseSysUser = model.NewBaseSysUser()
	)
	v.DBM(baseSysUser).Where("id=?", claims.UserId).Where("status=?", 1).Scan(&user)
	if user == nil {
		err = gerror.New("用户不存在")
		return
	}

	result, err = s.generateTokenByUser(ctx, user)
	return
}

// generateToken  生成token
func (*BaseSysLoginService) generateToken(ctx context.Context, user *model.BaseSysUser, roleIds []string, exprire uint, isRefresh bool) (token string) {
	err := v.CacheManager.Set(ctx, "admin:passwordVersion:"+gconv.String(user.ID), gconv.String(user.PasswordV), 0)
	if err != nil {
		g.Log().Error(ctx, "生成token失败", err)
	}

	claims := &v.Claims{
		IsRefresh:       isRefresh,
		RoleIds:         roleIds,
		Username:        user.Username,
		UserId:          user.ID,
		PasswordVersion: user.PasswordV,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(exprire) * time.Second)),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenClaims.SignedString([]byte(config.Config.Jwt.Secret))
	if err != nil {
		g.Log().Error(ctx, "生成token失败", err)
	}
	return
}

// 根据用户生成前端需要的Token信息
func (s *BaseSysLoginService) generateTokenByUser(ctx context.Context, user *model.BaseSysUser) (result *TokenResult, err error) {
	var (
		baseSysRoleService       = NewBaseSysRoleService()
		baseSysMenuService       = NewBaseSysMenuService()
		baseSysDepartmentService = NewBaseSysDepartmentService()
	)
	// 获取用户角色
	roleIds := baseSysRoleService.GetByUser(user.ID)
	// 如果没有角色，则报错
	if len(roleIds) == 0 {
		err = gerror.New("该用户未设置任何角色，无法登录~")
		return
	}

	// 生成token
	result = &TokenResult{}
	result.Expire = config.Config.Jwt.Token.Expire
	result.RefreshExpire = config.Config.Jwt.Token.RefreshExpire
	result.Token = s.generateToken(ctx, user, roleIds, result.Expire, false)
	result.RefreshToken = s.generateToken(ctx, user, roleIds, result.RefreshExpire, true)
	// 将用户相关信息保存到缓存
	perms := baseSysMenuService.GetPerms(roleIds)
	departments := baseSysDepartmentService.GetByRoleIds(roleIds, user.Username == "admin")
	v.CacheManager.Set(ctx, "admin:department:"+gconv.String(user.ID), departments, 0)
	v.CacheManager.Set(ctx, "admin:perms:"+gconv.String(user.ID), perms, 0)
	v.CacheManager.Set(ctx, "admin:token:"+gconv.String(user.ID), result.Token, 0)
	v.CacheManager.Set(ctx, "admin:token:refresh:"+gconv.String(user.ID), result.RefreshToken, 0)

	return
}

// NewBaseSysLoginService 创建一个新的BaseSysLoginService
func NewBaseSysLoginService() *BaseSysLoginService {
	return &BaseSysLoginService{}
}
