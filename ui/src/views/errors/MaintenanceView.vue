<template>
  <div class="maint-container">
    <img src="/assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 38px; margin-bottom: 1.5rem;" />
    <span class="cf-pill active" style="margin-bottom: 1rem;"><span class="status-dot"></span><span>Planned Maintenance Window</span></span>
    <h1 style="font-size: 2rem; font-weight: 700; color: var(--text-color-kumo-strong);">Scheduled Edge Infrastructure Upgrade</h1>
    <p style="font-size: 0.95rem; color: var(--text-color-kumo-subtle); max-width: 520px; line-height: 1.5; margin-top: 0.5rem;">
      We are performing kernel optimization and hardware cryptographic upgrades on sovereign cluster nodes. All live traffic continues to be served via Anycast failover.
    </p>

    <!-- Countdown Timer -->
    <div class="countdown-grid">
      <div class="countdown-box">
        <div class="countdown-num">{{ hours }}</div>
        <div class="countdown-label">Hours</div>
      </div>
      <div class="countdown-box">
        <div class="countdown-num">{{ minutes }}</div>
        <div class="countdown-label">Minutes</div>
      </div>
      <div class="countdown-box">
        <div class="countdown-num">{{ seconds }}</div>
        <div class="countdown-label">Seconds</div>
      </div>
    </div>

    <div style="font-size: 0.84rem; color: var(--text-color-kumo-subtle); margin-bottom: 2rem;">
      Live traffic status: <strong>100% Operational</strong> • Failover Region: <strong>FRA / AMS</strong>
    </div>

    <router-link to="/" class="btn-cf btn-cf-primary">Check Dashboard Status</router-link>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const hours = ref('00')
const minutes = ref('42')
const seconds = ref('18')

let interval = null

onMounted(() => {
  let sec = 18
  let min = 42
  interval = setInterval(() => {
    if (sec > 0) {
      sec--
    } else {
      sec = 59
      if (min > 0) min--
    }
    seconds.value = sec < 10 ? '0' + sec : String(sec)
    minutes.value = min < 10 ? '0' + min : String(min)
  }, 1000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<style scoped>
.maint-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 2.5rem;
  background: var(--color-kumo-canvas);
}
.countdown-grid {
  display: flex;
  gap: 1.5rem;
  margin: 2rem 0;
}
.countdown-box {
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-lg);
  padding: 1.25rem 1.75rem;
  min-width: 100px;
}
.countdown-num {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--cf-brand-orange);
  line-height: 1;
}
.countdown-label {
  font-size: 0.76rem;
  text-transform: uppercase;
  color: var(--text-color-kumo-subtle);
  margin-top: 0.35rem;
  font-weight: 600;
}
</style>
