<template>
  <div :class="['cf-skeleton-line', customClass]" :style="{ width: width || '100%', height: height || '0.5rem' }"></div>
</template>

<script setup>
defineProps({
  width: {
    type: String,
    default: '100%'
  },
  height: {
    type: String,
    default: '0.5rem'
  },
  customClass: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.cf-skeleton-line {
  background-color: #f3f4f6;
  border-radius: 2px;
  position: relative;
  overflow: hidden;
  display: inline-block;
}

.cf-skeleton-line::after {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0.6) 50%, rgba(255, 255, 255, 0) 100%);
  animation: cf-shimmer 1.5s infinite ease-in-out;
}

[data-theme="dark"] .cf-skeleton-line {
  background-color: rgba(255, 255, 255, 0.08);
}

[data-theme="dark"] .cf-skeleton-line::after {
  background: linear-gradient(90deg, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0.1) 50%, rgba(255, 255, 255, 0) 100%);
}

@keyframes cf-shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}
</style>
