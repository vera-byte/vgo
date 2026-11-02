package service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	pluginService "github.com/vera-byte/vgo/modules/plugin/service"
	sagf "github.com/sa-tokens/sa-token-go/integrations/gf"
	userModel "github.com/vera-byte/vgo/modules/user/model"
	"github.com/vera-byte/vgo/v"
)

// UserInfoService 用户信息服务（参照 cool-admin user 模块 UserInfoService）
type UserInfoService struct {
	*v.Service
}

// NewUserInfoService 创建用户信息服务
func NewUserInfoService() *UserInfoService {
	return &UserInfoService{
		Service: &v.Service{
			Model:              userModel.NewUserInfo(),
			InfoIgnoreProperty: "password",
			PageQueryOp: &v.QueryOp{
				KeyWordField: []string{"nickName", "phone"},
				AddOrderby:   g.MapStrStr{"createTime": "DESC"},
			},
		},
	}
}

// Person 获取个人信息（不含密码）
func (s *UserInfoService) Person(ctx context.Context, id uint) (map[string]any, error) {
	info, err := v.DBM(userModel.NewUserInfo()).Where("id", id).FieldsEx("password").One()
	if err != nil {
		return nil, err
	}
	if info.IsEmpty() {
		return nil, gerror.New("用户不存在")
	}
	return info.Map(), nil
}

// UpdatePerson 更新个人信息（白名单字段；头像变更经 upload 插件转存）
func (s *UserInfoService) UpdatePerson(ctx context.Context, id uint) error {
	r := g.RequestFromCtx(ctx)
	rMap := g.MapStrAny{}
	for _, field := range []string{"nickName", "avatarUrl", "gender", "description"} {
		if val := r.Get(field); !val.IsNil() {
			rMap[field] = val.Val()
		}
	}
	if len(rMap) == 0 {
		return nil
	}
	// 头像重新转存（无 upload 插件则保留原地址）
	if avatar := gconv.String(rMap["avatarUrl"]); avatar != "" {
		if raw, e := pluginService.Invoke(ctx, "upload", "downAndUpload", avatar, "avatar.png"); e == nil {
			rMap["avatarUrl"] = gconv.String(raw)
		}
	}
	if _, err := v.DBM(userModel.NewUserInfo()).Where("id", id).Data(rMap).Update(); err != nil {
		return gerror.New("更新失败，参数错误或者手机号已存在")
	}
	return nil
}

// UpdatePassword 修改密码（短信验证码校验）
func (s *UserInfoService) UpdatePassword(ctx context.Context, id uint, password, code string) error {
	user, err := v.DBM(userModel.NewUserInfo()).Where("id", id).One()
	if err != nil || user.IsEmpty() {
		return gerror.New("用户不存在")
	}
	if !NewUserSmsService().CheckCode(ctx, user["phone"].String(), code) {
		return gerror.New("验证码错误")
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}
	if _, err = v.DBM(userModel.NewUserInfo()).Where("id", id).Data(g.Map{"password": hashed}).Update(); err != nil {
		return err
	}
	// 密码变更后强制下线该用户全部会话
	_ = sagf.Logout(id, userDevice)
	return nil
}

// BindPhone 绑定手机号（短信验证码校验）
func (s *UserInfoService) BindPhone(ctx context.Context, id uint, phone, code string) error {
	if !NewUserSmsService().CheckCode(ctx, phone, code) {
		return gerror.New("验证码错误")
	}
	_, err := v.DBM(userModel.NewUserInfo()).Where("id", id).Data(g.Map{"phone": phone}).Update()
	return err
}

// MiniPhone 绑定小程序手机号
func (s *UserInfoService) MiniPhone(ctx context.Context, id uint, wxCode, encryptedData, iv string) (string, error) {
	phone, err := NewUserWxService().MiniPhone(ctx, wxCode, encryptedData, iv)
	if err != nil {
		return "", err
	}
	_, err = v.DBM(userModel.NewUserInfo()).Where("id", id).Data(g.Map{"phone": phone}).Update()
	return phone, err
}

// Logoff 注销（脱敏并标记已注销）
func (s *UserInfoService) Logoff(ctx context.Context, id uint) error {
	nick := "已注销-00" + gconv.String(id)
	_, err := v.DBM(userModel.NewUserInfo()).Where("id", id).Data(g.Map{
		"status":    2,
		"phone":     nil,
		"unionid":   nil,
		"nickName":  nick,
		"avatarUrl": nil,
	}).Update()
	return err
}
