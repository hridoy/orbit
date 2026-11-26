<template>
  <div class="kumo-empty-card">
    <!-- Cloudflare Dashed Cloud Illustration -->
    <div class="kumo-empty-svg-wrap">
      <svg class="kumo-empty-cloud-svg" viewBox="0 0 450 210" fill="none">
        <circle cx="114.45" cy="113.58" r="48.4" stroke="currentColor" stroke-width="3.5" stroke-dasharray="7 6" />
        <line x1="208.6" y1="207.9" x2="4.2" y2="207.9" stroke="currentColor" stroke-width="3.5" stroke-dasharray="7 6" />
        <path d="M449.299 208.419C452.346 199.411 453.999 189.761 453.999 179.725C453.999 130.292 413.926 90.2192 364.494 90.2192C315.061 90.2192 274.988 130.292 274.988 179.725C274.988 189.761 276.641 199.411 279.688 208.419" stroke="currentColor" stroke-width="3.5" stroke-dasharray="7 6" />
      </svg>
    </div>

    <h3 class="kumo-empty-title">{{ title }}</h3>
    <p v-if="description" class="kumo-empty-desc">{{ description }}</p>

    <div v-if="$slots.default || actionText" class="kumo-empty-action">
      <slot>
        <CfButton variant="primary" :icon="actionIcon" @click="$emit('action')">
          {{ actionText }}
        </CfButton>
      </slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    default: null
  },
  actionText: {
    type: String,
    default: null
  },
  actionIcon: {
    type: String,
    default: null
  }
})

defineEmits(['action'])
</script>

<style scoped>
.kumo-empty-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 3rem 1.5rem;
  background-color: var(--color-kumo-base, #ffffff);
  border: 1px solid var(--color-kumo-line, #e2e8f0);
  border-radius: 8px;
  max-width: 480px;
  margin: 0 auto;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

.kumo-empty-svg-wrap {
  width: 160px;
  margin-bottom: 1.25rem;
  color: #94a3b8;
}

.kumo-empty-cloud-svg {
  width: 100%;
  height: auto;
  display: block;
}

.kumo-empty-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color-kumo-default, #18181b);
  margin: 0 0 0.4rem 0;
}

[data-theme="dark"] .kumo-empty-title {
  color: #f8fafc;
}

.kumo-empty-desc {
  font-size: 13.5px;
  line-height: 1.45;
  color: var(--text-color-kumo-subtle, #71717a);
  margin: 0 0 1.25rem 0;
  max-width: 380px;
}

.kumo-empty-action {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
