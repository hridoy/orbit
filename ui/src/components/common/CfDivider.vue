<template>
  <div v-if="label" class="cf-divider-with-label" role="separator">
    <span class="cf-divider-line"></span>
    <span class="cf-divider-text">{{ label }}</span>
    <span class="cf-divider-line"></span>
  </div>
  <hr v-else class="cf-divider" :class="{ 'is-vertical': vertical }" role="separator" />
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: ''
  },
  vertical: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.cf-divider {
  border: none;
  border-top: 1px solid var(--color-kumo-line);
  margin: 0;
  width: 100%;
}

.cf-divider.is-vertical {
  border-top: none;
  border-left: 1px solid var(--color-kumo-line);
  width: 1px;
  height: 20px;
  display: inline-block;
}

.cf-divider-with-label {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.cf-divider-line {
  flex: 1;
  height: 1px;
  background-color: var(--color-kumo-line);
}

.cf-divider-text {
  font-size: 11.5px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-color-kumo-subtle);
}
</style>
