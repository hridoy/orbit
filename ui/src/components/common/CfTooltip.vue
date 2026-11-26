<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  text: {
    type: String,
    default: ''
  },
  position: {
    type: String,
    default: 'top', // 'top' | 'bottom' | 'left' | 'right'
    validator: (v) => ['top', 'bottom', 'left', 'right'].includes(v)
  },
  delay: {
    type: Number,
    default: 150
  }
})

const isVisible = ref(false)
let timer = null

function show() {
  clearTimeout(timer)
  timer = setTimeout(() => {
    isVisible.value = true
  }, props.delay)
}

function hide() {
  clearTimeout(timer)
  isVisible.value = false
}

onBeforeUnmount(() => {
  clearTimeout(timer)
})
</script>

<template>
  <div class="cf-tooltip-container" @mouseenter="show" @mouseleave="hide" @focusin="show" @focusout="hide">
    <slot />
    
    <Transition name="cf-tooltip-fade">
      <div 
        v-if="isVisible && (text || $slots.content)" 
        class="cf-tooltip-bubble" 
        :class="[`pos-${position}`]"
        role="tooltip"
      >
        <slot name="content">{{ text }}</slot>
        <span class="cf-tooltip-arrow"></span>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.cf-tooltip-container {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.cf-tooltip-bubble {
  position: absolute;
  z-index: 1000;
  padding: 6px 10px;
  background-color: var(--kumo-control-bg-active, #000000);
  color: #ffffff;
  font-size: 11px;
  font-weight: 500;
  line-height: 14px;
  border-radius: var(--kumo-radius-sm, 4px);
  white-space: nowrap;
  pointer-events: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
}

.pos-top {
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-6px);
}

.pos-bottom {
  top: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(6px);
}

.pos-left {
  right: 100%;
  top: 50%;
  transform: translateY(-50%) translateX(-6px);
}

.pos-right {
  left: 100%;
  top: 50%;
  transform: translateY(-50%) translateX(6px);
}

.cf-tooltip-arrow {
  position: absolute;
  width: 0;
  height: 0;
  border-style: solid;
}

.pos-top .cf-tooltip-arrow {
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 4px 4px 0 4px;
  border-color: var(--kumo-control-bg-active, #000000) transparent transparent transparent;
}

.pos-bottom .cf-tooltip-arrow {
  top: -4px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 0 4px 4px 4px;
  border-color: transparent transparent var(--kumo-control-bg-active, #000000) transparent;
}

.pos-left .cf-tooltip-arrow {
  right: -4px;
  top: 50%;
  transform: translateY(-50%);
  border-width: 4px 0 4px 4px;
  border-color: transparent transparent transparent var(--kumo-control-bg-active, #000000);
}

.pos-right .cf-tooltip-arrow {
  left: -4px;
  top: 50%;
  transform: translateY(-50%);
  border-width: 4px 4px 4px 0;
  border-color: transparent var(--kumo-control-bg-active, #000000) transparent transparent;
}

.cf-tooltip-fade-enter-active,
.cf-tooltip-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.cf-tooltip-fade-enter-from,
.cf-tooltip-fade-leave-to {
  opacity: 0;
}
</style>
