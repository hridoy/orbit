<template>
  <div ref="containerRef" class="cf-split-pane" :style="{ height }">
    <div class="cf-split-left" :style="{ width: `${leftPercent}%` }">
      <slot name="left" />
    </div>
    <div
      class="cf-split-handle"
      title="Drag to resize"
      @mousedown.prevent="startDrag"
    >
      <span class="cf-split-grip"></span>
    </div>
    <div class="cf-split-right">
      <slot name="right" />
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'

const props = defineProps({
  initialLeft: { type: Number, default: 34 },
  minLeft: { type: Number, default: 18 },
  maxLeft: { type: Number, default: 75 },
  height: { type: String, default: '240px' }
})

const leftPercent = ref(props.initialLeft)
const containerRef = ref(null)
let dragging = false

function startDrag() {
  dragging = true
  window.addEventListener('mousemove', onDrag)
  window.addEventListener('mouseup', stopDrag)
}

function onDrag(e) {
  if (!dragging || !containerRef.value) return
  const rect = containerRef.value.getBoundingClientRect()
  const pct = ((e.clientX - rect.left) / rect.width) * 100
  leftPercent.value = Math.min(props.maxLeft, Math.max(props.minLeft, pct))
}

function stopDrag() {
  dragging = false
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('mouseup', stopDrag)
}

onUnmounted(stopDrag)
</script>

<style scoped>
.cf-split-pane {
  display: flex;
  width: 100%;
  border: 1px solid var(--color-kumo-line);
  border-radius: var(--radius-lg, 8px);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-split-left {
  overflow: auto;
  padding: 12px;
  flex-shrink: 0;
  background-color: var(--color-kumo-canvas);
}

.cf-split-handle {
  width: 6px;
  cursor: col-resize;
  background-color: var(--color-kumo-hairline);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background-color 0.15s ease;
}

.cf-split-handle:hover {
  background-color: var(--color-kumo-brand);
}

.cf-split-grip {
  width: 2px;
  height: 24px;
  border-radius: 2px;
  background-color: var(--text-color-kumo-placeholder);
}

.cf-split-right {
  flex: 1;
  min-width: 0;
  overflow: auto;
  padding: 12px;
}
</style>
