<script setup>
import { useToast } from '../../composables/useToast'

const { toasts, remove } = useToast()

const iconMap = {
  success: 'check',
  error: 'x',
  warning: 'alert-triangle',
  info: 'info'
}
</script>

<template>
  <Teleport to="body">
    <div class="cf-toast-container" aria-live="polite">
      <TransitionGroup name="cf-toast-item">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="cf-toast"
          :class="[`type-${toast.type}`]"
          role="status"
        >
          <div class="cf-toast-icon-wrap">
            <CfIcon :name="iconMap[toast.type] || 'info'" size="16" />
          </div>

          <div class="cf-toast-content">
            <h5 v-if="toast.title" class="cf-toast-title">{{ toast.title }}</h5>
            <p v-if="toast.message" class="cf-toast-msg">{{ toast.message }}</p>
          </div>

          <button 
            type="button" 
            class="cf-toast-close" 
            aria-label="Dismiss toast"
            @click="remove(toast.id)"
          >
            <CfIcon name="x" size="14" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.cf-toast-container {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 400px;
  pointer-events: none;
}

.cf-toast {
  pointer-events: auto;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 14px;
  background-color: var(--kumo-surface-default, #ffffff);
  border-radius: var(--kumo-radius-md, 6px);
  border: 1px solid var(--kumo-border-default, #e4e4e7);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.15), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
  min-width: 280px;
}

.cf-toast-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 1px;
}

.type-success .cf-toast-icon-wrap {
  color: #059669;
}

.type-error .cf-toast-icon-wrap {
  color: #dc2626;
}

.type-warning .cf-toast-icon-wrap {
  color: #d97706;
}

.type-info .cf-toast-icon-wrap {
  color: #0051c3;
}

.cf-toast-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-toast-title {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--kumo-fg-default, #18181b);
  line-height: 18px;
}

.cf-toast-msg {
  margin: 0;
  font-size: 12px;
  color: var(--kumo-fg-muted, #71717a);
  line-height: 16px;
}

.cf-toast-close {
  background: transparent;
  border: none;
  color: var(--kumo-fg-muted, #71717a);
  cursor: pointer;
  padding: 2px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.cf-toast-close:hover {
  color: var(--kumo-fg-default, #18181b);
  background-color: var(--kumo-control-bg-hover, #f4f4f5);
}

.cf-toast-item-enter-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.cf-toast-item-leave-active {
  transition: all 0.2s cubic-bezier(0.7, 0, 0.84, 0);
}

.cf-toast-item-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.cf-toast-item-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.95);
}
</style>
