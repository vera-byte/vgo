package service

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/vera-byte/vgo/modules/base/model"
	"github.com/vera-byte/vgo/v"
)

// fakeDictProvider 假字典数据源：仅提供账号状态字典的测试数据，避免单测依赖数据库/缓存
type fakeDictProvider struct{}

func (fakeDictProvider) GetLabel(ctx context.Context, typeKey string, value string) (string, bool) {
	if typeKey != model.DictKeyUserStatus {
		return "", false
	}
	switch value {
	case "0":
		return "封禁", true
	case "3":
		return "空备注", true
	}
	return "", false
}

func (fakeDictProvider) GetValue(ctx context.Context, typeKey string, label string) (*gvar.Var, bool) {
	return nil, false
}

func (fakeDictProvider) GetOptions(ctx context.Context, typeKey string) []v.DictOption {
	return []v.DictOption{}
}

func (fakeDictProvider) GetTypes(ctx context.Context) []v.DictTypeInfo {
	return []v.DictTypeInfo{}
}

func (fakeDictProvider) GetInfo(ctx context.Context, typeKey string, label string) *v.DictEntry {
	if typeKey != model.DictKeyUserStatus {
		return nil
	}
	switch label {
	case "封禁":
		return &v.DictEntry{Name: label, Value: gvar.New("0"), Remark: "该账号已被封禁"}
	case "空备注":
		return &v.DictEntry{Name: label, Value: gvar.New("3")}
	}
	return nil
}

func TestCheckLoginStatus(t *testing.T) {
	v.RegisterDict(fakeDictProvider{})
	t.Cleanup(func() { v.RegisterDict(nil) })

	statusPtr := func(s int32) *int32 { return &s }

	cases := []struct {
		name    string
		status  *int32
		wantErr string // 空串表示放行
	}{
		{"正常状态放行", statusPtr(model.UserStatusNormal), ""},
		{"封禁取字典备注文案", statusPtr(model.UserStatusBanned), "该账号已被封禁"},
		{"未知状态默认文案", statusPtr(2), "当前账号状态不支持登录~"},
		{"空状态默认文案", nil, "当前账号状态不支持登录~"},
		{"条目备注为空回落默认文案", statusPtr(3), "当前账号状态不支持登录~"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkLoginStatus(context.Background(), c.status)
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("expected pass, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got pass", c.wantErr)
			}
			if err.Error() != c.wantErr {
				t.Fatalf("message mismatch: want %q, got %q", c.wantErr, err.Error())
			}
		})
	}
}
