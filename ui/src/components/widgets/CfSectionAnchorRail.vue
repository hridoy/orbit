<template>
  <nav class="cf-anchor-rail" aria-label="Section navigation">
    <div v-if="title" class="cf-rail-title">{{ title }}</div>
    <button
      v-for="item in items"
      :key="item.id"
      type="button"
      class="cf-rail-link"
      :class="{ 'is-active': modelValue === item.id }"
      @click="$emit('update:modelValue', item.id)"
    >
      {{ item.label }}
    </button>
  </nav>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'On this page' },
  modelValue: { type: String, required: true },
  items: { type: Array, required: true } // [{ id, label }]
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-anchor-rail {
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-left: 1px solid var(--color-kumo-line);
  padding-left: 12px;
}

.cf-rail-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-kumo-subtle);
  margin-bottom: 6px;
}

.cf-rail-link {
  position: relative;
  text-align: left;
  padding: 5px 8px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--text-color-kumo-subtle);
  font-size: 13px;
  cursor: pointer;
  transition: color 0.15s ease, background-color 0.15s ease;
}

.cf-rail-link:hover {
  color: var(--text-color-kumo-strong);
  background-color: var(--color-kumo-tint);
}

.cf-rail-link.is-active {
  color: var(--text-color-kumo-strong);
  font-weight: 600;
}

.cf-rail-link.is-active::before {
  content: '';
  position: absolute;
  left: -13px;
  top: 4px;
  bottom: 4px;
  width: 2px;
  background-color: var(--color-kumo-brand);
  border-radius: 2px;
}
</style>
