<template>
  <div class="cf-deploy-card">
    <div class="cf-deploy-head">
      <strong>{{ title }}</strong>
      <CfBadge variant="neutral" size="sm">{{ deployments.length }} versions</CfBadge>
    </div>

    <div class="cf-deploy-list">
      <div v-for="(d, idx) in deployments" :key="d.id" class="cf-deploy-row">
        <div class="cf-deploy-main">
          <div class="cf-deploy-ver-row">
            <code class="cf-deploy-id">{{ d.id }}</code>
            <CfBadge v-if="idx === 0" variant="success" size="sm" dot>Active</CfBadge>
            <span class="cf-deploy-msg">{{ d.message }}</span>
          </div>
          <div class="cf-deploy-meta">{{ d.author }} · {{ d.time }}</div>
        </div>
        <CfButton
          v-if="idx > 0"
          variant="secondary"
          size="sm"
          @click="$emit('rollback', d)"
        >
          Rollback
        </CfButton>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Deployment History' },
  deployments: {
    type: Array,
    required: true // [{ id, message, author, time }]
  }
})

defineEmits(['rollback'])
</script>

<style scoped>
.cf-deploy-card {
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-deploy-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-deploy-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
}

.cf-deploy-row:last-child {
  border-bottom: none;
}

.cf-deploy-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.cf-deploy-ver-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.cf-deploy-id {
  font-size: 11.5px;
}

.cf-deploy-msg {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-strong);
}

.cf-deploy-meta {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}
</style>
