package base

import (
	"context"
	"time"
	"vgo/app/admin/internal/dao"
	"vgo/app/admin/internal/model/entity"
	"vgo/app/admin/internal/service"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mojocn/base64Captcha"
	vck "github.com/vera-byte/vgo/vgo_core_kit"
	vck_config "github.com/vera-byte/vgo/vgo_core_kit/config"
)

func init() {
	service.RegisterBaseSysLoginLogic(NewBaseSysLoginLogic())
}

type sBaseSysLoginLogic struct {
	store base64Captcha.Store
}

func NewBaseSysLoginLogic() *sBaseSysLoginLogic {
	return &sBaseSysLoginLogic{
		store: base64Captcha.DefaultMemStore,
	}
}

// 生成验证码
func (c *sBaseSysLoginLogic) GenerateCaptcha(ctx context.Context, width int, height int) (id string, b64s string, answer string, err error) {

	driver := base64Captcha.NewDriverDigit(height, width, 4, 0, 10)
	captcha := base64Captcha.NewCaptcha(driver, c.store)
	id, b64s, answer, err = captcha.Generate()
	if err != nil {
		return "", "", "", err
	}
	return
}

// 验证验证码
func (c *sBaseSysLoginLogic) VerifyCaptcha(id, answer string) bool {
	return c.store.Verify(id, answer, true)
}

// 退出登录
func (c *sBaseSysLoginLogic) LoginOut(ctx context.Context) (err error) {
	var (
		admin = vck.GetAdminAtGateway(ctx)
	)
	if admin == nil {
		return gerror.New("用户不存在")
	}
	vck.CacheManager.Remove(ctx, "admin:passwordVersionadmin:passwordVersion:"+gconv.String(admin.UserId))
	vck.CacheManager.Remove(ctx, "admin:department:"+gconv.String(admin.UserId))
	vck.CacheManager.Remove(ctx, "admin:perms:"+gconv.String(admin.UserId))
	vck.CacheManager.Remove(ctx, "admin:token:"+gconv.String(admin.UserId))
	vck.CacheManager.Remove(ctx, "admin:token:refresh:"+gconv.String(admin.UserId))
	return nil
}

// 密码登录 此处只验证密码和验证码 Token由其他函数生成
func (c *sBaseSysLoginLogic) Login(ctx context.Context, captchaId string, password string, userName string, code string) (result *vck_config.TokenResult, err error) {
	var (
		md5password, _ = gmd5.Encrypt(password)
		user           *entity.BaseSysUser
	)

	if !c.VerifyCaptcha(captchaId, code) {
		return nil, gerror.New("验证码错误")
	}

	err = dao.BaseSysUser.Ctx(ctx).Where("username=?", userName).Scan(&user)
	if err != nil {
		return nil, gerror.New("系统异常!")
	}
	if user == nil {
		err = gerror.New("账户或密码不正确~")
		return
	}
	if user.Status == 0 || user.Password != md5password {
		err = gerror.New("账户或密码不正确~")
		return
	}
	result, err = c.generateTokenByUser(ctx, user)
	if err != nil {
		return
	}
	return
}

func (c *sBaseSysLoginLogic) RefreshToken(ctx context.Context, _refreshToken string) (result *vck_config.TokenResult, err error) {
	claims, err := vck_config.ParseToken(ctx, _refreshToken)
	if err != nil {
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
		user *entity.BaseSysUser
	)
	dao.BaseSysUser.Ctx(ctx).Where(dao.BaseSysUser.Columns().Id, claims.UserId).Where(dao.BaseSysUser.Columns().Status, 1).Scan(&user)
	if user == nil {
		err = gerror.New("用户不存在")
		return
	}

	result, err = c.generateTokenByUser(ctx, user)
	if err != nil {
		return
	}
	return nil, nil
}

// 根据用户生成前端需要的Token信息
func (c *sBaseSysLoginLogic) generateTokenByUser(ctx context.Context, user *entity.BaseSysUser) (result *vck_config.TokenResult, err error) {
	var (
		baseSysRoleService       = service.BaseSysRoleLogic()
		baseSysMenuService       = service.BaseSysMenuLogic()
		baseSysDepartmentService = service.BaseSysDepartmentLogic()
		_refreshExpire           = vck.GetAdminConfig.Jwt.Token.RefreshExpire
		_expire                  = vck.GetAdminConfig.Jwt.Token.Expire
	)
	// 获取用户角色
	roleIds, err := baseSysRoleService.GetByUser(ctx, int64(user.Id))
	if err != nil {
		return
	}
	// 如果没有角色，则报错
	if len(roleIds) == 0 {
		err = gerror.New("该用户未设置任何角色，无法登录~")
		return
	}
	// 生成token
	result = &vck_config.TokenResult{}
	result.Expire = _expire
	result.RefreshExpire = _refreshExpire
	result.Token = c.generateToken(ctx, int64(user.Id), roleIds, user.Username, user.PasswordV, result.Expire, false)
	result.RefreshToken = c.generateToken(ctx, int64(user.Id), roleIds, user.Username, user.PasswordV, result.Expire, true)

	// 将用户相关信息保存到缓存
	perms := baseSysMenuService.GetPerms(ctx, roleIds)
	departments := baseSysDepartmentService.GetByRoleIds(ctx, roleIds, user.Username == "admin")
	vck.CacheManager.Set(ctx, "admin:department:"+gconv.String(user.Id), departments, 0)
	vck.CacheManager.Set(ctx, "admin:perms:"+gconv.String(user.Id), perms, 0)
	vck.CacheManager.Set(ctx, "admin:token:"+gconv.String(user.Id), result.Token, 0)
	vck.CacheManager.Set(ctx, "admin:token:refresh:"+gconv.String(user.Id), result.RefreshToken, 0)

	return
}

// 生成token
func (c *sBaseSysLoginLogic) generateToken(ctx context.Context, userId int64, roleIds []string, username string, userPasswordV int, expire int64, isRefresh bool) string {
	err := vck.CacheManager.Set(ctx, "admin:passwordVersion:"+gconv.String(userId), gconv.String(userPasswordV), 0)
	if err != nil {
		g.Log().Error(ctx, "生成token失败", err)
	}
	claims := &vck_config.Claims{
		IsRefresh:       isRefresh,
		RoleIds:         roleIds,
		Username:        username,
		UserId:          userId,
		PasswordVersion: userPasswordV,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString([]byte(vck.GetAdminConfig.Jwt.Secret))
	if err != nil {
		g.Log().Error(ctx, "生成token失败", err)
		return ""
	}
	return token
}
