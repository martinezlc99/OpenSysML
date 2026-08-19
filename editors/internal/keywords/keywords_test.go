package keywords

import (
	"testing"

	"github.com/Open-MBEE/OpenSysML/internal/core/lexer"
)

func TestPartitionCoversTheLexerExactly(t *testing.T) {
	groups, rest, err := Partition()
	if err != nil {
		t.Fatalf("Partition() err = %v", err)
	}

	seen := map[string]bool{}
	for _, g := range groups {
		for _, kw := range g.Keywords {
			if seen[kw] {
				t.Errorf("keyword %q claimed twice", kw)
			}
			seen[kw] = true
		}
	}
	for _, kw := range rest {
		if seen[kw] {
			t.Errorf("keyword %q is grouped and in the leftovers", kw)
		}
		seen[kw] = true
	}

	for _, kw := range lexer.Keywords() {
		if !seen[kw] {
			t.Errorf("lexer keyword %q is neither grouped nor left over", kw)
		}
	}
	if got, want := len(seen), len(lexer.Keywords()); got != want {
		t.Errorf("partition names %d keywords, lexer has %d", got, want)
	}
}

func TestPartitionRejectsUnknownKeywords(t *testing.T) {
	bad := []Group{{Name: "control", Keywords: []string{"nosuchkeyword"}}}
	if _, _, err := partition(bad); err == nil {
		t.Error("partition with an unknown keyword succeeded, want an error")
	}
}

func TestPartitionRejectsDuplicateKeywords(t *testing.T) {
	bad := []Group{
		{Name: "control", Keywords: []string{"if"}},
		{Name: "modifier", Keywords: []string{"if"}},
	}
	if _, _, err := partition(bad); err == nil {
		t.Error("partition with a duplicated keyword succeeded, want an error")
	}
}
