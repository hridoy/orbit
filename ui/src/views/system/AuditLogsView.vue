<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Audit Logs</span>
    </div>
    <h1 class="page-title">Audit Logs</h1>
    <ComingSoon title="Not connected yet" text="Audit logs are not recorded yet. Actions on this server will appear here once logging is added." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
