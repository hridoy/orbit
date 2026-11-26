<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: {
    type: Number,
    default: 0
  },
  max: {
    type: Number,
    default: 100
  },
  variant: {
    type: String,
    default: 'brand', // 'brand' | 'success' | 'warning' | 'danger'
    validator: (v) => ['brand', 'success', 'warning', 'danger'].includes(v)
  },
  size: {
    type: String,
    default: 'md', // 'sm' | 'md' | 'lg'
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  },
  striped: {
    type: Boolean,
    default: false
  },
  animated: {
    type: Boolean,
    default: false
  },
  showLabel: {
    type: Boolean,
    default: false
  }
})

const percentage = computed(() => {
  if (props.max <= 0) return 0
  const pct = Math.round((props.value / props.max) * 100)
  return Math.min(100, Math.max(0, pct))
})
</script>

<template>
  <div class="cf-progress-wrap">
    <div v-if="showLabel" class="cf-progress-header">
      <slot name="label">
        <span class="cf-progress-label">Progress</span>
      </slot>
      <span class="cf-progress-value">{{ percentage }}%</span>
    </div>

    <div class="cf-progress-track" :class="[`size-${size}`]">
      <div 
        class="cf-progress-bar"
        :class="[`variant-${variant}`, { 'is-striped': striped, 'is-animated': animated }]"
        :style="{ width: `${percentage}%` }"
        role="progressbar"
        :aria-valuenow="value"
        aria-valuemin="0"
        :aria-valuemax="max"
      ></div>
    </div>
  </div>
</template>

<style scoped>
.cf-progress-wrap {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.cf-progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.cf-progress-label {
  font-weight: 500;
  color: var(--kumo-fg-default, #18181b);
}

.cf-progress-value {
  font-family: var(--font-mono, monospace);
  color: var(--kumo-fg-muted, #71717a);
}

.cf-progress-track {
  width: 100%;
  background-color: var(--kumo-surface-sunken, #f4f4f5);
  border-radius: 9999px;
  overflow: hidden;
}

.size-sm {
  height: 4px;
}

.size-md {
  height: 8px;
}

.size-lg {
  height: 12px;
}

.cf-progress-bar {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.variant-brand {
  background-color: #0051c3;
}

.variant-success {
  background-color: #059669;
}

.variant-warning {
  background-color: #d97706;
}

.variant-danger {
  background-color: #dc2626;
}

.is-striped {
  background-image: linear-gradient(
    45deg,
    rgba(255, 255, 255, 0.15) 25%,
    transparent 25%,
    transparent 50%,
    rgba(255, 255, 255, 0.15) 50%,
    rgba(255, 255, 255, 0.15) 75%,
    transparent 75%,
    transparent
  );
  background-size: 16px 16px;
}

.is-animated {
  animation: cf-progress-stripe 1s linear infinite;
}

@keyframes cf-progress-stripe {
  from {
    background-position: 16px 0;
  }
  to {
    background-position: 0 0;
  }
}
</style>
