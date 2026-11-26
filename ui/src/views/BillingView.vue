<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Billing</span>
    </div>
    <h1 class="page-title">Billing</h1>
    <ComingSoon title="Not connected yet" text="This server has no billing. Nothing is charged here." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
