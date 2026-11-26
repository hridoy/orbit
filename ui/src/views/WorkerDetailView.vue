<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/workers">Compute</router-link>
      <span class="divider">/</span>
      <span class="current">{{ worker.name }}</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.75rem;">
          <h1 class="page-title">⚡ {{ worker.name }}</h1>
          <span class="cf-pill active"><span class="status-dot"></span><span>{{ worker.status }}</span></span>
          <span class="badge-tag-beta">Production</span>
        </div>
        <p class="page-subtitle">Serverless Worker running on V8 Isolates with instant global cold starts & zero cold start latency.</p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <a :href="`https://${worker.custom_domain}`" target="_blank" class="btn-cf btn-cf-secondary">Visit Site ↗</a>
        <button class="btn-cf btn-cf-primary" @click="alert('Triggering zero-downtime blue/green deployment via wrangler...')">Deploy Version</button>
      </div>
    </div>

    <!-- Worker Top Tabs Navigation -->
    <div class="cf-tab-nav">
      <button class="cf-tab-btn" :class="{ active: activeTab === 'overview' }" @click="activeTab = 'overview'">Overview</button>
      <button class="cf-tab-btn" :class="{ active: activeTab === 'metrics' }" @click="activeTab = 'metrics'">Metrics</button>
      <button class="cf-tab-btn" :class="{ active: activeTab === 'deployments' }" @click="activeTab = 'deployments'">Deployments</button>
      <button class="cf-tab-btn" :class="{ active: activeTab === 'logs' }" @click="activeTab = 'logs'">Observability & Live Logs</button>
      <button class="cf-tab-btn" :class="{ active: activeTab === 'settings' }" @click="activeTab = 'settings'">Settings & Bindings</button>
    </div>

    <!-- Overview Tab -->
    <div v-if="activeTab === 'overview' || activeTab === 'metrics'">
      <div class="stats-overview-grid">
        <CfStatCard
          label="24h Invocations"
          :value="worker.invocations_24h"
          :delta="worker.invocations_delta"
          delta-type="up"
          tooltip="Total isolate invocations in past 24h"
          :sparkline-data="[15, 20, 18, 25, 30, 42]"
        />
        <CfStatCard
          label="Median CPU Time"
          :value="worker.median_cpu_time"
          :delta="worker.median_cpu_delta"
          delta-type="down"
          tooltip="p50 CPU runtime per invocation"
          :sparkline-data="[12, 10, 8, 6, 5, 3]"
        />
        <CfStatCard
          label="Active Version"
          :value="worker.current_version.split(' ')[0]"
          delta="Production"
          delta-type="up"
          :tooltip="`Git commit: ${worker.current_version.split(' ')[1]}`"
        />
        <CfStatCard
          label="Active Routes"
          :value="worker.custom_domain"
          delta="Anycast"
          delta-type="neutral"
          :tooltip="`Subdomain: ${worker.subdomain}`"
        />
      </div>

      <!-- Connected Bindings Summary Card -->
      <div class="kumo-table-container" style="margin-top: 1.5rem;">
        <div class="table-header-title">
          <span>Connected Resources & Bindings (D1, R2, Queues)</span>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="activeTab = 'settings'">Configure Bindings →</button>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Binding Variable</th>
              <th>Resource Type</th>
              <th>Connected Target</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in worker.bindings" :key="b.name">
              <td><code>env.{{ b.name }}</code></td>
              <td><span class="badge-tag-beta">{{ b.type }}</span></td>
              <td><strong>{{ b.target }}</strong></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Deployments Tab -->
    <div v-else-if="activeTab === 'deployments'" class="kumo-table-container">
      <div class="table-header-title">
        <span>Deployment Release History ({{ worker.deployments.length }})</span>
        <button class="btn-cf btn-cf-primary btn-cf-sm" @click="alert('Rollback triggered.')">Rollback to v2.8.3</button>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Version</th>
            <th>Git Commit</th>
            <th>Triggered By</th>
            <th>Timestamp (UTC)</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="dep in worker.deployments" :key="dep.version">
            <td><strong>{{ dep.version }}</strong></td>
            <td><code>{{ dep.commit }}</code></td>
            <td>{{ dep.author }}</td>
            <td>{{ dep.date }}</td>
            <td>
              <span class="cf-pill" :class="dep.status.includes('Active') ? 'active' : 'stopped'">
                <span class="status-dot"></span>
                <span>{{ dep.status }}</span>
              </span>
            </td>
            <td>
              <span v-if="dep.status.includes('Active')" style="color: var(--text-color-kumo-subtle); font-size: 0.8rem;">Current Live</span>
              <button v-else class="btn-cf btn-cf-secondary btn-cf-sm" @click="alert(`Rolling back to ${dep.version}`)">Rollback</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Observability & Live Logs Tab -->
    <div v-else-if="activeTab === 'logs'" class="d1-sql-console" style="height: 380px; display: flex; flex-direction: column;">
      <div class="d1-sql-console-head">
        <span style="color: #f8fafc; font-weight: 600;">⚡ Realtime Log Stream — {{ worker.name }}</span>
        <div style="display: flex; gap: 0.5rem; align-items: center;">
          <span class="status-dot" style="background: #10b981;"></span>
          <span style="color: #94a3b8; font-size: 0.75rem;">Listening on Anycast edges...</span>
        </div>
      </div>
      <div style="flex: 1; overflow-y: auto; font-family: var(--font-mono); font-size: 0.82rem; line-height: 1.6; color: #94a3b8; padding-top: 0.5rem;">
        <div>[2026-11-22 08:52:10 UTC] [INFO] GET https://codecollection.northlab.dev/ HTTP/2 200 OK (11.2ms) [IAD Edge]</div>
        <div>[2026-11-22 08:52:14 UTC] [INFO] D1 SQL Execution: SELECT * FROM code_snippets WHERE id = 1 -> 1 row returned (0.8ms)</div>
        <div>[2026-11-22 08:52:18 UTC] [INFO] R2 Object Stream: GET /thumbnails/app-hero-light.png -> 200 OK (Cache HIT)</div>
        <div>[2026-11-22 08:52:25 UTC] [INFO] GET https://codecollection.northlab.dev/api/snippets HTTP/2 200 OK (1.4ms) [FRA Edge]</div>
        <div style="color: #38bdf8;">[2026-11-22 08:52:32 UTC] [METRIC] Memory: 14.8 MB / 128 MB | CPU: 1.1ms | Isolate healthy</div>
      </div>
    </div>

    <!-- Settings & Bindings Tab -->
    <div v-else-if="activeTab === 'settings'" class="ui-kit-section">
      <h3 class="ui-kit-title">Runtime Environment Variables & Secrets</h3>
      <p class="ui-kit-desc">Encrypted variables injected directly into isolate memory at boot time.</p>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Variable Name</th>
            <th>Type</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in worker.secrets" :key="s.key">
            <td><strong>{{ s.key }}</strong></td>
            <td><span class="badge-tag-beta">{{ s.type }}</span></td>
            <td><code>{{ s.value }}</code></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const worker = mockData.workers[0] // code-collection
const activeTab = ref('overview')
</script>
