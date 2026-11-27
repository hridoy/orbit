<template>
  <div class="app-detail-page">
    <div class="breadcrumbs">
      <router-link to="/apps">Apps</router-link>
      <span class="divider">/</span>
      <span class="current">{{ app.name || appId }}</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap;">
          <h1 class="page-title">{{ app.name || appId }}</h1>
          <CfBadge :variant="app.status === 'running' ? 'success' : 'neutral'" dot>
            {{ app.status || 'stopped' }}
          </CfBadge>
          <CfBadge variant="info" size="sm">Port :{{ app.active_port || 8001 }}</CfBadge>
        </div>
        <p class="page-subtitle">
          App ID: <code>{{ app.app_id || appId }}</code>
          <span v-if="app.domain"> · Domain: <strong>{{ app.domain }}</strong></span>
        </p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <a v-if="app.domain" :href="`http://${app.domain}`" target="_blank" style="text-decoration: none;">
          <CfButton variant="secondary">Visit App ↗</CfButton>
        </a>
        <CfButton variant="primary" icon="refresh" @click="fetchApp">Refresh</CfButton>
      </div>
    </div>

    <CfTabs
      v-model="activeTab"
      :tabs="[
        { label: 'Overview', value: 'overview' },
        { label: 'Live Logs', value: 'logs' },
        { label: 'Settings', value: 'settings' }
      ]"
    />

    <!-- Overview Tab -->
    <div v-if="activeTab === 'overview'" style="display: flex; flex-direction: column; gap: 1.5rem;">
      <div class="stats-overview-grid">
        <CfStatCard
          label="Status"
          :value="app.status || 'stopped'"
          :delta="app.auto_restart ? 'Auto-restart' : 'Manual'"
          :delta-type="app.status === 'running' ? 'up' : 'neutral'"
        />
        <CfStatCard
          label="Active Port"
          :value="`:${app.active_port || 8001}`"
          :delta="`Idle :${app.idle_port || 8002}`"
          delta-type="neutral"
        />
        <CfStatCard
          label="RAM Limit"
          :value="`${app.ram_limit_mb || 256} MB`"
          :delta="stats?.memory_mb ? `${stats.memory_mb} MB used` : 'cgroups'"
          delta-type="up"
        />
        <CfStatCard
          label="CPU Quota"
          :value="`${app.cpu_quota_pct || 100}%`"
          :delta="`PIDs max: ${app.pids_limit || 128}`"
          delta-type="neutral"
        />
      </div>

      <CfCard title="Runtime & Deployment Configuration">
        <div class="app-kv-grid">
          <div class="app-kv-row">
            <span class="app-kv-label">Binary Entrypoint</span>
            <code class="app-kv-val">{{ app.binary_entrypoint || './server' }}</code>
          </div>
          <div class="app-kv-row">
            <span class="app-kv-label">Health Check Path</span>
            <code class="app-kv-val">{{ app.health_check_path || '/healthz' }}</code>
          </div>
          <div class="app-kv-row">
            <span class="app-kv-label">Current Version</span>
            <code class="app-kv-val">{{ app.current_version || 'Not deployed yet' }}</code>
          </div>
          <div class="app-kv-row">
            <span class="app-kv-label">Node ID</span>
            <code class="app-kv-val">{{ app.node_id || 'local-node' }}</code>
          </div>
        </div>
      </CfCard>
    </div>

    <!-- Live Logs Tab -->
    <div v-else-if="activeTab === 'logs'">
      <CfCard title="Application Logs" subtitle="Real-time stdout/stderr stream from the Orbit process runner">
        <template #actions>
          <div style="display: flex; gap: 0.5rem;">
            <CfButton variant="secondary" size="sm" icon="refresh" @click="fetchLogs">Reload</CfButton>
            <CfButton variant="outline" size="sm" @click="logs = []">Clear</CfButton>
          </div>
        </template>
        <div class="app-log-terminal">
          <div v-if="!logs.length" style="color: #94a3b8;">No logs recorded yet for {{ appId }}.</div>
          <div v-for="(line, i) in logs" :key="i" class="app-log-line">{{ line }}</div>
        </div>
      </CfCard>
    </div>

    <!-- Settings Tab -->
    <div v-else>
      <CfCard title="Process & Resource Limits" subtitle="Configured isolation and health check settings for this app">
        <div class="app-kv-grid">
          <div class="app-kv-row">
            <span class="app-kv-label">Exec Args</span>
            <code class="app-kv-val">{{ app.exec_args || '(none)' }}</code>
          </div>
          <div class="app-kv-row">
            <span class="app-kv-label">Health Timeout</span>
            <span class="app-kv-val">{{ app.health_timeout_sec || 10 }}s</span>
          </div>
          <div class="app-kv-row">
            <span class="app-kv-label">Linux User</span>
            <code class="app-kv-val">{{ app.linux_user || 'orbit' }}</code>
          </div>
        </div>
      </CfCard>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const appId = ref(route.params.id)
const activeTab = ref('overview')
const app = ref({})
const stats = ref(null)
const logs = ref([])
let sse = null

function authHeaders() {
  const token = localStorage.getItem('orbit_token') || ''
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function fetchApp() {
  try {
    const res = await fetch(`/api/v1/apps/${appId.value}`, { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      app.value = data.app || data
      stats.value = data.stats || null
    }
  } catch (e) {}
}

async function fetchLogs() {
  try {
    const res = await fetch(`/api/v1/apps/${appId.value}/logs?tail=100`, { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      logs.value = Array.isArray(data.logs) ? data.logs : (typeof data === 'string' ? data.split('\n') : [])
    }
  } catch (e) {}
}

function startLogStream() {
  if (sse) sse.close()
  const token = localStorage.getItem('orbit_token') || ''
  const url = token ? `/api/v1/apps/${appId.value}/logs/stream?token=${token}` : `/api/v1/apps/${appId.value}/logs/stream`
  sse = new EventSource(url)
  sse.onmessage = (ev) => {
    if (ev.data) logs.value.push(ev.data)
  }
}

watch(activeTab, (tab) => {
  if (tab === 'logs') {
    fetchLogs()
    startLogStream()
  } else if (sse) {
    sse.close()
    sse = null
  }
})

onMounted(fetchApp)
onBeforeUnmount(() => {
  if (sse) sse.close()
})
</script>

<style scoped>
.app-detail-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.app-kv-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.app-kv-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 0.6rem;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 13.5px;
}

.app-kv-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.app-kv-label {
  color: var(--text-color-kumo-subtle);
  font-weight: 500;
}

.app-kv-val {
  color: var(--text-color-kumo-strong);
  font-family: var(--font-mono);
  font-size: 12.5px;
}

.app-log-terminal {
  background: #0b0f17;
  color: #e2e8f0;
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  padding: 1rem;
  border-radius: 6px;
  min-height: 260px;
  max-height: 460px;
  overflow-y: auto;
}

.app-log-line {
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
