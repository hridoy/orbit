<template>
  <div class="db-page">
    <div class="breadcrumbs">
      <router-link to="/d1">Databases</router-link>
      <span class="divider">/</span>
      <span class="current">{{ info?.database?.name || id }}</span>
    </div>

    <p v-if="error" class="db-error" role="alert">{{ error }}</p>
    <p v-if="loading" class="db-muted">Loading…</p>

    <template v-if="info">
      <div class="page-title-row">
        <div class="db-title">
          <h1 class="page-title">{{ info.database.name }}</h1>
          <span :class="['db-status', info.database.status]">{{ info.database.status }}</span>
        </div>
      </div>

      <!-- Actions -->
      <div class="db-actions">
        <button class="btn-cf btn-cf-primary" :disabled="busy" @click="run('start')">Start</button>
        <button class="btn-cf btn-cf-secondary" :disabled="busy" @click="run('restart')">Restart</button>
        <button class="btn-cf btn-cf-secondary" :disabled="busy" @click="run('stop')">Stop</button>
        <button class="btn-cf btn-cf-secondary" :disabled="busy" @click="run('backup')">Backup now</button>
        <button class="btn-cf btn-cf-danger" :disabled="busy" @click="remove">Delete</button>
      </div>
      <p v-if="notice" class="db-notice">{{ notice }}</p>

      <!-- Details -->
      <div class="db-card">
        <h2 class="db-h2">Connection</h2>
        <dl class="db-dl">
          <dt>Engine</dt><dd>{{ info.database.engine }} {{ info.database.version }}</dd>
          <dt>Admin user</dt><dd>{{ info.database.admin_user }}</dd>
          <dt>Port (this server only)</dt><dd>{{ info.database.port }}</dd>
          <dt>Size</dt><dd>{{ info.database.ram_profile }}</dd>
          <dt>Password</dt><dd><code class="db-uri">{{ info.password }}</code></dd>
          <dt>Connection string</dt><dd><code class="db-uri">{{ info.tcp_uri }}</code></dd>
        </dl>
        <p class="db-muted">The database listens on 127.0.0.1 only. Apps on this server connect with the string above.</p>
      </div>

      <!-- Logs -->
      <div class="db-card">
        <h2 class="db-h2">Logs</h2>
        <pre class="db-logs">{{ logs.length ? logs.join('\n') : 'No log lines yet.' }}</pre>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '@/composables/api.js'

const route = useRoute()
const router = useRouter()
const id = computed(() => route.params.id)

const info = ref(null)
const logs = ref([])
const loading = ref(true)
const busy = ref(false)
const error = ref('')
const notice = ref('')

const load = async () => {
  try {
    info.value = await api(`/api/v1/databases/${id.value}`)
    error.value = ''
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

const loadLogs = async () => {
  try {
    const res = await api(`/api/v1/databases/${id.value}/logs?tail=200`)
    logs.value = res.lines || []
  } catch (e) {}
}

const run = async (action) => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    const res = await api(`/api/v1/databases/${id.value}/${action}`, { method: 'POST' })
    notice.value = res.path ? `Backup saved at ${res.path}` : `Done: ${action}`
    await load()
    await loadLogs()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const remove = async () => {
  if (!window.confirm(`Delete ${info.value?.database?.name}? Its data is removed and cannot be recovered.`)) return
  busy.value = true
  try {
    await api(`/api/v1/databases/${id.value}`, { method: 'DELETE' })
    router.push('/d1')
  } catch (e) {
    error.value = e.message
    busy.value = false
  }
}

let timer
onMounted(() => {
  load()
  loadLogs()
  timer = setInterval(() => { load(); loadLogs() }, 5000)
})
onUnmounted(() => clearInterval(timer))
watch(id, () => { loading.value = true; load(); loadLogs() })
</script>

<style scoped>
.db-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
.db-title { display: flex; flex-wrap: wrap; align-items: center; gap: 0.6rem; min-width: 0; }
.db-title .page-title { overflow-wrap: anywhere; }
.db-actions { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.db-card {
  background: var(--kumo-base, #fff);
  border: 1px solid var(--kumo-border, #e2e4e8);
  border-radius: 10px;
  padding: 1rem;
  box-sizing: border-box;
  min-width: 0;
}
.db-h2 { font-size: 15px; margin: 0 0 0.75rem; }
.db-dl { display: grid; grid-template-columns: minmax(0, 170px) minmax(0, 1fr); gap: 0.5rem 1rem; margin: 0; }
.db-dl dt { color: #6b7280; font-size: 13px; }
.db-dl dd { margin: 0; min-width: 0; }
.db-uri { font-size: 12px; word-break: break-all; white-space: pre-wrap; display: block; }
.db-logs {
  background: #0f1115; color: #d1d5db; padding: 0.75rem; border-radius: 8px;
  font-size: 12px; max-height: 320px; overflow: auto; white-space: pre-wrap; word-break: break-word; margin: 0;
}
.db-status { font-size: 12px; font-weight: 600; padding: 2px 8px; border-radius: 999px; background: #e5e7eb; color: #374151; }
.db-status.running { background: #dcfce7; color: #166534; }
.db-status.stopped { background: #fee2e2; color: #991b1b; }
.db-muted { font-size: 13px; color: #6b7280; margin: 0.5rem 0 0; }
.db-error { color: #b91c1c; font-size: 13px; margin: 0; }
.db-notice { font-size: 13px; color: #166534; margin: 0; overflow-wrap: anywhere; }

@media (max-width: 680px) {
  .db-dl { grid-template-columns: minmax(0, 1fr); gap: 0.2rem; }
  .db-dl dt { margin-top: 0.5rem; }
  .db-actions .btn-cf { flex: 1 1 calc(50% - 0.25rem); }
}
</style>
