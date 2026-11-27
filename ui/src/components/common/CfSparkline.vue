<template>
  <svg :viewBox="`0 0 ${width} ${height}`" class="cf-sparkline" preserveAspectRatio="none">
    <defs>
      <linearGradient :id="gradId" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" :stop-color="color" stop-opacity="0.25" />
        <stop offset="100%" :stop-color="color" stop-opacity="0" />
      </linearGradient>
    </defs>
    <path v-if="areaPath" :d="areaPath" :fill="`url(#${gradId})`" />
    <path v-if="linePath" :d="linePath" fill="none" :stroke="color" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" />
  </svg>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  values: { type: Array, default: () => [12, 18, 15, 24, 21, 30, 28] },
  color: { type: String, default: '#0051c3' },
  width: { type: Number, default: 120 },
  height: { type: Number, default: 36 }
})

const gradId = `cf-spark-${Math.random().toString(36).slice(2, 8)}`

const points = computed(() => {
  const vals = props.values
  if (!vals.length) return []
  const max = Math.max(...vals)
  const min = Math.min(...vals)
  const range = Math.max(1, max - min)
  return vals.map((v, i) => ({
    x: vals.length === 1 ? props.width / 2 : (i / (vals.length - 1)) * props.width,
    y: props.height - 4 - ((v - min) / range) * (props.height - 8)
  }))
})

const linePath = computed(() =>
  points.value.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(' ')
)

const areaPath = computed(() => {
  if (!points.value.length) return ''
  const last = points.value[points.value.length - 1]
  const first = points.value[0]
  return `${linePath.value} L ${last.x.toFixed(1)} ${props.height} L ${first.x.toFixed(1)} ${props.height} Z`
})
</script>

<style scoped>
.cf-sparkline {
  width: 100%;
  height: 36px;
  display: block;
}
</style>
