package cli

import (
	"fmt"
	"time"

	"orbit/internal/db"
	"orbit/internal/telemetry"
)

func RunStatus(client *Client) error {
	// 1. Get Node Info
	var node db.Node
	if err := client.GetJSON("/api/v1/node", &node); err != nil {
		return fmt.Errorf("failed to fetch node info: %w", err)
	}

	// 2. Get Telemetry Metrics
	var metrics telemetry.HostMetrics
	_ = client.GetJSON("/api/v1/telemetry/metrics", &metrics)

	// 3. Get Apps List
	var apps []*db.App
	_ = client.GetJSON("/api/v1/apps", &apps)

	// 4. Get Databases List
	var databases []*db.DatabaseInstance
	_ = client.GetJSON("/api/v1/databases", &databases)

	fmt.Println("================================================================================")
	fmt.Printf("🪐 ORBIT CLUSTER STATUS: Node %s (%s)\n", node.NodeID, node.Hostname)
	fmt.Println("================================================================================")
	fmt.Printf("Status:       %s\n", stringsToUpper(node.Status))
	fmt.Printf("IP Address:   %s\n", node.PublicIP)
	fmt.Printf("CPU Cores:    %d (Host Usage: %.1f%%)\n", node.CPUCores, metrics.CPUUsagePct)
	fmt.Printf("Memory:       %s / %s (Agent Alloc: %s)\n",
		formatBytes(metrics.RAMUsedBytes),
		formatBytes(metrics.RAMTotalBytes),
		formatBytes(metrics.AgentAllocRAM),
	)
	fmt.Printf("Last Seen:    %s\n", node.LastHeartbeatAt.Format(time.RFC3339))
	fmt.Println("--------------------------------------------------------------------------------")

	fmt.Printf("📦 WORKLOADS (%d Active):\n", len(apps))
	if len(apps) == 0 {
		fmt.Println("   (No applications deployed. Run 'orbit push' to launch your first service.)")
	} else {
		for _, a := range apps {
			domain := "(none)"
			if a.Domain != nil && *a.Domain != "" {
				domain = *a.Domain
			}
			version := "initial"
			if a.CurrentVersion != nil {
				version = *a.CurrentVersion
			}
			fmt.Printf("   • %-16s | Port: %d (Idle: %d) | RAM: %dMB | Status: %-8s | Domain: %s | Ver: %s\n",
				a.AppID, a.ActivePort, a.IdlePort, a.RAMLimitMB, a.Status, domain, version)
		}
	}
	fmt.Println("--------------------------------------------------------------------------------")

	fmt.Printf("🗄️  DEDICATED DATABASES (%d Running):\n", len(databases))
	if len(databases) == 0 {
		fmt.Println("   (No dedicated databases provisioned. Run 'orbit db create' to spin one up.)")
	} else {
		for _, d := range databases {
			fmt.Printf("   • %-16s | %-8s | Port: %d | Profile: %-8s | Status: %-8s | IPC Socket: %s\n",
				d.DBID, d.Engine, d.Port, d.RAMProfile, d.Status, d.SocketPath)
		}
	}
	fmt.Println("================================================================================")
	return nil
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func stringsToUpper(s string) string {
	return fmt.Sprintf("%s", s)
}
