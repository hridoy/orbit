<template>
  <div class="cs-page">
    <div class="breadcrumbs">
      <span class="current">Members</span>
    </div>
    <h1 class="page-title">Members</h1>
    <ComingSoon title="Not connected yet" text="Members and roles are not available yet. Only the dashboard login exists." />
  </div>
</template>

<script setup>
import ComingSoon from '@/components/ComingSoon.vue'
</script>

<style scoped>
.cs-page { display: flex; flex-direction: column; gap: 1rem; min-width: 0; }
</style>
