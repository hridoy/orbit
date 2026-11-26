<template>
  <div :class="['cf-loader', customClass]" :style="{ width: sizePx, height: sizePx }">
    <svg :width="size || 32" :height="size || 32" viewBox="0 0 40 40" class="cf-loader-svg">
      <defs>
        <pattern id="flare-hatch" width="4" height="4" patternTransform="rotate(45 0 0)" patternUnits="userSpaceOnUse">
          <line x1="0" y1="0" x2="0" y2="4" class="cf-hatch-line" />
        </pattern>
      </defs>
      <!-- Base guide circle -->
      <circle cx="20" cy="20" r="18" class="cf-guides" />
      <!-- Concentric animated dashes -->
      <circle cx="20" cy="20" r="18" class="cf-ring cf-ring-a" />
      <circle cx="20" cy="20" r="14" class="cf-ring cf-ring-b" />
      <circle cx="20" cy="20" r="10" class="cf-ring cf-ring-c" />
      <circle cx="20" cy="20" r="6" class="cf-ring cf-ring-d" />
    </svg>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: {
    type: [Number, String],
    default: 32
  },
  customClass: {
    type: String,
    default: ''
  }
})

const sizePx = computed(() => {
  return typeof props.size === 'number' ? `${props.size}px` : props.size
})
</script>

<style scoped>
.cf-loader {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  --cf-duration: 6.4s;
  --cf-play: running;
}

.cf-loader-svg {
  display: block;
}

.cf-guides {
  fill: none;
  stroke: var(--text-color-kumo-placeholder, #a1a1aa);
  stroke-width: 0.6px;
  stroke-dasharray: 1 2;
  opacity: 0.35;
}

.cf-hatch-line {
  stroke: var(--text-color-kumo-placeholder, #a1a1aa);
  stroke-width: 0.5px;
}

.cf-ring {
  fill: none;
  stroke: var(--text-color-kumo-link, #0051c3);
  stroke-width: 1.2px;
  stroke-dasharray: 4 6;
  stroke-linecap: round;
  animation-name: cf-circle-dashes;
  animation-duration: var(--cf-duration);
  animation-timing-function: linear;
  animation-iteration-count: infinite;
  animation-play-state: var(--cf-play);
}

.cf-ring-a {
  stroke: var(--cf-brand-orange, #f6821f);
  stroke-dasharray: 6 8;
}

.cf-ring-b {
  animation-delay: calc(var(--cf-duration) * -0.25);
  animation-direction: reverse;
  stroke: var(--cf-brand-blue, #0051c3);
}

.cf-ring-c {
  animation-delay: calc(var(--cf-duration) * -0.5);
  stroke: var(--cf-brand-orange, #f6821f);
}

.cf-ring-d {
  animation-delay: calc(var(--cf-duration) * -0.75);
  animation-direction: reverse;
  stroke: var(--cf-brand-blue, #0051c3);
}

@keyframes cf-circle-dashes {
  0% { stroke-dashoffset: 0; }
  100% { stroke-dashoffset: -40px; }
}
</style>
