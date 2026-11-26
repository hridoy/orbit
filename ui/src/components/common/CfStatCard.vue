<template>
  <div class="kumo-stat-card">
    <div class="kumo-stat-header">
      <span class="kumo-stat-label">
        <span>{{ label }}</span>
        <span v-if="tooltip" class="kumo-stat-info-icon" :title="tooltip">
          <CfIcon name="cis-circle-info" size="13" />
        </span>
      </span>
      <div v-if="$slots.action" class="kumo-stat-action">
        <slot name="action" />
      </div>
    </div>

    <div class="kumo-stat-body">
      <div class="kumo-stat-main">
        <span class="kumo-stat-value tabular-nums">{{ value }}</span>
        <span v-if="unit" class="kumo-stat-unit">{{ unit }}</span>
      </div>

      <div v-if="delta || deltaText" :class="['kumo-stat-delta', `is-${deltaType}`]">
        <span v-if="deltaType === 'up'">▲</span>
        <span v-else-if="deltaType === 'down'">▼</span>
        <span>{{ delta || deltaText }}</span>
      </div>
    </div>

    <!-- Optional Sparkline SVG -->
    <div v-if="sparklineData && sparklineData.length" class="kumo-stat-sparkline">
      <svg width="100%" height="24" viewBox="0 0 100 24" preserveAspectRatio="none">
        <path :d="sparklinePath" fill="none" :stroke="sparklineColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  unit: {
    type: String,
    default: null
  },
  tooltip: {
    type: String,
    default: null
  },
  delta: {
    type: String,
    default: null
  },
  deltaText: {
    type: String,
    default: null
  },
  deltaType: {
    type: String,
    default: 'up', // up | down | neutral
    validator: (v) => ['up', 'down', 'neutral'].includes(v)
  },
  sparklineData: {
    type: Array,
    default: () => []
  }
})

const sparklineColor = computed(() => {
  if (props.deltaType === 'down') return '#dc2626'
  return '#0051c3'
})

const sparklinePath = computed(() => {
  const pts = props.sparklineData
  if (!pts || pts.length < 2) return ''
  const min = Math.min(...pts)
  const max = Math.max(...pts)
  const range = max - min || 1
  return pts.map((val, i) => {
    const x = (i / (pts.length - 1)) * 100
    const y = 22 - ((val - min) / range) * 20
    return `${i === 0 ? 'M' : 'L'} ${x.toFixed(1)} ${y.toFixed(1)}`
  }).join(' ')
})
</script>

<style scoped>
.kumo-stat-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  background-color: var(--color-kumo-base, #ffffff);
  border: 1px solid var(--color-kumo-line, #e2e8f0);
  border-radius: 8px;
  box-shadow: var(--shadow-kumo-card, 0 1px 2px rgba(0, 0, 0, 0.04));
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  transition: all 0.15s ease;
  min-height: 100px;
}

.kumo-stat-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.kumo-stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.kumo-stat-label {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-subtle, #71717a);
}

.kumo-stat-info-icon {
  color: #a1a1aa;
  cursor: help;
  line-height: 1;
}

.kumo-stat-info-icon:hover {
  color: #71717a;
}

.kumo-stat-body {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.5rem;
}

.kumo-stat-main {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
}

.kumo-stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-color-kumo-default, #18181b);
  letter-spacing: -0.02em;
}

[data-theme="dark"] .kumo-stat-value {
  color: #f8fafc;
}

.kumo-stat-unit {
  font-size: 13px;
  color: #71717a;
  font-weight: 500;
}

.kumo-stat-delta {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 11.5px;
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 9999px;
  line-height: 1;
  white-space: nowrap;
  flex-shrink: 0;
}

.kumo-stat-delta.is-up {
  background-color: rgba(16, 185, 129, 0.1);
  color: #059669;
}

.kumo-stat-delta.is-down {
  background-color: rgba(220, 38, 38, 0.1);
  color: #dc2626;
}

.kumo-stat-delta.is-neutral {
  background-color: #f4f4f5;
  color: #71717a;
}

.kumo-stat-sparkline {
  margin-top: 0.75rem;
  padding-top: 0.25rem;
}
</style>
