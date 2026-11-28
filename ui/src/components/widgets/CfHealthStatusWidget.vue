<template>
  <div class="cf-health-card">
    <div class="cf-health-head">
      <div style="display: flex; align-items: center; gap: 8px;">
        <span class="cf-health-dot" :class="`is-${overall}`"></span>
        <strong class="cf-health-title">{{ title }}</strong>
      </div>
      <CfBadge :variant="overall === 'operational' ? 'success' : 'warning'" size="sm">
        {{ overallLabel }}
      </CfBadge>
    </div>

    <div class="cf-health-services">
      <div v-for="svc in services" :key="svc.name" class="cf-health-row">
        <span class="cf-health-svc-name">{{ svc.name }}</span>
        <span class="cf-health-svc-status">
          <span class="cf-health-dot sm" :class="`is-${svc.status || 'operational'}`"></span>
          <span>{{ svc.uptime || '99.99%' }}</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Global Edge Network' },
  overall: { type: String, default: 'operational' }, // operational | degraded
  overallLabel: { type: String, default: 'All Systems Operational' },
  services: {
    type: Array,
    default: () => [
      { name: 'Anycast DNS', status: 'operational', uptime: '100%' },
      { name: 'Workers & Pages', status: 'operational', uptime: '99.99%' },
      { name: 'R2 Object Storage', status: 'operational', uptime: '99.98%' },
      { name: 'Zero Trust Access', status: 'operational', uptime: '100%' }
    ]
  }
})
</script>

<style scoped>
.cf-health-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-health-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.cf-health-title {
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-health-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background-color: var(--color-kumo-success);
}

.cf-health-dot.sm {
  width: 7px;
  height: 7px;
}

.cf-health-dot.is-degraded {
  background-color: var(--color-kumo-warning);
}

.cf-health-services {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--color-kumo-hairline);
}

.cf-health-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12.5px;
}

.cf-health-svc-name {
  color: var(--text-color-kumo-default);
}

.cf-health-svc-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}
</style>
