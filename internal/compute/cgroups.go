package compute

import (
	"fmt"
)

// CgroupLimits configures Linux cgroups v2 resource limits.
type CgroupLimits struct {
	MemoryMaxBytes int64 `json:"memory_max_bytes"` // memory.max
	CPUQuotaUS     int64 `json:"cpu_quota_us"`      // cpu.max quota (e.g. 100000)
	CPUPeriodUS    int64 `json:"cpu_period_us"`     // cpu.max period (e.g. 100000)
	PIDsMax        int64 `json:"pids_max"`          // pids.max
}

// CgroupStats contains real-time resource utilization from cgroups v2.
type CgroupStats struct {
	MemoryCurrentBytes int64 `json:"memory_current_bytes"`
	PIDsCurrent        int64 `json:"pids_current"`
	CPUUsageUSec       int64 `json:"cpu_usage_usec"`
}

// CgroupController handles creation, limit setting, and process attachment in cgroups v2.
type CgroupController interface {
	Setup(appID string, limits CgroupLimits) error
	AttachPID(appID string, pid int) error
	GetStats(appID string) (*CgroupStats, error)
	Cleanup(appID string) error
}

// NewCgroupController creates a platform-appropriate cgroup controller.
func NewCgroupController() CgroupController {
	return newPlatformCgroupController()
}

// FormatCPULimit converts cpu_quota_pct (e.g. 100 = 1 core, 50 = 0.5 core, 200 = 2 cores)
// into cgroups v2 quota/period (period default: 100000us = 100ms).
func FormatCPULimit(cpuQuotaPct int) (quotaUS, periodUS int64) {
	periodUS = 100000 // 100ms default
	if cpuQuotaPct <= 0 {
		return -1, periodUS // max (no limit)
	}
	quotaUS = int64(cpuQuotaPct) * (periodUS / 100)
	return quotaUS, periodUS
}

// FormatMemoryLimit converts megabytes to bytes.
func FormatMemoryLimit(ramLimitMB int) int64 {
	if ramLimitMB <= 0 {
		return -1 // max
	}
	return int64(ramLimitMB) * 1024 * 1024
}

// FormatPIDsLimit returns max pids.
func FormatPIDsLimit(pids int) int64 {
	if pids <= 0 {
		return 128
	}
	return int64(pids)
}

func FormatCgroupLimits(ramMB, cpuPct, pids int) CgroupLimits {
	quota, period := FormatCPULimit(cpuPct)
	return CgroupLimits{
		MemoryMaxBytes: FormatMemoryLimit(ramMB),
		CPUQuotaUS:     quota,
		CPUPeriodUS:    period,
		PIDsMax:        FormatPIDsLimit(pids),
	}
}

var _ = fmt.Sprintf
