<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">AI Gateway</span>
    </div>
    <h1 class="page-title">AI Gateway</h1>
    <ComingSoon title="Not connected yet" text="AI gateway is not available on this server yet." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
