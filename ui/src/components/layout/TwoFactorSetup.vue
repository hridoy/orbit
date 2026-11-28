<template>
  <section class="tfa-card">
    <h3>Two-factor sign-in</h3>
    <p class="tfa-muted">
      Ask for a 6-digit code from an authenticator app (Google Authenticator, 1Password, Authy) each time you sign in.
    </p>

    <p v-if="error" class="tfa-error" role="alert">{{ error }}</p>
    <p v-if="notice" class="tfa-notice">{{ notice }}</p>

    <template v-if="status === null">
      <p class="tfa-muted">Loading…</p>
    </template>

    <template v-else-if="status.enabled">
      <p><strong>On.</strong> Sign-in asks for a code after your password.</p>
      <form class="tfa-form" @submit.prevent="disable">
        <CfInput v-model="disablePassword" type="password" label="Your password to turn it off" required autocomplete="current-password" />
        <button class="btn-cf btn-cf-secondary" type="submit" :disabled="busy">Turn off two-factor</button>
      </form>
    </template>

    <template v-else>
      <p v-if="!setup" class="tfa-muted">Off.</p>
      <button v-if="!setup" class="btn-cf btn-cf-primary" type="button" :disabled="busy" @click="startSetup">
        Set up two-factor
      </button>

      <div v-if="setup" class="tfa-form">
        <p>Add this account to your authenticator app. Enter the setup key, or open the link on the phone that has the app.</p>
        <CfCopyField :value="setup.secret" label="Setup key" />
        <a class="tfa-link" :href="setup.uri">Open in authenticator app</a>
        <form @submit.prevent="enable">
          <CfInput v-model="code" label="6-digit code from the app" inputmode="numeric" maxlength="6" autocomplete="one-time-code" required />
          <button class="btn-cf btn-cf-primary" type="submit" :disabled="busy">Turn on</button>
        </form>
      </div>
    </template>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/composables/api.js'

const status = ref(null)
const setup = ref(null)
const code = ref('')
const disablePassword = ref('')
const busy = ref(false)
const error = ref('')
const notice = ref('')

const load = async () => {
  try {
    status.value = await api('/api/v1/auth/2fa')
  } catch (e) {
    error.value = e.message
  }
}

const run = async (fn, done) => {
  busy.value = true
  error.value = ''
  notice.value = ''
  try {
    await fn()
    if (done) notice.value = done
  } catch (e) {
    error.value = e.message
  } finally {
    busy.value = false
  }
}

const startSetup = () => run(async () => {
  setup.value = await api('/api/v1/auth/2fa/setup', { method: 'POST' })
  code.value = ''
})

const enable = () => run(async () => {
  await api('/api/v1/auth/2fa/enable', { method: 'POST', body: { code: code.value.trim() } })
  setup.value = null
  code.value = ''
  status.value = { enabled: true, pending: false }
}, 'Two-factor is on.')

const disable = () => run(async () => {
  await api('/api/v1/auth/2fa/disable', { method: 'POST', body: { password: disablePassword.value } })
  disablePassword.value = ''
  status.value = { enabled: false, pending: false }
}, 'Two-factor is off.')

onMounted(load)
</script>

<style scoped>
.tfa-card { border: 1px solid var(--color-kumo-line); border-radius: 8px; padding: 1rem 1.25rem; display: flex; flex-direction: column; gap: 0.6rem; }
.tfa-card h3 { margin: 0; font-size: 1rem; }
.tfa-form { display: flex; flex-direction: column; gap: 0.6rem; max-width: 420px; }
.tfa-muted { color: var(--text-color-kumo-subtle); font-size: 0.85rem; margin: 0; }
.tfa-error { color: #b91c1c; font-size: 0.85rem; margin: 0; }
.tfa-notice { color: #15803d; font-size: 0.85rem; margin: 0; }
.tfa-link { font-size: 0.85rem; }
</style>
