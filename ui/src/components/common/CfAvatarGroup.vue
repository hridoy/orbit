<template>
  <div class="cf-avatar-group">
    <div
      v-for="(name, idx) in visibleNames"
      :key="name + idx"
      class="cf-avatar-ring"
      :title="name"
    >
      <CfAvatar :name="name" size="sm" />
    </div>
    <span v-if="overflow > 0" class="cf-avatar-overflow">+{{ overflow }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  names: { type: Array, required: true },
  max: { type: Number, default: 4 }
})

const visibleNames = computed(() => props.names.slice(0, props.max))
const overflow = computed(() => Math.max(0, props.names.length - props.max))
</script>

<style scoped>
.cf-avatar-group {
  display: inline-flex;
  align-items: center;
}

.cf-avatar-ring {
  margin-left: -8px;
  border-radius: 50%;
  border: 2px solid var(--color-kumo-base);
}

.cf-avatar-ring:first-child {
  margin-left: 0;
}

.cf-avatar-overflow {
  margin-left: -6px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid var(--color-kumo-base);
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-strong);
  font-size: 11px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>
