<template>
  <div class="orbit-shell">
    <!-- CLOUDFLARE STYLE SIDEBAR -->
    <aside class="cf-sidebar">
      <div class="sidebar-header">
        <a href="#" class="cf-logo" @click.prevent="currentTab = 'overview'">
          <!-- Cloudflare-style Orange Cloud SVG -->
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M37.5 19.5C36.8 13.5 31.8 9 25.5 9C20.2 9 15.8 12.2 14.1 16.8C13.4 16.6 12.7 16.5 12 16.5C7.6 16.5 4 20.1 4 24.5C4 28.9 7.6 32.5 12 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5Z" fill="#F38020"/>
            <path d="M32 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5C36.8 13.5 31.8 9 25.5 9C24.5 9 23.5 9.1 22.5 9.4C26.5 12.8 29.1 17.8 29.4 23.5C30.2 23.5 31.1 23.6 32 24C32.1 26.8 32 29.6 32 32.5Z" fill="#FAAD3F"/>
          </svg>
          <span>Orbit Sovereign</span>
        </a>
      </div>

      <div class="account-selector">
        <div style="display: flex; align-items: center; gap: 0.5rem; overflow: hidden;">
          <span style="width: 8px; height: 8px; border-radius: 50%; background: #f38020; flex-shrink: 0;"></span>
          <span style="font-weight: 500; text-overflow: ellipsis; overflow: hidden; white-space: nowrap;">Sovereign Cluster (Primary)</span>
        </div>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 9 6 6 6-6"/></svg>
      </div>

      <div class="sidebar-nav">
        <div class="nav-section-title">Core Services</div>
        <button class="cf-nav-item" :class="{ active: currentTab === 'overview' }" @click="currentTab = 'overview'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg>
          <span>Overview</span>
        </button>

        <button class="cf-nav-item" :class="{ active: currentTab === 'workers' }" @click="currentTab = 'workers'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
          <span>Workers & Compute</span>
        </button>

        <button class="cf-nav-item" :class="{ active: currentTab === 'databases' }" @click="currentTab = 'databases'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/><path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"/></svg>
          <span>D1 & Databases</span>
        </button>

        <button class="cf-nav-item" :class="{ active: currentTab === 'r2' }" @click="currentTab = 'r2'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.29 7 12 12 20.71 7"/><line x1="12" x2="12" y1="22" y2="12"/></svg>
          <span>R2 Storage (S3)</span>
        </button>

        <div class="nav-section-title">Host & Governance</div>
        <button class="cf-nav-item" :class="{ active: currentTab === 'files' }" @click="currentTab = 'files'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"/></svg>
          <span>Jailed Filesystem</span>
        </button>

        <button class="cf-nav-item" :class="{ active: currentTab === 'security' }" @click="currentTab = 'security'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
          <span>Zero Trust & Security</span>
        </button>
      </div>

      <div class="sidebar-footer">
        <div class="node-badge">
          <div class="status-dot"></div>
          <span style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ nodeInfo.node_id || 'node_sovereign' }}</span>
        </div>
      </div>
    </aside>

    <!-- MAIN PANE -->
    <div class="orbit-main">
      <!-- CLOUDFLARE TOP HEADER -->
      <header class="cf-header">
        <div class="header-left">
          <div class="search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
            <input type="text" placeholder="Search applications, databases, buckets... (⌘K)" />
          </div>
        </div>

        <div class="header-right">
          <div class="operational-pill">
            <span class="status-dot" style="width: 6px; height: 6px;"></span>
            <span>All Systems Operational</span>
          </div>
          <button class="theme-toggle-btn" @click="toggleTheme" :title="isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
            <span v-if="isDarkMode">☀️ Light</span>
            <span v-else>🌙 Dark</span>
          </button>
          <button class="btn-cf btn-cf-primary btn-cf-sm" @click="showNewAppModal = true">+ Create Application</button>
        </div>
      </header>

      <!-- CONTENT CANVAS -->
      <div class="cf-content">
        <!-- BREADCRUMBS -->
        <div class="breadcrumbs">
          <span>Sovereign Cloud</span>
          <span>/</span>
          <span>{{ nodeInfo.hostname || 'localhost' }}</span>
          <span>/</span>
          <span class="current">{{ capitalize(currentTab) }}</span>
        </div>

        <!-- OVERVIEW TAB -->
        <div v-if="currentTab === 'overview'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Sovereign Node Overview</h1>
              <p class="page-subtitle">Real-time metrics, cgroups v2 resource quotas, and zero-docker service telemetries.</p>
            </div>
            <div style="display: flex; gap: 0.5rem;">
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="fetchState">Refresh State</button>
            </div>
          </div>

          <!-- Cloudflare-style Metric Banners -->
          <div class="metrics-banner">
            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Host CPU Utilization</span>
                <span class="metric-tag green">Healthy</span>
              </div>
              <div class="metric-big">{{ metrics.cpu_usage_pct.toFixed(1) }}%</div>
              <div class="metric-foot">
                <span>Kernel cgroups v2 throttled</span>
              </div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Host Memory Allocated</span>
                <span class="metric-tag orange">{{ formatBytes(metrics.ram_used_bytes) }}</span>
              </div>
              <div class="metric-big">{{ formatBytes(metrics.ram_used_bytes) }}</div>
              <div class="metric-foot">
                <span>Total Available: {{ formatBytes(metrics.ram_total_bytes) }}</span>
              </div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Orbit Agent Footprint</span>
                <span class="metric-tag green">&lt; 25MB Limit</span>
              </div>
              <div class="metric-big" style="color: #38bdf8;">{{ formatBytes(metrics.agent_alloc_ram) }}</div>
              <div class="metric-foot">
                <span>Single static binary runtime</span>
              </div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Active Workloads</span>
                <span class="metric-tag green">{{ apps.length }} Active</span>
              </div>
              <div class="metric-big">{{ apps.length }}</div>
              <div class="metric-foot">
                <span>{{ databases.length }} Dedicated Native DBs</span>
              </div>
            </div>
          </div>

          <!-- Apps Table -->
          <div class="cf-table-wrapper">
            <div class="cf-table-header">
              <h3 style="font-size: 1rem; font-weight: 600;">Workloads & Applications</h3>
              <button class="btn-cf btn-cf-primary btn-cf-sm" @click="showNewAppModal = true">+ Deploy App</button>
            </div>
            <table class="cf-table">
              <thead>
                <tr>
                  <th>Application</th>
                  <th>Route / Domain</th>
                  <th>Port (Active / Idle)</th>
                  <th>RAM Quota</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="app in apps" :key="app.app_id">
                  <td>
                    <strong>{{ app.name }}</strong>
                    <div style="font-family: var(--font-mono); font-size: 0.78rem; color: var(--cf-text-dim);">{{ app.app_id }}</div>
                  </td>
                  <td><code>{{ app.domain || 'localhost:' + app.active_port }}</code></td>
                  <td>
                    <span class="cf-pill orange">:{{ app.active_port }}</span>
                    <span style="color: var(--cf-text-dim); margin-left: 0.35rem;">(idle :{{ app.idle_port }})</span>
                  </td>
                  <td>{{ app.ram_limit_mb }} MB</td>
                  <td><span class="cf-pill active">● {{ app.status }}</span></td>
                  <td>
                    <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="openLogs(app.app_id)">View Logs</button>
                  </td>
                </tr>
                <tr v-if="apps.length === 0">
                  <td colspan="6" style="text-align: center; color: var(--cf-text-muted); padding: 2.5rem;">
                    No applications deployed yet. Run <code>orbit push</code> from your local terminal to deploy.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- WORKERS & COMPUTE TAB -->
        <div v-if="currentTab === 'workers'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Workers & Compute Engine</h1>
              <p class="page-subtitle">Wrangler-style incremental delta deployment with zero-downtime Blue/Green port swapping.</p>
            </div>
            <button class="btn-cf btn-cf-primary" @click="showNewAppModal = true">+ Create Application</button>
          </div>

          <div class="cf-table-wrapper">
            <table class="cf-table">
              <thead>
                <tr>
                  <th>App Slug</th>
                  <th>Serving Route</th>
                  <th>Active Port</th>
                  <th>cgroups v2 Quota</th>
                  <th>Version</th>
                  <th>Status</th>
                  <th>Console</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="app in apps" :key="app.app_id">
                  <td><strong>{{ app.name }}</strong><br><span style="font-family: var(--font-mono); font-size: 0.75rem; color: var(--cf-text-dim);">{{ app.app_id }}</span></td>
                  <td><code>{{ app.domain || '127.0.0.1:' + app.active_port }}</code></td>
                  <td><span class="cf-pill orange">:{{ app.active_port }}</span></td>
                  <td>{{ app.ram_limit_mb }}MB RAM / {{ app.cpu_quota_pct }}% CPU</td>
                  <td><span style="font-family: var(--font-mono);">{{ app.current_version || 'initial' }}</span></td>
                  <td><span class="cf-pill active">● {{ app.status }}</span></td>
                  <td>
                    <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="openLogs(app.app_id)">Live Logs</button>
                  </td>
                </tr>
                <tr v-if="apps.length === 0">
                  <td colspan="7" style="text-align: center; padding: 2rem; color: var(--cf-text-muted);">No workloads found. Deploy with <code>orbit push</code>.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Live Terminal Viewer -->
          <div v-if="activeLogApp" class="cf-card" style="margin-top: 1.5rem;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
              <h3 style="font-size: 0.95rem; font-weight: 600;">Live Log Stream: {{ activeLogApp }}</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="activeLogApp = null">Close Terminal</button>
            </div>
            <div class="cf-terminal" ref="terminalRef">
              <div v-for="(line, idx) in appLogs" :key="idx">
                <span class="timestamp">[{{ formatTime(line.timestamp) }}]</span> {{ line.message }}
              </div>
              <div v-if="appLogs.length === 0" style="color: var(--cf-text-dim);">Listening for Server-Sent Events...</div>
            </div>
          </div>
        </div>

        <!-- D1 & DATABASES TAB -->
        <div v-if="currentTab === 'databases'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Dedicated Native Databases (D1 Style)</h1>
              <p class="page-subtitle">Isolated, dedicated native PostgreSQL 16 & MongoDB 7 with micro-RAM profiles and zero-network IPC sockets.</p>
            </div>
            <button class="btn-cf btn-cf-primary" @click="showNewDBModal = true">+ Provision Database</button>
          </div>

          <div class="cf-table-wrapper">
            <table class="cf-table">
              <thead>
                <tr>
                  <th>Instance ID & Name</th>
                  <th>Engine</th>
                  <th>RAM Profile</th>
                  <th>Port (Loopback)</th>
                  <th>IPC Domain Socket</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="d in databases" :key="d.db_id">
                  <td><strong>{{ d.name }}</strong><br><span style="font-family: var(--font-mono); font-size: 0.75rem; color: var(--cf-text-dim);">{{ d.db_id }}</span></td>
                  <td><span class="cf-pill orange">{{ d.engine }}</span></td>
                  <td><span class="cf-pill stopped">{{ d.ram_profile }} (&lt;32MB)</span></td>
                  <td><code>127.0.0.1:{{ d.port }}</code></td>
                  <td><span style="font-family: var(--font-mono); font-size: 0.8rem; color: var(--cf-text-muted);">{{ d.socket_path }}</span></td>
                  <td><span class="cf-pill active">● {{ d.status }}</span></td>
                </tr>
                <tr v-if="databases.length === 0">
                  <td colspan="6" style="text-align: center; padding: 2rem; color: var(--cf-text-muted);">No databases provisioned yet. Click "+ Provision Database" to create one.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- R2 OBJECT STORAGE (S3) TAB -->
        <div v-if="currentTab === 'r2'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">R2 Object Storage (S3 API)</h1>
              <p class="page-subtitle">Sovereign, built-in S3-compatible engine supporting AWS Signature v4, HTTP range GETs, and multipart uploads.</p>
            </div>
            <button class="btn-cf btn-cf-primary" @click="showNewBucketModal = true">+ Create Bucket</button>
          </div>

          <div class="metrics-banner">
            <div v-for="b in buckets" :key="b.bucket_id" class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">{{ b.bucket_id }}</span>
                <span class="cf-pill" :class="b.is_public ? 'stopped' : 'active'">{{ b.is_public ? 'Public' : 'Private (SigV4)' }}</span>
              </div>
              <div style="font-size: 0.9rem; color: var(--cf-text-secondary); margin-top: 0.5rem;">
                Endpoint: <code>/s3/{{ b.bucket_id }}/</code>
              </div>
              <div class="metric-foot">
                <span>Owner: {{ b.owner_id }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- JAILED FILESYSTEM TAB -->
        <div v-if="currentTab === 'files'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Jailed File Explorer</h1>
              <p class="page-subtitle">Kernel-hardened isolation via openat2 (RESOLVE_BENEATH) strictly confined to /opt/orbit.</p>
            </div>
            <div style="display: flex; gap: 0.5rem;">
              <select v-model="fileScope" class="cf-select" style="width: 140px;" @change="loadFiles">
                <option value="app">Apps Scope</option>
                <option value="storage">Storage Scope</option>
              </select>
              <input v-model="fileTargetID" placeholder="Enter App or Bucket ID..." class="cf-input" style="width: 220px;" @change="loadFiles" />
              <button class="btn-cf btn-cf-secondary" @click="loadFiles">Browse</button>
            </div>
          </div>

          <div class="cf-table-wrapper">
            <table class="cf-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Size</th>
                  <th>Last Modified</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in fileItems" :key="item.path">
                  <td>
                    <span v-if="item.is_dir">📁 </span>
                    <span v-else>📄 </span>
                    <strong>{{ item.name }}</strong>
                  </td>
                  <td>{{ item.is_dir ? 'Directory' : 'File' }}</td>
                  <td>{{ item.is_dir ? '—' : formatBytes(item.size) }}</td>
                  <td>{{ item.mod_time }}</td>
                </tr>
                <tr v-if="fileItems.length === 0">
                  <td colspan="4" style="text-align: center; padding: 2.5rem; color: var(--cf-text-muted);">
                    Enter an App ID or Bucket ID above to safely inspect files inside the sandbox.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- SECURITY TAB -->
        <div v-if="currentTab === 'security'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Zero Trust & Security Hardening</h1>
              <p class="page-subtitle">Host nftables firewall, AES-256-GCM encrypted secrets vault, and WireGuard mesh status.</p>
            </div>
          </div>

          <div class="metrics-banner">
            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Firewall</span>
                <span class="metric-tag green">nftables Active</span>
              </div>
              <div class="metric-big">Secured</div>
              <div class="metric-foot">Zero public administrative ports exposed</div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Secrets Vault</span>
                <span class="metric-tag green">AES-256-GCM</span>
              </div>
              <div class="metric-big">Encrypted</div>
              <div class="metric-foot">Injected directly into memory at runtime</div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Process Isolation</span>
                <span class="metric-tag green">Namespaces</span>
              </div>
              <div class="metric-big">cgroups v2</div>
              <div class="metric-foot">Non-root OS users per app</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: NEW APP -->
    <div v-if="showNewAppModal" class="cf-modal-backdrop">
      <div class="cf-modal-box">
        <div class="cf-modal-head">
          <h3 style="font-size: 1rem; font-weight: 600;">Deploy New Application</h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showNewAppModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>App Identifier (Slug)</label>
            <input v-model="newApp.app_id" placeholder="e.g. api-service" class="cf-input" />
          </div>
          <div class="cf-form-group">
            <label>Application Name</label>
            <input v-model="newApp.name" placeholder="e.g. Core API Service" class="cf-input" />
          </div>
          <div class="cf-form-group">
            <label>Custom Domain (Optional)</label>
            <input v-model="newApp.domain" placeholder="e.g. api.example.com" class="cf-input" />
          </div>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
            <div class="cf-form-group">
              <label>RAM Limit (MB)</label>
              <input v-model.number="newApp.ram_limit_mb" type="number" class="cf-input" />
            </div>
            <div class="cf-form-group">
              <label>CPU Quota (%)</label>
              <input v-model.number="newApp.cpu_quota_pct" type="number" class="cf-input" />
            </div>
          </div>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-secondary" @click="showNewAppModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="createApp">Create Application</button>
        </div>
      </div>
    </div>

    <!-- MODAL: NEW DATABASE -->
    <div v-if="showNewDBModal" class="cf-modal-backdrop">
      <div class="cf-modal-box">
        <div class="cf-modal-head">
          <h3 style="font-size: 1rem; font-weight: 600;">Provision Dedicated Database</h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showNewDBModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>Database Name</label>
            <input v-model="newDB.name" placeholder="e.g. production_db" class="cf-input" />
          </div>
          <div class="cf-form-group">
            <label>Engine</label>
            <select v-model="newDB.engine" class="cf-select">
              <option value="postgres">PostgreSQL 16 (Native Low-RAM)</option>
              <option value="mongodb">MongoDB 7 (Native Low-RAM)</option>
            </select>
          </div>
          <div class="cf-form-group">
            <label>RAM Profile</label>
            <select v-model="newDB.ram_profile" class="cf-select">
              <option value="micro">Micro (&lt;32MB RAM, 20 max conn)</option>
              <option value="standard">Standard (&lt;128MB RAM, 50 max conn)</option>
              <option value="heavy">Heavy (&lt;512MB RAM, 100 max conn)</option>
            </select>
          </div>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-secondary" @click="showNewDBModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="createDB">Provision</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

const currentTab = ref('overview')
const nodeInfo = ref({})
const apps = ref([])
const databases = ref([])
const buckets = ref([])
const fileItems = ref([])
const fileScope = ref('app')
const fileTargetID = ref('')

const metrics = ref({
  cpu_usage_pct: 0,
  ram_used_bytes: 0,
  ram_total_bytes: 0,
  agent_alloc_ram: 0,
})

const showNewAppModal = ref(false)
const newApp = ref({
  app_id: '',
  name: '',
  domain: '',
  ram_limit_mb: 256,
  cpu_quota_pct: 100,
})

const showNewDBModal = ref(false)
const newDB = ref({
  name: '',
  engine: 'postgres',
  ram_profile: 'micro',
})

const showNewBucketModal = ref(false)

const activeLogApp = ref(null)
const appLogs = ref([])
let logEventSource = null
let telemetryEventSource = null
const terminalRef = ref(null)

const capitalize = (s) => s.charAt(0).toUpperCase() + s.slice(1)

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatTime = (isoString) => {
  if (!isoString) return ''
  const d = new Date(isoString)
  return d.toLocaleTimeString()
}

const authToken = ref(localStorage.getItem('orbit_token') || '')

// Capture ?token=... from URL if provided (e.g. from dev.sh)
const urlParams = new URLSearchParams(window.location.search)
if (urlParams.get('token')) {
  authToken.value = urlParams.get('token')
  localStorage.setItem('orbit_token', authToken.value)
}

const apiFetch = async (url, options = {}) => {
  options.headers = options.headers || {}
  if (authToken.value) {
    options.headers['Authorization'] = 'Bearer ' + authToken.value
  }
  return fetch(url, options)
}

const fetchState = async () => {
  try {
    const [nodeRes, appsRes, dbsRes, bucketsRes] = await Promise.all([
      apiFetch('/api/v1/node').catch(() => null),
      apiFetch('/api/v1/apps').catch(() => null),
      apiFetch('/api/v1/databases').catch(() => null),
      apiFetch('/s3/').catch(() => null),
    ])
    if (nodeRes && nodeRes.ok) nodeInfo.value = await nodeRes.json()
    if (appsRes && appsRes.ok) apps.value = await appsRes.json()
    if (dbsRes && dbsRes.ok) databases.value = await dbsRes.json()
  } catch (err) {
    console.error('Fetch error:', err)
  }
}

const startTelemetrySSE = () => {
  const sseURL = authToken.value ? `/api/v1/telemetry/stream?token=${authToken.value}` : '/api/v1/telemetry/stream'
  telemetryEventSource = new EventSource(sseURL)
  telemetryEventSource.onmessage = (event) => {
    try {
      metrics.value = JSON.parse(event.data)
    } catch (e) {}
  }
}

const openLogs = (appID) => {
  if (logEventSource) logEventSource.close()
  activeLogApp.value = appID
  appLogs.value = []

  const sseURL = authToken.value ? `/api/v1/apps/${appID}/logs/stream?token=${authToken.value}` : `/api/v1/apps/${appID}/logs/stream`
  logEventSource = new EventSource(sseURL)
  logEventSource.onmessage = (event) => {
    try {
      const line = JSON.parse(event.data)
      appLogs.value.push(line)
      if (appLogs.value.length > 500) appLogs.value.shift()
      nextTick(() => {
        if (terminalRef.value) terminalRef.value.scrollTop = terminalRef.value.scrollHeight
      })
    } catch (e) {}
  }
}

const createApp = async () => {
  if (!newApp.value.app_id || !newApp.value.name) return
  try {
    const res = await apiFetch('/api/v1/apps', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newApp.value),
    })
    if (res.ok) {
      showNewAppModal.value = false
      newApp.value = { app_id: '', name: '', domain: '', ram_limit_mb: 256, cpu_quota_pct: 100 }
      fetchState()
    }
  } catch (e) {}
}

const createDB = async () => {
  if (!newDB.value.name) return
  try {
    const res = await apiFetch('/api/v1/databases', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newDB.value),
    })
    if (res.ok) {
      showNewDBModal.value = false
      newDB.value = { name: '', engine: 'postgres', ram_profile: 'micro' }
      fetchState()
    }
  } catch (e) {}
}

const loadFiles = async () => {
  if (!fileTargetID.value) return
  try {
    const res = await apiFetch(`/api/v1/jail/list?scope=${fileScope.value}&id=${fileTargetID.value}&path=`)
    if (res.ok) {
      const data = await res.json()
      fileItems.value = data.items || []
    }
  } catch (e) {}
}

const isDarkMode = ref(false)

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  if (isDarkMode.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
    localStorage.setItem('orbit_theme', 'dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
    localStorage.setItem('orbit_theme', 'light')
  }
}

onMounted(() => {
  const saved = localStorage.getItem('orbit_theme')
  if (saved === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    isDarkMode.value = false
    document.documentElement.removeAttribute('data-theme')
  }
  fetchState()
  startTelemetrySSE()
  setInterval(fetchState, 5000)
})

onUnmounted(() => {
  if (telemetryEventSource) telemetryEventSource.close()
  if (logEventSource) logEventSource.close()
})
</script>
