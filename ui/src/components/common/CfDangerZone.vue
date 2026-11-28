<template>
  <div class="cf-danger-zone">
    <div class="cf-danger-info">
      <div class="cf-danger-title">{{ title }}</div>
      <div v-if="description" class="cf-danger-desc">{{ description }}</div>
    </div>
    <div class="cf-danger-action">
      <slot>
        <CfButton variant="danger" size="sm" icon="trash" @click="$emit('action')">
          {{ buttonText }}
        </CfButton>
      </slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
  buttonText: { type: String, default: 'Delete' }
})

defineEmits(['action'])
</script>

<style scoped>
.cf-danger-zone {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid rgba(220, 38, 38, 0.35);
  background-color: var(--color-kumo-danger-tint);
  flex-wrap: wrap;
}

.cf-danger-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 220px;
}

.cf-danger-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-danger);
}

.cf-danger-desc {
  font-size: 13px;
  color: var(--text-color-kumo-default);
}
</style>
