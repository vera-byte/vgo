package service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	pluginService "github.com/vera-byte/vgo/modules/plugin/service"
	sagf "github.com/sa-tokens/sa-token-go/integrations/gf"
	userModel "github.com/vera-byte/vgo/modules/user/model"
	"github.com/vera-byte/vgo/modules/user/config"
	"github.com/vera-byte/vgo/v"
)

// UserTokenResult 登录返回（与管理端 TokenResult 同构，前端一致）
type UserTokenResult struct {
	Expire        uint   `json:"expire"`
	Token         string `json:"token"`
	RefreshExpire uint   `json:"refreshExpire"`
	RefreshToken  string `json:"refreshToken"`
}

// UserLoginService 用户登录服务（参照 cool-admin user 模块 UserLoginService）
type UserLoginService struct{}

// NewUserLoginService 创建用户登录服务
func NewUserLoginService() *UserLoginService {
	return &UserLoginService{}
}

// SmsCode 发送手机验证码（先校验图片验证码）
func (s *UserLoginService) SmsCode(ctx context.Context, phone, captchaId, code string) error {
	if !checkCaptcha(ctx, captchaId, code) {
		return gerror.New("图片验证码错误")
	}
	return NewUserSmsService().SendSms(ctx, phone)
}

// PhoneVerifyCode 手机验证码登录
func (s *UserLoginService) PhoneVerifyCode(ctx context.Context, phone, smsCode string) (*UserTokenResult, error) {
	if !NewUserSmsService().CheckCode(ctx, phone, smsCode) {
		return nil, gerror.New("验证码错误")
	}
	return s.Phone(ctx, phone)
}

// Phone 手机号登录（不存在则注册，昵称脱敏）
func (s *UserLoginService) Phone(ctx context.Context, phone string) (*UserTokenResult, error) {
	user, err := v.DBM(userModel.NewUserInfo()).Where("phone", phone).One()
	if err != nil {
		return nil, err
	}
	if user.IsEmpty() {
		nick := maskPhone(phone)
		unionid := phone
		loginType := 2
		if _, err = v.DBM(userModel.NewUserInfo()).Data(g.Map{
			"phone":     phone,
			"unionid":   unionid,
			"loginType": loginType,
			"nickName":  nick,
		}).Insert(); err != nil {
			return nil, err
		}
		if user, err = v.DBM(userModel.NewUserInfo()).Where("phone", phone).One(); err != nil {
			return nil, err
		}
	}
	if user["status"].Int() == 0 {
		return nil, gerror.New("账号已禁用")
	}
	return s.Token(user["id"].Uint())
}

// Password 密码登录（bcrypt，兼容旧 MD5 并自动升级）
func (s *UserLoginService) Password(ctx context.Context, phone, password string) (*UserTokenResult, error) {
	user, err := v.DBM(userModel.NewUserInfo()).Where("phone", phone).One()
	if err != nil {
		return nil, err
	}
	if user.IsEmpty() {
		return nil, gerror.New("账号或密码错误")
	}
	ok, legacy := checkPassword(password, user["password"].String())
	if !ok {
		return nil, gerror.New("账号或密码错误")
	}
	if user["status"].Int() == 0 {
		return nil, gerror.New("账号已禁用")
	}
	if legacy {
		if hashed, e := hashPassword(password); e == nil {
			if _, e = v.DBM(userModel.NewUserInfo()).Where("id", user["id"]).Data(g.Map{"password": hashed}).Update(); e != nil {
				g.Log().Warning(ctx, "密码哈希升级落库失败:", e)
			}
		}
	}
	return s.Token(user["id"].Uint())
}

// MiniPhone 小程序手机号登录
func (s *UserLoginService) MiniPhone(ctx context.Context, code, encryptedData, iv string) (*UserTokenResult, error) {
	phone, err := NewUserWxService().MiniPhone(ctx, code, encryptedData, iv)
	if err != nil {
		return nil, err
	}
	return s.Phone(ctx, phone)
}

// UniPhone 一键手机号登录（经 uniphone 插件）
func (s *UserLoginService) UniPhone(ctx context.Context, accessToken, openid, appId string) (*UserTokenResult, error) {
	raw, err := pluginService.Invoke(ctx, "uniphone", "getPhone", accessToken, openid, appId)
	if err != nil {
		return nil, gerror.Wrap(err, "获得手机号失败，请检查配置")
	}
	phone := gconv.String(raw)
	if phone == "" {
		return nil, gerror.New("获得手机号失败，请检查配置")
	}
	return s.Phone(ctx, phone)
}

// Mini 小程序登录
func (s *UserLoginService) Mini(ctx context.Context, code, encryptedData, iv string) (*UserTokenResult, error) {
	wxService := NewUserWxService()
	openid, sessionKey, unionid, err := wxService.Code2Session(ctx, code)
	if err != nil {
		return nil, err
	}
	info := g.Map{
		"openid":  openid,
		"unionid": unionid,
		"type":    0,
	}
	if encryptedData != "" && iv != "" {
		if data, e := wxService.DecryptData(sessionKey, encryptedData, iv); e == nil {
			info["nickName"] = gconv.String(data["nickName"])
			info["avatarUrl"] = gconv.String(data["avatarUrl"])
			info["gender"] = gconv.Int(data["gender"])
		}
	}
	return s.wxLoginToken(ctx, info)
}

// Mp 公众号登录
func (s *UserLoginService) Mp(ctx context.Context, code string) (*UserTokenResult, error) {
	wxInfo, err := NewUserWxService().OAuth2UserInfo(ctx, code)
	if err != nil {
		return nil, err
	}
	return s.wxLoginToken(ctx, g.Map{
		"openid":    wxInfo.Openid,
		"unionid":   wxInfo.Unionid,
		"nickName":  wxInfo.NickName,
		"avatarUrl": wxInfo.AvatarUrl,
		"gender":    wxInfo.Sex,
		"city":      wxInfo.City,
		"province":  wxInfo.Province,
		"country":   wxInfo.Country,
		"type":      1,
	})
}

// WxApp 微信 APP 授权登录
func (s *UserLoginService) WxApp(ctx context.Context, code string) (*UserTokenResult, error) {
	wxInfo, err := NewUserWxService().OAuth2UserInfo(ctx, code)
	if err != nil {
		return nil, err
	}
	return s.wxLoginToken(ctx, g.Map{
		"openid":    wxInfo.Openid,
		"unionid":   wxInfo.Unionid,
		"nickName":  wxInfo.NickName,
		"avatarUrl": wxInfo.AvatarUrl,
		"gender":    wxInfo.Sex,
		"city":      wxInfo.City,
		"province":  wxInfo.Province,
		"country":   wxInfo.Country,
		"type":      3,
	})
}

// wxLoginToken 微信用户入库并签发会话
func (s *UserLoginService) wxLoginToken(ctx context.Context, wxInfo g.Map) (*UserTokenResult, error) {
	saved, err := saveWxInfo(ctx, wxInfo)
	if err != nil {
		return nil, err
	}
	unionid := gconv.String(saved["unionid"])
	if unionid == "" {
		unionid = gconv.String(saved["openid"])
	}
	user, err := v.DBM(userModel.NewUserInfo()).Where("unionid", unionid).One()
	if err != nil {
		return nil, err
	}
	if user.IsEmpty() {
		// 头像转存走 upload 插件（无插件则保留原始地址）
		avatar := gconv.String(wxInfo["avatarUrl"])
		if avatar != "" {
			if raw, e := pluginService.Invoke(ctx, "upload", "downAndUpload", avatar, gtime.Now().Format("YmdHis")+".png"); e == nil {
				avatar = gconv.String(raw)
			}
		}
		gender := gconv.Int(wxInfo["gender"])
		loginType := gconv.Int(wxInfo["type"])
		if _, err = v.DBM(userModel.NewUserInfo()).Data(g.Map{
			"unionid":   unionid,
			"nickName":  wxInfo["nickName"],
			"avatarUrl": avatar,
			"gender":    gender,
			"loginType": loginType,
		}).Insert(); err != nil {
			return nil, err
		}
		if user, err = v.DBM(userModel.NewUserInfo()).Where("unionid", unionid).One(); err != nil {
			return nil, err
		}
	}
	return s.Token(user["id"].Uint())
}

// RefreshToken 刷新会话：校验旧 refreshToken 后作废重签（轮换）
func (s *UserLoginService) RefreshToken(ctx context.Context, refreshToken string) (*UserTokenResult, error) {
	pair, err := sagf.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, gerror.New("刷新token失败，请检查refreshToken是否正确或过期")
	}
	userId := gconv.Uint(pair.LoginID)
	user, err := v.DBM(userModel.NewUserInfo()).Where("id", userId).One()
	if err != nil || user.IsEmpty() || user["status"].Int() != 1 {
		_ = sagf.RevokeRefreshToken(refreshToken)
		return nil, gerror.New("刷新token失败，请检查refreshToken是否正确或过期")
	}
	// 旧 refreshToken 作废，重新签发
	_ = sagf.RevokeRefreshToken(refreshToken)
	return s.Token(userId)
}

// Token 签发 sa-token 会话（access/refresh token 对，返回结构与原契约一致）
func (s *UserLoginService) Token(id uint) (*UserTokenResult, error) {
	info, err := v.DBM(userModel.NewUserInfo()).Where("id", id).One()
	if err != nil || info.IsEmpty() {
		return nil, gerror.New("用户不存在")
	}
	pair, err := sagf.LoginWithRefreshToken(id, userDevice)
	if err != nil {
		return nil, err
	}
	if err = setUserSession(id, info["tenantId"].Uint()); err != nil {
		return nil, err
	}
	return &UserTokenResult{
		Expire:        config.Config.Token.Expire,
		Token:         pair.AccessToken,
		RefreshExpire: config.Config.Token.RefreshExpire,
		RefreshToken:  pair.RefreshToken,
	}, nil
}

// UserLoginID 从访问令牌解析用户ID（供中间件使用；refreshToken/无效令牌返回错误）
func UserLoginID(token string) (uint, error) {
	id, err := sagf.GetLoginID(token)
	if err != nil {
		return 0, gerror.New("登录失效~")
	}
	return gconv.Uint(id), nil
}

// checkCaptcha 校验图片验证码
func checkCaptcha(ctx context.Context, captchaId, code string) bool {
	val, err := v.CacheManager.Get(ctx, captchaId)
	if err != nil || val.IsNil() {
		return false
	}
	return val.String() == code
}

// maskPhone 手机号脱敏
func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// saveWxInfo 保存/更新微信用户信息
func saveWxInfo(ctx context.Context, wxInfo g.Map) (map[string]any, error) {
	openid := gconv.String(wxInfo["openid"])
	existing, err := v.DBM(userModel.NewUserWx()).Where("openid", openid).One()
	if err != nil {
		return nil, err
	}
	data := g.Map{
		"unionid":   wxInfo["unionid"],
		"openid":    openid,
		"nickName":  wxInfo["nickName"],
		"avatarUrl": wxInfo["avatarUrl"],
		"gender":    wxInfo["gender"],
		"city":      wxInfo["city"],
		"province":  wxInfo["province"],
		"country":   wxInfo["country"],
		"language":  wxInfo["language"],
		"type":      wxInfo["type"],
	}
	if existing.IsEmpty() {
		if _, err = v.DBM(userModel.NewUserWx()).Data(data).Insert(); err != nil {
			return nil, err
		}
		return data, nil
	}
	if _, err = v.DBM(userModel.NewUserWx()).Where("id", existing["id"]).Data(data).Update(); err != nil {
		return nil, err
	}
	data["id"] = existing["id"]
	return data, nil
}
