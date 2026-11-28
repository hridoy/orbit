<template>
  <span class="cf-proxy-status" :class="proxied ? 'is-proxied' : 'is-dns-only'">
    <svg width="20" height="12" viewBox="0 0 24 14" fill="currentColor" aria-hidden="true">
      <path d="M19.35 6.04A5.49 5.49 0 0 0 14 2c-2.22 0-4.12 1.32-4.96 3.22A3.99 3.99 0 0 0 5 9.13 3.5 3.5 0 0 0 5.5 14H19a3.5 3.5 0 0 0 .35-7.96Z" />
    </svg>
    <span>{{ proxied ? 'Proxied' : 'DNS only' }}</span>
  </span>
</template>

<script setup>
defineProps({
  proxied: { type: Boolean, default: true }
})
</script>

<style scoped>
.cf-proxy-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  font-weight: 500;
  white-space: nowrap;
}

.cf-proxy-status.is-proxied {
  color: var(--cf-brand-orange);
}

.cf-proxy-status.is-dns-only {
  color: var(--text-color-kumo-subtle);
}
</style>
