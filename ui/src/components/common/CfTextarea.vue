<script setup>
const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  rows: {
    type: Number,
    default: 4
  },
  hint: {
    type: String,
    default: ''
  },
  error: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  maxlength: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

function onInput(e) {
  emit('update:modelValue', e.target.value)
}
</script>

<template>
  <div class="cf-textarea-wrap" :class="{ 'is-disabled': disabled, 'has-error': !!error }">
    <div v-if="label || maxlength" class="cf-textarea-header">
      <label v-if="label" class="cf-textarea-label">{{ label }}</label>
      <span v-if="maxlength" class="cf-textarea-count">
        {{ (modelValue || '').length }} / {{ maxlength }}
      </span>
    </div>

    <textarea
      :value="modelValue"
      :placeholder="placeholder"
      :rows="rows"
      :disabled="disabled"
      :maxlength="maxlength"
      class="cf-textarea-field"
      @input="onInput"
    ></textarea>

    <p v-if="error" class="cf-textarea-msg is-error">{{ error }}</p>
    <p v-else-if="hint" class="cf-textarea-msg is-hint">{{ hint }}</p>
  </div>
</template>

<style scoped>
.cf-textarea-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-textarea-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cf-textarea-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--kumo-fg-default, #18181b);
}

.cf-textarea-count {
  font-size: 11px;
  color: var(--kumo-fg-muted, #71717a);
  font-family: var(--font-mono, monospace);
}

.cf-textarea-field {
  width: 100%;
  padding: 8px 12px;
  border-radius: var(--kumo-radius-md, 6px);
  background-color: var(--kumo-control-bg-default, #ffffff);
  border: 1px solid var(--kumo-control-border-default, #d4d4d8);
  font-family: inherit;
  font-size: 13px;
  line-height: 20px;
  color: var(--kumo-control-fg-default, #18181b);
  resize: vertical;
  outline: none;
  transition: all 0.15s ease;
}

.cf-textarea-field:hover:not(:disabled) {
  border-color: var(--kumo-control-border-hover, #a1a1aa);
}

.cf-textarea-field:focus {
  border-color: #0051c3;
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.15);
}

.has-error .cf-textarea-field {
  border-color: #dc2626 !important;
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12) !important;
}

.is-disabled {
  opacity: 0.6;
  pointer-events: none;
}

.cf-textarea-msg {
  margin: 0;
  font-size: 12px;
  line-height: 16px;
}

.cf-textarea-msg.is-error {
  color: #dc2626;
}

.cf-textarea-msg.is-hint {
  color: var(--kumo-fg-muted, #71717a);
}
</style>
