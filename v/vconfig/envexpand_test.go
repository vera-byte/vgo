package vconfig

import "testing"

func TestExpandEnv(t *testing.T) {
	t.Setenv("VGO_TEST_SECRET", "smoke-value")
	if got := ExpandEnv("pre-${VGO_TEST_SECRET}-post"); got != "pre-smoke-value-post" {
		t.Fatalf("ExpandEnv failed: %q", got)
	}
	if got := ExpandEnv("${VGO_TEST_UNSET_VAR_XYZ}"); got != "" {
		t.Fatalf("unset var should expand to empty: %q", got)
	}
	if got := ExpandEnv("no-placeholder"); got != "no-placeholder" {
		t.Fatalf("plain string changed: %q", got)
	}
	if got := ExpandEnv("$VGO_TEST_SECRET"); got != "$VGO_TEST_SECRET" {
		t.Fatalf("bare $VAR should not expand: %q", got)
	}
}
