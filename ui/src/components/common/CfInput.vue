<template>
  <div :class="['kumo-input-wrapper', { 'kumo-input-error': error, 'kumo-input-disabled': disabled }]">
    <label v-if="label" class="kumo-input-label">
      <span>{{ label }}</span>
      <span v-if="required" class="kumo-required-star">*</span>
    </label>

    <div class="kumo-input-container">
      <span v-if="$slots.prefix || prefixIcon" class="kumo-input-prefix">
        <slot name="prefix">
          <CfIcon v-if="prefixIcon" :name="prefixIcon" size="14" />
        </slot>
      </span>

      <input
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        class="kumo-input-native"
        @input="$emit('update:modelValue', $event.target.value)"
        @focus="$emit('focus', $event)"
        @blur="$emit('blur', $event)"
      />

      <span v-if="$slots.suffix || suffixIcon || clearable" class="kumo-input-suffix">
        <button
          v-if="clearable && modelValue"
          type="button"
          class="kumo-input-clear-btn"
          @click="$emit('update:modelValue', '')"
        >
          ✕
        </button>
        <slot name="suffix">
          <CfIcon v-if="suffixIcon" :name="suffixIcon" size="14" />
        </slot>
      </span>
    </div>

    <span v-if="hint && !error" class="kumo-input-hint">{{ hint }}</span>
    <span v-if="error" class="kumo-input-error-msg">{{ error }}</span>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  label: {
    type: String,
    default: null
  },
  placeholder: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  readonly: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: false
  },
  prefixIcon: {
    type: String,
    default: null
  },
  suffixIcon: {
    type: String,
    default: null
  },
  hint: {
    type: String,
    default: null
  },
  error: {
    type: String,
    default: null
  }
})

defineEmits(['update:modelValue', 'focus', 'blur'])
</script>

<style scoped>
.kumo-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  width: 100%;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

.kumo-input-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-default, #18181b);
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

[data-theme="dark"] .kumo-input-label {
  color: #f8fafc;
}

.kumo-required-star {
  color: #dc2626;
}

.kumo-input-container {
  display: flex;
  align-items: center;
  height: 36px;
  width: 100%;
  background-color: var(--kumo-control-bg-default, #ffffff);
  border: 1px solid var(--kumo-control-border-default, #d1d5db);
  border-radius: 6px;
  padding: 0 0.75rem;
  gap: 0.5rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
  box-sizing: border-box;
}

.kumo-input-container:focus-within {
  border-color: #0051c3;
  box-shadow: 0 0 0 1px #0051c3, 0 1px 2px rgba(0, 81, 195, 0.1);
}

.kumo-input-native {
  flex: 1;
  min-width: 0;
  height: 100%;
  border: none;
  background: transparent;
  outline: none;
  font-size: 13.5px;
  color: var(--kumo-control-fg-default, #18181b);
}

.kumo-input-native::placeholder {
  color: #a1a1aa;
}

.kumo-input-prefix,
.kumo-input-suffix {
  display: inline-flex;
  align-items: center;
  color: #71717a;
}

.kumo-input-clear-btn {
  border: none;
  background: transparent;
  color: #71717a;
  cursor: pointer;
  padding: 0 4px;
  font-size: 12px;
}

.kumo-input-hint {
  font-size: 12px;
  color: #71717a;
}

.kumo-input-error-msg {
  font-size: 12px;
  color: #dc2626;
}

.kumo-input-error .kumo-input-container {
  border-color: #dc2626;
}

.kumo-input-error .kumo-input-container:focus-within {
  box-shadow: 0 0 0 1px #dc2626;
}

.kumo-input-disabled .kumo-input-container {
  background-color: var(--kumo-surface-sunken, #f4f4f5);
  cursor: not-allowed;
  opacity: 0.7;
}
</style>
