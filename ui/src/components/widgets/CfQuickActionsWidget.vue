<template>
  <div class="cf-quick-actions-grid">
    <button
      v-for="item in actions"
      :key="item.title"
      type="button"
      class="cf-quick-tile"
      @click="$emit('action', item)"
    >
      <div class="cf-quick-icon">
        <CfIcon :name="item.icon" size="18" />
      </div>
      <div class="cf-quick-text">
        <div class="cf-quick-title">{{ item.title }}</div>
        <div v-if="item.subtitle" class="cf-quick-sub">{{ item.subtitle }}</div>
      </div>
      <span class="cf-quick-chevron">→</span>
    </button>
  </div>
</template>

<script setup>
defineProps({
  actions: {
    type: Array,
    required: true // [{ title, subtitle?, icon, route? }]
  }
})

defineEmits(['action'])
</script>

<style scoped>
.cf-quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.cf-quick-tile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  text-align: left;
  cursor: pointer;
  transition: border-color 0.15s ease, background-color 0.15s ease;
}

.cf-quick-tile:hover {
  border-color: var(--color-kumo-brand);
  background-color: var(--color-kumo-tint);
}

.cf-quick-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-strong);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.cf-quick-text {
  flex: 1;
  min-width: 0;
}

.cf-quick-title {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-quick-sub {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cf-quick-chevron {
  color: var(--text-color-kumo-subtle);
  font-size: 14px;
}
</style>
