<script setup>
import { computed, useAttrs } from 'vue'
defineOptions({ inheritAttrs: false })
const props = defineProps({
  modelValue: {
    type: [String, Number, Boolean],
    default: ''
  },
  options: {
    type: Array,
    default: () => []
    // string[] or { value, label, disabled }[]
  },
  placeholder: {
    type: String,
    default: 'Select an option...'
  },
  label: {
    type: String,
    default: ''
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
  size: {
    type: String,
    default: 'md', // 'sm' | 'md' | 'lg'
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

function onChange(e) {
  emit('update:modelValue', e.target.value)
  emit('change', e.target.value)
}

function getOptionValue(opt) {
  return typeof opt === 'object' && opt !== null ? opt.value : opt
}

function getOptionLabel(opt) {
  return typeof opt === 'object' && opt !== null ? opt.label : opt
}

function isOptionDisabled(opt) {
  return typeof opt === 'object' && opt !== null ? !!opt.disabled : false
}

const attrs = useAttrs()
const selectAttrs = computed(() => {
  const { class: _c, style: _s, ...rest } = attrs
  return rest
})
</script>

<template>
  <div class="cf-select-wrap" :class="{ 'is-disabled': disabled, 'has-error': !!error }">
    <label v-if="label" class="cf-select-label">
      {{ label }}
    </label>

    <div class="cf-select-box" :class="[`cf-size-${size}`]">
      <select
        :value="modelValue"
        :disabled="disabled"
        class="cf-select-native"
        v-bind="selectAttrs"
        @change="onChange"
      >
        <option v-if="placeholder" value="" disabled :selected="!modelValue">
          {{ placeholder }}
        </option>
        <template v-for="(opt, idx) in options" :key="idx">
          <option
            :value="getOptionValue(opt)"
            :disabled="isOptionDisabled(opt)"
          >
            {{ getOptionLabel(opt) }}
          </option>
        </template>
      </select>
      <svg class="cf-select-arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="m6 9 6 6 6-6" />
      </svg>
    </div>

    <p v-if="error" class="cf-select-msg is-error">{{ error }}</p>
    <p v-else-if="hint" class="cf-select-msg is-hint">{{ hint }}</p>
  </div>
</template>

<style scoped>
.cf-select-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-select-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--kumo-fg-default, #18181b);
}

.cf-select-box {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  border-radius: var(--kumo-radius-md, 6px);
  background-color: var(--kumo-control-bg-default, #ffffff);
  border: 1px solid var(--kumo-control-border-default, #d4d4d8);
  transition: all 0.15s ease;
}

.cf-select-box:hover:not(.is-disabled) {
  border-color: var(--kumo-control-border-hover, #a1a1aa);
}

.cf-select-box:focus-within {
  border-color: #0051c3;
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.15);
}

.has-error .cf-select-box {
  border-color: #dc2626 !important;
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12) !important;
}

.cf-size-sm {
  height: 28px;
  font-size: 12px;
}

.cf-size-md {
  height: 36px;
  font-size: 13px;
}

.cf-size-lg {
  height: 42px;
  font-size: 14px;
}

.cf-select-native {
  width: 100%;
  height: 100%;
  padding: 0 32px 0 12px;
  background: transparent;
  background-image: none;
  border: none;
  font-family: inherit;
  font-size: inherit;
  color: var(--kumo-control-fg-default, #18181b);
  cursor: pointer;
  -webkit-appearance: none !important;
  -moz-appearance: none !important;
  appearance: none !important;
  outline: none;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  box-sizing: border-box;
}

.cf-select-native::-ms-expand {
  display: none;
}

.cf-select-arrow {
  position: absolute;
  right: 10px;
  pointer-events: none;
  color: var(--kumo-fg-muted, #71717a);
}

.is-disabled {
  opacity: 0.6;
  pointer-events: none;
}

.cf-select-msg {
  margin: 0;
  font-size: 12px;
  line-height: 16px;
}

.cf-select-msg.is-error {
  color: #dc2626;
}

.cf-select-msg.is-hint {
  color: var(--kumo-fg-muted, #71717a);
}
</style>
