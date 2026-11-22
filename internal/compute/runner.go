package compute

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
)

type ProcessInstance struct {
	AppID      string
	Port       int
	PID        int
	Cmd        *exec.Cmd
	RingBuffer *RingBuffer
	Cancel     context.CancelFunc
	Exited     chan struct{}
	ExitErr    error
}

type ProcessSupervisor struct {
	mu        sync.RWMutex
	cfg       *config.Config
	cgroups   CgroupController
	processes map[string]*ProcessInstance // Key: appID
	buffers   map[string]*RingBuffer     // Key: appID
}

func NewProcessSupervisor(cfg *config.Config, cgroups CgroupController) *ProcessSupervisor {
	return &ProcessSupervisor{
		cfg:       cfg,
		cgroups:   cgroups,
		processes: make(map[string]*ProcessInstance),
		buffers:   make(map[string]*RingBuffer),
	}
}

// GetRingBuffer returns or initializes the log ring buffer for an app.
func (s *ProcessSupervisor) GetRingBuffer(appID string) *RingBuffer {
	s.mu.Lock()
	defer s.mu.Unlock()

	rb, exists := s.buffers[appID]
	if !exists {
		rb = NewRingBuffer(1000)
		s.buffers[appID] = rb
	}
	return rb
}

// StartProcess launches an application instance on the designated port.
func (s *ProcessSupervisor) StartProcess(
	ctx context.Context,
	app *db.App,
	port int,
	workDir string,
	envVars map[string]string,
) (*ProcessInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rb := s.GetRingBuffer(app.AppID)

	childCtx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(childCtx, app.BinaryEntrypoint)
	cmd.Dir = workDir

	// In-memory environment injection (Never written to disk!)
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("PORT=%d", port),
		fmt.Sprintf("ORBIT_APP_ID=%s", app.AppID),
		fmt.Sprintf("ORBIT_NODE_ID=%s", s.cfg.NodeID),
		"ORBIT_ENV=production",
	)
	for k, v := range envVars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	// Route stdout & stderr directly into the in-memory circular ring buffer
	cmd.Stdout = rb
	cmd.Stderr = rb

	// Apply platform-specific isolation (Linux cgroups + namespaces)
	configureSysProcAttr(cmd, app)

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start binary %s: %w", app.BinaryEntrypoint, err)
	}

	pid := cmd.Process.Pid
	rb.Append(fmt.Sprintf("[orbit] Process started (PID: %d, Port: %d)", pid, port))

	// Attach to cgroups v2
	_ = s.cgroups.Setup(app.AppID, FormatCgroupLimits(app.RAMLimitMB, app.CPUQuotaPct, app.PIDsLimit))
	_ = s.cgroups.AttachPID(app.AppID, pid)

	instance := &ProcessInstance{
		AppID:      app.AppID,
		Port:       port,
		PID:        pid,
		Cmd:        cmd,
		RingBuffer: rb,
		Cancel:     cancel,
		Exited:     make(chan struct{}),
	}

	// Monitor exit in background
	go func() {
		err := cmd.Wait()
		instance.ExitErr = err
		close(instance.Exited)
		rb.Append(fmt.Sprintf("[orbit] Process exited (PID: %d, Err: %v)", pid, err))
	}()

	return instance, nil
}

// StopProcess sends SIGTERM and waits up to gracePeriod before SIGKILL.
func (s *ProcessSupervisor) StopProcess(instance *ProcessInstance, gracePeriod time.Duration) error {
	if instance == nil || instance.Cmd == nil || instance.Cmd.Process == nil {
		return nil
	}

	select {
	case <-instance.Exited:
		return nil // Already exited
	default:
	}

	// Send SIGTERM
	_ = instance.Cmd.Process.Signal(syscall.SIGTERM)

	select {
	case <-instance.Exited:
		return nil
	case <-time.After(gracePeriod):
		// Force kill if lingering past grace period
		_ = instance.Cmd.Process.Kill()
		instance.Cancel()
		<-instance.Exited
		return nil
	}
}

// RegisterActive registers an instance as the active serving process.
func (s *ProcessSupervisor) RegisterActive(appID string, instance *ProcessInstance) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processes[appID] = instance
}

// GetActive returns the active instance for an app.
func (s *ProcessSupervisor) GetActive(appID string) *ProcessInstance {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.processes[appID]
}
