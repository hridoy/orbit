<template>
  <div class="cf-toggle-option-card">
    <div class="cf-toc-text">
      <div class="cf-toc-title">
        <span>{{ title }}</span>
        <CfBadge v-if="badge" variant="info" size="sm">{{ badge }}</CfBadge>
      </div>
      <div v-if="description" class="cf-toc-desc">{{ description }}</div>
    </div>
    <CfSwitch :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" />
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, required: true },
  description: { type: String, default: '' },
  badge: { type: String, default: '' }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-toggle-option-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-toc-text {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.cf-toc-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-toc-desc {
  font-size: 12.5px;
  color: var(--text-color-kumo-subtle);
}
</style>
