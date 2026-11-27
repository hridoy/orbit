<template>
  <div class="cf-copy-field">
    <label v-if="label" class="cf-copy-label">{{ label }}</label>
    <div class="cf-copy-box">
      <input
        :type="secret && !revealed ? 'password' : 'text'"
        :value="value"
        readonly
        class="cf-copy-input"
        @focus="$event.target.select()"
      />
      <button
        v-if="secret"
        type="button"
        class="cf-copy-action"
        :title="revealed ? 'Hide value' : 'Reveal value'"
        @click="revealed = !revealed"
      >
        {{ revealed ? 'Hide' : 'Show' }}
      </button>
      <button
        type="button"
        class="cf-copy-action is-primary"
        @click="copyValue"
      >
        {{ copied ? 'Copied' : 'Copy' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  value: {
    type: String,
    required: true
  },
  label: {
    type: String,
    default: ''
  },
  secret: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['copy'])
const revealed = ref(false)
const copied = ref(false)

async function copyValue() {
  try {
    await navigator.clipboard.writeText(props.value)
  } catch {
    // ignore in non-HTTPS
  }
  copied.value = true
  emit('copy', props.value)
  setTimeout(() => {
    copied.value = false
  }, 1800)
}
</script>

<style scoped>
.cf-copy-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-copy-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-copy-box {
  display: flex;
  align-items: stretch;
  height: 36px;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--color-kumo-recessed);
  overflow: hidden;
}

.cf-copy-input {
  flex: 1;
  min-width: 0;
  border: none;
  background: transparent;
  padding: 0 12px;
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: var(--text-color-kumo-strong);
  outline: none;
}

.cf-copy-action {
  border: none;
  border-left: 1px solid var(--kumo-control-border-default);
  background: var(--color-kumo-base);
  color: var(--text-color-kumo-default);
  font-family: var(--font-sans);
  font-size: 12px;
  font-weight: 500;
  padding: 0 12px;
  cursor: pointer;
  transition: background-color 0.15s ease;
  white-space: nowrap;
}

.cf-copy-action:hover {
  background: var(--color-kumo-tint);
}

.cf-copy-action.is-primary {
  color: var(--text-color-kumo-link);
  font-weight: 600;
}
</style>
