<template>
  <div class="cf-chart-card">
    <div v-if="title || subtitle" class="cf-chart-head">
      <div v-if="title" class="cf-chart-title">{{ title }}</div>
      <div v-if="subtitle" class="cf-chart-subtitle">{{ subtitle }}</div>
    </div>

    <div class="cf-donut-layout">
      <div class="cf-donut-svg-wrap">
        <svg viewBox="0 0 120 120" class="cf-donut-svg">
          <circle cx="60" cy="60" r="46" fill="none" stroke="var(--color-kumo-recessed)" stroke-width="14" />
          <circle
            v-for="(seg, idx) in computedSegments"
            :key="seg.label + idx"
            cx="60"
            cy="60"
            r="46"
            fill="none"
            :stroke="seg.color"
            stroke-width="14"
            :stroke-dasharray="`${seg.dash} ${circumference - seg.dash}`"
            :stroke-dashoffset="seg.offset"
            stroke-linecap="butt"
            style="transition: opacity 0.15s ease; cursor: pointer;"
            :style="{ opacity: hoveredIdx === -1 || hoveredIdx === idx ? 1 : 0.4 }"
            @mouseenter="hoveredIdx = idx"
            @mouseleave="hoveredIdx = -1"
          />
        </svg>
        <div class="cf-donut-center">
          <strong class="cf-donut-total">{{ centerValue || totalLabel }}</strong>
          <span class="cf-donut-caption">{{ centerLabel || 'Total' }}</span>
        </div>
      </div>

      <div class="cf-donut-legend">
        <div
          v-for="(seg, idx) in computedSegments"
          :key="seg.label + idx"
          class="cf-legend-row"
          :class="{ active: hoveredIdx === idx }"
          @mouseenter="hoveredIdx = idx"
          @mouseleave="hoveredIdx = -1"
        >
          <span class="cf-legend-dot" :style="{ backgroundColor: seg.color }"></span>
          <span class="cf-legend-label">{{ seg.label }}</span>
          <span class="cf-legend-val">{{ seg.percent }}%</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  segments: {
    type: Array,
    required: true // [{ label: string, value: number, color: string }]
  },
  centerValue: { type: String, default: '' },
  centerLabel: { type: String, default: '' }
})

const hoveredIdx = ref(-1)
const circumference = 2 * Math.PI * 46

const total = computed(() => Math.max(1, props.segments.reduce((sum, s) => sum + s.value, 0)))
const totalLabel = computed(() => total.value.toLocaleString())

const computedSegments = computed(() => {
  let acc = 0
  return props.segments.map((s) => {
    const frac = s.value / total.value
    const dash = frac * circumference
    // Start at top (-90deg = circumference * 0.25)
    const offset = circumference * 0.25 - acc
    acc += dash
    return {
      ...s,
      dash,
      offset,
      percent: Math.round(frac * 100)
    }
  })
})
</script>

<style scoped>
.cf-chart-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 14px 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-chart-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-chart-subtitle {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-donut-layout {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.cf-donut-svg-wrap {
  position: relative;
  width: 132px;
  height: 132px;
  flex-shrink: 0;
}

.cf-donut-svg {
  width: 100%;
  height: 100%;
  display: block;
}

.cf-donut-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.cf-donut-total {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-color-kumo-strong);
  line-height: 1.1;
}

.cf-donut-caption {
  font-size: 11px;
  color: var(--text-color-kumo-subtle);
}

.cf-donut-legend {
  flex: 1;
  min-width: 150px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cf-legend-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
  padding: 4px 6px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.cf-legend-row.active,
.cf-legend-row:hover {
  background-color: var(--color-kumo-tint);
}

.cf-legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.cf-legend-label {
  flex: 1;
  color: var(--text-color-kumo-default);
}

.cf-legend-val {
  font-family: var(--font-mono);
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}
</style>
