<template>
  <div class="kumo-segmented-tabs" role="tablist">
    <button
      v-for="tab in tabs"
      :key="tab.value || tab"
      type="button"
      role="tab"
      :aria-selected="(modelValue === (tab.value || tab)) ? 'true' : 'false'"
      :class="[
        'kumo-tab-btn',
        { 'is-active': modelValue === (tab.value || tab) }
      ]"
      @click="$emit('update:modelValue', tab.value || tab)"
    >
      <CfIcon v-if="tab.icon" :name="tab.icon" size="14" class="kumo-tab-icon" />
      <span>{{ tab.label || tab }}</span>
      <span v-if="tab.count !== undefined" class="kumo-tab-count">{{ tab.count }}</span>
    </button>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: [String, Number],
    required: true
  },
  tabs: {
    type: Array,
    required: true // array of strings or { label, value, icon, count }
  }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.kumo-segmented-tabs {
  display: inline-flex;
  align-items: center;
  background-color: #f4f4f5;
  height: 36px;
  padding: 2px;
  border-radius: 8px;
  gap: 2px;
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.05);
  box-sizing: border-box;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

[data-theme="dark"] .kumo-segmented-tabs {
  background-color: rgba(255, 255, 255, 0.08);
}

.kumo-tab-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.35rem;
  height: 32px;
  padding: 0 12px;
  font-size: 13.5px;
  font-weight: 400;
  color: #71717a;
  border-radius: 6px;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.kumo-tab-btn:hover {
  color: #18181b;
}

[data-theme="dark"] .kumo-tab-btn {
  color: #a1a1aa;
}

[data-theme="dark"] .kumo-tab-btn:hover {
  color: #f8fafc;
}

.kumo-tab-btn.is-active {
  background-color: #ffffff;
  color: #18181b;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.08);
  font-weight: 500;
}

[data-theme="dark"] .kumo-tab-btn.is-active {
  background-color: #18181b;
  color: #ffffff;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.4);
}

.kumo-tab-count {
  font-size: 11px;
  padding: 1px 5px;
  border-radius: 9999px;
  background-color: rgba(0, 0, 0, 0.06);
  color: #52525b;
}

[data-theme="dark"] .kumo-tab-count {
  background-color: rgba(255, 255, 255, 0.12);
  color: #d4d4d8;
}
</style>
