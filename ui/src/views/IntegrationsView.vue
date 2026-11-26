<template>
  <div class="int-page">
    <div class="breadcrumbs">
      <span>Server</span>
      <span class="divider">/</span>
      <span class="current">Credentials</span>
    </div>

    <h1 class="page-title">Credentials and integrations</h1>
    <p class="page-subtitle">Keys for third-party services. Orbit stores them encrypted on this server and never shows them again.</p>

    <p v-if="error" class="int-error" role="alert">{{ error }}</p>
    <p v-if="notice" class="int-notice">{{ notice }}</p>

    <div v-if="loading" class="int-muted">Loading…</div>

    <div v-for="item in items" :key="item.provider" class="int-card">
      <div class="int-head">
        <strong class="int-name">{{ item.name }}</strong>
        <span :class="['int-status', item.configured ? 'on' : 'off']">
          {{ item.configured ? 'Saved' : 'Not set' }}
        </span>
      </div>
      <p class="int-muted">{{ item.used_for }}</p>

      <div v-if="item.configured" class="int-saved">
        <div>{{ item.secret_label }}: <code>{{ item.hint }}</code></div>
        <div v-if="item.extra_label">{{ item.extra_label }}: <code>{{ item.extra_value }}</code></div>
        <div class="int-muted">Updated {{ fmtDate(item.updated_at) }}</div>
      </div>

      <form class="int-form" @submit.prevent="save(item)">
        <label>
          {{ item.secret_label }}
          <input v-model="drafts[item.provider].secret" type="password" autocomplete="off" required
                 :placeholder="item.configured ? 'Enter a new value to replace it' : 'Paste here'" />
        </label>
        <label v-if="item.extra_label">
          {{ item.extra_label }}
          <input v-model="drafts[item.provider].extra" autocomplete="off" :required="item.provider === 'zerossl'" />
        </label>
        <div class="int-actions">
          <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy">
            {{ item.configured ? 'Replace' : 'Save' }}
          </button>
          <button v-if="item.configured" class="btn-cf btn-cf-secondary" type="button" :disabled="busy" @click="remove(item)">
            Remove
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { api } from '@/composables/api.js'

const items = ref([])
const loading = ref(true)
const busy = ref(false)
const error = ref('')
const notice = ref('')
const drafts = reactive({})

const fmtDate = (iso) => (iso ? new Date(iso).toLocaleString() : '')

const load = async () => {
  try {
    items.value = await api('/api/v1/integrations')
    for (const it of items.value) {
      if (!drafts[it.provider]) drafts[it.provider] = { secret: '', extra: '' }
    }
    error.value = ''
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

const save = async (item) => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    const d = drafts[item.provider]
    await api(`/api/v1/integrations/${item.provider}`, {
      method: 'PUT',
      body: { secret: d.secret, extra: d.extra },
    })
    d.secret = ''
    d.extra = ''
    notice.value = `${item.name} saved.`
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const remove = async (item) => {
  if (!window.confirm(`Remove the saved ${item.name} credential? Certificates that depend on it keep working until they expire.`)) return
  busy.value = true
  try {
    await api(`/api/v1/integrations/${item.provider}`, { method: 'DELETE' })
    notice.value = `${item.name} removed.`
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.int-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
.int-card { background: var(--kumo-base, #fff); border: 1px solid var(--kumo-border, #e2e4e8); border-radius: 10px; padding: 1rem; box-sizing: border-box; min-width: 0; display: flex; flex-direction: column; gap: 0.6rem; }
.int-head { display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; gap: 0.5rem; }
.int-name { font-size: 16px; }
.int-status { font-size: 12px; font-weight: 600; padding: 2px 8px; border-radius: 999px; background: #e5e7eb; color: #374151; }
.int-status.on { background: #dcfce7; color: #166534; }
.int-saved { font-size: 13px; display: flex; flex-direction: column; gap: 2px; overflow-wrap: anywhere; }
.int-saved code { font-size: 12px; }
.int-form { display: flex; flex-direction: column; gap: 0.6rem; }
.int-form label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; min-width: 0; }
.int-form input { padding: 9px 10px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; background: #fff; color: inherit; width: 100%; box-sizing: border-box; }
.int-actions { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.int-muted { font-size: 13px; color: #6b7280; margin: 0; overflow-wrap: anywhere; }
.int-error { color: #b91c1c; font-size: 13px; margin: 0; overflow-wrap: anywhere; }
.int-notice { color: #166534; font-size: 13px; margin: 0; }
@media (max-width: 680px) {
  .int-actions .btn-cf { flex: 1 1 100%; }
}
</style>
