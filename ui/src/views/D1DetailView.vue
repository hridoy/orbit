<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/d1">D1 Databases</router-link>
      <span class="divider">/</span>
      <span class="current">{{ db.name }}</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.75rem;">
          <h1 class="page-title">🗄️ {{ db.name }}</h1>
          <span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span>
          <span class="badge-tag-beta">SQLite 3.45</span>
        </div>
        <p class="page-subtitle">UUID: <code>{{ db.id }}</code> (Region: {{ db.location }})</p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <button class="btn-cf btn-cf-secondary" @click="copyUuid">
          {{ copiedUuid ? '✓ Copied UUID' : 'Copy Database UUID' }}
        </button>
        <button class="btn-cf btn-cf-primary" @click="executeSql">Run SQL Query</button>
      </div>
    </div>

    <!-- D1 Top Tabs Navigation -->
    <div class="cf-tab-nav">
      <button class="cf-tab-btn" :class="{ active: subTab === 'console' }" @click="subTab = 'console'">SQL Console</button>
      <button class="cf-tab-btn" :class="{ active: subTab === 'tables' }" @click="subTab = 'tables'">Tables & Schema ({{ db.tables.length }})</button>
      <button class="cf-tab-btn" :class="{ active: subTab === 'metrics' }" @click="subTab = 'metrics'">Metrics & Units</button>
      <button class="cf-tab-btn" :class="{ active: subTab === 'time-travel' }" @click="subTab = 'time-travel'">Time Travel Backups</button>
    </div>

    <!-- SQL Console Tab -->
    <div v-if="subTab === 'console'">
      <div class="d1-sql-console">
        <div class="d1-sql-console-head">
          <span style="color: #f8fafc; font-weight: 600;">Interactive D1 SQL Runner</span>
          <div style="display: flex; gap: 0.5rem;">
            <button class="btn-cf btn-cf-primary btn-cf-sm" @click="executeSql">▶ Execute (Ctrl+Enter)</button>
          </div>
        </div>
        <textarea
          v-model="sqlQuery"
          class="d1-sql-editor"
          @keydown.ctrl.enter="executeSql"
        ></textarea>
      </div>

      <!-- Query Results -->
      <div class="kumo-table-container" style="margin-top: 1.5rem;">
        <div class="table-header-title">
          <span>Query Results ({{ results.length }} rows returned in {{ execDuration }}ms)</span>
          <span class="cf-pill active"><span class="status-dot"></span><span>200 OK</span></span>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>id</th>
              <th>title</th>
              <th>lang</th>
              <th>author</th>
              <th>downloads</th>
              <th>status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in results" :key="r.id">
              <td><code>{{ r.id }}</code></td>
              <td><strong>{{ r.title }}</strong></td>
              <td><span class="badge-tag-beta">{{ r.lang }}</span></td>
              <td>{{ r.author }}</td>
              <td>{{ r.downloads.toLocaleString() }}</td>
              <td><span class="cf-pill active"><span class="status-dot"></span><span>{{ r.status }}</span></span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Tables & Schema Tab -->
    <div v-else-if="subTab === 'tables'" class="kumo-table-container">
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Table Name</th>
            <th>Rows Count</th>
            <th>Columns</th>
            <th>Table Size</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in db.tables" :key="t.name">
            <td><strong>{{ t.name }}</strong></td>
            <td>{{ t.rows.toLocaleString() }}</td>
            <td>{{ t.columns }}</td>
            <td>{{ t.size }}</td>
            <td>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="queryTable(t.name)">SELECT *</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Metrics Tab -->
    <div v-else class="stats-overview-grid">
      <MetricCard label="Row Reads (24h)" :value="db.read_units_24h" delta="↗ 14.8%" strokeColor="#059669" footer="Anycast local edge reads" />
      <MetricCard label="Row Writes (24h)" :value="db.write_units_24h" delta="↘ 2.1%" strokeColor="var(--cf-brand-orange)" footer="Leader replicate write" />
      <MetricCard label="Total Queries" :value="db.queries_count" delta="↗ 8.2%" strokeColor="#0051c3" footer="Prepared statements" />
      <MetricCard label="Database Size" :value="db.size" strokeColor="#7c3aed" footer="Limit: 10 GB Included" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'
import MetricCard from '@/components/MetricCard.vue'

const db = mockData.databases[0]
const subTab = ref('console')
const sqlQuery = ref("SELECT id, title, lang, author, downloads, status FROM code_snippets WHERE status = 'published' LIMIT 10;")
const results = ref([...db.sampleQueryResults])
const execDuration = ref('0.8')
const copiedUuid = ref(false)

function executeSql() {
  execDuration.value = (Math.random() * 0.8 + 0.4).toFixed(1)
  alert(`Executing SQL on ${db.name}:\n${sqlQuery.value}\n5 rows returned in ${execDuration.value}ms`)
}

function queryTable(name) {
  sqlQuery.value = `SELECT * FROM ${name} LIMIT 10;`
  subTab.value = 'console'
}

function copyUuid() {
  if (navigator.clipboard) navigator.clipboard.writeText(db.id)
  copiedUuid.value = true
  setTimeout(() => { copiedUuid.value = false }, 2000)
}
</script>
