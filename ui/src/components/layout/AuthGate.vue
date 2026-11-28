<template>
  <div class="auth-split-wrapper">
    <!-- Left hero pane: Orbit identity (desktop only) -->
    <div class="auth-hero-pane">
      <div class="auth-brand">
        <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
          <path d="M37.5 19.5C36.8 13.5 31.8 9 25.5 9C20.2 9 15.8 12.2 14.1 16.8C13.4 16.6 12.7 16.5 12 16.5C7.6 16.5 4 20.1 4 24.5C4 28.9 7.6 32.5 12 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5Z" fill="#F38020"/>
          <path d="M32 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5C36.8 13.5 31.8 9 25.5 9C24.5 9 23.5 9.1 22.5 9.4C26.5 12.8 29.1 17.8 29.4 23.5C30.2 23.5 31.1 23.6 32 24C32.1 26.8 32 29.6 32 32.5Z" fill="#FAAD3F"/>
        </svg>
        <span>Orbit Sovereign</span>
      </div>

      <div>
        <h2 style="font-size: 2.25rem; font-weight: 700; line-height: 1.2; margin-bottom: 1rem;">
          Your own cloud, under your control.
        </h2>
        <p style="font-size: 1rem; color: #94a3b8; max-width: 480px; line-height: 1.5;">
          Deploy apps, databases, and object storage on a single server you own.
        </p>
      </div>

      <div style="font-size: 0.8rem; color: #64748b;">Orbit Sovereign PaaS</div>
    </div>

    <!-- Right card pane: sign in or first-run admin setup -->
    <div class="auth-card-pane">
      <form @submit.prevent="submit">
        <h1>{{ challenge ? 'Enter your code' : (setupRequired ? 'Create your admin login' : 'Sign in') }}</h1>
        <p v-if="challenge">Enter the 6-digit code from your authenticator app.</p>
        <p v-else>
          {{ setupRequired
            ? 'First run. Choose a username and password for this dashboard.'
            : 'Enter the username and password you set during setup.' }}
        </p>

        <template v-if="challenge">
          <div class="cf-form-group">
            <label for="auth-code">Authentication code</label>
            <input id="auth-code" v-model="code" class="cf-input" inputmode="numeric" maxlength="6" autocomplete="one-time-code" required autofocus />
          </div>
          <p v-if="error" class="auth-error" role="alert">{{ error }}</p>
          <button type="submit" class="btn-cf btn-cf-primary auth-submit" :disabled="busy">
            {{ busy ? 'Please wait…' : 'Verify' }}
          </button>
          <button type="button" class="btn-cf btn-cf-secondary auth-submit" :disabled="busy" @click="backToPassword">Back</button>
        </template>

        <template v-else>
        <div class="cf-form-group">
          <label for="auth-username">Username</label>
          <input id="auth-username" v-model="username" class="cf-input" type="text" autocomplete="username" maxlength="64" required autofocus />
        </div>

        <div class="cf-form-group">
          <label for="auth-password">Password</label>
          <input id="auth-password" v-model="password" class="cf-input" type="password" :autocomplete="setupRequired ? 'new-password' : 'current-password'" required />
        </div>

        <div v-if="setupRequired" class="cf-form-group">
          <label for="auth-confirm">Confirm password</label>
          <input id="auth-confirm" v-model="confirm" class="cf-input" type="password" autocomplete="new-password" required />
          <span class="auth-hint">At least 8 characters.</span>
        </div>

        <p v-if="error" class="auth-error" role="alert">{{ error }}</p>

        <button type="submit" class="btn-cf btn-cf-primary auth-submit" :disabled="busy">
          {{ busy ? 'Please wait…' : (setupRequired ? 'Create account' : 'Sign in') }}
        </button>
        </template>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  setupRequired: { type: Boolean, default: false },
})
const emit = defineEmits(['authed'])

const username = ref('')
const password = ref('')
const confirm = ref('')
const error = ref('')
const busy = ref(false)
const challenge = ref('')
const code = ref('')

const backToPassword = () => {
  challenge.value = ''
  code.value = ''
  error.value = ''
}

const submit = async () => {
  error.value = ''
  if (challenge.value) return verifyCode()
  if (props.setupRequired && password.value !== confirm.value) {
    error.value = 'Passwords do not match.'
    return
  }
  busy.value = true
  try {
    const path = props.setupRequired ? '/api/v1/auth/setup' : '/api/v1/auth/login'
    const res = await fetch(path, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value }),
    })
    const body = await res.json().catch(() => ({}))
    if (!res.ok) {
      error.value = body.error || `Request failed (HTTP ${res.status})`
      return
    }
    password.value = ''
    confirm.value = ''
    if (body.two_factor_required) {
      // The password was right. Ask for the authenticator code next.
      challenge.value = body.challenge
      return
    }
    emit('authed', body.token)
  } catch (e) {
    error.value = 'Could not reach the agent. Check the SSH tunnel.'
  } finally {
    busy.value = false
  }
}

const verifyCode = async () => {
  busy.value = true
  try {
    const res = await fetch('/api/v1/auth/login/2fa', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ challenge: challenge.value, code: code.value.trim() }),
    })
    const body = await res.json().catch(() => ({}))
    if (!res.ok) {
      error.value = body.error || `Request failed (HTTP ${res.status})`
      // An expired or used-up sign-in step must start over.
      if (res.status === 401 && /expired|too many|log in again/i.test(error.value)) backToPassword()
      return
    }
    challenge.value = ''
    code.value = ''
    emit('authed', body.token)
  } catch (e) {
    error.value = 'Could not reach the agent. Check the SSH tunnel.'
  } finally {
    busy.value = false
  }
}
</script>

<style scoped>
.auth-brand {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  font-size: 1.25rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
}
.auth-brand svg { width: 32px; height: 32px; }
.auth-card-pane form { display: flex; flex-direction: column; gap: 0.25rem; }
.auth-hint { font-size: 0.78rem; color: var(--text-color-kumo-subtle); }
.auth-error { font-size: 0.85rem; color: #b91c1c; margin: 0.25rem 0 0; }
.auth-submit { width: 100%; justify-content: center; height: 42px; margin-top: 0.75rem; }
</style>
