<template>
  <div class="error-container">
    <img src="/assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 38px; margin-bottom: 2rem;" />
    <div class="error-code">404</div>
    <h1 style="font-size: 1.75rem; font-weight: 700; color: var(--text-color-kumo-strong); margin-bottom: 0.5rem;">Page Not Found</h1>
    <p style="font-size: 0.95rem; color: var(--text-color-kumo-subtle); max-width: 460px; line-height: 1.5;">
      The requested URL or edge resource route was not found on this server or V8 worker.
    </p>

    <div class="ray-id-pill">
      <span>Cloudflare Ray ID: <strong>8e41a980b12c98ef</strong></span>
      <span>•</span>
      <span>Your IP: 198.51.100.24</span>
      <span>•</span>
      <span>Edge: IAD</span>
    </div>

    <div style="display: flex; gap: 0.75rem;">
      <router-link to="/" class="btn-cf btn-cf-primary">Return to Dashboard</router-link>
      <a href="https://cloudflarestatus.com" target="_blank" class="btn-cf btn-cf-secondary">System Status ↗</a>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.error-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 2rem;
  background: var(--color-kumo-canvas);
}
.error-code {
  font-size: 5.5rem;
  font-weight: 800;
  color: var(--cf-brand-orange);
  line-height: 1;
  margin-bottom: 0.5rem;
  letter-spacing: -0.04em;
}
.ray-id-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-full);
  padding: 0.35rem 0.85rem;
  font-family: var(--font-mono);
  font-size: 0.78rem;
  color: var(--text-color-kumo-subtle);
  margin: 1.5rem 0 2rem 0;
}
</style>
