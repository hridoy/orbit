#!/usr/bin/env bash
# One-command Orbit install or upgrade for Ubuntu/Debian Linux (amd64 or arm64).
#
#   sudo ORBIT_REPO=owner/repo ORBIT_PUBKEY=<hex> bash install.sh
#
# Downloads the latest signed release, verifies the manifest signature and binary checksum,
# installs the agent as a systemd service, and opens the firewall for SSH, HTTP, HTTPS, and RDP
# (only if xrdp is installed). Safe to re-run: it upgrades in place.
set -euo pipefail

ORBIT_REPO="${ORBIT_REPO:?set ORBIT_REPO=owner/repo}"
ORBIT_PUBKEY="${ORBIT_PUBKEY:?set ORBIT_PUBKEY to the release public key (hex)}"
ORBIT_VERSION="${ORBIT_VERSION:-latest}"
ROOT="/opt/orbit"

log() { printf '==> %s\n' "$*"; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

[ "$(id -u)" -eq 0 ] || die "run as root"
command -v openssl >/dev/null || die "openssl is required"
command -v curl >/dev/null || die "curl is required"

case "$(uname -m)" in
  x86_64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) die "unsupported architecture: $(uname -m)" ;;
esac

if [ "$ORBIT_VERSION" = "latest" ]; then
  BASE="https://github.com/${ORBIT_REPO}/releases/latest/download"
else
  BASE="https://github.com/${ORBIT_REPO}/releases/download/${ORBIT_VERSION}"
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

log "downloading release manifest"
curl -fsSL -o "$TMP/manifest.json" "$BASE/manifest.json"
curl -fsSL -o "$TMP/manifest.sig" "$BASE/manifest.sig"

log "verifying manifest signature"
# Ed25519 public key as DER (SubjectPublicKeyInfo) from the raw 32-byte hex key.
unhex() { python3 -c 'import sys,binascii;sys.stdout.buffer.write(binascii.unhexlify(sys.stdin.read().strip()))'; }
{ printf '302a300506032b6570032100'; printf '%s' "$ORBIT_PUBKEY"; } | unhex > "$TMP/pub.der"
unhex < "$TMP/manifest.sig" > "$TMP/manifest.sig.bin"
openssl pkeyutl -verify -pubin -keyform DER -inkey "$TMP/pub.der" -rawin \
  -in "$TMP/manifest.json" -sigfile "$TMP/manifest.sig.bin" >/dev/null \
  || die "manifest signature is INVALID; refusing to install"

KEY="linux/${ARCH}"
ASSET_URL="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1]))["assets"][sys.argv[2]]["url"])' "$TMP/manifest.json" "$KEY")"
ASSET_SHA="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1]))["assets"][sys.argv[2]]["sha256"])' "$TMP/manifest.json" "$KEY")"
VERSION="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1]))["version"])' "$TMP/manifest.json")"
case "$ASSET_URL" in
  "https://github.com/${ORBIT_REPO}/releases/download/"*) ;;
  *) die "asset URL is outside ${ORBIT_REPO} releases" ;;
esac

log "downloading orbit-agent ${VERSION} (${KEY})"
curl -fsSL -o "$TMP/orbit-agent" "$ASSET_URL"
echo "${ASSET_SHA}  $TMP/orbit-agent" | sha256sum -c --quiet || die "binary checksum mismatch"

log "installing to ${ROOT}"
install -d -m 0700 "$ROOT" "$ROOT/secrets" "$ROOT/data"
install -d -m 0755 "$ROOT/bin"
# Keep the previous binary for manual rollback.
[ -f "$ROOT/bin/orbit-agent" ] && cp -f "$ROOT/bin/orbit-agent" "$ROOT/bin/orbit-agent.prev"
install -m 0755 "$TMP/orbit-agent" "$ROOT/bin/orbit-agent"

log "writing systemd service"
cat > /etc/systemd/system/orbit-agent.service <<UNIT
[Unit]
Description=Orbit Sovereign PaaS Agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=${ROOT}
ExecStart=${ROOT}/bin/orbit-agent
Restart=always
RestartSec=3s
LimitNOFILE=65536
MemoryMax=512M

[Install]
WantedBy=multi-user.target
UNIT
systemctl daemon-reload
systemctl enable orbit-agent >/dev/null 2>&1 || true
systemctl restart orbit-agent

log "waiting for health check"
for _ in $(seq 1 20); do
  if curl -fsS -o /dev/null http://127.0.0.1:8787/healthz; then break; fi
  sleep 1
done
curl -fsS -o /dev/null http://127.0.0.1:8787/healthz || die "agent did not become healthy; see: journalctl -u orbit-agent"

log "scheduling daily state backup (keeps the 14 newest)"
install -d -m 0700 "$ROOT/backups"
cat > /etc/systemd/system/orbit-backup.service <<UNIT
[Unit]
Description=Orbit control-plane state backup

[Service]
Type=oneshot
ExecStart=/bin/sh -c 'umask 077 && ${ROOT}/bin/orbit-agent backup ${ROOT}/backups && ls -1dt ${ROOT}/backups/orbit-backup-* | tail -n +15 | xargs -r rm -rf'
UNIT
cat > /etc/systemd/system/orbit-backup.timer <<UNIT
[Unit]
Description=Daily Orbit state backup

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
UNIT
systemctl daemon-reload
systemctl enable --now orbit-backup.timer >/dev/null 2>&1 || true

log "installing rclone for off-host backups (Google Drive, Dropbox, OneDrive, S3)"
if ! command -v rclone >/dev/null; then
  DEBIAN_FRONTEND=noninteractive apt-get install -y -qq rclone >/dev/null || log "rclone install failed; cloud backups unavailable"
fi

log "installing fail2ban for SSH brute-force protection"
if ! command -v fail2ban-client >/dev/null; then
  DEBIAN_FRONTEND=noninteractive apt-get install -y -qq fail2ban >/dev/null || log "fail2ban install failed; continuing"
fi
if [ ! -f /etc/fail2ban/jail.d/orbit.local ] && command -v fail2ban-client >/dev/null; then
  install -d /etc/fail2ban/jail.d
  printf '[sshd]\nenabled = true\nmaxretry = 5\nfindtime = 10m\nbantime = 1h\n' > /etc/fail2ban/jail.d/orbit.local
  systemctl enable --now fail2ban >/dev/null 2>&1 || true
  systemctl restart fail2ban >/dev/null 2>&1 || true
fi

log "configuring firewall"
if command -v ufw >/dev/null; then
  ufw allow 22/tcp >/dev/null          # never lock out SSH
  ufw allow 80/tcp >/dev/null
  ufw allow 443/tcp >/dev/null
  if ss -ltn | grep -q ':3389 '; then
    ufw allow 3389/tcp >/dev/null      # keep RDP if xrdp is running
  fi
  ufw --force enable >/dev/null
else
  log "ufw not found; skipping firewall (install ufw to manage ports)"
fi

log "done: Orbit ${VERSION} is running"
echo "   dashboard (on the server):  curl http://127.0.0.1:8787/healthz"
echo "   from your computer:         ssh -L 8787:127.0.0.1:8787 root@<server-ip>"
echo "   admin token file:           ${ROOT}/secrets/admin.token"
echo "   public apps:                ports 80 and 443 are served by the agent; add a domain in Domains"
