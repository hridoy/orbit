<template>
  <div class="cf-metrics-page">
    <div class="breadcrumbs">
      <router-link to="/">System</router-link>
      <span class="divider">/</span>
      <span class="current">Metrics</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.65rem; flex-wrap: wrap;">
          <h1 class="page-title">System Metrics</h1>
          <CfBadge :variant="streamConnected ? 'success' : 'neutral'" size="sm" dot>
            {{ streamConnected ? 'Live Stream' : 'Polling' }}
          </CfBadge>
        </div>
        <p class="page-subtitle">
          Real-time CPU, memory, and runtime telemetry from <code>{{ nodeInfo.hostname || 'this VPS node' }}</code> via SSE.
        </p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <CfButton variant="secondary" icon="refresh" @click="fetchSnapshot">Refresh</CfButton>
      </div>
    </div>

    <!-- 4 Live Telemetry Stat Cards -->
    <div class="stats-overview-grid">
      <CfStatCard
        label="CPU Usage"
        :value="`${latest.cpu_usage_pct.toFixed(1)}%`"
        delta="Host CPU"
        delta-type="up"
        :sparkline-data="cpuHistory.map(p => p.value)"
      />
      <CfStatCard
        label="Memory Used"
        :value="formatBytes(latest.ram_used_bytes)"
        :delta="`of ${formatBytes(latest.ram_total_bytes)}`"
        delta-type="neutral"
        :sparkline-data="ramHistory.map(p => p.value)"
      />
      <CfStatCard
        label="Agent Heap (Go)"
        :value="formatBytes(latest.agent_alloc_ram)"
        :delta="`${latest.agent_num_gc} GC cycles`"
        delta-type="up"
      />
      <CfStatCard
        label="Active Goroutines"
        :value="latest.num_goroutines"
        delta="Runtime"
        delta-type="neutral"
      />
    </div>

    <!-- 2-Column Live Area Charts (Copied CfAreaChart from cloudflare-ui) -->
    <div class="cf-metrics-charts-grid">
      <CfAreaChart
        title="CPU Utilization (%)"
        subtitle="Real-time 2s interval stream"
        :data="cpuHistory"
        color="#f6821f"
        unit="%"
      />
      <CfAreaChart
        title="Memory Allocated (MB)"
        subtitle="System & process memory usage"
        :data="ramHistory"
        color="#0051c3"
        unit=" MB"
      />
    </div>

    <!-- Host & Runtime Details Card -->
    <CfCard title="Node & Runtime Details">
      <div class="cf-metrics-kv-list">
        <div class="cf-metrics-kv-row">
          <span class="cf-metrics-kv-label">Node ID</span>
          <code>{{ nodeInfo.node_id || 'local-node' }}</code>
        </div>
        <div class="cf-metrics-kv-row">
          <span class="cf-metrics-kv-label">Hostname</span>
          <span>{{ nodeInfo.hostname || 'localhost' }}</span>
        </div>
        <div class="cf-metrics-kv-row">
          <span class="cf-metrics-kv-label">OS / Architecture</span>
          <span>{{ nodeInfo.os || 'linux' }} / {{ nodeInfo.arch || 'amd64' }}</span>
        </div>
        <div class="cf-metrics-kv-row">
          <span class="cf-metrics-kv-label">Free Memory</span>
          <span>{{ formatBytes(latest.ram_free_bytes) }}</span>
        </div>
      </div>
    </CfCard>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'

const nodeInfo = ref({})
const streamConnected = ref(false)
const cpuHistory = ref([])
const ramHistory = ref([])

const latest = reactive({
  cpu_usage_pct: 0,
  ram_total_bytes: 0,
  ram_used_bytes: 0,
  ram_free_bytes: 0,
  agent_alloc_ram: 0,
  agent_num_gc: 0,
  num_goroutines: 0
})

let sse = null

function authHeaders() {
  const token = localStorage.getItem('orbit_token') || ''
  return token ? { Authorization: `Bearer ${token}` } : {}
}

function pushSample(m) {
  latest.cpu_usage_pct = Number(m.cpu_usage_pct || 0)
  latest.ram_total_bytes = Number(m.ram_total_bytes || 0)
  latest.ram_used_bytes = Number(m.ram_used_bytes || 0)
  latest.ram_free_bytes = Number(m.ram_free_bytes || 0)
  latest.agent_alloc_ram = Number(m.agent_alloc_ram || 0)
  latest.agent_num_gc = Number(m.agent_num_gc || 0)
  latest.num_goroutines = Number(m.num_goroutines || 0)

  const label = new Date(m.timestamp || Date.now()).toLocaleTimeString([], { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' })
  cpuHistory.value = [...cpuHistory.value.slice(-19), { label, value: Number(latest.cpu_usage_pct.toFixed(1)) }]
  const ramMb = Number((latest.ram_used_bytes / (1024 * 1024)).toFixed(1))
  ramHistory.value = [...ramHistory.value.slice(-19), { label, value: ramMb }]
}

async function fetchSnapshot() {
  try {
    const [nodeRes, metRes] = await Promise.all([
      fetch('/api/v1/node', { headers: authHeaders() }),
      fetch('/api/v1/telemetry/metrics', { headers: authHeaders() })
    ])
    if (nodeRes.ok) nodeInfo.value = await nodeRes.json()
    if (metRes.ok) {
      const m = await metRes.json()
      pushSample(m)
    }
  } catch (e) {}
}

function startStream() {
  const token = localStorage.getItem('orbit_token') || ''
  const url = token ? `/api/v1/telemetry/stream?token=${token}` : '/api/v1/telemetry/stream'
  sse = new EventSource(url)
  sse.onopen = () => {
    streamConnected.value = true
  }
  sse.onmessage = (ev) => {
    try {
      const data = JSON.parse(ev.data)
      streamConnected.value = true
      pushSample(data)
    } catch (e) {}
  }
  sse.onerror = () => {
    streamConnected.value = false
  }
}

function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

onMounted(() => {
  fetchSnapshot()
  startStream()
})

onBeforeUnmount(() => {
  if (sse) sse.close()
})
</script>

<style scoped>
.cf-metrics-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.cf-metrics-charts-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1.25rem;
}

@media (max-width: 900px) {
  .cf-metrics-charts-grid {
    grid-template-columns: 1fr;
  }
}

.cf-metrics-kv-list {
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}

.cf-metrics-kv-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 0.55rem;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 13px;
}

.cf-metrics-kv-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.cf-metrics-kv-label {
  color: var(--text-color-kumo-subtle);
  font-weight: 500;
}
</style>
