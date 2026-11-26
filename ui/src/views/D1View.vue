<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Database</router-link>
      <span class="divider">/</span>
      <span class="current">Databases</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Databases</h1>
        <p class="page-subtitle">Manage SQLite, PostgreSQL, and MySQL database instances running on this VPS node with integrated SQL console.</p>
      </div>
      <CfButton variant="primary" icon="plus" @click="alert('Create Database wizard.')">Create Database</CfButton>
    </div>

    <CfCard body-class="d1-table-body">
      <CfTable :columns="dbColumns" :data="mockData.databases" row-key="id">
        <template #col-name="{ row }">
          <router-link :to="`/d1/${row.name}`" style="font-weight: 600; color: var(--text-color-kumo-strong); text-decoration: none; display: inline-flex; align-items: center; gap: 0.45rem;">
            <CfIcon name="d1" size="15" style="color: #0051c3;" />
            <span>{{ row.name }}</span>
          </router-link>
        </template>
        <template #col-engine="{ val }">
          <CfBadge variant="info" size="sm">{{ val }}</CfBadge>
        </template>
        <template #col-reads="{ row }">
          <span>{{ row.read_units_24h }} / {{ row.write_units_24h }}</span>
        </template>
        <template #actions="{ row }">
          <router-link :to="`/d1/${row.name}`" style="text-decoration: none;">
            <CfButton variant="secondary" size="sm">Open Console →</CfButton>
          </router-link>
        </template>
      </CfTable>
    </CfCard>
  </div>
</template>

<script setup>
import { mockData } from '@/data/mockData.js'

const dbColumns = [
  { key: 'name', label: 'Database Name' },
  { key: 'engine', label: 'Engine', width: '130px' },
  { key: 'location', label: 'Node / Path', width: '140px' },
  { key: 'size', label: 'Size', width: '110px' },
  { key: 'tables_count', label: 'Tables', width: '90px' },
  { key: 'reads', label: '24h Reads / Writes', width: '160px' }
]
</script>

<style scoped>
:deep(.d1-table-body) {
  padding: 0;
}
</style>
