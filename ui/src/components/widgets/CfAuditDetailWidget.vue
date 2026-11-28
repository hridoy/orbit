<template>
  <div class="cf-audit-card">
    <div class="cf-audit-head">
      <div>
        <div class="cf-audit-title-row">
          <strong>{{ event.action }}</strong>
          <CfBadge :variant="event.status === 'Success' ? 'success' : 'danger'" size="sm" dot>
            {{ event.status }}
          </CfBadge>
        </div>
        <div class="cf-audit-sub">{{ event.actor }} · {{ event.timestamp }}</div>
      </div>
      <CfTabs v-model="viewMode" :tabs="['Structured', 'JSON']" />
    </div>

    <div v-if="viewMode === 'Structured'" class="cf-audit-fields">
      <div v-for="(val, key) in event.metadata" :key="key" class="cf-audit-field-row">
        <span class="cf-af-key">{{ key }}</span>
        <code class="cf-af-val">{{ val }}</code>
      </div>
    </div>

    <div v-else class="cf-audit-json">
      <CfCodeBlock :code="JSON.stringify(event, null, 2)" lang="json" filename="audit-event.json" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  event: {
    type: Object,
    required: true // { action, status, actor, timestamp, metadata: {} }
  }
})

const viewMode = ref('Structured')
</script>

<style scoped>
.cf-audit-card {
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-audit-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  flex-wrap: wrap;
}

.cf-audit-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-audit-sub {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
  margin-top: 2px;
}

.cf-audit-fields {
  display: flex;
  flex-direction: column;
}

.cf-audit-field-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 13px;
}

.cf-audit-field-row:last-child {
  border-bottom: none;
}

.cf-af-key {
  color: var(--text-color-kumo-subtle);
  font-weight: 500;
}

.cf-audit-json {
  padding: 12px;
}
</style>
