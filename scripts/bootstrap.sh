#!/usr/bin/env bash
# ==============================================================================
# Orbit Sovereign Node Bootstrap Script
# Hardens host Linux, installs nftables firewall, provisions cgroups v2,
# creates directory structure under /opt/orbit, and initializes systemd service.
# ==============================================================================

set -euo pipefail

if [[ $EUID -ne 0 ]]; then
   echo "❌ Orbit bootstrap must be run as root (or via sudo)." >&2
   exit 1
fi

echo "============================================================"
echo "🪐 Orbit Sovereign Host Bootstrap & Kernel Hardening"
echo "============================================================"

# 1. Directory Provisioning
echo "📁 Initializing /opt/orbit sovereign hierarchy..."
mkdir -p /opt/orbit/{bin,data,apps,databases,storage,secrets,run}
chmod 0700 /opt/orbit
chmod 0700 /opt/orbit/secrets

# 2. Linux Kernel & Host Hardening (sysctl)
echo "🔒 Applying kernel hardening sysctl params..."
cat << 'EOF' > /etc/sysctl.d/99-orbit-hardened.conf
# Orbit Sovereign Hardening
net.ipv4.ip_forward = 1
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.icmp_echo_ignore_broadcasts = 1
net.ipv4.conf.all.accept_redirects = 0
net.ipv6.conf.all.accept_redirects = 0
net.ipv4.conf.all.send_redirects = 0
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_max_syn_backlog = 4096
net.ipv4.tcp_fin_timeout = 15
fs.inotify.max_user_watches = 524288
fs.file-max = 2097152
EOF
sysctl --system > /dev/null 2>&1 || true

# 3. nftables Sovereign Firewall Configuration
echo "🛡️ Configuring nftables firewall..."
if command -v apt-get &> /dev/null; then
    apt-get update -qq && apt-get install -y -qq nftables curl > /dev/null 2>&1 || true
elif command -v dnf &> /dev/null; then
    dnf install -y -q nftables curl > /dev/null 2>&1 || true
fi

cat << 'EOF' > /etc/nftables.conf
#!/usr/sbin/nft -f

flush ruleset

table inet orbit_firewall {
    chain inbound {
        type filter hook input priority 0; policy drop;

        # Allow localhost loopback
        iif "lo" accept

        # Allow established and related connections
        ct state { established, related } accept
        ct state invalid drop

        # ICMP / Ping rate limited
        ip protocol icmp icmp type echo-request limit rate 5/second accept
        ip6 nexthdr ipv6-icmp icmpv6 type echo-request limit rate 5/second accept

        # Public Services: SSH, HTTP, HTTPS
        tcp dport 22 accept
        tcp dport { 80, 443 } accept

        # Sovereign WireGuard mesh
        udp dport 51820 accept

        # Reject all other inbound packets cleanly
        reject with icmpx type port-unreachable
    }

    chain forward {
        type filter hook forward priority 0; policy drop;
    }

    chain outbound {
        type filter hook output priority 0; policy accept;
    }
}
EOF

if command -v nft &> /dev/null; then
    nft -f /etc/nftables.conf || true
    systemctl enable --now nftables > /dev/null 2>&1 || true
fi

# 4. Master Key Generation (AES-256 / Ed25519)
if [[ ! -f /opt/orbit/secrets/master.key ]]; then
    echo "🔑 Generating 256-bit cryptographically secure master key..."
    head -c 32 /dev/urandom > /opt/orbit/secrets/master.key
    chmod 0600 /opt/orbit/secrets/master.key
fi

# 5. Node ID Assignment (Final Step)
if [[ ! -f /opt/orbit/data/node_id ]]; then
    NODE_UUID="node_$(od -x /dev/urandom | head -1 | awk '{print $2$3$4$5}')"
    echo "$NODE_UUID" > /opt/orbit/data/node_id
    chmod 0600 /opt/orbit/data/node_id
fi

NODE_ID=$(cat /opt/orbit/data/node_id)
MASTER_KEY_HEX=$(od -An -tx1 /opt/orbit/secrets/master.key | tr -d ' \n')

# 6. systemd Service Unit
echo "⚙️ Creating /etc/systemd/system/orbit-agent.service..."
cat << EOF > /etc/systemd/system/orbit-agent.service
[Unit]
Description=Orbit Sovereign PaaS Agent
After=network.target network-online.target nftables.service
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/orbit
ExecStart=/opt/orbit/bin/orbit-agent
Restart=always
RestartSec=3s
LimitNOFILE=65536
LimitNPROC=32768
MemoryMax=128M
CPUQuota=100%

# Sandboxing for agent itself
ProtectSystem=full
ProtectHome=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload > /dev/null 2>&1 || true

echo "============================================================"
echo "✅ Orbit Sovereign Node Bootstrap Complete!"
echo "   Node ID:    $NODE_ID"
echo "   Master Key: $MASTER_KEY_HEX"
echo "   Root Path:  /opt/orbit"
echo "============================================================"
