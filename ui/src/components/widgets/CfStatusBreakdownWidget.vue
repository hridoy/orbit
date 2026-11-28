<template>
  <div class="cf-breakdown-card">
    <div class="cf-breakdown-head">
      <strong class="cf-breakdown-title">{{ title }}</strong>
      <span class="cf-breakdown-total">{{ totalLabel }}</span>
    </div>

    <div class="cf-proportion-bar">
      <div
        v-for="(seg, idx) in normalized"
        :key="seg.label + idx"
        class="cf-proportion-seg"
        :style="{ width: `${seg.percent}%`, backgroundColor: seg.color }"
        :title="`${seg.label}: ${seg.percent}%`"
      ></div>
    </div>

    <div class="cf-breakdown-legend">
      <div v-for="(seg, idx) in normalized" :key="seg.label + idx" class="cf-bd-item">
        <span class="cf-bd-dot" :style="{ backgroundColor: seg.color }"></span>
        <span class="cf-bd-label">{{ seg.label }}</span>
        <strong class="cf-bd-pct">{{ seg.percent }}%</strong>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: { type: String, default: 'Traffic Breakdown' },
  totalLabel: { type: String, default: '' },
  segments: {
    type: Array,
    required: true // [{ label, value, color }]
  }
})

const total = computed(() => Math.max(1, props.segments.reduce((s, i) => s + i.value, 0)))

const normalized = computed(() =>
  props.segments.map(s => ({
    ...s,
    percent: Math.round((s.value / total.value) * 100)
  }))
)
</script>

<style scoped>
.cf-breakdown-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-breakdown-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cf-breakdown-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-breakdown-total {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-proportion-bar {
  display: flex;
  height: 10px;
  border-radius: 9999px;
  overflow: hidden;
  background-color: var(--color-kumo-recessed);
  gap: 2px;
}

.cf-proportion-seg {
  height: 100%;
  transition: width 0.25s ease;
}

.cf-breakdown-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
}

.cf-bd-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.cf-bd-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.cf-bd-label {
  color: var(--text-color-kumo-subtle);
}

.cf-bd-pct {
  font-family: var(--font-mono);
  color: var(--text-color-kumo-strong);
}
</style>
