<template>
  <div class="auth-split-wrapper">
    <!-- Left Hero Pane -->
    <div class="auth-hero-pane">
      <div style="display: flex; align-items: center; gap: 0.75rem;">
        <img src="/assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 32px; width: auto;" />
        <span style="font-size: 1.25rem; font-weight: 700; letter-spacing: -0.02em;">Cloudflare</span>
      </div>

      <div>
        <h2 style="font-size: 2.25rem; font-weight: 700; line-height: 1.2; margin-bottom: 1rem;">
          Secure Account Recovery
        </h2>
        <p style="font-size: 1rem; color: #94a3b8; max-width: 480px; line-height: 1.5;">
          Reset credentials with cryptographic verification and emergency recovery codes to regain access to your edge infrastructure.
        </p>
      </div>

      <div style="font-size: 0.8rem; color: #64748b;">
        © 2026 Cloudflare, Inc. All rights reserved.
      </div>
    </div>

    <!-- Right Password Reset Card Pane -->
    <div class="auth-card-pane">
      <h1>Reset password</h1>
      <p>Enter the email address registered with your account to receive instructions.</p>

      <div v-if="submitted" class="ui-alert-banner success" style="margin-bottom: 1.5rem;">
        <div>✓</div>
        <div><strong>Instructions Sent:</strong> Check your inbox for the password reset link. It expires in 15 minutes.</div>
      </div>

      <form v-if="!submitted" @submit.prevent="handleReset">
        <div class="cf-form-group">
          <label>Account email address</label>
          <input v-model="email" type="email" class="cf-input" placeholder="hridoy@duck.com" required />
        </div>

        <button type="submit" class="btn-cf btn-cf-primary" style="width: 100%; justify-content: center; height: 42px;">
          Send password reset email
        </button>
      </form>

      <div style="text-align: center; margin-top: 1.5rem; font-size: 0.84rem; color: var(--text-color-kumo-subtle);">
        Remembered your password? 
        <router-link to="/auth/login" style="color: var(--cf-brand-orange); font-weight: 500; text-decoration: none;">Back to log in</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const email = ref('hridoy@duck.com')
const submitted = ref(false)

function handleReset() {
  submitted.value = true
}
</script>
