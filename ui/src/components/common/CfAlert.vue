<template>
  <div :class="['kumo-alert', `kumo-alert-${variant}`]">
    <div class="kumo-alert-icon">
      <slot name="icon">
        <CfIcon :name="defaultIcon" size="16" />
      </slot>
    </div>
    <div class="kumo-alert-content">
      <span v-if="title" class="kumo-alert-title">{{ title }}</span>
      <div class="kumo-alert-message">
        <slot />
      </div>
    </div>
    <button
      v-if="dismissible"
      type="button"
      class="kumo-alert-dismiss"
      aria-label="Dismiss alert"
      @click="$emit('dismiss')"
    >
      ✕
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'info', // info | success | warning | danger
    validator: (v) => ['info', 'success', 'warning', 'danger'].includes(v)
  },
  title: {
    type: String,
    default: null
  },
  dismissible: {
    type: Boolean,
    default: false
  }
})

defineEmits(['dismiss'])

const defaultIcon = computed(() => {
  if (props.variant === 'success') return 'cis-checkmark-1-medium'
  if (props.variant === 'warning') return 'cis-clock-alert'
  if (props.variant === 'danger') return 'cis-trash-can'
  return 'cis-circle-info'
})
</script>

<style scoped>
.kumo-alert {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border-radius: 6px;
  font-family: var(--font-sans);
  font-size: 13.5px;
  line-height: 1.45;
  letter-spacing: -0.01em;
  box-sizing: border-box;
}

.kumo-alert-icon {
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.kumo-alert-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.kumo-alert-title {
  font-weight: 600;
  color: inherit;
}

.kumo-alert-message {
  color: inherit;
}

.kumo-alert-dismiss {
  border: none;
  background: transparent;
  color: inherit;
  opacity: 0.6;
  cursor: pointer;
  padding: 0;
  font-size: 13px;
  margin-top: 2px;
}

.kumo-alert-dismiss:hover {
  opacity: 1;
}

/* VARIANTS */
.kumo-alert-info {
  background-color: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #1e40af;
}

[data-theme="dark"] .kumo-alert-info {
  background-color: rgba(30, 64, 175, 0.15);
  border-color: rgba(30, 64, 175, 0.35);
  color: #93c5fd;
}

.kumo-alert-success {
  background-color: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #166534;
}

[data-theme="dark"] .kumo-alert-success {
  background-color: rgba(22, 101, 52, 0.15);
  border-color: rgba(22, 101, 52, 0.35);
  color: #86efac;
}

.kumo-alert-warning {
  background-color: #fffbeb;
  border: 1px solid #fde68a;
  color: #92400e;
}

[data-theme="dark"] .kumo-alert-warning {
  background-color: rgba(146, 64, 14, 0.15);
  border-color: rgba(146, 64, 14, 0.35);
  color: #fcd34d;
}

.kumo-alert-danger {
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}

[data-theme="dark"] .kumo-alert-danger {
  background-color: rgba(153, 27, 27, 0.15);
  border-color: rgba(153, 27, 27, 0.35);
  color: #fca5a5;
}
</style>
