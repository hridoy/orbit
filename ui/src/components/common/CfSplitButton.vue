<template>
  <div ref="rootRef" class="cf-split-btn-wrap">
    <div class="cf-split-btn">
      <button type="button" class="cf-split-main" @click="$emit('click')">
        {{ label }}
      </button>
      <button type="button" class="cf-split-caret" aria-label="More options" @click="open = !open">
        ▾
      </button>
    </div>

    <div v-if="open" class="cf-split-menu">
      <button
        v-for="opt in options"
        :key="opt.id"
        type="button"
        class="cf-split-menu-item"
        @click="selectOption(opt)"
      >
        {{ opt.label }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

defineProps({
  label: { type: String, required: true },
  options: { type: Array, default: () => [] } // [{ id, label }]
})

const emit = defineEmits(['click', 'select'])
const open = ref(false)
const rootRef = ref(null)

function selectOption(opt) {
  open.value = false
  emit('select', opt)
}

function onDocClick(e) {
  if (rootRef.value && !rootRef.value.contains(e.target)) open.value = false
}

onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))
</script>

<style scoped>
.cf-split-btn-wrap {
  position: relative;
  display: inline-block;
}

.cf-split-btn {
  display: inline-flex;
  height: 36px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #0043a4;
  background-color: #0051c3;
}

.cf-split-main,
.cf-split-caret {
  border: none;
  background: transparent;
  color: #ffffff;
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.cf-split-main {
  padding: 0 12px;
}

.cf-split-caret {
  padding: 0 9px;
  border-left: 1px solid rgba(255, 255, 255, 0.25);
}

.cf-split-main:hover,
.cf-split-caret:hover {
  background-color: #003d99;
}

.cf-split-menu {
  position: absolute;
  top: calc(100% + 4px);
  right: 0;
  min-width: 170px;
  z-index: 50;
  padding: 4px;
  border-radius: 8px;
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-elevated);
  box-shadow: var(--shadow-kumo-dialog);
}

.cf-split-menu-item {
  width: 100%;
  text-align: left;
  padding: 7px 10px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--text-color-kumo-default);
  font-size: 13px;
  cursor: pointer;
}

.cf-split-menu-item:hover {
  background-color: var(--color-kumo-tint);
}
</style>
