package telemetry

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"orbit/internal/proxy"
	"orbit/internal/sysinfo"
)

type HostMetrics struct {
	Timestamp     time.Time `json:"timestamp"`
	CPUUsagePct   float64   `json:"cpu_usage_pct"`
	RAMTotalBytes uint64    `json:"ram_total_bytes"`
	RAMUsedBytes  uint64    `json:"ram_used_bytes"`
	RAMFreeBytes  uint64    `json:"ram_free_bytes"`
	// OrbitRSSBytes is the memory used by the agent and everything it started (apps, databases).
	// RAMUsedBytes is the whole server, which includes other software on the same machine.
	OrbitRSSBytes   uint64 `json:"orbit_rss_bytes"`
	AgentAllocRAM   uint64 `json:"agent_alloc_ram"`
	AgentNumGC      uint32 `json:"agent_num_gc"`
	NumGoroutines   int    `json:"num_goroutines"`
	BandwidthInSec  uint64 `json:"bandwidth_in_sec"`
	BandwidthOutSec uint64 `json:"bandwidth_out_sec"`
}

type TelemetryCollector struct {
	meter        *proxy.BandwidthMeter
	lastCPUTotal uint64
	lastCPUIdle  uint64
}

func NewTelemetryCollector(meter *proxy.BandwidthMeter) *TelemetryCollector {
	return &TelemetryCollector{meter: meter}
}

// Collect returns a low-overhead snapshot of system and process metrics.
func (c *TelemetryCollector) Collect() HostMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ramTotal, ramUsed, ramFree := readLinuxMemory()
	if ramTotal == 0 {
		// Non-Linux or container fallback
		ramTotal = m.Sys
		ramUsed = m.Alloc
		ramFree = m.Sys - m.Alloc
	}

	cpuPct := c.readLinuxCPU()

	orbitRSS := sysinfo.ProcessTreeRSS(os.Getpid())
	if orbitRSS == 0 {
		// Non-Linux: fall back to what the Go runtime has allocated.
		orbitRSS = m.Sys
	}

	return HostMetrics{
		Timestamp:      time.Now().UTC(),
		CPUUsagePct:    cpuPct,
		RAMTotalBytes:  ramTotal,
		RAMUsedBytes:   ramUsed,
		RAMFreeBytes:   ramFree,
		OrbitRSSBytes:  orbitRSS,
		AgentAllocRAM:  m.Alloc,
		AgentNumGC:     m.NumGC,
		NumGoroutines:  runtime.NumGoroutine(),
		BandwidthInSec: 0,
	}
}

func readLinuxMemory() (total, used, free uint64) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0
	}
	defer f.Close()

	var memTotal, memFree, memAvail uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			val, _ := strconv.ParseUint(fields[1], 10, 64)
			memTotal = val * 1024
		case "MemFree:":
			val, _ := strconv.ParseUint(fields[1], 10, 64)
			memFree = val * 1024
		case "MemAvailable:":
			val, _ := strconv.ParseUint(fields[1], 10, 64)
			memAvail = val * 1024
		}
	}

	if memTotal > 0 {
		if memAvail > 0 {
			used = memTotal - memAvail
		} else {
			used = memTotal - memFree
		}
		return memTotal, used, memAvail
	}
	return 0, 0, 0
}

func (c *TelemetryCollector) readLinuxCPU() float64 {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0.0 // Non-Linux fallback
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return 0.0
	}

	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0.0
	}

	var user, nice, system, idle, iowait, irq, softirq uint64
	user, _ = strconv.ParseUint(fields[1], 10, 64)
	nice, _ = strconv.ParseUint(fields[2], 10, 64)
	system, _ = strconv.ParseUint(fields[3], 10, 64)
	idle, _ = strconv.ParseUint(fields[4], 10, 64)
	if len(fields) > 5 {
		iowait, _ = strconv.ParseUint(fields[5], 10, 64)
	}
	if len(fields) > 6 {
		irq, _ = strconv.ParseUint(fields[6], 10, 64)
	}
	if len(fields) > 7 {
		softirq, _ = strconv.ParseUint(fields[7], 10, 64)
	}

	idleTime := idle + iowait
	totalTime := user + nice + system + idle + iowait + irq + softirq

	deltaTotal := totalTime - c.lastCPUTotal
	deltaIdle := idleTime - c.lastCPUIdle

	c.lastCPUTotal = totalTime
	c.lastCPUIdle = idleTime

	if deltaTotal == 0 {
		return 0.0
	}

	usage := float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100.0
	if usage < 0.0 {
		usage = 0.0
	}
	if usage > 100.0 {
		usage = 100.0
	}
	return usage
}
