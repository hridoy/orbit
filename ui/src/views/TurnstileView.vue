<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Turnstile</span>
    </div>
    <h1 class="page-title">Turnstile</h1>
    <ComingSoon title="Not connected yet" text="Turnstile is not set up on this server." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
