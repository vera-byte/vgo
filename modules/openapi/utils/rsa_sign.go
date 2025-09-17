package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
)

// RSASignUtil RSA签名工具类
type RSASignUtil struct{}

// NewRSASignUtil 创建RSA签名工具实例
// 返回值: *RSASignUtil - RSA签名工具实例
func NewRSASignUtil() *RSASignUtil {
	return &RSASignUtil{}
}

// GenerateSignature 生成RSA签名
// 功能: 根据请求体、应用ID、时间戳、随机数和RSA私钥生成base64编码的签名
// 参数:
//   - jsonBody: JSON格式的请求体
//   - appId: 应用ID
//   - timestamp: 时间戳
//   - nonce: 随机数
//   - privateKeyPEM: PEM格式的RSA私钥
// 返回值:
//   - string: base64编码的签名结果
//   - error: 错误信息
func (r *RSASignUtil) GenerateSignature(jsonBody, appId string, timestamp int64, nonce, privateKeyPEM string) (string, error) {
	// 1. 构建待签名字符串: json-body + Appid + Timestamp + Nonce
	signData := jsonBody + appId + strconv.FormatInt(timestamp, 10) + nonce

	// 2. 解析RSA私钥
	privateKey, err := r.parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", gerror.Wrap(err, "解析RSA私钥失败")
	}

	// 3. 计算SHA256哈希
	hash := sha256.Sum256([]byte(signData))

	// 4. 使用RSA私钥进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", gerror.Wrap(err, "RSA签名失败")
	}

	// 5. 对签名结果进行base64编码
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	return signatureBase64, nil
}

// VerifySignature 验证RSA签名
// 功能: 验证签名是否正确
// 参数:
//   - jsonBody: JSON格式的请求体
//   - appId: 应用ID
//   - timestamp: 时间戳
//   - nonce: 随机数
//   - signature: base64编码的签名
//   - publicKeyPEM: PEM格式的RSA公钥
// 返回值:
//   - bool: 验证结果，true表示验证成功
//   - error: 错误信息
func (r *RSASignUtil) VerifySignature(jsonBody, appId string, timestamp int64, nonce, signature, publicKeyPEM string) (bool, error) {
	// 1. 构建待验证字符串
	signData := jsonBody + appId + strconv.FormatInt(timestamp, 10) + nonce

	// 2. 解析RSA公钥
	publicKey, err := r.parsePublicKey(publicKeyPEM)
	if err != nil {
		return false, gerror.Wrap(err, "解析RSA公钥失败")
	}

	// 3. 解码base64签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, gerror.Wrap(err, "base64解码签名失败")
	}

	// 4. 计算SHA256哈希
	hash := sha256.Sum256([]byte(signData))

	// 5. 验证签名
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return false, nil // 验证失败，但不是错误
	}

	return true, nil
}

// GenerateRSAKeyPair 生成RSA密钥对
// 功能: 生成指定位数的RSA密钥对
// 参数:
//   - bits: 密钥位数，建议2048或4096
// 返回值:
//   - string: PEM格式的私钥
//   - string: PEM格式的公钥
//   - error: 错误信息
func (r *RSASignUtil) GenerateRSAKeyPair(bits int) (string, string, error) {
	// 1. 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", gerror.Wrap(err, "生成RSA密钥对失败")
	}

	// 2. 编码私钥为PEM格式
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", gerror.Wrap(err, "编码私钥失败")
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 3. 编码公钥为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", gerror.Wrap(err, "编码公钥失败")
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// parsePrivateKey 解析PEM格式的RSA私钥
// 功能: 将PEM格式的私钥字符串解析为RSA私钥对象
// 参数:
//   - privateKeyPEM: PEM格式的私钥字符串
// 返回值:
//   - *rsa.PrivateKey: RSA私钥对象
//   - error: 错误信息
func (r *RSASignUtil) parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, gerror.New("无效的PEM格式私钥")
	}

	var privateKey interface{}
	var err error

	// 尝试不同的私钥格式
	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	default:
		return nil, gerror.New(fmt.Sprintf("不支持的私钥类型: %s", block.Type))
	}

	if err != nil {
		return nil, gerror.Wrap(err, "解析私钥失败")
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, gerror.New("私钥不是RSA格式")
	}

	return rsaPrivateKey, nil
}

// parsePublicKey 解析PEM格式的RSA公钥
// 功能: 将PEM格式的公钥字符串解析为RSA公钥对象
// 参数:
//   - publicKeyPEM: PEM格式的公钥字符串
// 返回值:
//   - *rsa.PublicKey: RSA公钥对象
//   - error: 错误信息
func (r *RSASignUtil) parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, gerror.New("无效的PEM格式公钥")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, gerror.Wrap(err, "解析公钥失败")
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, gerror.New("公钥不是RSA格式")
	}

	return rsaPublicKey, nil
}