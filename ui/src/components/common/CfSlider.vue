<template>
  <div class="cf-slider-wrap">
    <div v-if="label || showValue" class="cf-slider-head">
      <label v-if="label" class="cf-slider-label">{{ label }}</label>
      <span v-if="showValue" class="cf-slider-value">{{ modelValue }}{{ unit }}</span>
    </div>
    <input
      type="range"
      :min="min"
      :max="max"
      :step="step"
      :value="modelValue"
      :disabled="disabled"
      class="cf-slider-native"
      @input="$emit('update:modelValue', Number($event.target.value))"
    />
    <div v-if="showRange" class="cf-slider-bounds">
      <span>{{ min }}{{ unit }}</span>
      <span>{{ max }}{{ unit }}</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Number,
    default: 0
  },
  min: {
    type: Number,
    default: 0
  },
  max: {
    type: Number,
    default: 100
  },
  step: {
    type: Number,
    default: 1
  },
  label: {
    type: String,
    default: ''
  },
  unit: {
    type: String,
    default: ''
  },
  showValue: {
    type: Boolean,
    default: true
  },
  showRange: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-slider-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-slider-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cf-slider-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-slider-value {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-strong);
  border: 1px solid var(--color-kumo-line);
}

.cf-slider-native {
  width: 100%;
  accent-color: var(--color-kumo-brand);
  cursor: pointer;
}

.cf-slider-native:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.cf-slider-bounds {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text-color-kumo-subtle);
}
</style>
