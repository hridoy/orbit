//go:build linux

package compute

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const cgroupV2Root = "/sys/fs/cgroup"
const orbitCgroupBase = "/sys/fs/cgroup/orbit"

type LinuxCgroupController struct {
	available bool
}

func newPlatformCgroupController() CgroupController {
	available := false
	if info, err := os.Stat(cgroupV2Root); err == nil && info.IsDir() {
		// Test if cgroups v2 controllers file is present
		if _, err := os.Stat(filepath.Join(cgroupV2Root, "cgroup.controllers")); err == nil {
			available = true
			_ = os.MkdirAll(orbitCgroupBase, 0755)
		}
	}
	return &LinuxCgroupController{available: available}
}

func (c *LinuxCgroupController) cgroupPath(appID string) string {
	return filepath.Join(orbitCgroupBase, appID)
}

func (c *LinuxCgroupController) Setup(appID string, limits CgroupLimits) error {
	if !c.available {
		return nil // Graceful degradation if cgroups v2 is not mounted
	}

	cgPath := c.cgroupPath(appID)
	if err := os.MkdirAll(cgPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup %s: %w", cgPath, err)
	}

	// 1. Memory limit
	memFile := filepath.Join(cgPath, "memory.max")
	if limits.MemoryMaxBytes > 0 {
		_ = os.WriteFile(memFile, []byte(strconv.FormatInt(limits.MemoryMaxBytes, 10)), 0644)
	} else {
		_ = os.WriteFile(memFile, []byte("max"), 0644)
	}

	// 2. CPU quota limit
	cpuFile := filepath.Join(cgPath, "cpu.max")
	if limits.CPUQuotaUS > 0 && limits.CPUPeriodUS > 0 {
		val := fmt.Sprintf("%d %d", limits.CPUQuotaUS, limits.CPUPeriodUS)
		_ = os.WriteFile(cpuFile, []byte(val), 0644)
	} else {
		val := fmt.Sprintf("max %d", limits.CPUPeriodUS)
		_ = os.WriteFile(cpuFile, []byte(val), 0644)
	}

	// 3. PIDs limit
	pidsFile := filepath.Join(cgPath, "pids.max")
	if limits.PIDsMax > 0 {
		_ = os.WriteFile(pidsFile, []byte(strconv.FormatInt(limits.PIDsMax, 10)), 0644)
	} else {
		_ = os.WriteFile(pidsFile, []byte("max"), 0644)
	}

	return nil
}

func (c *LinuxCgroupController) AttachPID(appID string, pid int) error {
	if !c.available || pid <= 0 {
		return nil
	}

	cgPath := c.cgroupPath(appID)
	procsFile := filepath.Join(cgPath, "cgroup.procs")
	if err := os.WriteFile(procsFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to attach PID %d to cgroup %s: %w", pid, cgPath, err)
	}
	return nil
}

func (c *LinuxCgroupController) GetStats(appID string) (*CgroupStats, error) {
	if !c.available {
		return &CgroupStats{}, nil
	}

	cgPath := c.cgroupPath(appID)
	stats := &CgroupStats{}

	// Read memory.current
	if data, err := os.ReadFile(filepath.Join(cgPath, "memory.current")); err == nil {
		val, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		stats.MemoryCurrentBytes = val
	}

	// Read pids.current
	if data, err := os.ReadFile(filepath.Join(cgPath, "pids.current")); err == nil {
		val, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		stats.PIDsCurrent = val
	}

	// Read cpu.stat
	if f, err := os.Open(filepath.Join(cgPath, "cpu.stat")); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) == 2 && fields[0] == "usage_usec" {
				val, _ := strconv.ParseInt(fields[1], 10, 64)
				stats.CPUUsageUSec = val
				break
			}
		}
	}

	return stats, nil
}

func (c *LinuxCgroupController) Cleanup(appID string) error {
	if !c.available {
		return nil
	}
	cgPath := c.cgroupPath(appID)
	return os.Remove(cgPath)
}
