package service

import (
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword 使用 bcrypt 生成密码哈希。
func hashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", gerror.Wrap(err, "密码哈希失败")
	}
	return string(b), nil
}

// checkPassword 校验密码：新记录使用 bcrypt；兼容历史无盐 MD5 哈希。
// legacy 为 true 表示命中旧哈希格式，校验通过后应升级为 bcrypt 存储。
func checkPassword(plain, hash string) (ok bool, legacy bool) {
	if strings.HasPrefix(hash, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil, false
	}
	md5Plain, err := gmd5.Encrypt(plain)
	if err != nil {
		return false, true
	}
	return md5Plain == hash, true
}
