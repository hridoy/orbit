<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Domain overview</span>
    </div>
    <h1 class="page-title">Domain overview</h1>
    <ComingSoon title="Not connected yet" text="Per-domain analytics are not collected yet. Manage the domain certificate in SSL/TLS." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
