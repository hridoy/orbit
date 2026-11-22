#!/usr/bin/env bash
# ==============================================================================
# Orbit Sovereign PaaS — One-Click Local Development Runner
# Runs Orbit in zero-docker local mode with embedded Cloudflare dashboard UI.
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PORT="${ORBIT_PORT:-8787}"
DEV_ROOT="$SCRIPT_DIR/.orbit-local"

echo "========================================================================"
echo "🪐  ORBIT SOVEREIGN PaaS — LOCAL DEV ENVIRONMENT"
echo "    100% Sovereign | Zero Docker | Cloudflare-Style Dashboard"
echo "========================================================================"

# 1. Check Prerequisites
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed." >&2; exit 1; }

# 2. Build Frontend UI if dist is missing
if [ ! -d "ui/dist" ] || [ ! -f "ui/dist/index.html" ]; then
    echo "📦 Building Cloudflare-style Vue 3 PWA dashboard..."
    (cd ui && npm install --silent && npm run build)
fi

# 3. Build Static Single Binaries
echo "⚡ Compiling static orbit-agent and orbit CLI..."
mkdir -p bin
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/orbit-agent ./cmd/orbit-agent
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/orbit ./cmd/orbit

# 4. Prepare Local Data Directory (No root required)
mkdir -p "$DEV_ROOT"/{apps,databases,storage,secrets,data,run}

export ORBIT_ROOT="$DEV_ROOT"
export ORBIT_LISTEN_ADDR="127.0.0.1:${PORT}"
export ORBIT_SERVER="http://127.0.0.1:${PORT}"

# Trap SIGINT / SIGTERM for clean teardown
cleanup() {
    echo ""
    echo "🛑 Stopping Orbit Agent..."
    if [ -n "${AGENT_PID:-}" ]; then
        kill -TERM "$AGENT_PID" 2>/dev/null || true
        wait "$AGENT_PID" 2>/dev/null || true
    fi
    echo "👋 Orbit local dev stopped cleanly."
    exit 0
}
trap cleanup SIGINT SIGTERM EXIT

# 5. Launch Orbit Agent in background
echo "🚀 Starting Orbit Sovereign Agent on http://127.0.0.1:${PORT}..."
./bin/orbit-agent &
AGENT_PID=$!

# Wait for agent to come online
for i in {1..20}; do
    if curl -s "http://127.0.0.1:${PORT}/healthz" >/dev/null 2>&1; then
        break
    fi
    sleep 0.25
done

# Read generated master key
MASTER_KEY=""
if [ -f "$DEV_ROOT/secrets/master.key" ]; then
    MASTER_KEY=$(od -An -tx1 "$DEV_ROOT/secrets/master.key" | tr -d ' \n')
fi
export ORBIT_TOKEN="$MASTER_KEY"

# 6. Seed Sample Workload & Dedicated DB (if first launch)
if [ ! -f "$DEV_ROOT/.seeded" ]; then
    echo "🌱 Seeding sample demo workload and micro-RAM database..."
    
    # Create sample demo app
    curl -s -X POST "http://127.0.0.1:${PORT}/api/v1/apps" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $MASTER_KEY" \
        -d '{"app_id":"demo-worker","name":"Edge API Worker","domain":"api.local.orbit","ram_limit_mb":128,"cpu_quota_pct":100}' >/dev/null 2>&1 || true

    # Provision sample dedicated database
    curl -s -X POST "http://127.0.0.1:${PORT}/api/v1/databases" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $MASTER_KEY" \
        -d '{"name":"main_d1_db","engine":"postgres","ram_profile":"micro"}' >/dev/null 2>&1 || true

    # Create sample S3 bucket
    curl -s -X PUT "http://127.0.0.1:${PORT}/s3/assets" \
        -H "Authorization: Bearer $MASTER_KEY" >/dev/null 2>&1 || true

    touch "$DEV_ROOT/.seeded"
fi

echo ""
echo "========================================================================"
echo "🎉 Orbit Sovereign PaaS is LIVE!"
echo ""
echo "🌐 Cloudflare Dashboard:  http://localhost:${PORT}"
echo "🔑 Master Admin Key:      ${MASTER_KEY}"
echo "📁 Local Storage Root:    ${DEV_ROOT}"
echo ""
echo "Try running these in another terminal tab:"
echo "   ./bin/orbit status"
echo "   ./bin/orbit db create analytics_db --engine=postgres"
echo "   ./bin/orbit push ./my-app --app=demo-worker"
echo "========================================================================"
echo "Press Ctrl+C to stop."
echo ""

# Attempt to open browser on Mac / Linux
DASHBOARD_URL="http://localhost:${PORT}/?token=${MASTER_KEY}"
if command -v open >/dev/null 2>&1; then
    open "$DASHBOARD_URL" >/dev/null 2>&1 || true
elif command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$DASHBOARD_URL" >/dev/null 2>&1 || true
fi

# Wait for process
wait "$AGENT_PID"
