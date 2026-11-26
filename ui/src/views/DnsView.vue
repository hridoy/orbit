<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">DNS records</span>
    </div>
    <h1 class="page-title">DNS records</h1>
    <ComingSoon title="Not connected yet" text="DNS records are not managed here. Point your domain at this server's IP in your DNS provider." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
