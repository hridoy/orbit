<template>
  <div class="cf-app-detail-page">
    <!-- 1. TOP BREADCRUMB ROW -->
    <div class="cf-ad-breadcrumb-row">
      <div class="cf-ad-breadcrumbs">
        <CfIcon name="edgeworker" size="15" style="color: #f6821f;" />
        <router-link to="/apps" class="cf-ad-bc-link">Apps</router-link>
        <span class="cf-ad-bc-sep">›</span>
        <span class="cf-ad-bc-current">{{ app.name || appId }}</span>
        <span class="cf-ad-bc-sep">›</span>
        <span class="cf-ad-env-pill">
          <CfIcon name="cis-globe" size="12" />
          <span>Production</span>
        </span>
      </div>
    </div>

    <!-- 2. SUB-HEADER NAVIGATION & ACTION ROW (Exact Cloudflare Worker Header) -->
    <div class="cf-ad-nav-action-row">
      <CfTabs
        v-model="activeTab"
        :tabs="[
          { label: 'Overview', value: 'overview' },
          { label: 'Metrics', value: 'metrics' },
          { label: 'Deployments', value: 'deployments' },
          { label: 'Logs', value: 'logs' },
          { label: 'Domains', value: 'domains' },
          { label: 'Settings', value: 'settings' }
        ]"
      />

      <div class="cf-ad-actions">
        <CfButton variant="secondary" size="md" icon="refresh" @click="fetchApp">Refresh</CfButton>
        <CfButton variant="secondary" size="md" @click="activeTab = 'logs'">View logs</CfButton>
        <a
          :href="app.domain ? `http://${app.domain}` : `http://127.0.0.1:${app.active_port || 8001}`"
          target="_blank"
          style="text-decoration: none;"
        >
          <CfButton variant="primary" size="md">Visit ↗</CfButton>
        </a>
      </div>
    </div>

    <!-- ================================================================ -->
    <!-- TAB 1: OVERVIEW (Exact Cloudflare Worker Overview Structure)      -->
    <!-- ================================================================ -->
    <div v-if="activeTab === 'overview'" class="cf-ad-stack">
      <!-- Top Summary Bar -->
      <CfCard>
        <div class="cf-ad-summary-bar">
          <div class="cf-ad-summary-left">
            <CfIcon name="cis-globe" size="15" style="color: var(--text-color-kumo-subtle);" />
            <a
              :href="app.domain ? `http://${app.domain}` : `http://127.0.0.1:${app.active_port || 8001}`"
              target="_blank"
              class="cf-ad-domain-link"
            >
              {{ app.domain || `127.0.0.1:${app.active_port || 8001}` }}
            </a>
            <CfBadge :variant="app.status === 'running' ? 'success' : 'neutral'" size="sm" dot>
              {{ app.status || 'stopped' }}
            </CfBadge>
          </div>
          <div class="cf-ad-summary-right">
            <span class="cf-ad-repo-pill">
              <CfIcon name="server" size="12" />
              <span>{{ app.app_id || appId }}</span>
            </span>
            <span class="cf-ad-muted-meta">Port :{{ app.active_port || 8001 }}</span>
            <span class="cf-ad-muted-meta">{{ formatRelative(app.updated_at) }}</span>
          </div>
        </div>
      </CfCard>

      <!-- Interactive Architecture & Bindings Diagram Card -->
      <CfCard title="Architecture & Bindings" subtitle="Traffic ingress, runtime process, and connected VPS resources">
        <div class="cf-arch-diagram">
          <!-- Left Column: Inputs / Triggers -->
          <div class="cf-arch-col">
            <div class="cf-arch-col-label">Inputs</div>
            <div class="cf-arch-pill-card" @click="activeTab = 'domains'">
              <span>Domains & Routes</span>
              <CfBadge variant="neutral" size="sm">{{ app.domain ? 1 : 0 }}</CfBadge>
            </div>
            <div class="cf-arch-pill-card">
              <span>Reverse Proxy Port</span>
              <CfBadge variant="info" size="sm">:{{ app.active_port || 8001 }}</CfBadge>
            </div>
            <div class="cf-arch-pill-card">
              <span>Health Probe</span>
              <CfBadge variant="neutral" size="sm">{{ app.health_check_path || '/healthz' }}</CfBadge>
            </div>
          </div>

          <!-- Connector -->
          <div class="cf-arch-connector">
            <div class="cf-arch-line"></div>
          </div>

          <!-- Center Node: App Process Card -->
          <div class="cf-arch-center-card">
            <div class="cf-arch-center-head">
              <div style="display: flex; align-items: center; gap: 0.5rem;">
                <CfIcon name="edgeworker" size="18" style="color: #f6821f;" />
                <strong>{{ app.name || appId }}</strong>
              </div>
              <span class="cf-arch-status-dot" :class="{ 'is-live': app.status === 'running' }"></span>
            </div>
            <div class="cf-arch-center-meta">
              <code>{{ app.binary_entrypoint || './server' }}</code>
            </div>
            <div class="cf-arch-center-badges">
              <div class="cf-arch-obs-row" @click="activeTab = 'logs'">
                <span>Process Logs</span>
                <CfBadge variant="success" size="sm">Enabled</CfBadge>
              </div>
              <div class="cf-arch-obs-row" @click="activeTab = 'settings'">
                <span>Auto-Restart</span>
                <CfBadge :variant="app.auto_restart ? 'success' : 'neutral'" size="sm">
                  {{ app.auto_restart ? 'Enabled' : 'Manual' }}
                </CfBadge>
              </div>
            </div>
          </div>

          <!-- Connector -->
          <div class="cf-arch-connector">
            <div class="cf-arch-line"></div>
          </div>

          <!-- Right Column: Bindings -->
          <div class="cf-arch-col">
            <div class="cf-arch-col-head">
              <span class="cf-arch-col-label">Bindings</span>
              <button type="button" class="cf-arch-add-link" @click="activeTab = 'settings'">+ Add binding</button>
            </div>
            <div v-for="b in bindings" :key="b.name" class="cf-arch-pill-card" @click="activeTab = 'settings'">
              <span style="display: inline-flex; align-items: center; gap: 0.4rem;">
                <CfBadge variant="neutral" size="sm">{{ b.type }}</CfBadge>
                <span style="font-weight: 500;">{{ b.name }}</span>
              </span>
              <span style="font-size: 11.5px; color: var(--text-color-kumo-subtle);">{{ b.target }}</span>
            </div>
          </div>
        </div>
      </CfCard>

      <!-- Metrics Summary Card (3-column sparkline grid) -->
      <CfCard>
        <template #title>
          <div style="display: flex; align-items: center; gap: 0.5rem;">
            <span style="font-size: 15px; font-weight: 600;">Metrics</span>
            <CfBadge variant="neutral" size="sm">Last 24 hours</CfBadge>
          </div>
        </template>
        <template #actions>
          <CfButton variant="ghost" size="sm" @click="activeTab = 'metrics'">View metrics →</CfButton>
        </template>
        <div class="cf-ad-metrics-3col">
          <CfStatCard
            label="RAM Allocation"
            :value="`${app.ram_limit_mb || 256} MB`"
            :delta="stats?.memory_mb ? `${stats.memory_mb} MB used` : 'cgroups'"
            delta-type="up"
            :sparkline-data="[42, 48, 45, 52, 50, 56]"
          />
          <CfStatCard
            label="CPU Quota"
            :value="`${app.cpu_quota_pct || 100}%`"
            delta="Healthy"
            delta-type="up"
            :sparkline-data="[12, 16, 14, 18, 15, 14]"
          />
          <CfStatCard
            label="PIDs Limit"
            :value="`${app.pids_limit || 128}`"
            delta="Zero OOM"
            delta-type="neutral"
            :sparkline-data="[4, 5, 4, 6, 5, 5]"
          />
        </div>
      </CfCard>

      <!-- Deployments Summary Card -->
      <CfCard title="Deployments" body-class="cf-ad-table-body">
        <template #actions>
          <CfButton variant="ghost" size="sm" @click="activeTab = 'deployments'">View all →</CfButton>
        </template>
        <CfTable :columns="deployColumns" :data="deployments" row-key="id">
          <template #col-version="{ val }">
            <div style="display: inline-flex; align-items: center; gap: 0.35rem;">
              <code class="cf-ad-code-pill">{{ val }}</code>
              <CfButton variant="ghost" size="sm" icon="clipboard" icon-only title="Copy version" @click="copyText(val)" />
            </div>
          </template>
          <template #col-status="{ val }">
            <CfBadge :variant="val === 'Ready' ? 'success' : 'neutral'" size="sm" dot>{{ val }}</CfBadge>
          </template>
          <template #col-source="{ val }">
            <CfBadge variant="neutral" size="sm">{{ val }}</CfBadge>
          </template>
        </CfTable>
      </CfCard>

      <!-- Bottom Grid: Domains & Next Steps -->
      <div class="cf-ad-bottom-2col">
        <CfCard title="Domains and routes">
          <template #actions>
            <CfButton variant="ghost" size="sm" @click="activeTab = 'domains'">Manage →</CfButton>
          </template>
          <div class="cf-ad-kv-list">
            <div class="cf-ad-kv-row">
              <span class="cf-ad-kv-label">Primary Domain</span>
              <span class="cf-ad-kv-val">{{ app.domain || 'None configured' }}</span>
            </div>
            <div class="cf-ad-kv-row">
              <span class="cf-ad-kv-label">Active Upstream</span>
              <code>127.0.0.1:{{ app.active_port || 8001 }}</code>
            </div>
            <div class="cf-ad-kv-row">
              <span class="cf-ad-kv-label">Idle / Blue-Green Port</span>
              <code>127.0.0.1:{{ app.idle_port || 8002 }}</code>
            </div>
          </div>
        </CfCard>

        <CfCard title="Next steps">
          <div class="cf-ad-kv-list">
            <div class="cf-ad-next-row" @click="activeTab = 'logs'">
              <span>Stream real-time stdout/stderr logs</span>
              <CfIcon name="chevron-right" size="14" />
            </div>
            <div class="cf-ad-next-row" @click="activeTab = 'settings'">
              <span>Configure environment variables & secrets</span>
              <CfIcon name="chevron-right" size="14" />
            </div>
            <div class="cf-ad-next-row" @click="activeTab = 'domains'">
              <span>Bind custom domain & Auto-SSL</span>
              <CfIcon name="chevron-right" size="14" />
            </div>
          </div>
        </CfCard>
      </div>
    </div>

    <!-- ================================================================ -->
    <!-- TAB 2: METRICS                                                    -->
    <!-- ================================================================ -->
    <div v-else-if="activeTab === 'metrics'" class="cf-ad-stack">
      <div class="cf-ad-controls-bar">
        <div style="width: 220px;">
          <CfSelect v-model="selectedVersion" :options="['All deployed versions', app.current_version || 'Current release']" />
        </div>
        <div style="display: flex; align-items: center; gap: 0.5rem;">
          <div style="width: 170px;">
            <CfSelect v-model="metricsRange" :options="['Last 30 minutes', 'Last 1 hour', 'Last 24 hours', 'Last 7 days']" />
          </div>
          <CfButton variant="secondary" icon="refresh" icon-only title="Refresh metrics" @click="fetchApp" />
        </div>
      </div>

      <div class="stats-overview-grid">
        <CfStatCard
          label="Status"
          :value="app.status || 'stopped'"
          :delta="app.auto_restart ? 'Auto-restart' : 'Manual'"
          :delta-type="app.status === 'running' ? 'up' : 'neutral'"
        />
        <CfStatCard
          label="Active Port"
          :value="`:${app.active_port || 8001}`"
          :delta="`Idle :${app.idle_port || 8002}`"
          delta-type="neutral"
        />
        <CfStatCard
          label="RAM Limit"
          :value="`${app.ram_limit_mb || 256} MB`"
          :delta="stats?.memory_mb ? `${stats.memory_mb} MB used` : 'cgroups'"
          delta-type="up"
          :sparkline-data="[38, 42, 40, 44, 46, 45]"
        />
        <CfStatCard
          label="CPU Quota"
          :value="`${app.cpu_quota_pct || 100}%`"
          :delta="`PIDs max: ${app.pids_limit || 128}`"
          delta-type="neutral"
          :sparkline-data="[10, 14, 12, 18, 15, 13]"
        />
      </div>

      <CfCard title="Active Deployment">
        <div class="cf-ad-active-deploy-grid">
          <div>
            <div class="cf-ad-kv-label">Version ID</div>
            <code class="cf-ad-code-pill">{{ app.current_version || 'v1.0.0' }}</code>
          </div>
          <div>
            <div class="cf-ad-kv-label">Traffic %</div>
            <div style="margin-top: 4px; width: 120px;">
              <CfProgress :value="100" variant="brand" size="sm" show-label />
            </div>
          </div>
          <div>
            <div class="cf-ad-kv-label">Health Check</div>
            <div style="font-weight: 600;">{{ app.health_check_path || '/healthz' }}</div>
          </div>
          <div>
            <div class="cf-ad-kv-label">Timeout</div>
            <div style="font-weight: 600;">{{ app.health_timeout_sec || 10 }}s</div>
          </div>
        </div>
      </CfCard>
    </div>

    <!-- ================================================================ -->
    <!-- TAB 3: DEPLOYMENTS                                                -->
    <!-- ================================================================ -->
    <div v-else-if="activeTab === 'deployments'" class="cf-ad-stack">
      <CfCard title="Active deployment">
        <div class="cf-ad-active-deploy-grid">
          <div>
            <div class="cf-ad-kv-label">Version</div>
            <code class="cf-ad-code-pill">{{ app.current_version || 'Initial' }}</code>
          </div>
          <div>
            <div class="cf-ad-kv-label">Port</div>
            <div style="font-weight: 600;">:{{ app.active_port || 8001 }}</div>
          </div>
          <div>
            <div class="cf-ad-kv-label">Entrypoint</div>
            <code>{{ app.binary_entrypoint || './server' }}</code>
          </div>
          <div>
            <div class="cf-ad-kv-label">Status</div>
            <CfBadge :variant="app.status === 'running' ? 'success' : 'neutral'" dot>
              {{ app.status || 'stopped' }}
            </CfBadge>
          </div>
        </div>
      </CfCard>

      <CfCard title="Deployment history" body-class="cf-ad-table-body">
        <template #actions>
          <CfButton variant="secondary" size="sm" icon="refresh" @click="fetchDeployments">Refresh</CfButton>
        </template>
        <CfTable
          :columns="deployColumns"
          :data="filteredDeployments"
          row-key="id"
          searchable
          search-placeholder="Filter deployments..."
          column-picker
          density-toggle
        >
          <template #toolbar>
            <CfTabs v-model="deployStatusTab" :tabs="['All', 'Active', 'Previous']" />
          </template>
          <template #col-version="{ val }">
            <div style="display: inline-flex; align-items: center; gap: 0.35rem;">
              <code class="cf-ad-code-pill">{{ val }}</code>
              <CfButton variant="ghost" size="sm" icon="clipboard" icon-only title="Copy Version ID" @click="copyText(val)" />
            </div>
          </template>
          <template #col-traffic="{ val }">
            <div style="display: flex; align-items: center; gap: 8px; min-width: 110px;">
              <div style="flex: 1; height: 6px; border-radius: 9999px; background: var(--color-kumo-recessed, #f4f4f5); overflow: hidden;">
                <div :style="{ width: `${val}%`, height: '100%', background: 'var(--color-kumo-brand, #f6821f)' }"></div>
              </div>
              <span class="font-mono" style="font-size: 0.78rem;">{{ val }}%</span>
            </div>
          </template>
          <template #col-cpu="{ val }">
            <div style="width: 90px;">
              <CfSparkline :values="val" :height="24" />
            </div>
          </template>
          <template #col-status="{ val }">
            <CfBadge :variant="val === 'Active' || val === 'Ready' ? 'success' : 'neutral'" size="sm" dot>{{ val }}</CfBadge>
          </template>
          <template #col-source="{ val }">
            <CfBadge variant="neutral" size="sm">{{ val }}</CfBadge>
          </template>
          <template #actions="{ row }">
            <CfDropdownMenu :items="deployMenuItems" @select="handleDeployAction($event, row)">
              <template #trigger>
                <CfButton variant="ghost" size="sm" icon-only title="Row actions">⋯</CfButton>
              </template>
            </CfDropdownMenu>
          </template>
        </CfTable>
      </CfCard>
    </div>

    <!-- ================================================================ -->
    <!-- TAB 4: LOGS (Observability using CfLogViewer from cloudflare-ui)  -->
    <!-- ================================================================ -->
    <div v-else-if="activeTab === 'logs'" class="cf-ad-stack">
      <div style="display: flex; justify-content: flex-end; gap: 0.5rem;">
        <CfButton variant="secondary" size="sm" icon="refresh" @click="fetchLogs">Reload</CfButton>
      </div>
      <CfLogViewer
        :title="`${app.name || appId} — Live Stream`"
        :logs="formattedLogs"
        @clear="logs = []"
      />
    </div>

    <!-- ================================================================ -->
    <!-- TAB 5: DOMAINS                                                    -->
    <!-- ================================================================ -->
    <div v-else-if="activeTab === 'domains'" class="cf-ad-stack">
      <CfCard title="App URL & Reverse Proxy">
        <div class="cf-ad-kv-list">
          <div class="cf-ad-kv-row">
            <span class="cf-ad-kv-label">Production Domain</span>
            <span class="cf-ad-kv-val">{{ app.domain || 'Not assigned' }}</span>
          </div>
          <div class="cf-ad-kv-row">
            <span class="cf-ad-kv-label">Local Upstream</span>
            <code>http://127.0.0.1:{{ app.active_port || 8001 }}</code>
          </div>
        </div>
      </CfCard>

      <CfCard title="Custom Domains and Routes" body-class="cf-ad-table-body">
        <template #actions>
          <router-link to="/ssl" style="text-decoration: none;">
            <CfButton variant="primary" size="sm" icon="plus">Manage SSL & Domains</CfButton>
          </router-link>
        </template>
        <CfTable :columns="domainColumns" :data="domainRows" row-key="name">
          <template #col-type="{ val }">
            <CfBadge variant="info" size="sm">{{ val }}</CfBadge>
          </template>
        </CfTable>
      </CfCard>
    </div>

    <!-- ================================================================ -->
    <!-- TAB 6: SETTINGS (Exact Cloudflare 2-Col + Sticky Anchor Sidebar)  -->
    <!-- ================================================================ -->
    <div v-else class="cf-ad-settings-layout">
      <!-- Left Scrollable Sections -->
      <div class="cf-ad-settings-main">
        <section id="sec-vars">
          <CfCard title="Runtime variables and secrets" subtitle="Environment variables injected into the app process">
            <template #actions>
              <CfButton variant="secondary" size="sm" icon="plus" @click="addEnvVar">Add variable</CfButton>
            </template>
            <CfTable :columns="varColumns" :data="envVars" row-key="name">
              <template #col-type="{ val }">
                <CfBadge :variant="val === 'Secret' ? 'warning' : 'neutral'" size="sm">{{ val }}</CfBadge>
              </template>
              <template #col-value="{ val }">
                <code>{{ val }}</code>
              </template>
              <template #actions="{ row }">
                <CfButton variant="danger-outline" size="sm" icon="trash" icon-only @click="removeEnvVar(row.name)" />
              </template>
            </CfTable>
          </CfCard>
        </section>

        <section id="sec-bindings">
          <CfCard title="Bindings" subtitle="Connected Databases and S3 Object Storage buckets">
            <CfTable :columns="bindingColumns" :data="bindings" row-key="name">
              <template #col-type="{ val }">
                <CfBadge variant="info" size="sm">{{ val }}</CfBadge>
              </template>
            </CfTable>
          </CfCard>
        </section>

        <section id="sec-runtime">
          <CfCard title="Runtime & Resource Limits" subtitle="cgroups isolation and health check configuration">
            <div class="cf-ad-kv-list">
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">RAM Limit</span>
                <span class="cf-ad-kv-val">{{ app.ram_limit_mb || 256 }} MB</span>
              </div>
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">CPU Quota</span>
                <span class="cf-ad-kv-val">{{ app.cpu_quota_pct || 100 }}%</span>
              </div>
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">Health Check Path</span>
                <code>{{ app.health_check_path || '/healthz' }}</code>
              </div>
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">Auto-Restart on Crash</span>
                <CfSwitch :model-value="!!app.auto_restart" disabled />
              </div>
            </div>
          </CfCard>
        </section>

        <section id="sec-builds">
          <CfCard title="Build & Entrypoint" subtitle="Executable command and Git deployment settings">
            <div class="cf-ad-kv-list">
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">Binary Entrypoint</span>
                <code>{{ app.binary_entrypoint || './server' }}</code>
              </div>
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">Exec Args</span>
                <code>{{ app.exec_args || '(none)' }}</code>
              </div>
              <div class="cf-ad-kv-row">
                <span class="cf-ad-kv-label">Linux User</span>
                <code>{{ app.linux_user || 'orbit' }}</code>
              </div>
            </div>
          </CfCard>
        </section>

        <section id="sec-general">
          <CfCard title="General" subtitle="Application identity on this VPS node">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
              <CfInput :model-value="app.name || appId" label="Name" readonly />
              <CfInput :model-value="app.app_id || appId" label="App ID" readonly />
            </div>
          </CfCard>
        </section>

        <section id="sec-danger">
          <CfCard title="Danger zone" subtitle="Irreversible actions for this application">
            <div style="display: flex; align-items: center; justify-content: space-between;">
              <div>
                <div style="font-weight: 600; font-size: 13.5px; color: var(--text-color-kumo-strong);">Delete App</div>
                <div style="font-size: 12.5px; color: var(--text-color-kumo-subtle);">
                  Stop the running process and remove its configuration from this node.
                </div>
              </div>
              <CfButton variant="danger" icon="trash" @click="alert('Delete confirmation')">Delete</CfButton>
            </div>
          </CfCard>
        </section>
      </div>

      <!-- Right Sticky Anchor Sidebar -->
      <aside class="cf-ad-settings-rail">
        <a href="#sec-vars" class="cf-rail-link" @click.prevent="scrollTo('sec-vars')">Variables and secrets</a>
        <a href="#sec-bindings" class="cf-rail-link" @click.prevent="scrollTo('sec-bindings')">Bindings</a>
        <a href="#sec-runtime" class="cf-rail-link" @click.prevent="scrollTo('sec-runtime')">Runtime</a>
        <a href="#sec-builds" class="cf-rail-link" @click.prevent="scrollTo('sec-builds')">Builds</a>
        <a href="#sec-general" class="cf-rail-link" @click.prevent="scrollTo('sec-general')">General</a>
        <a href="#sec-danger" class="cf-rail-link is-danger" @click.prevent="scrollTo('sec-danger')">Danger zone</a>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const appId = ref(route.params.id)
const activeTab = ref('overview')
const app = ref({})
const stats = ref(null)
const logs = ref([])
const selectedVersion = ref('All deployed versions')
const metricsRange = ref('Last 24 hours')
let sse = null

const bindings = ref([
  { type: 'Database', name: 'DB', target: 'main_d1_db' },
  { type: 'S3 Bucket', name: 'FILES', target: 'assets' }
])

const envVars = ref([
  { type: 'Variable', name: 'PORT', value: '8001' },
  { type: 'Secret', name: 'ORBIT_APP_KEY', value: 'Value encrypted' }
])

const deployStatusTab = ref('All')
const rawDeployments = ref([])

const deployColumns = [
  { key: 'version', label: 'Version ID', width: '165px', sortable: true },
  { key: 'traffic', label: 'Traffic', width: '140px', sortable: true },
  { key: 'cpu', label: 'CPU Trend', width: '115px' },
  { key: 'status', label: 'Status', width: '105px', sortable: true },
  { key: 'source', label: 'Source', width: '110px' },
  { key: 'time', label: 'Deployed', width: '110px' }
]

const deployMenuItems = [
  { label: 'View Logs', value: 'logs' },
  { label: 'Copy Version ID', value: 'copy' },
  { label: 'Redeploy', value: 'redeploy' }
]

const domainColumns = [
  { key: 'name', label: 'Name' },
  { key: 'type', label: 'Type', width: '140px' },
  { key: 'target', label: 'Upstream Target', width: '180px' }
]

const varColumns = [
  { key: 'type', label: 'Type', width: '100px' },
  { key: 'name', label: 'Name', width: '180px' },
  { key: 'value', label: 'Value' }
]

const bindingColumns = [
  { key: 'type', label: 'Type', width: '120px' },
  { key: 'name', label: 'Name', width: '140px' },
  { key: 'target', label: 'Resource' }
]

const formattedLogs = computed(() =>
  logs.value.map(line => ({
    time: new Date().toLocaleTimeString([], { hour12: false }),
    level: line.toLowerCase().includes('error') ? 'error' : 'info',
    message: line
  }))
)

const deployments = computed(() => {
  if (rawDeployments.value.length > 0) {
    return rawDeployments.value.map((d, idx) => ({
      id: d.deployment_id || idx + 1,
      version: d.version_hash || 'initial',
      traffic: idx === 0 && app.value.status === 'running' ? 100 : 0,
      cpu: idx === 0 ? [12, 18, 16, 22, 19, 24, 21] : [8, 10, 9, 11, 10, 9, 8],
      status: idx === 0 && app.value.status === 'running' ? 'Active' : (d.status || 'Previous'),
      source: d.deployed_by || 'Orbit CLI',
      time: formatRelative(d.created_at)
    }))
  }
  return [
    {
      id: 1,
      version: app.value.current_version || 'initial-seed',
      traffic: app.value.status === 'running' ? 100 : 0,
      cpu: [14, 19, 17, 23, 20, 26, 22],
      status: app.value.status === 'running' ? 'Active' : 'Stopped',
      source: 'Orbit CLI',
      time: formatRelative(app.value.updated_at)
    }
  ]
})

const filteredDeployments = computed(() => {
  if (deployStatusTab.value === 'All') return deployments.value
  if (deployStatusTab.value === 'Active') return deployments.value.filter(d => d.status === 'Active' || d.traffic > 0)
  return deployments.value.filter(d => d.status !== 'Active' && d.traffic === 0)
})

async function fetchDeployments() {
  try {
    const res = await fetch(`/api/v1/apps/${appId.value}/deployments`, { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      rawDeployments.value = data.deployments || []
    }
  } catch (e) {}
}

function handleDeployAction(action, row) {
  if (action === 'logs') activeTab.value = 'logs'
  else if (action === 'copy') copyText(row.version)
  else if (action === 'redeploy') fetchApp()
}

const domainRows = computed(() => {
  const rows = []
  if (app.value.domain) {
    rows.push({ name: app.value.domain, type: 'Custom Domain', target: `127.0.0.1:${app.value.active_port || 8001}` })
  }
  rows.push({ name: `127.0.0.1:${app.value.active_port || 8001}`, type: 'Local Port', target: 'Active' })
  return rows
})

function authHeaders() {
  const token = localStorage.getItem('orbit_token') || ''
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function fetchApp() {
  try {
    const res = await fetch(`/api/v1/apps/${appId.value}`, { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      app.value = data.app || data
      stats.value = data.stats || null
    }
  } catch (e) {}
}

async function fetchLogs() {
  try {
    const res = await fetch(`/api/v1/apps/${appId.value}/logs?tail=100`, { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      logs.value = Array.isArray(data.logs) ? data.logs : (typeof data === 'string' ? data.split('\n') : [])
    }
  } catch (e) {}
}

function startLogStream() {
  if (sse) sse.close()
  const token = localStorage.getItem('orbit_token') || ''
  const url = token ? `/api/v1/apps/${appId.value}/logs/stream?token=${token}` : `/api/v1/apps/${appId.value}/logs/stream`
  sse = new EventSource(url)
  sse.onmessage = (ev) => {
    if (ev.data) logs.value.push(ev.data)
  }
}

function addEnvVar() {
  const name = prompt('Variable name (e.g. NODE_ENV):')
  if (!name) return
  const val = prompt('Variable value:') || ''
  envVars.value.push({ type: 'Variable', name: name.trim(), value: val })
}

function removeEnvVar(name) {
  envVars.value = envVars.value.filter(v => v.name !== name)
}

function scrollTo(id) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function copyText(text) {
  navigator.clipboard?.writeText(text)
}

function formatRelative(ts) {
  if (!ts) return 'Just now'
  const diff = Math.floor((Date.now() - new Date(ts).getTime()) / 1000)
  if (diff < 60) return 'Just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

watch(activeTab, (tab) => {
  if (tab === 'logs') {
    fetchLogs()
    startLogStream()
  } else if (sse) {
    sse.close()
    sse = null
  }
})

onMounted(() => {
  fetchApp()
  fetchDeployments()
})
onBeforeUnmount(() => {
  if (sse) sse.close()
})
</script>

<style scoped>
.cf-app-detail-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.cf-ad-breadcrumb-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cf-ad-breadcrumbs {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  font-size: 13.5px;
}

.cf-ad-bc-link {
  color: var(--text-color-kumo-subtle);
  text-decoration: none;
}

.cf-ad-bc-link:hover {
  color: var(--text-color-kumo-strong);
}

.cf-ad-bc-sep {
  color: var(--text-color-kumo-inactive);
}

.cf-ad-bc-current {
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-ad-env-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 2px 8px;
  border-radius: 6px;
  border: 1px solid var(--color-kumo-line);
  background: var(--color-kumo-base);
  font-size: 12px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-ad-nav-action-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.cf-ad-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.cf-ad-stack {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.cf-ad-summary-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.cf-ad-summary-left,
.cf-ad-summary-right {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.cf-ad-domain-link {
  font-size: 13.5px;
  font-weight: 500;
  color: #0051c3;
  text-decoration: underline;
}

.cf-ad-repo-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 2px 8px;
  border-radius: 9999px;
  border: 1px solid var(--color-kumo-line);
  background: var(--color-kumo-tint);
  font-size: 12px;
  color: var(--text-color-kumo-default);
}

.cf-ad-muted-meta {
  font-size: 12.5px;
  color: var(--text-color-kumo-subtle);
}

/* Architecture Diagram */
.cf-arch-diagram {
  display: grid;
  grid-template-columns: 1fr 48px 1.15fr 48px 1fr;
  align-items: center;
  gap: 0;
}

@media (max-width: 900px) {
  .cf-arch-diagram {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  .cf-arch-connector {
    display: none;
  }
}

.cf-arch-col {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.cf-arch-col-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cf-arch-col-label {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--text-color-kumo-subtle);
}

.cf-arch-add-link {
  background: none;
  border: none;
  color: #0051c3;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
}

.cf-arch-pill-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.6rem 0.85rem;
  border-radius: 8px;
  border: 1px solid var(--color-kumo-line);
  background: var(--color-kumo-canvas);
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s ease;
}

.cf-arch-pill-card:hover {
  border-color: #0051c3;
}

.cf-arch-connector {
  display: flex;
  align-items: center;
  justify-content: center;
}

.cf-arch-line {
  width: 100%;
  height: 1px;
  border-top: 2px dashed var(--color-kumo-line);
}

.cf-arch-center-card {
  border: 1px solid var(--color-kumo-line);
  border-radius: 10px;
  padding: 1rem;
  background: var(--color-kumo-base);
  box-shadow: var(--shadow-kumo-card);
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.cf-arch-center-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cf-arch-status-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #9ca3af;
}

.cf-arch-status-dot.is-live {
  background: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
}

.cf-arch-center-meta {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-arch-center-badges {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  padding-top: 0.5rem;
  border-top: 1px solid var(--color-kumo-hairline);
}

.cf-arch-obs-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12.5px;
  cursor: pointer;
}

.cf-ad-metrics-3col {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

@media (max-width: 768px) {
  .cf-ad-metrics-3col {
    grid-template-columns: 1fr;
  }
}

:deep(.cf-ad-table-body) {
  padding: 0;
}

.cf-ad-code-pill {
  font-family: var(--font-mono);
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--color-kumo-recessed);
}

.cf-ad-bottom-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
}

@media (max-width: 768px) {
  .cf-ad-bottom-2col {
    grid-template-columns: 1fr;
  }
}

.cf-ad-kv-list {
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}

.cf-ad-kv-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 0.55rem;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 13px;
}

.cf-ad-kv-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.cf-ad-kv-label {
  font-size: 12.5px;
  color: var(--text-color-kumo-subtle);
  font-weight: 500;
}

.cf-ad-kv-val {
  color: var(--text-color-kumo-strong);
  font-weight: 500;
}

.cf-ad-next-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.55rem 0.25rem;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 13px;
  cursor: pointer;
}

.cf-ad-next-row:hover {
  color: #0051c3;
}

.cf-ad-controls-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.cf-ad-active-deploy-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
  align-items: center;
}

.app-log-terminal {
  background: #0b0f17;
  color: #e2e8f0;
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  padding: 1rem;
  border-radius: 6px;
  min-height: 260px;
  max-height: 460px;
  overflow-y: auto;
}

.app-log-line {
  white-space: pre-wrap;
  word-break: break-all;
}

/* Settings 2-Column Sticky Anchor Layout */
.cf-ad-settings-layout {
  display: grid;
  grid-template-columns: 1fr 220px;
  gap: 1.5rem;
  align-items: start;
}

@media (max-width: 900px) {
  .cf-ad-settings-layout {
    grid-template-columns: 1fr;
  }
  .cf-ad-settings-rail {
    display: none;
  }
}

.cf-ad-settings-main {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  min-width: 0;
}

.cf-ad-settings-rail {
  position: sticky;
  top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding-left: 1rem;
  border-left: 1px solid var(--color-kumo-line);
}

.cf-rail-link {
  font-size: 13px;
  color: var(--text-color-kumo-subtle);
  text-decoration: none;
  padding: 0.35rem 0.5rem;
  border-radius: 4px;
}

.cf-rail-link:hover {
  background: var(--color-kumo-fill);
  color: var(--text-color-kumo-strong);
}

.cf-rail-link.is-danger {
  color: #dc2626;
}
</style>
