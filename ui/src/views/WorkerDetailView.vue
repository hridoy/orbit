<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">App details</span>
    </div>
    <h1 class="page-title">App details</h1>
    <ComingSoon title="Not connected yet" text="Open an app from the Apps list on the home page." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
