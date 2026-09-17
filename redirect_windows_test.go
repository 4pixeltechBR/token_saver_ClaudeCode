//go:build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestWindowsJunctionWriteRefused(t *testing.T) {
	o := fixture(t)
	outside := filepath.Join(filepath.Dir(o.Project), "outside")
	if err := os.MkdirAll(outside, 0700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(o.Project, ".claude")
	if b, err := exec.Command("cmd.exe", "/c", "mklink", "/J", link, outside).CombinedOutput(); err != nil {
		t.Fatalf("create junction: %s: %v", b, err)
	}
	o.Concise = true
	a, err := audit(o)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := makePlan(o, a); err == nil {
		t.Fatal("junction accepted as a normal directory")
	}
	entries, err := os.ReadDir(outside)
	if err != nil || len(entries) != 0 {
		t.Fatal("redirected target changed")
	}
}
