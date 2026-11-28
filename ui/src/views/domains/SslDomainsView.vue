<template>
  <div class="tls-page">
    <div class="breadcrumbs">
      <span>SSL/TLS</span>
      <span class="divider">/</span>
      <span class="current">Domains</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Domains and HTTPS</h1>
        <p class="page-subtitle">Each domain gets its own certificate. Port 443 serves it, and port 80 answers the certificate checks.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showForm = !showForm">
        {{ showForm ? 'Close' : '+ Add domain' }}
      </button>
    </div>

    <p v-if="error" class="tls-error" role="alert">{{ error }}</p>
    <p v-if="notice" class="tls-notice">{{ notice }}</p>

    <form v-if="showForm" class="tls-card tls-form" @submit.prevent="add">
      <CfInput v-model="form.domain" label="Domain" required placeholder="app.example.com" @change="detect">
        <template #suffix>
          <button class="btn-cf btn-cf-secondary" type="button" :disabled="!form.domain || detecting" @click="detect">
            {{ detecting ? 'Checking…' : 'Check' }}
          </button>
        </template>
      </CfInput>

      <div v-if="detection" class="tls-detect">
        <span v-if="detection.cloudflare">Cloudflare detected. Recommended: <strong>Cloudflare Origin Certificate</strong>.</span>
        <span v-else>Not behind Cloudflare. Recommended: <strong>Let's Encrypt</strong>.</span>
        <span class="tls-muted">Resolves to {{ detection.ips.join(', ') }}</span>
      </div>

      <CfSelect
        v-model="form.provider"
        label="Certificate provider"
        :options="[
          { value: 'letsencrypt', label: 'Let’s Encrypt (free, automatic)' },
          { value: 'zerossl', label: 'ZeroSSL (free, automatic, needs EAB keys)' },
          { value: 'cloudflare_origin', label: 'Cloudflare Origin Certificate (paste)' },
        ]"
      />

      <template v-if="form.provider === 'zerossl'">
        <CfInput v-model="form.eab_kid" label="EAB key ID" required autocomplete="off" />
        <CfInput v-model="form.eab_hmac" label="EAB HMAC key" required autocomplete="off" />
      </template>

      <template v-if="form.provider === 'cloudflare_origin'">
        <CfInput v-model="form.cf_api_token" type="password" label="Cloudflare API token" autocomplete="off" placeholder="Paste token" />
        <p class="tls-muted">Create it in Cloudflare under My Profile, API Tokens, with Zone, SSL and Certificates, Edit. Orbit issues and renews the certificate itself.</p>
        <details>
          <summary class="tls-muted">Or paste the origin certificate yourself</summary>
          <label>Origin certificate (PEM)
          <textarea v-model="form.cert_pem" rows="5" required placeholder="-----BEGIN CERTIFICATE-----"></textarea>
        </label>
        <label>Private key (PEM)
          <textarea v-model="form.key_pem" rows="5" placeholder="-----BEGIN PRIVATE KEY-----"></textarea>
        </label>
        </details>
      </template>

      <p v-if="form.provider !== 'cloudflare_origin'" class="tls-muted">
        The domain's DNS record must point at this server, and port 80 must be open, for the certificate to issue.
      </p>

      <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy">
        {{ busy ? 'Saving…' : 'Add domain' }}
      </button>
    </form>

    <p v-if="loading" class="tls-muted">Loading…</p>
    <p v-else-if="!domains.length" class="tls-muted">No domains yet.</p>

    <div class="tls-list">
      <div v-for="d in domains" :key="d.domain" class="tls-card tls-row">
        <div class="tls-row-main">
          <strong class="tls-domain">{{ d.domain }}</strong>
          <span :class="['tls-status', d.status]">{{ d.status }}</span>
        </div>
        <div class="tls-row-meta">
          <span class="tls-muted">{{ providerLabel(d.provider) }}</span>
          <span v-if="d.expires_at" class="tls-muted">expires {{ fmtDate(d.expires_at) }}</span>
        </div>
        <p v-if="d.last_error" class="tls-error">{{ d.last_error }}</p>
        <div class="tls-row-actions">
          <button class="btn-cf btn-cf-secondary" :disabled="busy" @click="remove(d.domain)">Remove</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { api } from '@/composables/api.js'

const domains = ref([])
const loading = ref(true)
const busy = ref(false)
const detecting = ref(false)
const error = ref('')
const notice = ref('')
const showForm = ref(false)
const detection = ref(null)

const blankForm = () => ({ domain: '', provider: 'letsencrypt', eab_kid: '', eab_hmac: '', cert_pem: '', key_pem: '', cf_api_token: '' })
const form = ref(blankForm())

const providerLabel = (p) => ({
  letsencrypt: "Let's Encrypt",
  zerossl: 'ZeroSSL',
  cloudflare_origin: 'Cloudflare Origin Certificate',
}[p] || p)

const fmtDate = (iso) => new Date(iso).toLocaleDateString()

const load = async () => {
  try {
    domains.value = await api('/api/v1/tls/domains')
    error.value = ''
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

// Check DNS and Cloudflare, then pre-select the recommended provider.
const detect = async () => {
  const d = form.value.domain.trim()
  if (!d) return
  detecting.value = true
  error.value = ''
  try {
    detection.value = await api(`/api/v1/tls/detect?domain=${encodeURIComponent(d)}`)
    form.value.provider = detection.value.recommended
  } catch (e) {
    detection.value = null
    error.value = e.message
  } finally {
    detecting.value = false
  }
}

const add = async () => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    const body = { ...form.value }
    if (body.provider !== 'zerossl') { body.eab_kid = ''; body.eab_hmac = '' }
    if (body.provider !== 'cloudflare_origin') { body.cert_pem = ''; body.key_pem = ''; body.cf_api_token = '' }
    const rec = await api('/api/v1/tls/domains', { method: 'POST', body })
    notice.value = rec.status === 'active'
      ? `${rec.domain} is active.`
      : `${rec.domain} is being issued. This takes a minute.`
    form.value = blankForm()
    detection.value = null
    showForm.value = false
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const remove = async (domain) => {
  if (!window.confirm(`Remove ${domain}? Its HTTPS certificate is deleted.`)) return
  busy.value = true
  try {
    await api(`/api/v1/tls/domains/${encodeURIComponent(domain)}`, { method: 'DELETE' })
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
  timer = setInterval(load, 10000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.tls-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
.tls-card {
  background: var(--kumo-base, #fff);
  border: 1px solid var(--kumo-border, #e2e4e8);
  border-radius: 10px;
  padding: 1rem;
  box-sizing: border-box;
  min-width: 0;
}
.tls-form { display: flex; flex-direction: column; gap: 0.75rem; }
.tls-form label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; min-width: 0; }
.tls-form input, .tls-form select, .tls-form textarea {
  padding: 9px 10px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px;
  background: #fff; color: inherit; width: 100%; box-sizing: border-box; font-family: inherit;
}
.tls-form textarea { font-family: ui-monospace, monospace; font-size: 12px; }
.tls-inline { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.tls-inline input { flex: 1 1 180px; min-width: 0; }
.tls-detect { font-size: 13px; background: #f3f4f6; padding: 0.6rem 0.75rem; border-radius: 6px; display: flex; flex-direction: column; gap: 2px; overflow-wrap: anywhere; }
.tls-list { display: flex; flex-direction: column; gap: 0.6rem; }
.tls-row-main, .tls-row-meta { display: flex; flex-wrap: wrap; gap: 0.5rem 0.75rem; align-items: center; justify-content: space-between; }
.tls-domain { overflow-wrap: anywhere; }
.tls-row-actions { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-top: 0.5rem; }
.tls-muted { font-size: 13px; color: #6b7280; margin: 0; overflow-wrap: anywhere; }
.tls-status { font-size: 12px; font-weight: 600; padding: 2px 8px; border-radius: 999px; background: #e5e7eb; color: #374151; }
.tls-status.active { background: #dcfce7; color: #166534; }
.tls-status.pending { background: #fef9c3; color: #854d0e; }
.tls-status.failed { background: #fee2e2; color: #991b1b; }
.tls-error { color: #b91c1c; font-size: 13px; margin: 0; overflow-wrap: anywhere; }
.tls-notice { font-size: 13px; color: #166534; margin: 0; }

@media (max-width: 680px) {
  .tls-inline .btn-cf { flex: 1 1 100%; }
  .tls-row-actions .btn-cf { flex: 1 1 100%; }
}
</style>
