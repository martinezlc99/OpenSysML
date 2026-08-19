package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Open-MBEE/OpenSysML/internal/core/lexer"
)

// TestCommittedSyntaxIsCurrent is the drift gate: the syntax files the plugin
// ships must be what the current keyword list generates.
func TestCommittedSyntaxIsCurrent(t *testing.T) {
	for _, s := range Syntaxes() {
		want, err := Render(s)
		if err != nil {
			t.Fatalf("Render(%s) err = %v", s.File, err)
		}
		path := filepath.Join("..", "..", "syntax", s.File)
		got, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if string(got) != string(want) {
			t.Errorf("%s is stale; regenerate it with `make nvim-syntax`", path)
		}
	}
}

func TestEveryKeywordIsHighlighted(t *testing.T) {
	data, err := Render(Syntaxes()[0])
	if err != nil {
		t.Fatalf("Render err = %v", err)
	}

	// `syn keyword <group> kw kw ...` lines carry only keywords and item
	// options (nextgroup=..., skipwhite) after the group name.
	highlighted := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != "syn" || fields[1] != "keyword" {
			continue
		}
		for _, kw := range fields[3:] {
			if strings.Contains(kw, "=") || kw == "skipwhite" {
				continue
			}
			highlighted[kw] = true
		}
	}

	for _, kw := range lexer.Keywords() {
		if !highlighted[kw] {
			t.Errorf("keyword %q is not in any generated syn keyword group", kw)
		}
	}
}

func TestGroupsAreScopedPerFiletype(t *testing.T) {
	for _, s := range Syntaxes() {
		data, err := Render(s)
		if err != nil {
			t.Fatalf("Render(%s) err = %v", s.File, err)
		}
		text := string(data)
		if !strings.Contains(text, "syn keyword "+s.Filetype+"Declaration ") {
			t.Errorf("%s does not prefix its groups with %q", s.File, s.Filetype)
		}
		if !strings.Contains(text, "let b:current_syntax = '"+s.Filetype+"'") {
			t.Errorf("%s does not set b:current_syntax to %q", s.File, s.Filetype)
		}
	}
}
