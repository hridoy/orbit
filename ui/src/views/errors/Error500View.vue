<template>
  <div class="error-container">
    <img src="/assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 38px; margin-bottom: 1.5rem;" />
    <h1 style="font-size: 1.85rem; font-weight: 700; color: var(--text-color-kumo-strong);">Web Server Is Returning An Unknown Error</h1>
    <p style="font-size: 0.95rem; color: var(--text-color-kumo-subtle); margin-top: 0.25rem;">There is an issue between Cloudflare's edge cache and your origin server.</p>

    <!-- Cloudflare Diagnostic Connection Diagram -->
    <div class="flow-diagram">
      <div class="flow-step">
        <div class="flow-icon">💻</div>
        <strong style="font-size: 0.9rem;">Browser</strong>
        <span style="color: #059669; font-size: 0.78rem; font-weight: 600; margin-top: 0.25rem;">Working</span>
      </div>
      <div style="font-size: 1.5rem; color: var(--text-color-kumo-subtle);">→</div>
      <div class="flow-step">
        <div class="flow-icon">☁️</div>
        <strong style="font-size: 0.9rem;">Cloudflare</strong>
        <span style="color: #059669; font-size: 0.78rem; font-weight: 600; margin-top: 0.25rem;">Working (IAD)</span>
      </div>
      <div style="font-size: 1.5rem; color: #dc2626;">⇸</div>
      <div class="flow-step error">
        <div class="flow-icon">🖥️</div>
        <strong style="font-size: 0.9rem;">Host Server</strong>
        <span style="color: #dc2626; font-size: 0.78rem; font-weight: 600; margin-top: 0.25rem;">Error 520</span>
      </div>
    </div>

    <div class="ray-id-pill" style="margin-bottom: 2rem;">
      Cloudflare Ray ID: <strong>8e41bb109a4f2100</strong> • IP: 198.51.100.24 • Origin: origin.northlab.dev
    </div>

    <router-link to="/" class="btn-cf btn-cf-primary">Return to Cloudflare Dashboard</router-link>
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
  padding: 2.5rem;
  background: var(--color-kumo-canvas);
}
.flow-diagram {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin: 2rem 0;
}
.flow-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-lg);
  padding: 1.25rem 1.5rem;
  min-width: 140px;
}
.flow-step.error {
  border-color: #dc2626;
  background: rgba(220, 38, 38, 0.04);
}
.flow-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}
.ray-id-pill {
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-md);
  padding: 0.5rem 1rem;
  font-family: var(--font-mono);
  font-size: 0.78rem;
  color: var(--text-color-kumo-subtle);
}
</style>
