# 🪐 Orbit — 100% Sovereign, Zero-Docker, Single-Binary PaaS

**Orbit** is a 100% Sovereign, Zero-Docker, Anti-Bloat, Single-Binary Platform as a Service (PaaS) engineered in Go to run with **<25MB idle RAM** and scale from a single cheap VPS ($3-5/mo) to multi-region bare-metal clusters.

---

## ⚡ Core Tenets & Non-Negotiable Rules

1. **ZERO DOCKER:** Zero Docker, Docker Compose, or Kubernetes overhead. Process isolation and resource quotas are enforced natively via Linux `cgroups v2`, kernel `namespaces` (`CLONE_NEWPID`, `CLONE_NEWNS`, `CLONE_NEWUTS`, `CLONE_NEWIPC`), and system user privilege separation.
2. **ZERO THIRD-PARTY PAID TOOLS / BLOAT:** No MinIO, no Prometheus, no Grafana, no external SaaS. All S3 object storage, metrics, routing, logging, and state are built natively inside a single static Go binary.
3. **SUB-25MB IDLE FOOTPRINT:** Verified runtime idle memory usage of **~1.1 MB** (Go runtime + embedded SQLite + S3 Engine + Reverse Proxy + Web Dashboard).
4. **PURE STATIC BINARIES:** Compiles with `CGO_ENABLED=0` using pure Go SQLite (`modernc.org/sqlite`) into standalone single binaries (`orbit-agent` and `orbit` CLI).
5. **EMBEDDED MOBILE PWA:** 100% native CSS (Zero Tailwind, Zero Bootstrap) Vue 3 SPA built into a mobile-first installable PWA and embedded directly into the binary via `go:embed`.

---

## 🏗️ Repository Architecture

```text
orbit/
├── cmd/
│   ├── orbit-agent/             # Sovereign node daemon (runs under systemd)
│   └── orbit/                   # Developer & operator CLI (push, status, logs, db)
├── internal/
│   ├── compute/                 # cgroups v2, namespaces, delta deployer, port swapper
│   ├── database/                # Dedicated native PostgreSQL & MongoDB micro-RAM orchestrator
│   ├── db/                      # Embedded SQLite (WAL mode, ~8MB cache, DDL migrations)
│   ├── jail/                    # Jailed File Browser (Linux openat2 RESOLVE_BENEATH)
│   ├── proxy/                   # Native reverse proxy, ACME Let's Encrypt TLS, bandwidth meter
│   ├── s3/                      # Sovereign S3-compatible engine (AWS SigV4, Multipart, Range GETs)
│   ├── security/                # AES-256-GCM vault, Ed25519 token manager, auth middleware
│   ├── telemetry/               # Low-overhead procfs / sysfs scraper (<0.1% CPU) & SSE broadcaster
│   ├── server/                  # HTTP API multiplexer & recovery middleware
│   └── cli/                     # CLI command implementations
├── ui/                          # Pure Vue 3 PWA (embedded into agent binary via embed.FS)
├── scripts/
│   └── bootstrap.sh             # Host hardening, nftables firewall, node ID provisioning
└── Makefile                     # Static build & test automation
```

---

## 🚀 Quickstart

### 1. Build Binaries
```bash
make all
# Outputs statically linked binaries:
#   bin/orbit-agent (11MB - contains complete PaaS + UI)
#   bin/orbit       (9.7MB - sovereign CLI)
```

### 2. Run All Unit & Integration Tests
```bash
make test
# Runs comprehensive test suite across compute, database, s3, jail, proxy, and security with -race
```

### 3. Bootstrap Host Node (Bare Metal / VPS)
```bash
sudo ./scripts/bootstrap.sh
```
The bootstrap script:
- Provisions `/opt/orbit/{apps,databases,storage,secrets,data,run}` (0700 permissions).
- Hardens Linux kernel sysctl (`net.ipv4.tcp_syncookies = 1`, `rp_filter = 1`).
- Configures `nftables` firewall (drops all inbound ports except 22, 80, 443, WireGuard 51820; zero WAN admin ports).
- Generates 256-bit AES-GCM Master Key (`/opt/orbit/secrets/master.key`).
- Assigns cryptographic Node ID and installs `/etc/systemd/system/orbit-agent.service`.

---

## 🛠️ Developer CLI (`orbit`)

### Cluster & Host Status
```bash
orbit status
```
```text
================================================================================
🪐 ORBIT CLUSTER STATUS: Node node_b6ea0459 (vps-sfo-01)
================================================================================
Status:       online
IP Address:   127.0.0.1
CPU Cores:    8 (Host Usage: 2.1%)
Memory:       1.1 MB / 3.8 GB (Agent Alloc: 1.1 MB)
--------------------------------------------------------------------------------
📦 WORKLOADS (1 Active):
   • api-service      | Port: 8001 (Idle: 8002) | RAM: 256MB | Status: healthy
--------------------------------------------------------------------------------
🗄️  DEDICATED DATABASES (1 Running):
   • db_postgres_8c3  | postgres | Port: 5433 | Profile: micro | IPC Socket: /opt/orbit/...
================================================================================
```

### Wrangler-Style Incremental Delta Deploy
Uploads only changed files or binaries using content-addressable SHA256 diffing, verifies candidate health on `idle_port`, and executes an atomic zero-downtime port swap:
```bash
orbit push ./my-app --app=api-service
```

### Dedicated Native Database Provisioning
Spins up an isolated, dedicated native PostgreSQL or MongoDB instance with micro low-RAM tuning (`<32MB RAM` target, 20 max connections, socket & loopback binding):
```bash
orbit db create my_postgres --engine=postgres --profile=micro
```
Outputs:
- **Loopback TCP URI:** `postgres://orbit_admin:<pass>@127.0.0.1:5433/my_postgres?sslmode=disable`
- **Zero-Network IPC Socket:** `postgres://orbit_admin:<pass>@/my_postgres?host=/opt/orbit/databases/<id>/run`

### Real-Time Live Log Streaming (SSE)
```bash
orbit logs api-service -f
```

---

## 🔒 Security Architecture

- **Path-Traversal Proof Jailed File Browser:** Dual-layer protection using Linux `openat2` syscall with `RESOLVE_BENEATH | RESOLVE_NO_MAGICLINKS`, making TOCTOU symlink race attacks mathematically impossible.
- **In-Memory Encrypted Vault:** Secrets are encrypted on disk via AES-256-GCM. Plaintext is never stored in `.env` files; secrets are decrypted and injected directly into child process memory at runtime.
- **Dedicated Non-Root OS Users:** Applications run under dedicated system users (`orbit-app-<id>`).
- **Internal Only Databases:** Database engines never bind to `0.0.0.0`; all instances bind strictly to `127.0.0.1` and local Unix domain sockets for zero-latency IPC.
- **Encrypted Off-Host Backups:** Backups include `master.key`, so the cloud copy must be encrypted. Create an rclone `crypt` remote that wraps your Drive, Dropbox, OneDrive, or S3 remote (`rclone config`, type `crypt`, with a strong password), then use that crypt remote name as the Cloud backup remote path (for example `secure:orbit-backups`). Store the crypt password safely; without it the cloud copies cannot be restored.
- **Token on the CLI:** Pass the token with `ORBIT_TOKEN`, not `--token`. The flag is kept for compatibility, but it exposes the token in the process list and shell history.
