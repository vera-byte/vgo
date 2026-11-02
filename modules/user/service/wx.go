package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/vera-byte/vgo/modules/user/config"
)

// UserWxService 微信开放平台服务（参照 cool-admin user 模块 UserWxService）
type UserWxService struct{}

// NewUserWxService 创建微信服务实例
func NewUserWxService() *UserWxService {
	return &UserWxService{}
}

// WxUserInfo 微信用户信息
type WxUserInfo struct {
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	NickName   string `json:"nickname"`
	AvatarUrl  string `json:"headimgurl"`
	Sex        int    `json:"sex"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	Privilege  []any  `json:"privilege"`
	SessionKey string `json:"-"`
}

// getJSON GET 请求并解析 JSON
func getJSON(ctx context.Context, uri string, out any) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

// Code2Session 小程序登录凭证换取会话（jscode2session）
func (s *UserWxService) Code2Session(ctx context.Context, code string) (openid, sessionKey, unionid string, err error) {
	uri := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		url.QueryEscape(config.Config.Wx.AppId), url.QueryEscape(config.Config.Wx.AppSecret), url.QueryEscape(code))
	var result struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		Unionid    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err = getJSON(ctx, uri, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		return "", "", "", gerror.Newf("微信登录失败: %s", result.ErrMsg)
	}
	return result.Openid, result.SessionKey, result.Unionid, nil
}

// OAuth2UserInfo 网页授权（公众号/APP）换取用户信息
func (s *UserWxService) OAuth2UserInfo(ctx context.Context, code string) (*WxUserInfo, error) {
	tokenURI := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		url.QueryEscape(config.Config.Wx.AppId), url.QueryEscape(config.Config.Wx.AppSecret), url.QueryEscape(code))
	var token struct {
		AccessToken string `json:"access_token"`
		Openid      string `json:"openid"`
		Unionid     string `json:"unionid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := getJSON(ctx, tokenURI, &token); err != nil {
		return nil, err
	}
	if token.ErrCode != 0 {
		return nil, gerror.Newf("微信授权失败: %s", token.ErrMsg)
	}
	userURI := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		url.QueryEscape(token.AccessToken), url.QueryEscape(token.Openid))
	info := &WxUserInfo{}
	if err := getJSON(ctx, userURI, info); err != nil {
		return nil, err
	}
	if info.Unionid == "" {
		info.Unionid = token.Unionid
	}
	return info, nil
}

// DecryptData 解密微信加密数据（AES-256-CBC，sessionKey/base64(iv)）
func (s *UserWxService) DecryptData(sessionKey, encryptedData, iv string) (map[string]any, error) {
	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, gerror.Wrap(err, "sessionKey 解码失败")
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, gerror.Wrap(err, "encryptedData 解码失败")
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, gerror.Wrap(err, "iv 解码失败")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) < block.BlockSize() {
		return nil, gerror.New("加密数据长度非法")
	}
	plain := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, ivBytes).CryptBlocks(plain, cipherText)
	plain = pkcs7Unpad(plain)
	var out map[string]any
	if err = json.Unmarshal(plain, &out); err != nil {
		return nil, gerror.Wrap(err, "解密数据解析失败")
	}
	return out, nil
}

// sha1Hex 计算 SHA1 十六进制摘要（JS-SDK 签名）
func sha1Hex(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// pkcs7Unpad 去除 PKCS7 填充
func pkcs7Unpad(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	pad := int(data[len(data)-1])
	if pad <= 0 || pad > len(data) {
		return data
	}
	return data[:len(data)-pad]
}

// MiniPhone 小程序手机号解密（encryptedData 解出 phoneNumber）
func (s *UserWxService) MiniPhone(ctx context.Context, code, encryptedData, iv string) (string, error) {
	_, sessionKey, _, err := s.Code2Session(ctx, code)
	if err != nil {
		return "", err
	}
	data, err := s.DecryptData(sessionKey, encryptedData, iv)
	if err != nil {
		return "", err
	}
	phone := gconv.String(data["phoneNumber"])
	if phone == "" {
		return "", gerror.New("获得手机号失败，请检查配置")
	}
	return phone, nil
}

// WxMpConfig 公众号 JS-SDK 配置（签名）
func (s *UserWxService) WxMpConfig(ctx context.Context, pageUrl string) (map[string]any, error) {
	var token struct {
		AccessToken string `json:"access_token"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	tkURI := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		url.QueryEscape(config.Config.Wx.AppId), url.QueryEscape(config.Config.Wx.AppSecret))
	if err := getJSON(ctx, tkURI, &token); err != nil {
		return nil, err
	}
	if token.ErrCode != 0 {
		return nil, gerror.Newf("获取 access_token 失败: %s", token.ErrMsg)
	}
	var ticket struct {
		Ticket  string `json:"ticket"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	tk2 := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi",
		url.QueryEscape(token.AccessToken))
	if err := getJSON(ctx, tk2, &ticket); err != nil {
		return nil, err
	}
	if ticket.ErrCode != 0 {
		return nil, gerror.Newf("获取 jsapi_ticket 失败: %s", ticket.ErrMsg)
	}
	nonceStr := grand.S(16, false)
	timestamp := gtime.Now().Unix()
	signSrc := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket.Ticket, nonceStr, timestamp, pageUrl)
	signature := sha1Hex(signSrc)
	return map[string]any{
		"appId":     config.Config.Wx.AppId,
		"timestamp": timestamp,
		"nonceStr":  nonceStr,
		"signature": signature,
	}, nil
}
