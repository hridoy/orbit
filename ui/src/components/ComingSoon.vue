<template>
  <div class="cs-card">
    <h2 class="cs-title">{{ title }}</h2>
    <p class="cs-text">{{ text }}</p>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  text: { type: String, default: 'This part is not connected to the server yet, so nothing is shown here.' },
})
</script>

<style scoped>
.cs-card {
  background: var(--kumo-base, #fff);
  border: 1px dashed var(--kumo-border, #cbd1d9);
  border-radius: 10px;
  padding: 1.5rem 1rem;
  text-align: center;
  box-sizing: border-box;
}
.cs-title { font-size: 16px; margin: 0 0 0.5rem; }
.cs-text { font-size: 13px; color: #6b7280; margin: 0; }
</style>
