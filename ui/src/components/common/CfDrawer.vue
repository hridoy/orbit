<script setup>
import { watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  width: {
    type: String,
    default: '460px'
  }
})

const emit = defineEmits(['update:modelValue', 'close'])

function close() {
  emit('update:modelValue', false)
  emit('close')
}

// Lock body scroll when open
watch(() => props.modelValue, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition name="cf-drawer-fade">
      <div v-if="modelValue" class="cf-drawer-overlay" @click.self="close">
        <Transition name="cf-drawer-slide">
          <aside 
            v-if="modelValue" 
            class="cf-drawer-sheet" 
            :style="{ width: width }"
            role="dialog"
            aria-modal="true"
          >
            <!-- Header -->
            <div class="cf-drawer-header">
              <div class="cf-drawer-title-group">
                <slot name="title">
                  <h3 class="cf-drawer-title">{{ title }}</h3>
                </slot>
                <p v-if="subtitle" class="cf-drawer-subtitle">{{ subtitle }}</p>
              </div>
              <button class="cf-drawer-close" type="button" aria-label="Close" @click="close">
                <CfIcon name="x" size="16" />
              </button>
            </div>

            <!-- Body -->
            <div class="cf-drawer-body">
              <slot />
            </div>

            <!-- Footer -->
            <div v-if="$slots.footer" class="cf-drawer-footer">
              <slot name="footer" :close="close" />
            </div>
          </aside>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.cf-drawer-overlay {
  position: fixed;
  inset: 0;
  z-index: 1050;
  background-color: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(2px);
  display: flex;
  justify-content: flex-end;
}

.cf-drawer-sheet {
  height: 100%;
  max-width: 100vw;
  background-color: var(--kumo-surface-default, #ffffff);
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  outline: none;
}

.cf-drawer-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--kumo-border-default, #e4e4e7);
  background-color: var(--kumo-surface-default, #ffffff);
}

.cf-drawer-title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cf-drawer-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--kumo-fg-default, #18181b);
}

.cf-drawer-subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--kumo-fg-muted, #71717a);
}

.cf-drawer-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--kumo-radius-sm, 4px);
  color: var(--kumo-fg-muted, #71717a);
  cursor: pointer;
  transition: all 0.15s ease;
}

.cf-drawer-close:hover {
  background-color: var(--kumo-control-bg-hover, #f4f4f5);
  color: var(--kumo-fg-default, #18181b);
}

.cf-drawer-body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.cf-drawer-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--kumo-border-default, #e4e4e7);
  background-color: var(--kumo-surface-secondary, #fafafa);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

/* Transitions */
.cf-drawer-fade-enter-active,
.cf-drawer-fade-leave-active {
  transition: opacity 0.25s ease;
}

.cf-drawer-fade-enter-from,
.cf-drawer-fade-leave-to {
  opacity: 0;
}

.cf-drawer-slide-enter-active {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.cf-drawer-slide-leave-active {
  transition: transform 0.2s cubic-bezier(0.7, 0, 0.84, 0);
}

.cf-drawer-slide-enter-from,
.cf-drawer-slide-leave-to {
  transform: translateX(100%);
}
</style>
