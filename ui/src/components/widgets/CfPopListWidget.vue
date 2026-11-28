<template>
  <div class="cf-pop-card">
    <div class="cf-pop-head">
      <div>
        <strong class="cf-pop-title">{{ title }}</strong>
        <div class="cf-pop-sub">Top edge locations by traffic</div>
      </div>
      <CfTabs v-model="unit" :tabs="['Count', 'RPS']" />
    </div>

    <div class="cf-pop-list">
      <div v-for="pop in pops" :key="pop.code" class="cf-pop-row">
        <div class="cf-pop-info">
          <code class="cf-pop-code">{{ pop.code }}</code>
          <span class="cf-pop-city">{{ pop.city }}</span>
          <strong class="cf-pop-val">{{ unit === 'Count' ? pop.count : pop.rps }}</strong>
        </div>
        <div class="cf-pop-bar-track">
          <div class="cf-pop-bar-fill" :style="{ width: `${pop.pct}%` }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  title: { type: String, default: 'Request Distribution by PoP' },
  pops: {
    type: Array,
    default: () => [
      { code: 'SIN', city: 'Singapore', count: '18.4k', rps: '214/s', pct: 86 },
      { code: 'IAD', city: 'Ashburn, US', count: '12.1k', rps: '140/s', pct: 62 },
      { code: 'LHR', city: 'London, UK', count: '8.7k', rps: '98/s', pct: 44 },
      { code: 'FRA', city: 'Frankfurt, DE', count: '5.2k', rps: '61/s', pct: 28 }
    ]
  }
})

const unit = ref('Count')
</script>

<style scoped>
.cf-pop-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-pop-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.cf-pop-title {
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-pop-sub {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-pop-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.cf-pop-row {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.cf-pop-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
}

.cf-pop-code {
  font-weight: 700;
}

.cf-pop-city {
  flex: 1;
  color: var(--text-color-kumo-subtle);
}

.cf-pop-val {
  font-family: var(--font-mono);
  color: var(--text-color-kumo-strong);
}

.cf-pop-bar-track {
  height: 6px;
  border-radius: 9999px;
  background-color: var(--color-kumo-recessed);
  overflow: hidden;
}

.cf-pop-bar-fill {
  height: 100%;
  border-radius: 9999px;
  background-color: var(--color-kumo-brand);
}
</style>
