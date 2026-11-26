<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  indeterminate: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

function onChange(e) {
  emit('update:modelValue', e.target.checked)
  emit('change', e.target.checked)
}
</script>

<template>
  <label class="cf-checkbox-wrap" :class="{ 'is-disabled': disabled }">
    <span class="cf-checkbox-control">
      <input
        type="checkbox"
        :checked="modelValue"
        :disabled="disabled"
        :indeterminate="indeterminate"
        class="cf-checkbox-input"
        @change="onChange"
      />
      <span class="cf-checkbox-box" :class="{ 'is-checked': modelValue, 'is-indeterminate': indeterminate }">
        <CfIcon v-if="indeterminate" name="minus" size="12" class="cf-check-icon" />
        <CfIcon v-else-if="modelValue" name="check" size="12" class="cf-check-icon" />
      </span>
    </span>

    <span v-if="label || $slots.default || description" class="cf-checkbox-text">
      <span class="cf-checkbox-label">
        <slot>{{ label }}</slot>
      </span>
      <span v-if="description" class="cf-checkbox-desc">{{ description }}</span>
    </span>
  </label>
</template>

<style scoped>
.cf-checkbox-wrap {
  display: inline-flex;
  align-items: flex-start;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.cf-checkbox-control {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-top: 2px;
}

.cf-checkbox-input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  margin: 0;
  cursor: pointer;
}

.cf-checkbox-box {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 1px solid var(--kumo-control-border-default, #d4d4d8);
  background-color: var(--kumo-control-bg-default, #ffffff);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.12s ease;
}

.cf-checkbox-wrap:hover .cf-checkbox-box:not(.is-checked):not(.is-indeterminate) {
  border-color: var(--kumo-control-border-hover, #a1a1aa);
}

.cf-checkbox-input:focus-visible + .cf-checkbox-box {
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.2);
}

.cf-checkbox-box.is-checked,
.cf-checkbox-box.is-indeterminate {
  background-color: #0051c3;
  border-color: #0051c3;
  color: #ffffff;
}

.cf-check-icon {
  stroke-width: 2.5;
}

.cf-checkbox-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-checkbox-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--kumo-fg-default, #18181b);
  line-height: 20px;
}

.cf-checkbox-desc {
  font-size: 12px;
  color: var(--kumo-fg-muted, #71717a);
  line-height: 16px;
}

.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}
</style>
