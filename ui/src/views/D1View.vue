<template>
  <div class="db-page">
    <div class="breadcrumbs">
      <span>Database</span>
      <span class="divider">/</span>
      <span class="current">Databases</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Databases</h1>
        <p class="page-subtitle">Dedicated database instances on this server. Each one runs isolated, with its own password and data folder.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showCreate = !showCreate">
        {{ showCreate ? 'Close' : '+ Create database' }}
      </button>
    </div>

    <p v-if="error" class="db-error" role="alert">{{ error }}</p>

    <!-- Create form -->
    <form v-if="showCreate" class="db-card db-form" @submit.prevent="create">
      <label>
        Name
        <input v-model="form.name" required maxlength="63" placeholder="my-app-db" />
      </label>
      <label>
        Engine
        <select v-model="form.engine">
          <option v-for="e in engines" :key="e.value" :value="e.value">{{ e.label }}</option>
        </select>
      </label>
      <label>
        Size
        <select v-model="form.ram_profile">
          <option v-for="p in profiles" :key="p.value" :value="p.value">{{ p.label }}</option>
        </select>
      </label>
      <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy">
        {{ busy ? 'Creating…' : 'Create database' }}
      </button>
      <div v-if="created" class="db-card db-created">
        <strong>Saved. Copy this password now, it is shown only here.</strong>
        <code class="db-uri">{{ created.password }}</code>
        <code class="db-uri">{{ created.tcp_uri }}</code>
      </div>
    </form>

    <p v-if="loading" class="db-muted">Loading…</p>
    <p v-else-if="!databases.length" class="db-muted">No databases yet. Create one above.</p>

    <!-- Cards work the same on phones and desktops -->
    <div class="db-list">
      <router-link
        v-for="d in databases"
        :key="d.db_id"
        :to="`/d1/${d.db_id}`"
        class="db-card db-row"
      >
        <div class="db-row-main">
          <strong class="db-name">{{ d.name }}</strong>
          <span class="db-engine">{{ d.engine }} {{ d.version }}</span>
        </div>
        <div class="db-row-meta">
          <span :class="['db-status', d.status]">{{ d.status }}</span>
          <span class="db-muted">{{ d.ram_profile }} · port {{ d.port }}</span>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { api } from '@/composables/api.js'

const databases = ref([])
const loading = ref(true)
const error = ref('')
const showCreate = ref(false)
const busy = ref(false)
const created = ref(null)

const engines = [
  { value: 'postgres', label: 'PostgreSQL' },
  { value: 'mysql', label: 'MySQL' },
  { value: 'mariadb', label: 'MariaDB' },
  { value: 'mongodb', label: 'MongoDB' },
  { value: 'redis', label: 'Redis' },
  { value: 'libsql', label: 'LibSQL (SQLite)' },
]
const profiles = [
  { value: 'micro', label: 'Micro (about 32 MB)' },
  { value: 'standard', label: 'Standard (about 128 MB)' },
  { value: 'heavy', label: 'Heavy (about 512 MB)' },
  { value: 'unlimited', label: 'Unlimited (no cap)' },
]
const form = ref({ name: '', engine: 'postgres', ram_profile: 'micro' })

const load = async () => {
  try {
    databases.value = await api('/api/v1/databases')
    error.value = ''
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

const create = async () => {
  busy.value = true
  error.value = ''
  created.value = null
  try {
    const res = await api('/api/v1/databases', { method: 'POST', body: form.value })
    created.value = res
    form.value.name = ''
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

let timer
onMounted(() => {
  load()
  timer = setInterval(load, 5000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.db-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
.db-card {
  background: var(--kumo-base, #fff);
  border: 1px solid var(--kumo-border, #e2e4e8);
  border-radius: 10px;
  padding: 1rem;
  box-sizing: border-box;
  min-width: 0;
}
.db-form { display: flex; flex-direction: column; gap: 0.75rem; }
.db-form label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
.db-form input, .db-form select {
  padding: 9px 10px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px;
  background: #fff; color: inherit; width: 100%; box-sizing: border-box;
}
.db-list { display: flex; flex-direction: column; gap: 0.6rem; }
.db-row { display: flex; flex-direction: column; gap: 0.4rem; text-decoration: none; color: inherit; }
.db-row-main, .db-row-meta { display: flex; flex-wrap: wrap; gap: 0.5rem 0.75rem; align-items: center; justify-content: space-between; }
.db-name { overflow-wrap: anywhere; }
.db-engine, .db-muted { font-size: 13px; color: #6b7280; }
.db-status { font-size: 12px; font-weight: 600; padding: 2px 8px; border-radius: 999px; background: #e5e7eb; color: #374151; }
.db-status.running { background: #dcfce7; color: #166534; }
.db-status.stopped { background: #fee2e2; color: #991b1b; }
.db-error { color: #b91c1c; font-size: 13px; margin: 0; }
.db-created { display: flex; flex-direction: column; gap: 0.5rem; background: #fffbeb; }
.db-uri { font-size: 12px; word-break: break-all; white-space: pre-wrap; }
</style>
