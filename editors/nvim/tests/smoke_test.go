// Package tests runs the Neovim plugin against a real nvim when one is
// available.
package tests

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// TestPluginSmoke drives a headless Neovim through filetype detection, the
// generated syntax and the LSP client (testdata/smoke.lua). It skips without
// nvim 0.10+ on PATH; OPENSYSML_REQUIRE_NVIM=1 fails instead.
func TestPluginSmoke(t *testing.T) {
	nvim, err := exec.LookPath("nvim")
	if err != nil {
		requireOrSkip(t, "nvim is not on PATH")
	}
	version, err := exec.Command(nvim, "--version").Output()
	if err != nil {
		t.Fatalf("nvim --version: %v", err)
	}
	if !atLeast010(version) {
		requireOrSkip(t, "nvim is older than 0.10")
	}

	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, nvim, "--headless", "--clean",
		"-l", filepath.Join("testdata", "smoke.lua"), root)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("smoke test failed: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("all checks passed")) {
		t.Fatalf("smoke test did not pass:\n%s", out)
	}
}

func requireOrSkip(t *testing.T, reason string) {
	t.Helper()
	if os.Getenv("OPENSYSML_REQUIRE_NVIM") != "" {
		t.Fatalf("%s and OPENSYSML_REQUIRE_NVIM is set", reason)
	}
	t.Skip(reason)
}

var versionPattern = regexp.MustCompile(`NVIM v(\d+)\.(\d+)`)

func atLeast010(version []byte) bool {
	m := versionPattern.FindSubmatch(version)
	if m == nil {
		return false
	}
	major, _ := strconv.Atoi(string(m[1]))
	minor, _ := strconv.Atoi(string(m[2]))
	return major > 0 || minor >= 10
}
