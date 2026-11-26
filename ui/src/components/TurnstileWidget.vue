<template>
  <div class="turnstile-simulator-card">
    <div
      class="turnstile-box"
      @click="triggerSolve"
    >
      <div style="display: flex; align-items: center; gap: 0.85rem;">
        <div
          class="turnstile-checkbox"
          :class="{ solved: isSolved, solving: isSolving }"
        >
          <span v-if="isSolved" style="color: #ffffff; font-size: 14px; font-weight: 700;">✓</span>
          <div v-else-if="isSolving" class="turnstile-spinner"></div>
        </div>
        <span class="turnstile-text" :class="{ 'text-solved': isSolved }">
          {{ labelText }}
        </span>
      </div>
      <div style="display: flex; flex-direction: column; align-items: flex-end;">
        <img src="/assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 18px;" />
        <span style="font-size: 0.65rem; color: var(--text-color-kumo-subtle); margin-top: 2px;">Cloudflare Turnstile</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isSolved = ref(false)
const isSolving = ref(false)
const labelText = ref('Verify you are human')

function triggerSolve() {
  if (isSolved.value || isSolving.value) return
  isSolving.value = true
  labelText.value = 'Verifying...'

  setTimeout(() => {
    isSolving.value = false
    isSolved.value = true
    labelText.value = 'Verification successful (Human visitor verified)'
  }, 450)
}
</script>

<style scoped>
.turnstile-simulator-card {
  background: var(--color-kumo-canvas);
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  max-width: 380px;
}
.turnstile-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-md);
  padding: 0.85rem 1.15rem;
  cursor: pointer;
  user-select: none;
  transition: var(--kumo-transition);
}
.turnstile-box:hover {
  border-color: rgba(0, 0, 0, 0.2);
}
.turnstile-checkbox {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-kumo-line);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-kumo-base);
  transition: var(--kumo-transition);
}
.turnstile-checkbox.solved {
  background: #059669;
  border-color: #059669;
}
.turnstile-text {
  font-size: 0.88rem;
  font-weight: 500;
  color: var(--text-color-kumo-strong);
}
.text-solved {
  color: #059669 !important;
}
.turnstile-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid var(--cf-brand-orange);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spinSlow 0.6s linear infinite;
}
</style>
