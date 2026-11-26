<template>
  <div>
    <div class="breadcrumbs">
      <span>Manage Account</span>
      <span class="divider">/</span>
      <span class="current">Audit Logs</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Audit Logs & Activity Stream</h1>
        <p class="page-subtitle">Immutable security trail tracking all account configurations, DNS mutations, database queries, and worker deployments.</p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <button class="btn-cf btn-cf-secondary" @click="exportJson">Export JSON</button>
      </div>
    </div>

    <!-- Filter Bar -->
    <div style="display: flex; gap: 0.75rem; margin-bottom: 1.25rem; flex-wrap: wrap;">
      <input 
        v-model="searchQuery" 
        type="text" 
        class="cf-input" 
        placeholder="Filter by event, actor, or subject..." 
        style="max-width: 320px;"
      />
      <select v-model="selectedZone" class="cf-select" style="max-width: 200px;">
        <option value="All">All Zone Scopes</option>
        <option value="northlab.dev">northlab.dev</option>
        <option value="Global">Global</option>
      </select>
    </div>

    <!-- Audit Logs Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Security Audit Trail ({{ filteredLogs.length }} events)</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Retention: 90 days</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Timestamp (UTC)</th>
            <th>Action / Event</th>
            <th>Subject / Resource</th>
            <th>Actor</th>
            <th>IP Address</th>
            <th>Zone Scope</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(log, idx) in filteredLogs" :key="idx">
            <td>
              <code style="font-size: 0.78rem;">{{ log.timestamp }}</code>
            </td>
            <td>
              <strong>{{ log.event }}</strong>
            </td>
            <td>
              <code style="font-size: 0.8rem; color: var(--cf-brand-orange);">{{ log.subject }}</code>
            </td>
            <td>{{ log.actor }}</td>
            <td>
              <code style="font-size: 0.78rem;">{{ log.ip }}</code>
            </td>
            <td>
              <span class="badge-tag-beta">{{ log.zone }}</span>
            </td>
            <td>
              <span class="cf-pill active">
                <span class="status-dot"></span>
                <span>{{ log.status }}</span>
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { mockData } from '@/data/mockData.js'

const logs = ref([...mockData.auditLogs])
const searchQuery = ref('')
const selectedZone = ref('All')

const filteredLogs = computed(() => {
  return logs.value.filter(l => {
    const matchesZone = selectedZone.value === 'All' || l.zone === selectedZone.value
    const q = searchQuery.value.toLowerCase()
    const matchesQuery = !q || 
      l.event.toLowerCase().includes(q) || 
      l.actor.toLowerCase().includes(q) || 
      l.subject.toLowerCase().includes(q)
    return matchesZone && matchesQuery
  })
})

function exportJson() {
  const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(filteredLogs.value, null, 2))
  const downloadAnchor = document.createElement('a')
  downloadAnchor.setAttribute('href', dataStr)
  downloadAnchor.setAttribute('download', `cloudflare-audit-logs-${Date.now()}.json`)
  document.body.appendChild(downloadAnchor)
  downloadAnchor.click()
  downloadAnchor.remove()
}
</script>
