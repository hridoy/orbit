<template>
  <div class="mem-page">
    <div class="breadcrumbs">
      <span>Manage Account</span>
      <span class="divider">/</span>
      <span class="current">Members</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Members & roles</h1>
        <p class="page-subtitle">People who can sign in to this dashboard. Admins manage everything, operators deploy, and viewers can only look.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showInvite = !showInvite">
        {{ showInvite ? 'Close' : '+ Add member' }}
      </button>
    </div>

    <p v-if="error" class="mem-error" role="alert">{{ error }}</p>
    <p v-if="notice" class="mem-notice">{{ notice }}</p>

    <!-- One-time password is shown only here -->
    <div v-if="issued" class="mem-card mem-issued">
      <strong>{{ issued.username }} was added as {{ issued.role }}.</strong>
      <span class="mem-muted">Share this password now. It is not shown again.</span>
      <code class="mem-code">{{ issued.temporary_password }}</code>
      <button class="btn-cf btn-cf-secondary" @click="issued = null">Done</button>
    </div>

    <form v-if="showInvite" class="mem-card mem-form" @submit.prevent="invite">
      <label>
        Username
        <input v-model="form.username" required minlength="2" maxlength="64" autocomplete="off" placeholder="e.g. ravi" />
      </label>
      <label>
        Role
        <select v-model="form.role">
          <option value="admin">Admin: full access, manages members</option>
          <option value="operator">Operator: deploys apps and databases</option>
          <option value="viewer">Viewer: read only</option>
        </select>
      </label>
      <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy">{{ busy ? 'Adding…' : 'Add member' }}</button>
    </form>

    <p v-if="loading" class="mem-muted">Loading…</p>

    <div class="mem-list">
      <div v-for="m in members" :key="m.username" class="mem-card mem-row">
        <div class="mem-row-main">
          <strong class="mem-name">{{ m.username }}</strong>
          <span :class="['mem-role', m.role]">{{ m.role }}</span>
          <span v-if="m.disabled" class="mem-role disabled">disabled</span>
        </div>
        <div class="mem-muted">Added {{ m.created_at }}</div>
        <div class="mem-actions">
          <select :value="m.role" :disabled="busy" @change="changeRole(m, $event.target.value)" aria-label="Change role">
            <option value="admin">Admin</option>
            <option value="operator">Operator</option>
            <option value="viewer">Viewer</option>
          </select>
          <button class="btn-cf btn-cf-secondary" :disabled="busy" @click="toggleDisabled(m)">
            {{ m.disabled ? 'Enable' : 'Disable' }}
          </button>
          <button class="btn-cf btn-cf-secondary" :disabled="busy" @click="remove(m)">Remove</button>
        </div>
      </div>
    </div>

    <form class="mem-card mem-form" @submit.prevent="changeOwnPassword">
      <strong>Change my password</strong>
      <label>Current password <input v-model="pw.current" type="password" required autocomplete="current-password" /></label>
      <label>New password (8+ characters) <input v-model="pw.new" type="password" required minlength="8" autocomplete="new-password" /></label>
      <button class="btn-cf btn-cf-secondary" type="submit" :disabled="busy">Change password</button>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/composables/api.js'

const members = ref([])
const loading = ref(true)
const busy = ref(false)
const error = ref('')
const notice = ref('')
const showInvite = ref(false)
const issued = ref(null)
const form = ref({ username: '', role: 'viewer' })
const pw = ref({ current: '', new: '' })

const load = async () => {
  try {
    members.value = await api('/api/v1/users')
    error.value = ''
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

const run = async (fn, done) => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    await fn()
    if (done) notice.value = done
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const invite = () => run(async () => {
  issued.value = await api('/api/v1/users', { method: 'POST', body: form.value })
  form.value = { username: '', role: 'viewer' }
  showInvite.value = false
})

const changeRole = (m, role) => run(
  () => api(`/api/v1/users/${encodeURIComponent(m.username)}`, { method: 'PUT', body: { role } }),
  `${m.username} is now ${role}.`,
)

const toggleDisabled = (m) => run(
  () => api(`/api/v1/users/${encodeURIComponent(m.username)}`, { method: 'PUT', body: { disabled: !m.disabled } }),
  m.disabled ? `${m.username} can sign in again.` : `${m.username} is disabled and signed out.`,
)

const remove = (m) => {
  if (!window.confirm(`Remove ${m.username}? They lose access right away.`)) return
  run(
    () => api(`/api/v1/users/${encodeURIComponent(m.username)}`, { method: 'DELETE' }),
    `${m.username} was removed.`,
  )
}

const changeOwnPassword = () => run(async () => {
  await api('/api/v1/users/me/password', { method: 'POST', body: pw.value })
  pw.value = { current: '', new: '' }
}, 'Password changed.')

onMounted(load)
</script>

<style scoped>
.mem-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
.mem-card { background: var(--kumo-base, #fff); border: 1px solid var(--kumo-border, #e2e4e8); border-radius: 10px; padding: 1rem; box-sizing: border-box; min-width: 0; display: flex; flex-direction: column; gap: 0.6rem; }
.mem-form label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; min-width: 0; }
.mem-form input, .mem-form select, .mem-actions select {
  padding: 9px 10px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; background: #fff; color: inherit; width: 100%; box-sizing: border-box;
}
.mem-actions select { width: auto; flex: 0 1 auto; }
.mem-list { display: flex; flex-direction: column; gap: 0.6rem; }
.mem-row-main { display: flex; flex-wrap: wrap; gap: 0.5rem; align-items: center; }
.mem-name { overflow-wrap: anywhere; }
.mem-role { font-size: 11px; font-weight: 600; padding: 2px 8px; border-radius: 999px; background: #e5e7eb; color: #374151; text-transform: capitalize; }
.mem-role.admin { background: #fee2e2; color: #991b1b; }
.mem-role.operator { background: #dbeafe; color: #1e40af; }
.mem-role.viewer { background: #f3f4f6; color: #374151; }
.mem-role.disabled { background: #e5e7eb; color: #6b7280; }
.mem-actions { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.mem-issued { background: #fffbeb; }
.mem-code { font-size: 14px; word-break: break-all; }
.mem-muted { font-size: 13px; color: #6b7280; margin: 0; overflow-wrap: anywhere; }
.mem-error { color: #b91c1c; font-size: 13px; margin: 0; overflow-wrap: anywhere; }
.mem-notice { color: #166534; font-size: 13px; margin: 0; }
@media (max-width: 680px) {
  .mem-actions .btn-cf, .mem-actions select { flex: 1 1 100%; width: 100%; }
}
</style>
