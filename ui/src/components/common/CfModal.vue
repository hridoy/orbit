<template>
  <Teleport to="body">
    <Transition name="kumo-modal-fade">
      <div v-if="modelValue" class="kumo-modal-backdrop" @click.self="onBackdropClick">
        <div :class="['kumo-modal-dialog', `kumo-modal-${size}`]" role="dialog" aria-modal="true">
          <!-- Modal Header -->
          <div class="kumo-modal-header">
            <div class="kumo-modal-title-row">
              <slot name="title">
                <h2 class="kumo-modal-title">{{ title }}</h2>
              </slot>
            </div>
            <button
              v-if="showClose"
              type="button"
              class="kumo-modal-close-btn"
              aria-label="Close modal"
              @click="$emit('update:modelValue', false)"
            >
              ✕
            </button>
          </div>

          <!-- Modal Body -->
          <div class="kumo-modal-body">
            <slot />
          </div>

          <!-- Modal Footer -->
          <div v-if="$slots.footer" class="kumo-modal-footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  size: {
    type: String,
    default: 'md', // sm (420px) | md (560px) | lg (720px) | xl (900px)
    validator: (v) => ['sm', 'md', 'lg', 'xl'].includes(v)
  },
  showClose: {
    type: Boolean,
    default: true
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue'])

const onBackdropClick = () => {
  if (props.closeOnBackdrop) {
    emit('update:modelValue', false)
  }
}

const onKeyDown = (e) => {
  if (e.key === 'Escape' && props.modelValue) {
    emit('update:modelValue', false)
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeyDown)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeyDown)
})
</script>

<style scoped>
.kumo-modal-backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 1rem;
  box-sizing: border-box;
}

.kumo-modal-dialog {
  display: flex;
  flex-direction: column;
  width: 100%;
  background-color: #ffffff;
  border: 1px solid var(--color-kumo-line, #e2e8f0);
  border-radius: 8px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.15), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  box-sizing: border-box;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  max-height: 90vh;
}

[data-theme="dark"] .kumo-modal-dialog {
  background-color: #0c0f17;
  border-color: #1e293b;
}

/* SIZES */
.kumo-modal-sm { max-width: 420px; }
.kumo-modal-md { max-width: 560px; }
.kumo-modal-lg { max-width: 720px; }
.kumo-modal-xl { max-width: 900px; }

.kumo-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--color-kumo-line, #f4f4f5);
  gap: 1rem;
}

[data-theme="dark"] .kumo-modal-header {
  border-bottom-color: #1e293b;
}

.kumo-modal-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color-kumo-default, #18181b);
  margin: 0;
  letter-spacing: -0.01em;
}

[data-theme="dark"] .kumo-modal-title {
  color: #f8fafc;
}

.kumo-modal-close-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: #71717a;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.15s ease;
}

.kumo-modal-close-btn:hover {
  background-color: #f4f4f5;
  color: #18181b;
}

[data-theme="dark"] .kumo-modal-close-btn:hover {
  background-color: #151a26;
  color: #f8fafc;
}

.kumo-modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  font-size: 13.5px;
  line-height: 1.5;
  color: var(--text-color-kumo-default, #18181b);
}

[data-theme="dark"] .kumo-modal-body {
  color: #f8fafc;
}

.kumo-modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  background-color: #fafafa;
  border-top: 1px solid var(--color-kumo-line, #f4f4f5);
}

[data-theme="dark"] .kumo-modal-footer {
  background-color: #070a12;
  border-top-color: #1e293b;
}

/* ANIMATION */
.kumo-modal-fade-enter-active,
.kumo-modal-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.kumo-modal-fade-enter-from,
.kumo-modal-fade-leave-to {
  opacity: 0;
  transform: scale(0.97);
}
</style>
