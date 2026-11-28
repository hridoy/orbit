<template>
  <div class="srv-page">
    <div class="breadcrumbs">
      <span>Server</span>
      <span class="divider">/</span>
      <span class="current">Settings</span>
    </div>

    <h1 class="page-title">Server settings</h1>
    <p class="page-subtitle">Set a domain for this dashboard so you can open it with a name instead of the server IP.</p>

    <p v-if="error" class="srv-error" role="alert">{{ error }}</p>
    <p v-if="notice" class="srv-notice">{{ notice }}</p>

    <form class="srv-card srv-form" @submit.prevent="save">
      <CfInput v-model="form.panel_domain" label="Panel domain" placeholder="panel.example.com" @change="detect">
        <template #suffix>
          <button class="btn-cf btn-cf-secondary" type="button" :disabled="!form.panel_domain || checking" @click="detect">
            {{ checking ? 'Checking…' : 'Check' }}
          </button>
        </template>
      </CfInput>

      <div v-if="detection" class="srv-detect">
        <span v-if="detection.cloudflare">Cloudflare detected. Recommended: <strong>Cloudflare Origin Certificate</strong>.</span>
        <span v-else>Not behind Cloudflare. Recommended: <strong>Let's Encrypt</strong>.</span>
      </div>

      <CfSelect
        v-model="form.provider"
        label="Certificate provider"
        :options="[
          { value: 'letsencrypt', label: 'Let’s Encrypt' },
          { value: 'zerossl', label: 'ZeroSSL (needs EAB keys)' },
          { value: 'cloudflare_origin', label: 'Cloudflare Origin Certificate (paste)' },
        ]"
      />

      <template v-if="form.provider === 'zerossl'">
        <CfInput v-model="form.eab_kid" label="EAB key ID" autocomplete="off" />
        <CfInput v-model="form.eab_hmac" label="EAB HMAC key" autocomplete="off" />
      </template>

      <template v-if="form.provider === 'cloudflare_origin'">
        <CfInput v-model="form.cf_api_token" type="password" label="Cloudflare API token" autocomplete="off" placeholder="Paste token (not stored in the browser)" />
        <p class="srv-muted">Create it in Cloudflare under My Profile, API Tokens, with the permission Zone, SSL and Certificates, Edit. Orbit issues the certificate and renews it itself.</p>
        <details>
          <summary class="srv-muted">Or paste the origin certificate yourself</summary>
          <label>Origin certificate (PEM)
          <textarea v-model="form.cert_pem" rows="5" placeholder="-----BEGIN CERTIFICATE-----"></textarea>
        </label>
        <label>Private key (PEM)
          <textarea v-model="form.key_pem" rows="5" placeholder="-----BEGIN PRIVATE KEY-----"></textarea>
        </label>
        </details>
      </template>

      <p class="srv-muted">Once saved, this domain serves the dashboard over HTTPS. Its DNS record must point at this server.</p>

      <div class="srv-actions">
        <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy || !form.panel_domain">
          {{ busy ? 'Saving…' : 'Save panel domain' }}
        </button>
        <button class="btn-cf btn-cf-secondary" type="button" :disabled="busy || !current" @click="clear">
          Remove panel domain
        </button>
      </div>
    </form>

    <div v-if="current" class="srv-card">
      <strong>Current panel domain</strong>
      <p class="srv-muted"><a :href="'https://' + current" target="_blank" rel="noopener">https://{{ current }}</a></p>
      <label class="srv-check">
        <input type="checkbox" v-model="lockOn" :disabled="busy" @change="saveLock" />
        <span>
          <strong>Only allow access through this domain</strong>
          <span class="srv-muted">The server IP stops opening the dashboard. Your SSH tunnel still works, so you can always get back in.</span>
        </span>
      </label>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/composables/api.js'

const current = ref('')
const lockOn = ref(false)
const busy = ref(false)
const checking = ref(false)
const error = ref('')
const notice = ref('')
const detection = ref(null)
const form = ref({ panel_domain: '', provider: 'letsencrypt', eab_kid: '', eab_hmac: '', cert_pem: '', key_pem: '', cf_api_token: '' })

const load = async () => {
  try {
    const res = await api('/api/v1/server/settings')
    current.value = res.panel_domain || ''
    lockOn.value = !!res.lock_to_panel_domain
    if (current.value) form.value.panel_domain = current.value
  } catch (e) {
    error.value = e.message
  }
}

const detect = async () => {
  const d = form.value.panel_domain.trim()
  if (!d) return
  checking.value = true
  error.value = ''
  try {
    detection.value = await api(`/api/v1/tls/detect?domain=${encodeURIComponent(d)}`)
    form.value.provider = detection.value.recommended
  } catch (e) {
    detection.value = null
    error.value = e.message
  } finally {
    checking.value = false
  }
}

const save = async () => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    const body = { ...form.value }
    if (body.provider !== 'zerossl') { body.eab_kid = ''; body.eab_hmac = '' }
    if (body.provider !== 'cloudflare_origin') { body.cert_pem = ''; body.key_pem = ''; body.cf_api_token = '' }
    const res = await api('/api/v1/server/settings', { method: 'PUT', body })
    current.value = res.panel_domain
    notice.value = res.certificate?.status === 'active'
      ? `Saved. https://${res.panel_domain} is live.`
      : `Saved. The certificate for ${res.panel_domain} is being issued; this can take a minute.`
    form.value.cert_pem = ''
    form.value.key_pem = ''
    form.value.cf_api_token = ''
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const saveLock = async () => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    const res = await api('/api/v1/server/settings', {
      method: 'PUT',
      body: { panel_domain: current.value, provider: '', lock_to_panel_domain: lockOn.value },
    })
    lockOn.value = !!res.lock_to_panel_domain
    notice.value = lockOn.value
      ? 'Locked. The dashboard now answers only on the panel domain.'
      : 'Lock turned off. The server IP works again.'
  } catch (e) {
    lockOn.value = !lockOn.value
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const clear = async () => {
  if (!window.confirm('Remove the panel domain? The dashboard stays available on the server IP.')) return
  busy.value = true
  try {
    await api('/api/v1/server/settings', { method: 'PUT', body: { panel_domain: '' } })
    current.value = ''
    form.value.panel_domain = ''
    notice.value = 'Panel domain removed.'
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.srv-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
.srv-card { background: var(--kumo-base, #fff); border: 1px solid var(--kumo-border, #e2e4e8); border-radius: 10px; padding: 1rem; box-sizing: border-box; min-width: 0; display: flex; flex-direction: column; gap: 0.5rem; }
.srv-form { gap: 0.75rem; }
.srv-form label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; min-width: 0; }
.srv-form input, .srv-form select, .srv-form textarea {
  padding: 9px 10px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px;
  background: #fff; color: inherit; width: 100%; box-sizing: border-box; font-family: inherit;
}
.srv-form textarea { font-family: ui-monospace, monospace; font-size: 12px; }
.srv-inline { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.srv-inline input { flex: 1 1 180px; min-width: 0; }
.srv-detect { font-size: 13px; background: #f3f4f6; padding: 0.6rem 0.75rem; border-radius: 6px; }
.srv-actions { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.srv-muted { font-size: 13px; color: #6b7280; margin: 0; overflow-wrap: anywhere; }
.srv-check { display: flex; gap: 0.6rem; align-items: flex-start; font-size: 14px; cursor: pointer; }
.srv-check input { margin-top: 3px; }
.srv-check .srv-muted { display: block; }
.srv-error { color: #b91c1c; font-size: 13px; margin: 0; overflow-wrap: anywhere; }
.srv-notice { color: #166534; font-size: 13px; margin: 0; }
@media (max-width: 680px) {
  .srv-inline .btn-cf, .srv-actions .btn-cf { flex: 1 1 100%; }
}
</style>
