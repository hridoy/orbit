<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Analytics</span>
    </div>
    <h1 class="page-title">Analytics</h1>
    <ComingSoon title="Not connected yet" text="Traffic analytics are not collected yet." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
