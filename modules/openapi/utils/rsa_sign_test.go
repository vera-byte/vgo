package utils

import (
	"strings"
	"testing"
)

func TestSignDataFormats(t *testing.T) {
	canonical := signDataCanonical(`{"a":1}`, "app1", 1700000000, "n1")
	legacy := signDataLegacy(`{"a":1}`, "app1", 1700000000, "n1")
	if canonical == legacy {
		t.Fatal("canonical and legacy formats should differ")
	}
	if strings.Count(canonical, "\n") != 3 {
		t.Fatalf("canonical should have 3 separators: %q", canonical)
	}
	// 边界歧义：旧格式下 body 拼接 appId 可产生相同串，规范格式必须区分
	c1 := signDataCanonical(`{"a":1}`, "app", 1, "n")
	c2 := signDataCanonical(`{"a":1}x`, "pp", 1, "n")
	if c1 == c2 {
		t.Fatal("ambiguous concatenation not prevented")
	}
}

func TestGenerateAndVerifySignature(t *testing.T) {
	u := &RSASignUtil{}
	priv, pub, err := u.GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	sig, err := u.GenerateSignature(`{"k":"v"}`, "app1", 1700000000, "nonce1", priv)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	ok, err := u.VerifySignature(`{"k":"v"}`, "app1", 1700000000, "nonce1", sig, pub)
	if err != nil || !ok {
		t.Fatalf("verify failed: %v %v", ok, err)
	}
	// 篡改 nonce 必须验签失败
	ok, _ = u.VerifySignature(`{"k":"v"}`, "app1", 1700000000, "nonce2", sig, pub)
	if ok {
		t.Fatal("tampered payload verified")
	}
}
