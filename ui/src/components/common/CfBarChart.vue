<template>
  <div class="cf-chart-card">
    <div v-if="title || subtitle" class="cf-chart-head">
      <div>
        <div v-if="title" class="cf-chart-title">{{ title }}</div>
        <div v-if="subtitle" class="cf-chart-subtitle">{{ subtitle }}</div>
      </div>
      <div v-if="hoveredBar" class="cf-chart-readout">
        <span class="cf-chart-readout-label">{{ hoveredBar.label }}:</span>
        <strong class="cf-chart-readout-val">{{ hoveredBar.value }}{{ unit }}</strong>
      </div>
    </div>

    <div class="cf-bars-grid" @mouseleave="hoveredIdx = -1">
      <div
        v-for="(item, idx) in data"
        :key="item.label + idx"
        class="cf-bar-col"
        @mouseenter="hoveredIdx = idx"
      >
        <div class="cf-bar-track">
          <div
            class="cf-bar-fill"
            :style="{
              height: `${Math.max(4, (item.value / maxVal) * 100)}%`,
              backgroundColor: item.color || color,
              opacity: hoveredIdx === -1 || hoveredIdx === idx ? 1 : 0.45
            }"
          ></div>
        </div>
        <span class="cf-bar-label">{{ item.label }}</span>
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
    required: true // [{ label: string, value: number, color?: string }]
  },
  color: { type: String, default: '#f6821f' },
  unit: { type: String, default: '' }
})

const hoveredIdx = ref(-1)
const maxVal = computed(() => Math.max(1, ...props.data.map(d => d.value)))
const hoveredBar = computed(() => (hoveredIdx.value >= 0 ? props.data[hoveredIdx.value] : null))
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

.cf-bars-grid {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  height: 160px;
  padding-top: 8px;
}

.cf-bar-col {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.cf-bar-track {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  background-color: var(--color-kumo-recessed);
  border-radius: 6px;
  overflow: hidden;
}

.cf-bar-fill {
  width: 100%;
  border-radius: 4px 4px 0 0;
  transition: height 0.25s ease, opacity 0.15s ease;
}

.cf-bar-label {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-color-kumo-subtle);
  white-space: nowrap;
}
</style>
