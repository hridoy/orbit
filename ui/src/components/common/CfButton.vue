<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'kumo-btn',
      `kumo-btn-${variant}`,
      `kumo-btn-${size}`,
      { 'kumo-btn-loading': loading, 'kumo-btn-icon-only': iconOnly }
    ]"
    @click="$emit('click', $event)"
  >
    <CfLoader v-if="loading" size="14" class="kumo-btn-loader" />
    <span v-else-if="$slots.icon || icon" class="kumo-btn-icon">
      <slot name="icon">
        <CfIcon v-if="icon" :name="icon" :size="iconSize" />
      </slot>
    </span>
    <span v-if="!iconOnly && $slots.default" class="kumo-btn-label">
      <slot />
    </span>
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary', // primary | secondary | outline | danger | danger-outline | ghost | link
    validator: (v) => ['primary', 'secondary', 'outline', 'danger', 'danger-outline', 'ghost', 'link'].includes(v)
  },
  size: {
    type: String,
    default: 'md', // sm (28px) | md (36px) | lg (40px)
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  },
  type: {
    type: String,
    default: 'button'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  icon: {
    type: String,
    default: null
  },
  iconOnly: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])

const iconSize = computed(() => {
  if (props.size === 'sm') return 13
  if (props.size === 'lg') return 16
  return 14
})
</script>

<style scoped>
.kumo-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  font-family: var(--font-sans);
  font-weight: 500;
  letter-spacing: -0.01em;
  border-radius: 6px;
  cursor: pointer;
  white-space: nowrap;
  user-select: none;
  transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  text-decoration: none;
  outline: none;
  box-sizing: border-box;
}

.kumo-btn:focus-visible {
  box-shadow: 0 0 0 2px #ffffff, 0 0 0 4px #0051c3;
}

.kumo-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
  pointer-events: none;
}

/* SIZES */
.kumo-btn-sm {
  height: 28px;
  padding: 0 10px;
  font-size: 12px;
}

.kumo-btn-md {
  height: 36px;
  padding: 0 14px;
  font-size: 13.5px;
}

.kumo-btn-lg {
  height: 40px;
  padding: 0 18px;
  font-size: 14px;
}

.kumo-btn-icon-only.kumo-btn-sm {
  width: 28px;
  padding: 0;
}

.kumo-btn-icon-only.kumo-btn-md {
  width: 36px;
  padding: 0;
}

.kumo-btn-icon-only.kumo-btn-lg {
  width: 40px;
  padding: 0;
}

/* VARIANTS */
/* 1. Primary Blue (Signature Cloudflare Gradient Gloss) */
.kumo-btn-primary {
  color: #ffffff !important;
  background-color: #0051c3;
  background-image: linear-gradient(180deg, rgba(255, 255, 255, 0.12), transparent);
  border: 1px solid #0043a4;
  box-shadow: 0 1px 2px rgba(0, 81, 195, 0.2), inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.kumo-btn-primary:hover:not(:disabled) {
  background-color: #003d99;
  border-color: #003681;
}

/* 2. Secondary White / Neutral */
.kumo-btn-secondary,
.kumo-btn-outline {
  color: var(--kumo-control-fg-default, #18181b);
  background-color: var(--kumo-control-bg-default, #ffffff);
  border: 1px solid var(--kumo-control-border-default, #d1d5db);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.kumo-btn-secondary:hover:not(:disabled),
.kumo-btn-outline:hover:not(:disabled) {
  background-color: var(--kumo-control-bg-hover, #f4f4f5);
  border-color: var(--kumo-control-border-hover, #9ca3af);
}

/* 3. Danger Solid Red */
.kumo-btn-danger {
  color: #ffffff !important;
  background-color: #dc2626;
  border: 1px solid #b91c1c;
  box-shadow: 0 1px 2px rgba(220, 38, 38, 0.2);
}

.kumo-btn-danger:hover:not(:disabled) {
  background-color: #b91c1c;
  border-color: #991b1b;
}

/* 3b. Danger Outline (Cloudflare Purge Everything / Destructive Secondary) */
.kumo-btn-danger-outline {
  color: #dc2626;
  background-color: var(--kumo-control-bg-default, #ffffff);
  border: 1px solid #fca5a5;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.kumo-btn-danger-outline:hover:not(:disabled) {
  background-color: rgba(220, 38, 38, 0.06);
  border-color: #dc2626;
}

/* 4. Ghost / Transparent */
.kumo-btn-ghost {
  color: var(--text-color-kumo-subtle, #52525b);
  background: transparent;
  border: 1px solid transparent;
}

.kumo-btn-ghost:hover:not(:disabled) {
  background-color: var(--color-kumo-fill, rgba(0, 0, 0, 0.05));
  color: var(--text-color-kumo-strong, #18181b);
}

[data-theme="dark"] .kumo-btn-ghost {
  color: #a1a1aa;
}

[data-theme="dark"] .kumo-btn-ghost:hover:not(:disabled) {
  background-color: rgba(255, 255, 255, 0.08);
  color: #f8fafc;
}

/* 5. Link */
.kumo-btn-link {
  color: #0051c3;
  background: transparent;
  border: none;
  padding: 0 !important;
  height: auto !important;
}

.kumo-btn-link:hover:not(:disabled) {
  text-decoration: underline;
}

.kumo-btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
</style>
