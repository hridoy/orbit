<template>
  <ol class="cf-stepper" role="list">
    <li
      v-for="(step, idx) in steps"
      :key="idx"
      class="cf-step-item"
      :class="{
        'is-done': idx < currentStep,
        'is-active': idx === currentStep,
        'is-clickable': clickable
      }"
      @click="clickable && $emit('step-click', idx)"
    >
      <div class="cf-step-indicator">
        <span v-if="idx < currentStep">✓</span>
        <span v-else>{{ idx + 1 }}</span>
      </div>
      <div class="cf-step-meta">
        <div class="cf-step-title">{{ step.title || step }}</div>
        <div v-if="step.description" class="cf-step-desc">{{ step.description }}</div>
      </div>
      <div v-if="idx < steps.length - 1" class="cf-step-connector" aria-hidden="true"></div>
    </li>
  </ol>
</template>

<script setup>
defineProps({
  steps: {
    type: Array,
    required: true // string[] or { title, description }[]
  },
  currentStep: {
    type: Number,
    default: 0
  },
  clickable: {
    type: Boolean,
    default: false
  }
})

defineEmits(['step-click'])
</script>

<style scoped>
.cf-stepper {
  display: flex;
  align-items: flex-start;
  list-style: none;
  margin: 0;
  padding: 0;
  gap: 12px;
  flex-wrap: wrap;
}

.cf-step-item {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 150px;
}

.cf-step-item.is-clickable {
  cursor: pointer;
}

.cf-step-indicator {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
  background-color: var(--color-kumo-recessed);
  border: 1px solid var(--color-kumo-line);
  color: var(--text-color-kumo-subtle);
  transition: all 0.15s ease;
}

.cf-step-item.is-active .cf-step-indicator {
  background-color: var(--color-kumo-brand);
  border-color: var(--color-kumo-brand);
  color: #ffffff;
}

.cf-step-item.is-done .cf-step-indicator {
  background-color: var(--color-kumo-success-tint);
  border-color: var(--color-kumo-success);
  color: var(--text-color-kumo-success);
}

.cf-step-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-step-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color-kumo-subtle);
}

.cf-step-item.is-active .cf-step-title,
.cf-step-item.is-done .cf-step-title {
  color: var(--text-color-kumo-strong);
}

.cf-step-desc {
  font-size: 11.5px;
  color: var(--text-color-kumo-subtle);
}

.cf-step-connector {
  flex: 1;
  height: 1px;
  min-width: 16px;
  background-color: var(--color-kumo-line);
}

.cf-step-item.is-done .cf-step-connector {
  background-color: var(--color-kumo-success);
}
</style>
