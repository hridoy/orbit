<template>
  <div :class="['kumo-card', { 'kumo-card-interactive': interactive }]">
    <!-- Optional Header -->
    <div v-if="$slots.header || title" class="kumo-card-header">
      <div class="kumo-card-title-group">
        <slot name="title">
          <h3 class="kumo-card-title">{{ title }}</h3>
        </slot>
        <span v-if="subtitle" class="kumo-card-subtitle">{{ subtitle }}</span>
      </div>
      <div v-if="$slots.actions" class="kumo-card-actions">
        <slot name="actions" />
      </div>
    </div>

    <!-- Main Card Body -->
    <div :class="['kumo-card-body', bodyClass]">
      <slot />
    </div>

    <!-- Optional Footer -->
    <div v-if="$slots.footer" class="kumo-card-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: null
  },
  subtitle: {
    type: String,
    default: null
  },
  interactive: {
    type: Boolean,
    default: false
  },
  bodyClass: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.kumo-card {
  display: flex;
  flex-direction: column;
  width: 100%;
  background-color: var(--color-kumo-base, #ffffff);
  border: 1px solid var(--color-kumo-line, #e2e8f0);
  border-radius: 8px;
  box-shadow: var(--shadow-kumo-card, 0 1px 2px rgba(0, 0, 0, 0.04));
  box-sizing: border-box;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.kumo-card-interactive:hover {
  border-color: #94a3b8;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
}

.kumo-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1.25rem;
  background-color: var(--color-kumo-base, #ffffff);
  border-bottom: 1px solid var(--color-kumo-line, #f4f4f5);
  border-top-left-radius: 7px;
  border-top-right-radius: 7px;
  gap: 1rem;
}

.kumo-card-title-group {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.kumo-card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-color-kumo-default, #18181b);
  margin: 0;
  letter-spacing: -0.01em;
}

[data-theme="dark"] .kumo-card-title {
  color: #f8fafc;
}

.kumo-card-subtitle {
  font-size: 13px;
  color: var(--text-color-kumo-subtle, #71717a);
}

.kumo-card-body {
  padding: 1.25rem;
}

.kumo-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1.25rem;
  background-color: #fafafa;
  border-top: 1px solid var(--color-kumo-line, #f4f4f5);
  font-size: 13px;
}

[data-theme="dark"] .kumo-card-footer {
  background-color: #070a12;
  border-top-color: #1e293b;
}
</style>
