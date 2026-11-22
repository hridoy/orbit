package compute

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// CheckHealth polls the candidate port until it responds or context times out.
func CheckHealth(ctx context.Context, port int, path string, timeoutSec int) error {
	if timeoutSec <= 0 {
		timeoutSec = 10
	}

	deadlineCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	targetURL := fmt.Sprintf("http://127.0.0.1:%d%s", port, path)
	if path == "" {
		targetURL = fmt.Sprintf("http://127.0.0.1:%d/", port)
	}

	for {
		select {
		case <-deadlineCtx.Done():
			return fmt.Errorf("health check timed out after %ds for %s", timeoutSec, targetURL)
		case <-ticker.C:
			// Try HTTP GET first
			req, err := http.NewRequestWithContext(deadlineCtx, "GET", targetURL, nil)
			if err == nil {
				resp, err := client.Do(req)
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode < 500 {
						// 2xx, 3xx, or 4xx indicates application is online and bound
						return nil
					}
				}
			}

			// Fallback: Test raw TCP dial if HTTP path fails
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 500*time.Millisecond)
			if err == nil {
				conn.Close()
				// If no specific health path is strictly required, socket listening is considered healthy
				if path == "" || path == "/" {
					return nil
				}
			}
		}
	}
}
