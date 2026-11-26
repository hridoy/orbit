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
          Connect, protect, and build everywhere.
        </h2>
        <p style="font-size: 1rem; color: #94a3b8; max-width: 480px; line-height: 1.5;">
          Protect applications from threats, accelerate global web performance, and deploy serverless compute to 330+ cities worldwide with zero egress fees.
        </p>
      </div>

      <div style="font-size: 0.8rem; color: #64748b;">
        © 2026 Cloudflare, Inc. All rights reserved. • Privacy Policy • Terms of Use
      </div>
    </div>

    <!-- Right Login Card Pane -->
    <div class="auth-card-pane">
      <h1>Log in</h1>
      <p>Log in to your Cloudflare dashboard to manage your edge services.</p>

      <button class="btn-cf btn-cf-secondary" style="width: 100%; justify-content: center; height: 42px; margin-bottom: 1rem;" @click="handleSocialLogin">
        <svg width="18" height="18" viewBox="0 0 24 24" style="margin-right: 0.5rem;"><path fill="#EA4335" d="M12 5c1.6 0 3 .6 4.1 1.7l3.1-3.1C17.3 1.8 14.8 1 12 1 7.5 1 3.7 3.6 1.9 7.3l3.7 2.9C6.5 7.4 9 5 12 5z"/><path fill="#4285F4" d="M23.5 12.3c0-.8-.1-1.6-.2-2.3H12v4.5h6.5c-.3 1.5-1.1 2.8-2.4 3.7l3.7 2.9c2.2-2 3.7-5 3.7-8.8z"/><path fill="#FBBC05" d="M5.6 14.8c-.2-.7-.4-1.5-.4-2.3s.2-1.6.4-2.3L1.9 7.3C.7 9.7 0 12.3 0 15s.7 5.3 1.9 7.7l3.7-2.9z"/><path fill="#34A853" d="M12 24c3.2 0 6-1.1 8-3l-3.7-2.9c-1.1.7-2.5 1.2-4.3 1.2-3 0-5.5-2.4-6.4-5.2L1.9 17C3.7 20.7 7.5 24 12 24z"/></svg>
        Continue with Google
      </button>

      <div class="auth-divider"><span>Or continue with email</span></div>

      <form @submit.prevent="handleLogin">
        <div class="cf-form-group">
          <label>Email address</label>
          <input v-model="email" type="email" class="cf-input" required />
        </div>

        <div class="cf-form-group">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <label>Password</label>
            <router-link to="/auth/forgot-password" style="font-size: 0.78rem; color: var(--cf-brand-orange); text-decoration: none;">Forgot password?</router-link>
          </div>
          <input v-model="password" type="password" class="cf-input" required />
        </div>

        <div style="display: flex; align-items: center; gap: 0.5rem; margin-bottom: 1.5rem;">
          <input type="checkbox" id="remember-me" v-model="remember" />
          <label for="remember-me" style="font-size: 0.82rem; margin: 0; cursor: pointer;">Remember this device for 30 days</label>
        </div>

        <button type="submit" class="btn-cf btn-cf-primary" style="width: 100%; justify-content: center; height: 42px;">
          Log in
        </button>
      </form>

      <div style="text-align: center; margin-top: 1.5rem; font-size: 0.84rem; color: var(--text-color-kumo-subtle);">
        Don't have an account? 
        <router-link to="/auth/register" style="color: var(--cf-brand-orange); font-weight: 500; text-decoration: none;">Sign up</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('hridoy@duck.com')
const password = ref('••••••••••••')
const remember = ref(true)

function handleLogin() {
  router.push('/')
}

function handleSocialLogin() {
  router.push('/')
}
</script>
