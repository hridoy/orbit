<template>
  <div class="kumo-empty-card">
    <!-- Cloudflare Dashed Cloud Illustration -->
    <div class="kumo-empty-svg-wrap">
      <svg class="kumo-empty-cloud-svg" viewBox="0 0 200 110" fill="none" aria-hidden="true">
        <path
          d="M52 90H150C168.778 90 184 74.7777 184 56C184 38.3085 170.491 23.773 153.214 22.1502C146.766 9.57544 133.681 1 118.5 1C99.3895 1 83.4896 14.6064 79.8763 32.6689C76.2225 30.9557 72.1424 30 67.8333 30C51.9091 30 39 42.9091 39 58.8333C39 59.9114 39.0591 60.9756 39.1743 62.0226C26.1859 64.0842 16.25 75.331 16.25 88.9167C16.25 89.2817 16.2572 89.643 16.2715 90H52Z"
          stroke="currentColor"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-dasharray="6 5"
        />
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
