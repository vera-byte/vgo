package service

import (
	"testing"

	"github.com/gogf/gf/v2/crypto/gmd5"
)

func TestHashAndCheckPassword(t *testing.T) {
	hashed, err := hashPassword("s3cret-Pass")
	if err != nil {
		t.Fatalf("hashPassword: %v", err)
	}
	ok, legacy := checkPassword("s3cret-Pass", hashed)
	if !ok || legacy {
		t.Fatalf("bcrypt verify failed: ok=%v legacy=%v", ok, legacy)
	}
	ok, _ = checkPassword("wrong", hashed)
	if ok {
		t.Fatal("wrong password accepted")
	}
}

func TestCheckPasswordLegacyMD5(t *testing.T) {
	md5Hash, _ := gmd5.Encrypt("old-pass")
	ok, legacy := checkPassword("old-pass", md5Hash)
	if !ok || !legacy {
		t.Fatalf("legacy md5 verify failed: ok=%v legacy=%v", ok, legacy)
	}
	ok, _ = checkPassword("wrong", md5Hash)
	if ok {
		t.Fatal("wrong legacy password accepted")
	}
}
