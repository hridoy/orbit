package database

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

// The forwarder must carry bytes both ways between a public port and the loopback backend.
func TestForwarderEchoesThroughPublicPort(t *testing.T) {
	backend, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()
	go func() {
		for {
			c, err := backend.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = io.Copy(c, c)
			}(c)
		}
	}()

	public := freePort(t)
	ex, err := startForwarder(public, backend.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer ex.close()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", public), 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	msg := []byte("ping-through-forwarder")
	if _, err := conn.Write(msg); err != nil {
		t.Fatal(err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	got := make([]byte, len(msg))
	if _, err := io.ReadFull(conn, got); err != nil {
		t.Fatal(err)
	}
	if string(got) != string(msg) {
		t.Fatalf("got %q", got)
	}
}

func TestForwarderClosePortIsFreedAfterClose(t *testing.T) {
	public := freePort(t)
	ex, err := startForwarder(public, "127.0.0.1:1")
	if err != nil {
		t.Fatal(err)
	}
	ex.close()
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", public))
	if err != nil {
		t.Fatalf("port not released: %v", err)
	}
	_ = ln.Close()
}
