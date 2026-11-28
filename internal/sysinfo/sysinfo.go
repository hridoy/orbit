// Package sysinfo reads the host facts the node record shows: the address the server
// uses for outbound traffic, total memory, and total disk on the root filesystem.
package sysinfo

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// PrimaryIP returns the local address the kernel picks for outbound traffic. On a VPS with a
// public interface that is the public address. It returns "" when the host has no route.
// No packet is sent: a UDP "connect" only asks the kernel which source address it would use.
func PrimaryIP() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || addr.IP.IsLoopback() || addr.IP.IsUnspecified() {
		return ""
	}
	return addr.IP.String()
}

// MemoryTotal returns MemTotal from /proc/meminfo in bytes, or 0 when unavailable.
func MemoryTotal() int64 {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) >= 2 && fields[0] == "MemTotal:" {
			kb, _ := strconv.ParseInt(fields[1], 10, 64)
			return kb * 1024
		}
	}
	return 0
}

// DiskTotal returns the size of the filesystem holding path, in bytes, or 0 when unavailable.
func DiskTotal(path string) int64 {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return 0
	}
	return int64(st.Blocks) * int64(st.Bsize)
}

// DiskUsedPercent returns how full the filesystem holding path is, from 0 to 100, or -1 when
// unavailable.
func DiskUsedPercent(path string) int {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil || st.Blocks == 0 {
		return -1
	}
	used := st.Blocks - st.Bavail
	return int(float64(used) / float64(st.Blocks) * 100)
}

// MemoryAvailablePercent returns MemAvailable as a share of MemTotal, or -1 when unavailable.
func MemoryAvailablePercent() int {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return -1
	}
	defer f.Close()
	var total, avail int64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 2 {
			continue
		}
		v, _ := strconv.ParseInt(fields[1], 10, 64)
		switch fields[0] {
		case "MemTotal:":
			total = v
		case "MemAvailable:":
			avail = v
		}
	}
	if total <= 0 {
		return -1
	}
	return int(float64(avail) / float64(total) * 100)
}
