package file

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	vck "github.com/vera-byte/vgo/vgo_core_kit"
)

type AliyunOssSignResult struct {
	Policy                string `json:"policy"`
	OSSAccessKeyId        string `json:"OSSAccessKeyId"`
	Signature             string `json:"Signature"`
	Host                  string `json:"host"`
	Key                   string `json:"key"`
	SUCCESS_ACTION_STATUS int32  `json:"success_action_status"`
}
type File struct{}

func (*File) SignUrl(ctx context.Context, dir string, uid int64) (res *AliyunOssSignResult, err error) {
	var (
		CachePrefix = "file:oss:sign:"

		Key = fmt.Sprintf("%s/%d/", dir, uid)
	)
	// 查看缓存是否存在
	cache, err := vck.CacheManager.Get(ctx, fmt.Sprintf("%s%s", CachePrefix, Key))
	cache.Scan(&res)

	if err != nil {
		g.Log().Error(ctx, "读取缓存出错", err)
	}
	if res != nil {
		g.Log().Info(ctx, "当前签名信息数据来源RedisCache")
		return res, nil
	}
	if vck.GetAdminConfig.File.Oss.AccessKeyId == "" {

		g.Log().Error(ctx, fmt.Errorf("ACCESS_KEY_ID Is Empty %s", "请设置环境变量ACCESS_KEY_ID"))

	}
	if vck.GetAdminConfig.File.Oss.AccessKeySecret == "" {
		g.Log().Error(ctx, fmt.Errorf("ACCESS_KEY_SECRET Is Empty %s", "请设置环境变量ACCESS_KEY_SECRET"))
	}
	// 获取当前时间并加上环境变量的过期秒作为过期时间
	OssExpiresTimestamp := time.Now().UTC().Add(time.Duration(vck.GetAdminConfig.File.Oss.Expires) * time.Second)
	// 构建策略
	expiration := OssExpiresTimestamp.Format("2006-01-02T15:04:05.000Z")
	policy := map[string]interface{}{
		"expiration": expiration,
		"conditions": []interface{}{
			[]interface{}{"content-length-range", 0, vck.GetAdminConfig.File.Oss.ContentLength},
		},
	}

	// 将策略转换为 JSON 并编码为 Base64
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		g.Log().Error(ctx, "将策略转换为 JSON 并编码为 Base64 失败", err)
	}
	policyBase64 := base64.StdEncoding.EncodeToString(policyJSON)

	// 计算签名
	h := hmac.New(sha1.New, []byte(vck.GetAdminConfig.File.Oss.AccessKeySecret))
	h.Write([]byte(policyBase64))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// 返回RPC结果
	resp := &AliyunOssSignResult{
		Policy:                policyBase64,
		OSSAccessKeyId:        vck.GetAdminConfig.File.Oss.AccessKeyId,
		Signature:             signature,
		Host:                  vck.GetAdminConfig.File.Oss.Host,
		Key:                   Key,
		SUCCESS_ACTION_STATUS: vck.GetAdminConfig.File.Oss.SuccessActionStatus,
	}
	// 缓存当前签名
	err = vck.CacheManager.Set(ctx, fmt.Sprintf("%s%s", CachePrefix, Key), resp, time.Duration(vck.GetAdminConfig.File.Oss.Expires)*time.Second)
	if err != nil {
		g.Log().Error(ctx, "签名信息存储失败")
	}
	return resp, nil

}
