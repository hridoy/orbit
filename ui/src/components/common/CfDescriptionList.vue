<template>
  <div class="cf-desc-list">
    <div v-for="(item, idx) in items" :key="item.label + idx" class="cf-desc-row">
      <div class="cf-desc-left">
        <div class="cf-desc-label">
          <span>{{ item.label }}</span>
          <CfBadge v-if="item.badge" :variant="item.badgeVariant || 'neutral'" size="sm" dot>
            {{ item.badge }}
          </CfBadge>
        </div>
        <div v-if="item.description" class="cf-desc-hint">{{ item.description }}</div>
      </div>
      <div class="cf-desc-right">
        <slot :name="`row-${idx}`" :item="item">
          <span class="cf-desc-val" :class="{ 'font-mono': item.mono }">{{ item.value }}</span>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  items: {
    type: Array,
    required: true // [{ label, value?, description?, badge?, badgeVariant?, mono? }]
  }
})
</script>

<style scoped>
.cf-desc-list {
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-lg, 8px);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-desc-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  flex-wrap: wrap;
}

.cf-desc-row:last-child {
  border-bottom: none;
}

.cf-desc-left {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 180px;
}

.cf-desc-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text-color-kumo-strong);
}

.cf-desc-hint {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-desc-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cf-desc-val {
  font-size: 13px;
  color: var(--text-color-kumo-default);
}
</style>
