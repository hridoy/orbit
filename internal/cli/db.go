package cli

import (
	"fmt"

	"orbit/internal/database"
)

func RunDBCreate(client *Client, name, engine, profile string) error {
	if name == "" {
		return fmt.Errorf("database name is required")
	}
	if engine == "" {
		engine = "postgres"
	}
	if profile == "" {
		profile = string(database.ProfileMicro)
	}

	payload := map[string]string{
		"name":        name,
		"engine":      engine,
		"ram_profile": profile,
	}

	var resp struct {
		Database  map[string]any `json:"database"`
		Password  string         `json:"password"`
		TCPURI    string         `json:"tcp_uri"`
		SocketURI string         `json:"socket_uri"`
	}

	fmt.Printf("📦 Provisioning isolated native %s instance ('%s') with profile '%s'...\n", engine, name, profile)
	err := client.PostJSON("/api/v1/databases", payload, &resp)
	if err != nil {
		return fmt.Errorf("failed to provision database: %w", err)
	}

	dbID := resp.Database["db_id"]
	port := resp.Database["port"]

	fmt.Println("✅ Dedicated Database Provisioned Successfully!")
	fmt.Printf("   DB ID:       %v\n", dbID)
	fmt.Printf("   Engine:      %s\n", engine)
	fmt.Printf("   RAM Profile: %s (<32MB low-RAM tuned)\n", profile)
	fmt.Printf("   Port:        %v (bound strictly to 127.0.0.1)\n", port)
	fmt.Printf("   Admin User:  orbit_admin\n")
	fmt.Printf("   Password:    %s (AES-256-GCM encrypted in Orbit Vault)\n", resp.Password)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("🔗 Connection Strings:")
	fmt.Printf("   • TCP URI (Loopback):   %s\n", resp.TCPURI)
	fmt.Printf("   • Socket URI (IPC):     %s\n", resp.SocketURI)
	fmt.Println("================================================================================")
	return nil
}
