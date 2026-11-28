package sysinfo

import "testing"

func TestMemoryAndDiskAreReported(t *testing.T) {
	if MemoryTotal() <= 0 {
		t.Skip("no /proc/meminfo on this platform")
	}
	if DiskTotal("/") <= 0 {
		t.Fatal("disk total should be positive for /")
	}
}
