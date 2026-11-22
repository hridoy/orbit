package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"orbit/internal/compute"
)

func RunLogs(client *Client, appID string, follow bool) error {
	if appID == "" {
		return fmt.Errorf("app ID is required")
	}

	if !follow {
		// Non-follow tail
		var resp struct {
			AppID string            `json:"app_id"`
			Lines []compute.LogLine `json:"lines"`
		}
		err := client.GetJSON(fmt.Sprintf("/api/v1/apps/%s/logs?tail=100", appID), &resp)
		if err != nil {
			return err
		}

		for _, l := range resp.Lines {
			fmt.Printf("[%s] %s\n", l.Timestamp.Format("15:04:05"), l.Message)
		}
		return nil
	}

	// Follow live SSE stream
	fmt.Printf("📡 Streaming live logs for '%s' (Ctrl+C to stop)...\n", appID)
	resp, err := client.Do("GET", fmt.Sprintf("/api/v1/apps/%s/logs/stream", appID), nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			payload := strings.TrimPrefix(line, "data: ")
			var logLine compute.LogLine
			if err := json.Unmarshal([]byte(payload), &logLine); err == nil {
				ts := logLine.Timestamp.Format("15:04:05")
				if logLine.Timestamp.IsZero() {
					ts = time.Now().Format("15:04:05")
				}
				fmt.Printf("[%s] %s\n", ts, logLine.Message)
			}
		}
	}

	return scanner.Err()
}
