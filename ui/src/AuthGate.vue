<template>
  <div class="auth-gate">
    <form class="auth-card" @submit.prevent="submit">
      <div class="auth-logo">
        <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
          <path d="M37.5 19.5C36.8 13.5 31.8 9 25.5 9C20.2 9 15.8 12.2 14.1 16.8C13.4 16.6 12.7 16.5 12 16.5C7.6 16.5 4 20.1 4 24.5C4 28.9 7.6 32.5 12 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5Z" fill="#F38020"/>
          <path d="M32 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5C36.8 13.5 31.8 9 25.5 9C24.5 9 23.5 9.1 22.5 9.4C26.5 12.8 29.1 17.8 29.4 23.5C30.2 23.5 31.1 23.6 32 24C32.1 26.8 32 29.6 32 32.5Z" fill="#FAAD3F"/>
        </svg>
        <span>Orbit Sovereign</span>
      </div>

      <h2>{{ setupRequired ? 'Create your admin login' : 'Sign in' }}</h2>
      <p class="auth-hint">
        {{ setupRequired
          ? 'First run. Choose a username and password for this dashboard.'
          : 'Enter the username and password you set during setup.' }}
      </p>

      <label>
        Username
        <input v-model="username" type="text" autocomplete="username" maxlength="64" required autofocus />
      </label>
      <label>
        Password
        <input v-model="password" type="password" :autocomplete="setupRequired ? 'new-password' : 'current-password'" required />
      </label>
      <label v-if="setupRequired">
        Confirm password
        <input v-model="confirm" type="password" autocomplete="new-password" required />
      </label>
      <p v-if="setupRequired" class="auth-hint">At least 8 characters.</p>

      <p v-if="error" class="auth-error" role="alert">{{ error }}</p>

      <button type="submit" :disabled="busy">
        {{ busy ? 'Please wait…' : (setupRequired ? 'Create account' : 'Sign in') }}
      </button>
    </form>
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

const submit = async () => {
  error.value = ''
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
    emit('authed', body.token)
  } catch (e) {
    error.value = 'Could not reach the agent. Check the SSH tunnel.'
  } finally {
    busy.value = false
  }
}
</script>

<style scoped>
.auth-gate {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: #f6f7f9;
  font-family: inherit;
}
.auth-card {
  width: 100%;
  max-width: 380px;
  background: #fff;
  border: 1px solid #e2e4e8;
  border-radius: 10px;
  padding: 28px 24px;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-sizing: border-box;
}
.auth-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #1d1f23;
}
.auth-logo svg { width: 28px; height: 28px; }
h2 { margin: 8px 0 0; font-size: 18px; color: #1d1f23; }
.auth-hint { margin: 0; font-size: 13px; color: #6b7280; }
label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; color: #374151; }
input {
  padding: 9px 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  background: #fff;
  color: #111827;
}
input:focus { outline: 2px solid #f38020; outline-offset: 0; border-color: #f38020; }
button {
  margin-top: 4px;
  padding: 10px;
  border: 0;
  border-radius: 6px;
  background: #f38020;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}
button:disabled { opacity: 0.6; cursor: default; }
.auth-error { margin: 0; font-size: 13px; color: #b91c1c; }
@media (prefers-color-scheme: dark) {
  .auth-gate { background: #0f1115; }
  .auth-card { background: #171a21; border-color: #2a2f3a; }
  .auth-logo, h2 { color: #f3f4f6; }
  label { color: #d1d5db; }
  input { background: #0f1115; border-color: #374151; color: #f3f4f6; }
}
</style>
