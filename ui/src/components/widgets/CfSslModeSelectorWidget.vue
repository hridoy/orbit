<template>
  <div class="cf-ssl-grid">
    <button
      v-for="mode in modes"
      :key="mode.id"
      type="button"
      class="cf-ssl-card"
      :class="{ 'is-selected': modelValue === mode.id }"
      @click="$emit('update:modelValue', mode.id)"
    >
      <div class="cf-ssl-top">
        <strong class="cf-ssl-title">{{ mode.label }}</strong>
        <span class="cf-ssl-radio" :class="{ checked: modelValue === mode.id }">
          <span v-if="modelValue === mode.id" class="cf-ssl-dot"></span>
        </span>
      </div>
      <p class="cf-ssl-desc">{{ mode.description }}</p>
      <CfBadge v-if="mode.recommended" variant="success" size="sm">Recommended</CfBadge>
    </button>
  </div>
</template>

<script setup>
defineProps({
  modelValue: { type: String, default: 'full_strict' },
  modes: {
    type: Array,
    default: () => [
      { id: 'off', label: 'Off (not secure)', description: 'No encryption applied between browser, Cloudflare, or origin.' },
      { id: 'flexible', label: 'Flexible', description: 'Encrypts traffic between browser and Cloudflare only.' },
      { id: 'full', label: 'Full', description: 'Encrypts end-to-end using a self-signed origin certificate.' },
      { id: 'full_strict', label: 'Full (strict)', description: 'Encrypts end-to-end with a valid CA-signed origin certificate.', recommended: true }
    ]
  }
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.cf-ssl-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.cf-ssl-card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding: 14px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  text-align: left;
  cursor: pointer;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.cf-ssl-card:hover {
  border-color: var(--color-kumo-brand);
}

.cf-ssl-card.is-selected {
  border-color: var(--color-kumo-brand);
  box-shadow: 0 0 0 2px rgba(0, 81, 195, 0.15);
  background-color: var(--color-kumo-info-tint);
}

.cf-ssl-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.cf-ssl-title {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-ssl-radio {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1.5px solid var(--kumo-control-border-default);
  display: flex;
  align-items: center;
  justify-content: center;
}

.cf-ssl-radio.checked {
  border-color: var(--color-kumo-brand);
  background-color: var(--color-kumo-brand);
}

.cf-ssl-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #ffffff;
}

.cf-ssl-desc {
  margin: 0;
  font-size: 12px;
  line-height: 1.45;
  color: var(--text-color-kumo-subtle);
}
</style>
