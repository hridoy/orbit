<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
    // item: { id, label, icon, shortcut, danger, divider, disabled, action }
  },
  align: {
    type: String,
    default: 'left', // 'left' | 'right'
    validator: (v) => ['left', 'right'].includes(v)
  }
})

const emit = defineEmits(['select'])

const isOpen = ref(false)
const dropdownRef = ref(null)

function toggle() {
  isOpen.value = !isOpen.value
}

function close() {
  isOpen.value = false
}

function handleSelect(item) {
  if (item.disabled) return
  if (item.action) {
    item.action(item)
  }
  emit('select', item)
  close()
}

function handleClickOutside(e) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target)) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="dropdownRef" class="cf-dropdown-wrap">
    <div class="cf-dropdown-trigger" @click.stop="toggle">
      <slot name="trigger" :isOpen="isOpen">
        <button class="cf-dropdown-btn" type="button">
          <span>Options</span>
          <CfIcon name="chevron-down" size="14" class="cf-dropdown-chevron" :class="{ 'is-open': isOpen }" />
        </button>
      </slot>
    </div>

    <Transition name="cf-dropdown-pop">
      <div 
        v-if="isOpen" 
        class="cf-dropdown-menu" 
        :class="[`align-${align}`]"
        role="menu"
      >
        <slot :close="close">
          <template v-for="(item, idx) in items" :key="item.id || idx">
            <div v-if="item.divider" class="cf-dropdown-divider" role="separator"></div>
            <button
              v-else
              type="button"
              class="cf-dropdown-item"
              :class="{ 'is-danger': item.danger, 'is-disabled': item.disabled }"
              :disabled="item.disabled"
              role="menuitem"
              @click="handleSelect(item)"
            >
              <span class="cf-dropdown-item-left">
                <CfIcon v-if="item.icon" :name="item.icon" size="14" class="cf-dropdown-item-icon" />
                <span class="cf-dropdown-item-label">{{ item.label }}</span>
              </span>
              <span v-if="item.shortcut" class="cf-dropdown-shortcut">{{ item.shortcut }}</span>
            </button>
          </template>
        </slot>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.cf-dropdown-wrap {
  position: relative;
  display: inline-block;
}

.cf-dropdown-trigger {
  cursor: pointer;
  display: inline-flex;
}

.cf-dropdown-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 12px;
  font-size: 13px;
  font-weight: 500;
  color: var(--kumo-control-fg-default, #1d1d1f);
  background-color: var(--kumo-control-bg-default, #f4f4f5);
  border: 1px solid var(--kumo-control-border-default, #d4d4d8);
  border-radius: var(--kumo-radius-md, 6px);
  cursor: pointer;
  transition: all 0.15s ease;
}

.cf-dropdown-btn:hover {
  background-color: var(--kumo-control-bg-hover, #e4e4e7);
  border-color: var(--kumo-control-border-hover, #a1a1aa);
}

.cf-dropdown-chevron {
  transition: transform 0.2s ease;
  color: var(--kumo-fg-muted, #71717a);
}

.cf-dropdown-chevron.is-open {
  transform: rotate(180deg);
}

.cf-dropdown-menu {
  position: absolute;
  top: calc(100% + 4px);
  z-index: 999;
  min-width: 180px;
  padding: 4px;
  background-color: var(--kumo-surface-default, #ffffff);
  border: 1px solid var(--kumo-border-default, #e4e4e7);
  border-radius: var(--kumo-radius-md, 6px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
}

.align-left {
  left: 0;
}

.align-right {
  right: 0;
}

.cf-dropdown-divider {
  height: 1px;
  background-color: var(--kumo-border-default, #e4e4e7);
  margin: 4px -4px;
}

.cf-dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  height: 30px;
  padding: 0 8px;
  font-size: 13px;
  font-weight: 400;
  color: var(--kumo-control-fg-default, #1d1d1f);
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.12s ease;
}

.cf-dropdown-item:hover:not(:disabled) {
  background-color: var(--kumo-control-bg-hover, #f4f4f5);
}

.cf-dropdown-item.is-danger {
  color: #dc2626;
}

.cf-dropdown-item.is-danger:hover:not(:disabled) {
  background-color: #fee2e2;
}

.cf-dropdown-item:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.cf-dropdown-item-left {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.cf-dropdown-item-icon {
  color: var(--kumo-fg-muted, #71717a);
}

.cf-dropdown-item.is-danger .cf-dropdown-item-icon {
  color: #dc2626;
}

.cf-dropdown-shortcut {
  font-size: 11px;
  font-family: var(--font-mono, monospace);
  color: var(--kumo-fg-subtle, #a1a1aa);
  letter-spacing: 0.5px;
}

.cf-dropdown-pop-enter-active,
.cf-dropdown-pop-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.cf-dropdown-pop-enter-from,
.cf-dropdown-pop-leave-to {
  opacity: 0;
  transform: translateY(-4px) scale(0.98);
}
</style>
