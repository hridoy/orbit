package database

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Exposure is one public port that forwards to a database on loopback.
type Exposure struct {
	mu      sync.Mutex
	ln      net.Listener
	port    int
	backend string
	conns   map[net.Conn]bool
}

// suggestedBase is where the free-port suggestions start for each engine.
var suggestedBase = map[string]int{
	"postgres": 15432,
	"mysql":    13306,
	"mariadb":  13306,
	"mongodb":  37017,
	"redis":    16379,
	"libsql":   18080,
}

// exposures holds the running forwarders by database ID.
type exposures struct {
	mu  sync.Mutex
	run map[string]*Exposure
}

func newExposures() *exposures {
	return &exposures{run: map[string]*Exposure{}}
}

// Expose starts forwarding publicPort on all interfaces to the database's loopback port.
// It also opens the port in the firewall.
func (m *DatabaseManager) Expose(dbID string, publicPort int) error {
	if publicPort < 1024 || publicPort > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535")
	}
	instance, err := m.dbStore.GetDatabase(dbID)
	if err != nil || instance == nil {
		return fmt.Errorf("database not found")
	}
	if instance.Port == publicPort {
		return fmt.Errorf("port %d is the database's own port", publicPort)
	}
	if err := m.portAvailable(publicPort); err != nil {
		return err
	}

	m.exp.mu.Lock()
	defer m.exp.mu.Unlock()
	if old := m.exp.run[dbID]; old != nil {
		old.close()
		delete(m.exp.run, dbID)
	}
	ex, err := startForwarder(publicPort, fmt.Sprintf("127.0.0.1:%d", instance.Port))
	if err != nil {
		return fmt.Errorf("could not open port %d: %w", publicPort, err)
	}
	m.exp.run[dbID] = ex
	if err := m.dbExposures.Set(dbID, publicPort); err != nil {
		ex.close()
		delete(m.exp.run, dbID)
		return err
	}
	firewallPort(publicPort, true)
	return nil
}

// Unexpose stops the public forwarder and closes the port.
func (m *DatabaseManager) Unexpose(dbID string) error {
	m.exp.mu.Lock()
	defer m.exp.mu.Unlock()
	cur, err := m.dbExposures.Get(dbID)
	if err != nil {
		return err
	}
	if cur != nil {
		firewallPort(cur.PublicPort, false)
	}
	if ex := m.exp.run[dbID]; ex != nil {
		ex.close()
		delete(m.exp.run, dbID)
	}
	return m.dbExposures.Delete(dbID)
}

// RestoreExposures reopens the public ports that were open before the agent restarted.
func (m *DatabaseManager) RestoreExposures() {
	list, err := m.dbExposures.List()
	if err != nil {
		return
	}
	for _, e := range list {
		if err := m.Expose(e.DBID, e.PublicPort); err != nil {
			log.Printf("⚠️ database %s: could not reopen port %d: %v", e.DBID, e.PublicPort, err)
		}
	}
}

// SuggestedPorts returns up to n free public ports for an engine, starting at its base.
func (m *DatabaseManager) SuggestedPorts(engine string, n int) []int {
	base, ok := suggestedBase[engine]
	if !ok {
		return nil
	}
	var out []int
	for p := base; p < base+400 && len(out) < n; p++ {
		if m.portAvailable(p) == nil {
			out = append(out, p)
		}
	}
	return out
}

// portAvailable reports whether a public port is free for a new forwarder.
func (m *DatabaseManager) portAvailable(port int) error {
	if list, err := m.dbStore.ListDatabases(); err == nil {
		for _, d := range list {
			if d.Port == port {
				return fmt.Errorf("port %d belongs to a database", port)
			}
		}
	}
	if list, err := m.dbExposures.List(); err == nil {
		for _, e := range list {
			if e.PublicPort == port {
				return fmt.Errorf("port %d is already open for another database", port)
			}
		}
	}
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return fmt.Errorf("port %d is in use on this server", port)
	}
	_ = ln.Close()
	return nil
}

func startForwarder(port int, backend string) (*Exposure, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}
	ex := &Exposure{ln: ln, port: port, backend: backend, conns: map[net.Conn]bool{}}
	go ex.accept()
	return ex, nil
}

func (ex *Exposure) accept() {
	for {
		c, err := ex.ln.Accept()
		if err != nil {
			return
		}
		go ex.pipe(c)
	}
}

func (ex *Exposure) pipe(client net.Conn) {
	backend, err := net.DialTimeout("tcp", ex.backend, 5*time.Second)
	if err != nil {
		_ = client.Close()
		return
	}
	ex.track(client, true)
	ex.track(backend, true)
	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(backend, client); done <- struct{}{} }()
	go func() { _, _ = io.Copy(client, backend); done <- struct{}{} }()
	<-done
	_ = client.Close()
	_ = backend.Close()
	ex.track(client, false)
	ex.track(backend, false)
}

func (ex *Exposure) track(c net.Conn, add bool) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	if add {
		ex.conns[c] = true
	} else {
		delete(ex.conns, c)
	}
}

func (ex *Exposure) close() {
	_ = ex.ln.Close()
	ex.mu.Lock()
	defer ex.mu.Unlock()
	for c := range ex.conns {
		_ = c.Close()
	}
}

// firewallPort opens or closes a TCP port with ufw when it is installed. Failure is logged,
// not fatal, because the forwarder still works on hosts without ufw.
func firewallPort(port int, open bool) {
	if _, err := exec.LookPath("ufw"); err != nil {
		return
	}
	args := []string{"allow", fmt.Sprintf("%d/tcp", port), "comment", "orbit database"}
	if !open {
		args = []string{"delete", "allow", fmt.Sprintf("%d/tcp", port)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "ufw", args...).CombinedOutput()
	if err != nil {
		log.Printf("⚠️ ufw %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
}
