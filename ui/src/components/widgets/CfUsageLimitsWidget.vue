<template>
  <div class="cf-usage-card">
    <div class="cf-usage-head">
      <div>
        <strong class="cf-usage-title">{{ title }}</strong>
        <div class="cf-usage-sub">{{ period === 'Today' ? 'Resets daily at 00:00 UTC' : 'Current billing cycle' }}</div>
      </div>
      <div style="display: flex; align-items: center; gap: 8px;">
        <CfTabs v-model="period" :tabs="['Today', 'This month']" />
        <CfButton variant="secondary" size="sm" @click="$emit('upgrade')">Upgrade</CfButton>
      </div>
    </div>

    <div class="cf-usage-rows">
      <div v-for="item in activeItems" :key="item.label" class="cf-usage-row">
        <svg width="36" height="36" viewBox="0 0 36 36" class="cf-ring-svg">
          <circle cx="18" cy="18" r="14" fill="none" stroke="var(--color-kumo-recessed)" stroke-width="3.5" />
          <circle
            cx="18"
            cy="18"
            r="14"
            fill="none"
            stroke="var(--color-kumo-brand)"
            stroke-width="3.5"
            :stroke-dasharray="`${ringDash(item.used, item.limit)} 88`"
            stroke-dashoffset="22"
            stroke-linecap="round"
          />
        </svg>
        <div class="cf-usage-meta">
          <span class="cf-usage-label">{{ item.label }}</span>
          <strong class="cf-usage-val">{{ item.used.toLocaleString() }} / {{ item.limit.toLocaleString() }}</strong>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: { type: String, default: 'Usage & Limits' },
  dailyItems: {
    type: Array,
    default: () => [
      { label: 'Requests today', used: 14200, limit: 100000 },
      { label: 'Observability events today', used: 48500, limit: 200000 }
    ]
  },
  monthlyItems: {
    type: Array,
    default: () => [
      { label: 'Monthly requests', used: 420000, limit: 3000000 },
      { label: 'Workers build minutes', used: 420, limit: 3000 }
    ]
  }
})

defineEmits(['upgrade'])
const period = ref('Today')
const activeItems = computed(() => (period.value === 'Today' ? props.dailyItems : props.monthlyItems))

function ringDash(used, limit) {
  const pct = Math.min(1, Math.max(0.03, used / Math.max(1, limit)))
  return (pct * 88).toFixed(1)
}
</script>

<style scoped>
.cf-usage-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-usage-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.cf-usage-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-usage-sub {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-usage-rows {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.cf-usage-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  background-color: var(--color-kumo-canvas);
  border: 1px solid var(--color-kumo-hairline);
}

.cf-ring-svg {
  flex-shrink: 0;
}

.cf-usage-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-usage-label {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-usage-val {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-color-kumo-strong);
}
</style>
