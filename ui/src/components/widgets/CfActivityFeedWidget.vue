<template>
  <div class="cf-activity-card">
    <div v-if="title" class="cf-activity-head">
      <strong>{{ title }}</strong>
      <slot name="action" />
    </div>
    <div class="cf-activity-list">
      <div v-for="(item, idx) in items" :key="idx" class="cf-activity-row">
        <CfAvatar :name="item.actor" size="sm" />
        <div class="cf-activity-content">
          <div class="cf-activity-line">
            <strong>{{ item.actor }}</strong>
            <span>{{ item.action }}</span>
            <code v-if="item.resource">{{ item.resource }}</code>
          </div>
          <div class="cf-activity-time">{{ item.time }}</div>
        </div>
        <CfBadge v-if="item.badge" :variant="item.badgeVariant || 'neutral'" size="sm">{{ item.badge }}</CfBadge>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Recent Activity' },
  items: {
    type: Array,
    required: true // [{ actor, action, resource?, time, badge?, badgeVariant? }]
  }
})
</script>

<style scoped>
.cf-activity-card {
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-activity-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-activity-list {
  display: flex;
  flex-direction: column;
}

.cf-activity-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
}

.cf-activity-row:last-child {
  border-bottom: none;
}

.cf-activity-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-activity-line {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 5px;
  font-size: 13px;
  color: var(--text-color-kumo-default);
}

.cf-activity-time {
  font-size: 11.5px;
  color: var(--text-color-kumo-subtle);
}
</style>
