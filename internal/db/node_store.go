package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Node struct {
	NodeID          string    `json:"node_id"`
	Hostname        string    `json:"hostname"`
	PublicIP        string    `json:"public_ip"`
	WireguardIP     string    `json:"wireguard_ip,omitempty"`
	WireguardPubkey string    `json:"wireguard_pubkey,omitempty"`
	Role            string    `json:"role"`
	CPUCores        int       `json:"cpu_cores"`
	RAMTotalBytes   int64     `json:"ram_total_bytes"`
	DiskTotalBytes  int64     `json:"disk_total_bytes"`
	AgentVersion    string    `json:"agent_version"`
	Status          string    `json:"status"`
	LastHeartbeatAt time.Time `json:"last_heartbeat_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type NodeStore struct {
	db *DB
}

func NewNodeStore(db *DB) *NodeStore {
	return &NodeStore{db: db}
}

func (s *NodeStore) UpsertNode(n *Node) error {
	query := `
	INSERT INTO nodes (
		node_id, hostname, public_ip, wireguard_ip, wireguard_pubkey,
		role, cpu_cores, ram_total_bytes, disk_total_bytes,
		agent_version, status, last_heartbeat_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(node_id) DO UPDATE SET
		hostname = excluded.hostname,
		public_ip = excluded.public_ip,
		wireguard_ip = excluded.wireguard_ip,
		wireguard_pubkey = excluded.wireguard_pubkey,
		cpu_cores = excluded.cpu_cores,
		ram_total_bytes = excluded.ram_total_bytes,
		disk_total_bytes = excluded.disk_total_bytes,
		agent_version = excluded.agent_version,
		status = excluded.status,
		last_heartbeat_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err := s.db.Exec(query,
		n.NodeID, n.Hostname, n.PublicIP, n.WireguardIP, n.WireguardPubkey,
		n.Role, n.CPUCores, n.RAMTotalBytes, n.DiskTotalBytes,
		n.AgentVersion, n.Status,
	)
	if err != nil {
		return fmt.Errorf("upsert node %s: %w", n.NodeID, err)
	}
	return nil
}

func (s *NodeStore) GetNode(nodeID string) (*Node, error) {
	query := `
	SELECT node_id, hostname, public_ip, COALESCE(wireguard_ip, ''), COALESCE(wireguard_pubkey, ''),
	       role, cpu_cores, ram_total_bytes, disk_total_bytes,
	       agent_version, status, last_heartbeat_at, created_at, updated_at
	FROM nodes WHERE node_id = ?;
	`
	row := s.db.QueryRow(query, nodeID)
	var n Node
	err := row.Scan(
		&n.NodeID, &n.Hostname, &n.PublicIP, &n.WireguardIP, &n.WireguardPubkey,
		&n.Role, &n.CPUCores, &n.RAMTotalBytes, &n.DiskTotalBytes,
		&n.AgentVersion, &n.Status, &n.LastHeartbeatAt, &n.CreatedAt, &n.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get node %s: %w", nodeID, err)
	}
	return &n, nil
}
