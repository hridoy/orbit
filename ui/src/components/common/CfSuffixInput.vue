<template>
  <div class="cf-suffix-wrap">
    <label v-if="label" class="cf-suffix-label">{{ label }}</label>
    <div class="cf-suffix-box">
      <span v-if="prefix" class="cf-affix is-prefix">{{ prefix }}</span>
      <input
        :value="modelValue"
        type="text"
        :placeholder="placeholder"
        class="cf-suffix-input"
        @input="$emit('update:modelValue', $event.target.value)"
      />
      <span v-if="suffix" class="cf-affix is-suffix">{{ suffix }}</span>
    </div>
    <span v-if="hint" class="cf-suffix-hint">{{ hint }}</span>
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  prefix: { type: String, default: '' },
  suffix: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  hint: { type: String, default: '' }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-suffix-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-suffix-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-suffix-box {
  display: flex;
  align-items: stretch;
  height: 36px;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  overflow: hidden;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.cf-suffix-box:focus-within {
  border-color: var(--color-kumo-brand);
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.15);
}

.cf-suffix-input {
  flex: 1;
  min-width: 0;
  border: none;
  background: transparent;
  padding: 0 12px;
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-color-kumo-strong);
  outline: none;
}

.cf-affix {
  display: inline-flex;
  align-items: center;
  padding: 0 10px;
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-subtle);
  font-family: var(--font-mono);
  font-size: 12px;
  white-space: nowrap;
  user-select: none;
}

.cf-affix.is-prefix {
  border-right: 1px solid var(--kumo-control-border-default);
}

.cf-affix.is-suffix {
  border-left: 1px solid var(--kumo-control-border-default);
}

.cf-suffix-hint {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}
</style>
