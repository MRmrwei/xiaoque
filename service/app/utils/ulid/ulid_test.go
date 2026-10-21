package ulid

import (
	"strings"
	"testing"
)

func TestNewGeneratesValidUniqueULID(t *testing.T) {
	seen := make(map[string]struct{}, 200)
	for index := 0; index < 200; index++ {
		value, err := New()
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if len(value) != 26 {
			t.Fatalf("ULID 长度 = %d，期望 26：%s", len(value), value)
		}
		if value[0] < '0' || value[0] > '7' {
			t.Fatalf("ULID 首字符越界：%s", value)
		}
		for _, char := range value {
			if !strings.ContainsRune(encoding, char) {
				t.Fatalf("ULID 包含非法字符 %q：%s", char, value)
			}
		}
		if _, exists := seen[value]; exists {
			t.Fatalf("生成了重复 ULID：%s", value)
		}
		seen[value] = struct{}{}
	}
}
