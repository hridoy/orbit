<template>
  <div class="cf-quota-card">
    <div class="cf-quota-head">
      <div>
        <div class="cf-quota-title">{{ title }}</div>
        <div v-if="subtitle" class="cf-quota-sub">{{ subtitle }}</div>
      </div>
      <CfBadge :variant="percent >= 85 ? 'warning' : 'neutral'" size="sm">{{ percent }}%</CfBadge>
    </div>

    <CfProgress :value="percent" :variant="percent >= 85 ? 'warning' : 'brand'" />

    <div class="cf-quota-foot">
      <span><strong>{{ usedLabel }}</strong> of {{ limitLabel }}</span>
      <span v-if="resetLabel" class="cf-quota-reset">{{ resetLabel }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
  used: { type: Number, required: true },
  limit: { type: Number, required: true },
  usedLabel: { type: String, required: true },
  limitLabel: { type: String, required: true },
  resetLabel: { type: String, default: '' }
})

const percent = computed(() => Math.min(100, Math.round((props.used / Math.max(1, props.limit)) * 100)))
</script>

<style scoped>
.cf-quota-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-quota-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.cf-quota-title {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-quota-sub {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-quota-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-quota-foot strong {
  color: var(--text-color-kumo-strong);
  font-family: var(--font-mono);
}
</style>
