<template>
  <nav class="cf-breadcrumbs" aria-label="Breadcrumb">
    <template v-for="(item, idx) in items" :key="idx">
      <router-link
        v-if="item.to && idx < items.length - 1"
        :to="item.to"
        class="cf-breadcrumb-link"
      >
        {{ item.label }}
      </router-link>
      <span
        v-else
        class="cf-breadcrumb-item"
        :class="{ 'is-current': idx === items.length - 1 }"
      >
        {{ item.label }}
      </span>
      <span v-if="idx < items.length - 1" class="cf-breadcrumb-sep" aria-hidden="true">/</span>
    </template>
  </nav>
</template>

<script setup>
defineProps({
  items: {
    type: Array,
    required: true // [{ label: string, to?: string }]
  }
})
</script>

<style scoped>
.cf-breadcrumbs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  font-family: var(--font-sans);
  font-size: 13px;
  line-height: 1.4;
}

.cf-breadcrumb-link {
  color: var(--text-color-kumo-subtle);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.15s ease;
}

.cf-breadcrumb-link:hover {
  color: var(--text-color-kumo-strong);
  text-decoration: underline;
}

.cf-breadcrumb-item {
  color: var(--text-color-kumo-subtle);
  font-weight: 500;
}

.cf-breadcrumb-item.is-current {
  color: var(--text-color-kumo-strong);
  font-weight: 600;
}

.cf-breadcrumb-sep {
  color: var(--text-color-kumo-placeholder);
  user-select: none;
}
</style>
