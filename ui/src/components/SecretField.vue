<template>
  <div class="secret">
    <code class="secret-value">{{ shown ? value : mask }}</code>
    <button type="button" class="secret-btn" :aria-label="shown ? 'Hide' : 'Show'" :title="shown ? 'Hide' : 'Show'" @click="shown = !shown">
      <svg v-if="!shown" viewBox="0 0 24 24" width="16" height="16" aria-hidden="true"><path fill="currentColor" d="M12 5c-7 0-10 7-10 7s3 7 10 7 10-7 10-7-3-7-10-7zm0 11a4 4 0 110-8 4 4 0 010 8zm0-6a2 2 0 100 4 2 2 0 000-4z"/></svg>
      <svg v-else viewBox="0 0 24 24" width="16" height="16" aria-hidden="true"><path fill="currentColor" d="M2 4.3L3.3 3 21 20.7 19.7 22l-3-3C15.4 19.6 13.8 20 12 20c-7 0-10-7-10-7a18 18 0 014.2-5.2L2 4.3zM12 7c3.5 0 6 2.3 7.3 3.7.4.4.7.8.9 1.3-.3.4-.6.8-1 1.2l-2.6-2.6A4 4 0 0012 7z"/></svg>
    </button>
    <button type="button" class="secret-btn" :title="copied ? 'Copied' : 'Copy'" @click="copy">
      <span v-if="!copied" class="secret-text">Copy</span>
      <span v-else class="secret-text">Copied</span>
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { copyText } from '@/composables/clipboard.js'

const props = defineProps({
  value: { type: String, default: '' },
  // Show the value by default. Passwords start hidden.
  revealed: { type: Boolean, default: true },
})
const shown = ref(props.revealed)
const copied = ref(false)
const mask = '••••••••••••'

const copy = async () => {
  const ok = await copyText(props.value)
  if (ok) {
    copied.value = true
    setTimeout(() => (copied.value = false), 1500)
  }
}
</script>

<style scoped>
.secret { display: inline-flex; flex-wrap: wrap; align-items: center; gap: 0.4rem; max-width: 100%; }
.secret-value { font-size: 12px; word-break: break-all; white-space: pre-wrap; max-width: 100%; }
.secret-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 4px;
  min-height: 28px; min-width: 28px; padding: 4px 8px; border: 1px solid #d1d5db; border-radius: 6px;
  background: #fff; color: inherit; cursor: pointer; font-size: 12px;
}
.secret-btn:hover { background: #f3f4f6; }
.secret-text { font-weight: 600; }
</style>
