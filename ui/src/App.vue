<template>
  <div class="app-container">
    <header>
      <div class="brand">
        <svg viewBox="0 0 512 512">
          <circle cx="256" cy="256" r="100" fill="none" stroke="#38bdf8" stroke-width="24"/>
          <circle cx="256" cy="256" r="48" fill="#38bdf8"/>
        </svg>
        <span>ORBIT PAAS</span>
        <span class="badge badge-success" style="margin-left: 0.5rem;">Sovereign Node</span>
      </div>

      <nav>
        <ul class="nav-tabs">
          <li><button class="nav-tab" :class="{ active: currentTab === 'dashboard' }" @click="currentTab = 'dashboard'">Dashboard</button></li>
          <li><button class="nav-tab" :class="{ active: currentTab === 'apps' }" @click="currentTab = 'apps'">Workloads ({{ apps.length }})</button></li>
          <li><button class="nav-tab" :class="{ active: currentTab === 'databases' }" @click="currentTab = 'databases'">Databases ({{ databases.length }})</button></li>
          <li><button class="nav-tab" :class="{ active: currentTab === 'storage' }" @click="currentTab = 'storage'">S3 Storage</button></li>
          <li><button class="nav-tab" :class="{ active: currentTab === 'files' }" @click="currentTab = 'files'">File Browser</button></li>
        </ul>
      </nav>
    </header>

    <!-- DASHBOARD VIEW -->
    <main v-if="currentTab === 'dashboard'">
      <div class="grid-3">
        <div class="card">
          <div class="card-title">Host CPU Usage</div>
          <div class="stat-value">{{ metrics.cpu_usage_pct.toFixed(1) }}%</div>
          <div class="stat-desc">Kernel & cgroups v2 throttled</div>
        </div>
        <div class="card">
          <div class="card-title">Host Memory Used</div>
          <div class="stat-value">{{ formatBytes(metrics.ram_used_bytes) }}</div>
          <div class="stat-desc">Total: {{ formatBytes(metrics.ram_total_bytes) }}</div>
        </div>
        <div class="card">
          <div class="card-title">Agent Footprint</div>
          <div class="stat-value" style="color: #38bdf8;">{{ formatBytes(metrics.agent_alloc_ram) }}</div>
          <div class="stat-desc">&lt;25MB target idle footprint</div>
        </div>
      </div>

      <div class="grid-2">
        <div class="card">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
            <h3 style="font-size: 1.1rem;">Active Workloads</h3>
            <span class="badge badge-blue">{{ apps.length }} Apps</span>
          </div>
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>App ID</th>
                  <th>Active Port</th>
                  <th>RAM Cap</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="app in apps" :key="app.app_id">
                  <td><strong>{{ app.name }}</strong><br><span style="color: var(--text-muted); font-size: 0.75rem;">{{ app.app_id }}</span></td>
                  <td><span class="badge badge-blue">:{{ app.active_port }}</span></td>
                  <td>{{ app.ram_limit_mb }} MB</td>
                  <td><span class="badge badge-success">{{ app.status }}</span></td>
                </tr>
                <tr v-if="apps.length === 0">
                  <td colspan="4" style="text-align: center; color: var(--text-muted); padding: 2rem;">No workloads deployed yet. Use <code>orbit push</code> to launch.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="card">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
            <h3 style="font-size: 1.1rem;">Dedicated Databases</h3>
            <span class="badge badge-blue">{{ databases.length }} Instances</span>
          </div>
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>Instance</th>
                  <th>Engine</th>
                  <th>Port</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="db in databases" :key="db.db_id">
                  <td><strong>{{ db.name }}</strong><br><span style="color: var(--text-muted); font-size: 0.75rem;">{{ db.db_id }}</span></td>
                  <td><span class="badge badge-blue">{{ db.engine }}</span></td>
                  <td>:{{ db.port }}</td>
                  <td><span class="badge badge-success">{{ db.status }}</span></td>
                </tr>
                <tr v-if="databases.length === 0">
                  <td colspan="4" style="text-align: center; color: var(--text-muted); padding: 2rem;">No databases created yet.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </main>

    <!-- APPS VIEW -->
    <main v-if="currentTab === 'apps'">
      <div class="card" style="margin-bottom: 1.5rem;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
          <h2>Application Workloads</h2>
          <button class="btn btn-primary" @click="showNewAppModal = true">+ New App</button>
        </div>
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>App Name & ID</th>
                <th>Domain</th>
                <th>Port (Active / Idle)</th>
                <th>Limits (RAM / CPU)</th>
                <th>Version</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="app in apps" :key="app.app_id">
                <td><strong>{{ app.name }}</strong><br><span style="font-family: var(--font-mono); font-size: 0.8rem; color: var(--text-muted);">{{ app.app_id }}</span></td>
                <td>{{ app.domain || '(none)' }}</td>
                <td><span class="badge badge-blue">{{ app.active_port }}</span> / {{ app.idle_port }}</td>
                <td>{{ app.ram_limit_mb }}MB / {{ app.cpu_quota_pct }}%</td>
                <td><span style="font-family: var(--font-mono);">{{ app.current_version || 'initial' }}</span></td>
                <td><span class="badge badge-success">{{ app.status }}</span></td>
                <td>
                  <button class="btn btn-secondary" style="padding: 0.25rem 0.6rem; font-size: 0.8rem;" @click="openLogs(app.app_id)">Logs</button>
                </td>
              </tr>
              <tr v-if="apps.length === 0">
                <td colspan="7" style="text-align: center; padding: 2rem; color: var(--text-muted);">No workloads found. Run <code>orbit push</code> to deploy.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Live Terminal -->
      <div v-if="activeLogApp" class="card">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
          <h3>Live Telemetry Logs: {{ activeLogApp }}</h3>
          <button class="btn btn-secondary" style="padding: 0.2rem 0.5rem;" @click="activeLogApp = null">Close</button>
        </div>
        <div class="log-terminal" ref="terminalRef">
          <div v-for="(line, idx) in appLogs" :key="idx">
            <span style="color: #64748b;">[{{ formatTime(line.timestamp) }}]</span> {{ line.message }}
          </div>
          <div v-if="appLogs.length === 0" style="color: #64748b;">Waiting for log stream...</div>
        </div>
      </div>
    </main>

    <!-- DATABASES VIEW -->
    <main v-if="currentTab === 'databases'">
      <div class="card">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
          <div>
            <h2>Dedicated Native Databases</h2>
            <p style="color: var(--text-muted); font-size: 0.85rem;">Zero-Docker, dedicated native PostgreSQL and MongoDB instances with internal socket & loopback binding.</p>
          </div>
          <button class="btn btn-primary" @click="showNewDBModal = true">+ Provision DB</button>
        </div>
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Database ID & Name</th>
                <th>Engine</th>
                <th>RAM Profile</th>
                <th>Internal Port</th>
                <th>IPC Socket</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="d in databases" :key="d.db_id">
                <td><strong>{{ d.name }}</strong><br><span style="font-family: var(--font-mono); font-size: 0.8rem; color: var(--text-muted);">{{ d.db_id }}</span></td>
                <td><span class="badge badge-blue">{{ d.engine }}</span></td>
                <td><span class="badge badge-warning">{{ d.ram_profile }} (&lt;32MB)</span></td>
                <td>127.0.0.1:{{ d.port }}</td>
                <td><span style="font-family: var(--font-mono); font-size: 0.8rem;">{{ d.socket_path }}</span></td>
                <td><span class="badge badge-success">{{ d.status }}</span></td>
              </tr>
              <tr v-if="databases.length === 0">
                <td colspan="6" style="text-align: center; padding: 2rem; color: var(--text-muted);">No databases provisioned yet.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <!-- S3 STORAGE VIEW -->
    <main v-if="currentTab === 'storage'">
      <div class="card">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
          <div>
            <h2>Sovereign S3 Object Storage</h2>
            <p style="color: var(--text-muted); font-size: 0.85rem;">Native Go S3-compatible engine with AWS Signature v4, range requests, and multipart upload.</p>
          </div>
          <button class="btn btn-primary" @click="showNewBucketModal = true">+ New Bucket</button>
        </div>
        <div class="grid-3" style="margin-top: 1rem;">
          <div v-for="b in buckets" :key="b.bucket_id" class="card" style="border-color: var(--border-hover);">
            <div style="display: flex; justify-content: space-between; align-items: flex-start;">
              <div>
                <h3 style="font-size: 1.1rem; color: var(--accent-blue);">{{ b.bucket_id }}</h3>
                <span class="stat-desc">Owner: {{ b.owner_id }}</span>
              </div>
              <span class="badge" :class="b.is_public ? 'badge-warning' : 'badge-success'">{{ b.is_public ? 'Public' : 'Private (SigV4)' }}</span>
            </div>
            <div style="margin-top: 1rem; font-size: 0.85rem; color: var(--text-muted);">
              Endpoint: <code>/s3/{{ b.bucket_id }}/</code>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- JAILED FILE BROWSER VIEW -->
    <main v-if="currentTab === 'files'">
      <div class="card">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
          <div>
            <h2>Jailed File Browser</h2>
            <p style="color: var(--text-muted); font-size: 0.85rem;">Kernel-hardened sandboxed explorer (openat2 RESOLVE_BENEATH) confined strictly to /opt/orbit.</p>
          </div>
          <div style="display: flex; gap: 0.5rem;">
            <select v-model="fileScope" class="form-select" style="width: 140px;" @change="loadFiles">
              <option value="app">Apps Scope</option>
              <option value="storage">Storage Scope</option>
            </select>
            <input v-model="fileTargetID" placeholder="App/Bucket ID..." class="form-input" style="width: 160px;" @change="loadFiles" />
          </div>
        </div>

        <div class="table-container">
          <table>
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
                <td>{{ item.is_dir ? '-' : formatBytes(item.size) }}</td>
                <td>{{ item.mod_time }}</td>
              </tr>
              <tr v-if="fileItems.length === 0">
                <td colspan="4" style="text-align: center; padding: 2rem; color: var(--text-muted);">Select or enter an App ID above to inspect jailed filesystem.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <!-- NEW DATABASE MODAL -->
    <div v-if="showNewDBModal" class="modal-backdrop">
      <div class="modal-content">
        <h3 style="margin-bottom: 1rem;">Provision Dedicated Database</h3>
        <div class="form-group">
          <label>Database Name</label>
          <input v-model="newDB.name" placeholder="e.g. production_db" class="form-input" />
        </div>
        <div class="form-group">
          <label>Engine</label>
          <select v-model="newDB.engine" class="form-select">
            <option value="postgres">PostgreSQL 16 (Native)</option>
            <option value="mongodb">MongoDB 7 (Native)</option>
          </select>
        </div>
        <div class="form-group">
          <label>RAM Tuning Profile</label>
          <select v-model="newDB.ram_profile" class="form-select">
            <option value="micro">Micro (&lt;32MB RAM, 20 max conn)</option>
            <option value="standard">Standard (&lt;128MB RAM, 50 max conn)</option>
            <option value="heavy">Heavy (&lt;512MB RAM, 100 max conn)</option>
          </select>
        </div>
        <div style="display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 1.5rem;">
          <button class="btn btn-secondary" @click="showNewDBModal = false">Cancel</button>
          <button class="btn btn-primary" @click="createDB">Provision</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

const currentTab = ref('dashboard')
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

const showNewDBModal = ref(false)
const newDB = ref({ name: '', engine: 'postgres', ram_profile: 'micro' })

const activeLogApp = ref(null)
const appLogs = ref([])
let logEventSource = null
let telemetryEventSource = null
const terminalRef = ref(null)

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

const fetchState = async () => {
  try {
    const [appsRes, dbsRes] = await Promise.all([
      fetch('/api/v1/apps'),
      fetch('/api/v1/databases'),
    ])
    if (appsRes.ok) apps.value = await appsRes.json()
    if (dbsRes.ok) databases.value = await dbsRes.json()
  } catch (err) {
    console.error('Fetch error:', err)
  }
}

const startTelemetrySSE = () => {
  telemetryEventSource = new EventSource('/api/v1/telemetry/stream')
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

  logEventSource = new EventSource(`/api/v1/apps/${appID}/logs/stream`)
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

const createDB = async () => {
  if (!newDB.value.name) return
  try {
    const res = await fetch('/api/v1/databases', {
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
    const res = await fetch(`/api/v1/jail/list?scope=${fileScope.value}&id=${fileTargetID.value}&path=`)
    if (res.ok) {
      const data = await res.json()
      fileItems.value = data.items || []
    }
  } catch (e) {}
}

onMounted(() => {
  fetchState()
  startTelemetrySSE()
  setInterval(fetchState, 5000)
})

onUnmounted(() => {
  if (telemetryEventSource) telemetryEventSource.close()
  if (logEventSource) logEventSource.close()
})
</script>
