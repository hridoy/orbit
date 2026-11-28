<template>
  <div class="cf-timerange">
    <CfIcon name="cis-clock" size="14" class="cf-timerange-icon" />
    <select
      :value="modelValue"
      class="cf-timerange-select"
      @change="$emit('update:modelValue', $event.target.value)"
    >
      <option v-for="opt in presets" :key="opt.value" :value="opt.value">
        {{ opt.label }}
      </option>
    </select>
    <svg class="cf-timerange-caret" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="m6 9 6 6 6-6" />
    </svg>
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: String, default: '24h' },
  presets: {
    type: Array,
    default: () => [
      { value: '30m', label: 'Last 30 minutes' },
      { value: '6h', label: 'Last 6 hours' },
      { value: '24h', label: 'Last 24 hours' },
      { value: '7d', label: 'Last 7 days' },
      { value: '30d', label: 'Last 30 days' }
    ]
  }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-timerange {
  position: relative;
  display: inline-flex;
  align-items: center;
  height: 36px;
  padding: 0 28px 0 32px;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  color: var(--text-color-kumo-default);
}

.cf-timerange:hover {
  border-color: var(--kumo-control-border-hover);
  background-color: var(--kumo-control-bg-hover);
}

.cf-timerange-icon {
  position: absolute;
  left: 10px;
  color: var(--text-color-kumo-subtle);
  pointer-events: none;
}

.cf-timerange-select {
  appearance: none;
  -webkit-appearance: none;
  border: none;
  background: transparent;
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 500;
  color: inherit;
  cursor: pointer;
  outline: none;
}

.cf-timerange-caret {
  position: absolute;
  right: 9px;
  color: var(--text-color-kumo-subtle);
  pointer-events: none;
}
</style>
