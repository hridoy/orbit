<template>
  <div ref="rootRef" class="cf-popover-wrap">
    <div class="cf-popover-trigger" @click="open = !open">
      <slot name="trigger">
        <CfButton variant="secondary" size="sm">Open Popover</CfButton>
      </slot>
    </div>

    <div v-if="open" class="cf-popover-card" :class="`align-${align}`" :style="{ width }">
      <div v-if="title" class="cf-popover-header">
        <span>{{ title }}</span>
        <button type="button" class="cf-popover-close" @click="open = false">✕</button>
      </div>
      <div class="cf-popover-body">
        <slot :close="() => (open = false)" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

defineProps({
  title: { type: String, default: '' },
  width: { type: String, default: '280px' },
  align: { type: String, default: 'left' } // left | right
})

const open = ref(false)
const rootRef = ref(null)

function onDocClick(e) {
  if (rootRef.value && !rootRef.value.contains(e.target)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))
</script>

<style scoped>
.cf-popover-wrap {
  position: relative;
  display: inline-block;
}

.cf-popover-trigger {
  display: inline-block;
}

.cf-popover-card {
  position: absolute;
  top: calc(100% + 6px);
  z-index: 60;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-elevated);
  box-shadow: var(--shadow-kumo-dialog);
  overflow: hidden;
}

.cf-popover-card.align-left { left: 0; }
.cf-popover-card.align-right { right: 0; }

.cf-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-popover-close {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-subtle);
  cursor: pointer;
  font-size: 12px;
}

.cf-popover-body {
  padding: 12px;
  font-size: 13px;
  color: var(--text-color-kumo-default);
}
</style>
