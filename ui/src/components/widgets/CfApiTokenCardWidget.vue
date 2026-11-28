<template>
  <div class="cf-token-card">
    <div class="cf-token-head">
      <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
        <CfIcon name="cis-key-1" size="16" style="color: var(--text-color-kumo-subtle);" />
        <strong class="cf-token-name">{{ name }}</strong>
        <CfBadge :variant="status === 'Active' ? 'success' : 'danger'" size="sm" dot>{{ status }}</CfBadge>
      </div>
      <div style="display: flex; gap: 6px;">
        <CfButton variant="secondary" size="sm" @click="$emit('roll', name)">Roll</CfButton>
        <CfButton variant="danger-outline" size="sm" @click="$emit('revoke', name)">Revoke</CfButton>
      </div>
    </div>

    <div class="cf-token-scopes">
      <span v-for="scope in scopes" :key="scope" class="cf-scope-pill">{{ scope }}</span>
    </div>

    <div class="cf-token-foot">
      <span>Last used: {{ lastUsed }}</span>
      <span>Expires: {{ expires }}</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  name: { type: String, required: true },
  status: { type: String, default: 'Active' },
  scopes: { type: Array, default: () => [] },
  lastUsed: { type: String, default: 'Today' },
  expires: { type: String, default: 'Never' }
})

defineEmits(['roll', 'revoke'])
</script>

<style scoped>
.cf-token-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px 16px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-token-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.cf-token-name {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-token-scopes {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.cf-scope-pill {
  font-family: var(--font-mono);
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 9999px;
  background-color: var(--color-kumo-recessed);
  border: 1px solid var(--color-kumo-line);
  color: var(--text-color-kumo-default);
}

.cf-token-foot {
  display: flex;
  justify-content: space-between;
  font-size: 11.5px;
  color: var(--text-color-kumo-subtle);
  padding-top: 8px;
  border-top: 1px solid var(--color-kumo-hairline);
}
</style>
