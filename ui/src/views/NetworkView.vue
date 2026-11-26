<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Network</span>
    </div>
    <h1 class="page-title">Network</h1>
    <ComingSoon title="Not connected yet" text="Network settings are not configured yet." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
