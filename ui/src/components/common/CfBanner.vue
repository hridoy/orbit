<template>
  <div v-if="!dismissed" class="cf-banner" :class="`cf-banner-${variant}`" role="status">
    <div class="cf-banner-left">
      <CfIcon v-if="icon" :name="icon" size="16" class="cf-banner-icon" />
      <div class="cf-banner-text">
        <strong v-if="title" class="cf-banner-title">{{ title }}</strong>
        <span class="cf-banner-body"><slot /></span>
      </div>
    </div>
    <div class="cf-banner-actions">
      <slot name="action" />
      <button
        v-if="dismissible"
        type="button"
        class="cf-banner-close"
        aria-label="Dismiss banner"
        @click="dismissed = true; $emit('dismiss')"
      >
        ✕
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  variant: {
    type: String,
    default: 'info', // info | warning | brand | success | danger
    validator: (v) => ['info', 'warning', 'brand', 'success', 'danger'].includes(v)
  },
  title: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: 'cis-circle-info'
  },
  dismissible: {
    type: Boolean,
    default: false
  }
})

defineEmits(['dismiss'])
const dismissed = ref(false)
</script>

<style scoped>
.cf-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  font-family: var(--font-sans);
  font-size: 13px;
  line-height: 1.45;
  flex-wrap: wrap;
}

.cf-banner-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 200px;
}

.cf-banner-icon {
  flex-shrink: 0;
}

.cf-banner-text {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
}

.cf-banner-title {
  font-weight: 600;
}

.cf-banner-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cf-banner-close {
  border: none;
  background: transparent;
  color: inherit;
  opacity: 0.65;
  cursor: pointer;
  padding: 2px 4px;
  font-size: 13px;
}

.cf-banner-close:hover {
  opacity: 1;
}

.cf-banner-info {
  background-color: var(--color-kumo-banner-info);
  color: var(--text-color-kumo-info);
}

.cf-banner-warning {
  background-color: var(--color-kumo-banner-warning);
  color: var(--text-color-kumo-warning);
}

.cf-banner-brand {
  background-color: var(--cf-brand-orange-subtle);
  border-color: rgba(246, 130, 31, 0.3);
  color: var(--text-color-kumo-strong);
}

.cf-banner-success {
  background-color: var(--color-kumo-success-tint);
  color: var(--text-color-kumo-success);
}

.cf-banner-danger {
  background-color: var(--color-kumo-danger-tint);
  color: var(--text-color-kumo-danger);
}
</style>
