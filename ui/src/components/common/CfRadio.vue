<script setup>
const props = defineProps({
  modelValue: {
    type: [String, Number, Boolean],
    default: ''
  },
  value: {
    type: [String, Number, Boolean],
    required: true
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
  name: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

function onChange() {
  emit('update:modelValue', props.value)
  emit('change', props.value)
}
</script>

<template>
  <label class="cf-radio-wrap" :class="{ 'is-disabled': disabled }">
    <span class="cf-radio-control">
      <input
        type="radio"
        :name="name"
        :value="value"
        :checked="modelValue === value"
        :disabled="disabled"
        class="cf-radio-input"
        @change="onChange"
      />
      <span class="cf-radio-circle" :class="{ 'is-checked': modelValue === value }">
        <span v-if="modelValue === value" class="cf-radio-dot"></span>
      </span>
    </span>

    <span v-if="label || $slots.default || description" class="cf-radio-text">
      <span class="cf-radio-label">
        <slot>{{ label }}</slot>
      </span>
      <span v-if="description" class="cf-radio-desc">{{ description }}</span>
    </span>
  </label>
</template>

<style scoped>
.cf-radio-wrap {
  display: inline-flex;
  align-items: flex-start;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.cf-radio-control {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-top: 2px;
}

.cf-radio-input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  margin: 0;
  cursor: pointer;
}

.cf-radio-circle {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1px solid var(--kumo-control-border-default, #d4d4d8);
  background-color: var(--kumo-control-bg-default, #ffffff);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.12s ease;
}

.cf-radio-wrap:hover .cf-radio-circle:not(.is-checked) {
  border-color: var(--kumo-control-border-hover, #a1a1aa);
}

.cf-radio-input:focus-visible + .cf-radio-circle {
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.2);
}

.cf-radio-circle.is-checked {
  border-color: #0051c3;
  background-color: #0051c3;
}

.cf-radio-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #ffffff;
}

.cf-radio-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-radio-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--kumo-fg-default, #18181b);
  line-height: 20px;
}

.cf-radio-desc {
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
