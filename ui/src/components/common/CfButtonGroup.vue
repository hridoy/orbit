<template>
  <div class="cf-btn-group" role="group">
    <button
      v-for="item in items"
      :key="item.value"
      type="button"
      class="cf-btn-group-item"
      :class="{ 'is-active': modelValue === item.value }"
      @click="$emit('update:modelValue', item.value)"
    >
      {{ item.label }}
    </button>
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: [String, Number], required: true },
  items: { type: Array, required: true } // [{ label, value }]
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-btn-group {
  display: inline-flex;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  overflow: hidden;
}

.cf-btn-group-item {
  height: 34px;
  padding: 0 12px;
  border: none;
  border-right: 1px solid var(--kumo-control-border-default);
  background: transparent;
  color: var(--text-color-kumo-default);
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.cf-btn-group-item:last-child {
  border-right: none;
}

.cf-btn-group-item:hover {
  background-color: var(--kumo-control-bg-hover);
}

.cf-btn-group-item.is-active {
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-strong);
  font-weight: 600;
}
</style>
