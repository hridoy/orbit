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
        <select v-model="form.engine" @change="pickEngine">
          <option v-for="e in engines" :key="e.engine" :value="e.engine">{{ e.name }}</option>
        </select>
      </label>
      <label>
        Version
        <select v-model="form.version" :disabled="!currentVersions.length">
          <option v-for="v in currentVersions" :key="v.version" :value="v.version">
            {{ v.version }}{{ v.status === 'installed' ? '' : ' (' + statusLabel(v) + ')' }}
          </option>
        </select>
        <span v-if="selectedVersion && selectedVersion.status !== 'installed'" class="db-install">
          <span class="db-muted">{{ selectedVersion.status === 'installing' ? 'Installing now. This can take several minutes.' : selectedVersion.error || 'Not installed on this server yet.' }}</span>
          <button v-if="selectedVersion.status !== 'installing'" type="button" class="btn-cf btn-cf-secondary" :disabled="installing" @click="installSelected">
            Install {{ selectedVersion.version }}
          </button>
        </span>
      </label>
      <label>
        Username
        <input v-model="form.username" maxlength="32" autocomplete="off"
               placeholder="Leave empty to generate one" />
      </label>
      <label>
        Password
        <span class="db-pass">
          <input v-model="form.password" :type="showPass ? 'text' : 'password'" maxlength="128" autocomplete="new-password"
                 placeholder="Leave empty to generate one" />
          <button type="button" class="db-eye" :aria-label="showPass ? 'Hide password' : 'Show password'" @click="showPass = !showPass">
            {{ showPass ? 'Hide' : 'Show' }}
          </button>
        </span>
      </label>
      <label>
        Size
        <select v-model="form.ram_profile">
          <option v-for="p in profiles" :key="p.value" :value="p.value">{{ p.label }}</option>
        </select>
      </label>
      <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy || !selectedInstalled">
        {{ busy ? 'Creating…' : 'Create database' }}
      </button>
      <div v-if="created" class="db-card db-created">
        <strong>Saved. Copy these now, the password is shown only here.</strong>
        <div class="db-kv"><span>Username</span><SecretField :value="created.database.admin_user" /></div>
        <div class="db-kv"><span>Password</span><SecretField :value="created.password" :revealed="false" /></div>
        <div class="db-kv"><span>Connection</span><SecretField :value="created.tcp_uri" :revealed="false" /></div>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { api } from '@/composables/api.js'
import SecretField from '@/components/SecretField.vue'

const databases = ref([])
const loading = ref(true)
const error = ref('')
const showCreate = ref(false)
const busy = ref(false)
const showPass = ref(false)
const created = ref(null)

const engines = ref([])
const profiles = [
  { value: 'micro', label: 'Micro (about 32 MB)' },
  { value: 'standard', label: 'Standard (about 128 MB)' },
  { value: 'heavy', label: 'Heavy (about 512 MB)' },
  { value: 'unlimited', label: 'Unlimited (no cap)' },
]
const form = ref({ name: '', engine: '', version: '', username: '', password: '', ram_profile: 'micro' })
const currentVersions = computed(() => engines.value.find((e) => e.engine === form.value.engine)?.versions || [])
const selectedVersion = computed(() => currentVersions.value.find((v) => v.version === form.value.version) || null)
const selectedInstalled = computed(() => !!selectedVersion.value && selectedVersion.value.status === 'installed')
const installing = ref(false)

const statusLabel = (v) => ({ installing: 'installing', failed: 'install failed', available: 'not installed' }[v.status] || v.status)

const pickEngine = () => {
  const first = currentVersions.value.find((v) => v.status === 'installed') || currentVersions.value[0]
  form.value.version = first ? first.version : ''
}

const installSelected = async () => {
  installing.value = true
  error.value = ''
  try {
    await api(`/api/v1/databases/engines/${form.value.engine}/${form.value.version}/install`, { method: 'POST' })
    await loadEngines()
  } catch (e) {
    error.value = e.message
  } finally {
    installing.value = false
  }
}

const loadEngines = async () => {
  try {
    engines.value = await api('/api/v1/databases/engines')
    if (!form.value.engine && engines.value.length) {
      form.value.engine = engines.value[0].engine
      pickEngine()
    }
  } catch (e) {
    error.value = e.message
  }
}

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
    const body = { ...form.value }
    if (!body.username) delete body.username
    if (!body.password) delete body.password
    const res = await api('/api/v1/databases', { method: 'POST', body })
    created.value = res
    form.value.name = ''
    form.value.username = ''
    form.value.password = ''
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

let timer
onMounted(() => {
  loadEngines()
  load()
  timer = setInterval(load, 5000)
})
onUnmounted(() => { clearInterval(timer); clearInterval(enginesTimer) })
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
.db-form select:disabled { opacity: 0.6; }
.db-install { display: flex; flex-wrap: wrap; gap: 0.5rem; align-items: center; }
.db-pass { display: flex; gap: 0.5rem; }
.db-pass input { flex: 1 1 auto; min-width: 0; }
.db-eye { padding: 0 12px; border: 1px solid #d1d5db; border-radius: 6px; background: #fff; color: inherit; cursor: pointer; font-size: 13px; }
.db-kv { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
.db-kv > span:first-child { color: #6b7280; }
</style>
