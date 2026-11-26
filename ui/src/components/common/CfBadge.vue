<template>
  <span
    :class="[
      'kumo-badge',
      `kumo-badge-${variant}`,
      `kumo-badge-${size}`,
      { 'kumo-badge-dot': dot }
    ]"
  >
    <span v-if="dot" class="kumo-badge-dot-indicator"></span>
    <CfIcon v-if="icon" :name="icon" :size="iconSize" class="kumo-badge-icon" />
    <slot />
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'neutral', // neutral | success | warning | danger | brand | info | beta
    validator: (v) => ['neutral', 'success', 'warning', 'danger', 'brand', 'info', 'beta'].includes(v)
  },
  size: {
    type: String,
    default: 'md', // sm | md
    validator: (v) => ['sm', 'md'].includes(v)
  },
  dot: {
    type: Boolean,
    default: false
  },
  icon: {
    type: String,
    default: null
  }
})

const iconSize = computed(() => (props.size === 'sm' ? 11 : 12))
</script>

<style scoped>
.kumo-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-family: var(--font-sans);
  font-weight: 500;
  letter-spacing: -0.01em;
  border-radius: 9999px;
  line-height: 1;
  white-space: nowrap;
  box-sizing: border-box;
}

.kumo-badge-sm {
  padding: 2px 6px;
  font-size: 11px;
}

.kumo-badge-md {
  padding: 3px 8px;
  font-size: 12px;
}

.kumo-badge-dot-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: currentColor;
}

/* VARIANTS */
/* Neutral / Default */
.kumo-badge-neutral {
  background-color: #f4f4f5;
  color: #3f3f46;
  border: 1px solid #e4e4e7;
}

[data-theme="dark"] .kumo-badge-neutral {
  background-color: #18181b;
  color: #d4d4d8;
  border-color: #27272a;
}

/* Success / Active */
.kumo-badge-success {
  background-color: rgba(16, 185, 129, 0.1);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.25);
}

[data-theme="dark"] .kumo-badge-success {
  background-color: rgba(16, 185, 129, 0.15);
  color: #34d399;
}

/* Warning */
.kumo-badge-warning {
  background-color: rgba(245, 158, 11, 0.1);
  color: #b45309;
  border: 1px solid rgba(245, 158, 11, 0.25);
}

/* Danger */
.kumo-badge-danger {
  background-color: rgba(239, 68, 68, 0.1);
  color: #dc2626;
  border: 1px solid rgba(239, 68, 68, 0.25);
}

/* Brand Orange */
.kumo-badge-brand {
  background-color: rgba(246, 130, 31, 0.1);
  color: #f6821f;
  border: 1px solid rgba(246, 130, 31, 0.25);
}

/* Info Blue */
.kumo-badge-info {
  background-color: rgba(0, 81, 195, 0.08);
  color: #0051c3;
  border: 1px solid rgba(0, 81, 195, 0.2);
}

/* Beta / Indigo */
.kumo-badge-beta {
  background-color: rgba(99, 102, 241, 0.1);
  color: #4f46e5;
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.kumo-badge-icon {
  display: inline-flex;
  align-items: center;
}
</style>
