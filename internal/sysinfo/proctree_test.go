package sysinfo

import (
	"os"
	"os/exec"
	"testing"
)

func TestProcessTreeRSSIncludesChildren(t *testing.T) {
	if _, err := os.Stat("/proc/self/status"); err != nil {
		t.Skip("no /proc on this platform")
	}
	self := ProcessTreeRSS(os.Getpid())
	if self == 0 {
		t.Fatal("own process reported zero RSS")
	}
	child := exec.Command("sleep", "5")
	if err := child.Start(); err != nil {
		t.Skip("sleep not available:", err)
	}
	defer func() { _ = child.Process.Kill(); _ = child.Wait() }()
	withChild := ProcessTreeRSS(os.Getpid())
	if withChild <= self {
		t.Fatalf("tree with a child (%d) should exceed own RSS (%d)", withChild, self)
	}
	if got := ProcessTreeRSS(child.Process.Pid); got == 0 {
		t.Fatal("child RSS should be non-zero")
	}
}
