<template>
  <div :class="['kumo-switch-wrapper', { 'is-disabled': disabled }]">
    <label class="kumo-switch-label">
      <input
        type="checkbox"
        :checked="modelValue"
        :disabled="disabled"
        class="kumo-switch-native"
        @change="$emit('update:modelValue', $event.target.checked)"
      />
      <span class="kumo-switch-track">
        <span class="kumo-switch-thumb"></span>
      </span>
      <span v-if="label || $slots.default" class="kumo-switch-text">
        <slot>{{ label }}</slot>
      </span>
    </label>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: null
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.kumo-switch-wrapper {
  display: inline-flex;
  align-items: center;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

.kumo-switch-label {
  display: inline-flex;
  align-items: center;
  gap: 0.65rem;
  cursor: pointer;
  user-select: none;
}

.kumo-switch-native {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.kumo-switch-track {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
  background-color: #e4e4e7;
  border-radius: 9999px;
  transition: background-color 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

[data-theme="dark"] .kumo-switch-track {
  background-color: #27272a;
}

.kumo-switch-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  background-color: #ffffff;
  border-radius: 50%;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.kumo-switch-native:checked + .kumo-switch-track {
  background-color: #0051c3;
}

.kumo-switch-native:checked + .kumo-switch-track .kumo-switch-thumb {
  transform: translateX(16px);
}

.kumo-switch-native:focus-visible + .kumo-switch-track {
  box-shadow: 0 0 0 2px #ffffff, 0 0 0 4px #0051c3;
}

.kumo-switch-text {
  font-size: 13.5px;
  color: var(--text-color-kumo-default, #18181b);
}

[data-theme="dark"] .kumo-switch-text {
  color: #f8fafc;
}

.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.is-disabled .kumo-switch-label {
  cursor: not-allowed;
}
</style>
