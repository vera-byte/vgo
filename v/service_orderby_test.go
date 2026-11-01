package v

import "testing"

func TestSafeOrderBy(t *testing.T) {
	field, dir, err := safeOrderBy("createTime", "DESC")
	if err != nil || field != "createTime" || dir != "desc" {
		t.Fatalf("valid order rejected: %q %q %v", field, dir, err)
	}
	if _, _, err = safeOrderBy("a.userId", "asc"); err != nil {
		t.Fatalf("dotted field rejected: %v", err)
	}
	// 注入类输入必须拒绝
	for _, bad := range []string{"id; DROP TABLE t", "id--", "", "1", "(select 1)"} {
		if _, _, err = safeOrderBy(bad, "asc"); err == nil {
			t.Fatalf("injection accepted: %q", bad)
		}
	}
	if _, _, err = safeOrderBy("id", "asc; DROP TABLE t"); err == nil {
		t.Fatal("injection direction accepted")
	}
	if _, _, err = safeOrderBy("id", "random"); err == nil {
		t.Fatal("invalid direction accepted")
	}
}
