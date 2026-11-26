<template>
  <div class="metric-card">
    <div class="metric-header">
      <span class="metric-label">{{ label }}</span>
      <span v-if="delta" class="metric-delta" :class="isPositive ? 'positive' : 'negative'">{{ delta }}</span>
    </div>
    <div class="metric-value">{{ value }}</div>

    <svg v-if="sparkline" width="100%" height="24" viewBox="0 0 100 24" fill="none" style="margin: 4px 0;">
      <path :d="sparkline" :stroke="strokeColor || 'var(--cf-brand-orange)'" stroke-width="2" fill="none" stroke-linecap="round" />
    </svg>

    <div v-if="footer" class="metric-footer">{{ footer }}</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: String,
  value: [String, Number],
  delta: String,
  sparkline: String,
  strokeColor: String,
  footer: String
})

const isPositive = computed(() => {
  return props.delta && props.delta.includes('↗')
})
</script>
