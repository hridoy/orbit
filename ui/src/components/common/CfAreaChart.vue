<template>
  <div class="cf-chart-card">
    <div v-if="title || subtitle" class="cf-chart-head">
      <div>
        <div v-if="title" class="cf-chart-title">{{ title }}</div>
        <div v-if="subtitle" class="cf-chart-subtitle">{{ subtitle }}</div>
      </div>
      <div v-if="hoveredPoint" class="cf-chart-readout">
        <span class="cf-chart-readout-label">{{ hoveredPoint.label }}</span>
        <strong class="cf-chart-readout-val">{{ hoveredPoint.value }}{{ unit }}</strong>
      </div>
    </div>

    <div class="cf-chart-wrap" @mouseleave="hoveredIdx = -1">
      <svg :viewBox="`0 0 ${width} ${height}`" class="cf-chart-svg" preserveAspectRatio="none">
        <defs>
          <linearGradient :id="gradId" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" :stop-color="color" stop-opacity="0.28" />
            <stop offset="100%" :stop-color="color" stop-opacity="0.01" />
          </linearGradient>
        </defs>

        <!-- Horizontal Gridlines -->
        <g class="cf-chart-grid">
          <line
            v-for="g in 4"
            :key="g"
            :x1="padL"
            :x2="width - padR"
            :y1="padT + ((g - 1) / 3) * plotH"
            :y2="padT + ((g - 1) / 3) * plotH"
          />
        </g>

        <!-- Area & Line -->
        <path v-if="areaPath" :d="areaPath" :fill="`url(#${gradId})`" />
        <path v-if="linePath" :d="linePath" fill="none" :stroke="color" stroke-width="2" stroke-linejoin="round" stroke-linecap="round" />

        <!-- Hover Crosshair & Dot -->
        <line
          v-if="hoveredPoint"
          :x1="hoveredPoint.x"
          :x2="hoveredPoint.x"
          :y1="padT"
          :y2="padT + plotH"
          stroke="var(--color-kumo-line)"
          stroke-dasharray="3 3"
        />
        <circle
          v-if="hoveredPoint"
          :cx="hoveredPoint.x"
          :cy="hoveredPoint.y"
          r="4.5"
          :fill="color"
          stroke="var(--color-kumo-base)"
          stroke-width="2"
        />

        <!-- Invisible Hover Columns -->
        <rect
          v-for="(pt, idx) in points"
          :key="idx"
          :x="pt.x - colW / 2"
          :y="padT"
          :width="colW"
          :height="plotH"
          fill="transparent"
          @mouseenter="hoveredIdx = idx"
        />
      </svg>

      <!-- X-Axis Labels -->
      <div class="cf-chart-xaxis">
        <span v-for="(pt, idx) in xLabels" :key="idx">{{ pt }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  data: {
    type: Array,
    required: true // [{ label: string, value: number }]
  },
  color: { type: String, default: '#0051c3' },
  unit: { type: String, default: '' }
})

const gradId = `cf-area-grad-${Math.random().toString(36).slice(2, 8)}`
const hoveredIdx = ref(-1)

const width = 560
const height = 180
const padL = 8
const padR = 8
const padT = 12
const padB = 8
const plotW = width - padL - padR
const plotH = height - padT - padB

const maxVal = computed(() => Math.max(1, ...props.data.map(d => d.value)))
const colW = computed(() => (props.data.length > 1 ? plotW / (props.data.length - 1) : plotW))

const points = computed(() => {
  const len = props.data.length
  if (!len) return []
  return props.data.map((d, i) => {
    const x = len === 1 ? padL + plotW / 2 : padL + (i / (len - 1)) * plotW
    const y = padT + plotH - (d.value / maxVal.value) * (plotH - 6)
    return { ...d, x, y }
  })
})

const linePath = computed(() => {
  if (!points.value.length) return ''
  return points.value.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(' ')
})

const areaPath = computed(() => {
  if (!points.value.length) return ''
  const first = points.value[0]
  const last = points.value[points.value.length - 1]
  const baseY = padT + plotH
  return `${linePath.value} L ${last.x.toFixed(1)} ${baseY} L ${first.x.toFixed(1)} ${baseY} Z`
})

const hoveredPoint = computed(() => (hoveredIdx.value >= 0 ? points.value[hoveredIdx.value] : null))

const xLabels = computed(() => {
  const len = props.data.length
  if (len <= 4) return props.data.map(d => d.label)
  return [props.data[0].label, props.data[Math.floor(len / 2)].label, props.data[len - 1].label]
})
</script>

<style scoped>
.cf-chart-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-chart-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  min-height: 36px;
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

.cf-chart-readout {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 8px;
  border-radius: 6px;
  background-color: var(--color-kumo-recessed);
  border: 1px solid var(--color-kumo-line);
  font-size: 12px;
}

.cf-chart-readout-label {
  color: var(--text-color-kumo-subtle);
}

.cf-chart-readout-val {
  font-family: var(--font-mono);
  color: var(--text-color-kumo-strong);
}

.cf-chart-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cf-chart-svg {
  width: 100%;
  height: 160px;
  display: block;
}

.cf-chart-grid line {
  stroke: var(--color-kumo-hairline);
  stroke-width: 1;
}

.cf-chart-xaxis {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text-color-kumo-subtle);
  font-family: var(--font-mono);
}
</style>
