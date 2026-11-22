//go:build !linux

package compute

import "sync"

type GenericCgroupController struct {
	mu     sync.RWMutex
	limits map[string]CgroupLimits
	pids   map[string]int
}

func newPlatformCgroupController() CgroupController {
	return &GenericCgroupController{
		limits: make(map[string]CgroupLimits),
		pids:   make(map[string]int),
	}
}

func (c *GenericCgroupController) Setup(appID string, limits CgroupLimits) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.limits[appID] = limits
	return nil
}

func (c *GenericCgroupController) AttachPID(appID string, pid int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pids[appID] = pid
	return nil
}

func (c *GenericCgroupController) GetStats(appID string) (*CgroupStats, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return &CgroupStats{
		MemoryCurrentBytes: 16 * 1024 * 1024, // Mock 16MB
		PIDsCurrent:        1,
		CPUUsageUSec:       1000,
	}, nil
}

func (c *GenericCgroupController) Cleanup(appID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.limits, appID)
	delete(c.pids, appID)
	return nil
}
