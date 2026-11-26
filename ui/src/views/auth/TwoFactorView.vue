<template>
  <div class="auth-split-wrapper">
    <!-- Left Hero Pane -->
    <div class="auth-hero-pane">
      <div style="display: flex; align-items: center; gap: 0.75rem;">
        <img src="/assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 32px; width: auto;" />
        <span style="font-size: 1.25rem; font-weight: 700;">Cloudflare</span>
      </div>

      <div>
        <h2 style="font-size: 2.25rem; font-weight: 700; line-height: 1.2; margin-bottom: 1rem;">
          Two-Factor Authentication Enforced
        </h2>
        <p style="font-size: 1rem; color: #94a3b8; max-width: 480px; line-height: 1.5;">
          Protecting your sovereign edge cluster, DNS records, and V8 worker isolates with hardware-backed TOTP verification.
        </p>
      </div>

      <div style="font-size: 0.8rem; color: #64748b;">
        © 2026 Cloudflare, Inc. All rights reserved.
      </div>
    </div>

    <!-- Right 2FA Card Pane -->
    <div class="auth-card-pane">
      <h1>Two-Factor Auth</h1>
      <p>Enter the 6-digit verification code from your authenticator app or hardware security key.</p>

      <form @submit.prevent="handleVerify">
        <div class="cf-form-group">
          <label>Authentication Code (TOTP)</label>
          <input 
            v-model="totpCode" 
            type="text" 
            class="cf-input" 
            placeholder="000 000" 
            maxlength="7"
            style="font-family: var(--font-mono); font-size: 1.4rem; letter-spacing: 0.25em; text-align: center; height: 50px;" 
            required 
            autofocus
          />
        </div>

        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem;">
          <button type="button" class="btn-cf btn-cf-secondary btn-cf-sm" @click="useSecurityKey">
            🔑 Use Hardware Key (FIDO2)
          </button>
          <a href="#" @click.prevent="useBackupCode" style="font-size: 0.78rem; color: var(--cf-brand-orange); text-decoration: none;">Use backup code</a>
        </div>

        <button type="submit" class="btn-cf btn-cf-primary" style="width: 100%; justify-content: center; height: 42px;">
          Verify & Sign In
        </button>
      </form>

      <div style="text-align: center; margin-top: 1.5rem; font-size: 0.84rem; color: var(--text-color-kumo-subtle);">
        Having trouble? 
        <router-link to="/auth/login" style="color: var(--cf-brand-orange); font-weight: 500; text-decoration: none;">Return to log in</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const totpCode = ref('849 102')

function handleVerify() {
  router.push('/')
}

function useSecurityKey() {
  alert('Simulating WebAuthn / FIDO2 hardware challenge prompt...')
  setTimeout(() => {
    router.push('/')
  }, 1000)
}

function useBackupCode() {
  const code = prompt('Enter 12-character emergency recovery code:')
  if (code) {
    router.push('/')
  }
}
</script>
