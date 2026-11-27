<template>
  <div class="cf-workers-page-real">
    <div class="cf-workers-container">
      <!-- 1. PAGE HEADER ROW -->
      <div class="cf-wp-header-row">
        <div class="cf-wp-title-block">
          <div class="cf-wp-h1-row">
            <h1 class="cf-wp-h1">Apps</h1>
            <span class="cf-wp-view-docs-pill">
              <span>{{ apps.length }} {{ apps.length === 1 ? 'app' : 'apps' }}</span>
            </span>
          </div>
          <p class="cf-wp-subtitle">
            Deploy, manage, and monitor applications running on this VPS node.
          </p>
        </div>

        <button class="cf-wp-create-btn" @click="$router.push('/apps/create')">
          <span>Create application</span>
        </button>
      </div>

      <!-- 2. NODE RESOURCES PANEL (ProductUsagePanel) -->
      <div class="cf-wp-usage-panel">
        <div class="cf-wp-usage-header">
          <span class="cf-wp-usage-title">
            <span>Node Resources</span>
            <CfIcon name="cis-circle-info" size="13" />
          </span>
          <div class="cf-wp-usage-actions">
            <button class="cf-wp-time-pill" type="button" @click="fetchApps">
              <span>Live</span>
              <CfIcon name="refresh" size="11" />
            </button>
          </div>
        </div>

        <div class="cf-wp-usage-body">
          <div class="cf-wp-usage-row">
            <div class="cf-wp-usage-left">
              <div class="cf-wp-donut-mini"></div>
              <span>Active apps</span>
            </div>
            <span class="cf-wp-usage-right">{{ runningCount }} / {{ apps.length }} running</span>
          </div>

          <div class="cf-wp-usage-row">
            <div class="cf-wp-usage-left">
              <div class="cf-wp-donut-mini" style="border-top-color: #f6821f;"></div>
              <span>Allocated RAM</span>
            </div>
            <span class="cf-wp-usage-right">{{ totalRamMB }} MB</span>
          </div>
        </div>
      </div>

      <!-- 3. UNIFIED JOINED TOOLBAR (Cloudflare Kumo 1:1) -->
      <div class="cf-wp-toolbar-joined" role="toolbar">
        <label class="cf-wp-tb-search">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #71717a; flex-shrink: 0;">
            <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            class="cf-wp-tb-input"
            placeholder="Search apps"
          />
        </label>

        <div class="cf-wp-tb-divider"></div>

        <CfDropdownMenu :items="filterMenuItems" align="right" @select="filterStatus = $event.value">
          <template #trigger>
            <button type="button" class="cf-wp-tb-combo" role="combobox">
              <span>{{ filterLabel }}</span>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#71717a" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                <path d="m7 15 5 5 5-5"/><path d="m7 9 5-5 5 5"/>
              </svg>
            </button>
          </template>
        </CfDropdownMenu>

        <div class="cf-wp-tb-divider"></div>

        <CfDropdownMenu :items="sortMenuItems" align="right" @select="sortBy = $event.value">
          <template #trigger>
            <button type="button" class="cf-wp-tb-combo" role="combobox">
              <span>{{ sortLabel }}</span>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#71717a" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                <path d="m7 15 5 5 5-5"/><path d="m7 9 5-5 5 5"/>
              </svg>
            </button>
          </template>
        </CfDropdownMenu>

        <div class="cf-wp-tb-divider"></div>

        <button class="cf-wp-tb-btn" title="Refresh" @click="fetchApps">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/>
          </svg>
        </button>
      </div>

      <!-- 4. TWO-TIER APPLICATION CARDS -->
      <div v-if="loading" style="padding: 2rem; display: flex; justify-content: center;">
        <CfLoader size="28" />
      </div>

      <CfEmptyState
        v-else-if="!filteredApps.length"
        title="No apps found"
        description="Deploy your first application or connect a Git repository to run it on this server."
        action-text="Deploy App"
        action-icon="plus"
        @action="showCreateModal = true"
      />

      <div v-else class="cf-wp-card-stream">
        <router-link
          v-for="app in filteredApps"
          :key="app.app_id"
          :to="`/apps/${app.app_id}`"
          class="cf-wp-app-card-two-tier"
        >
          <!-- Tier 1: Main Header -->
          <div class="cf-wp-card-tier-main">
            <div class="cf-wp-card-left">
              <div class="cf-wp-logo-bare">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m8 3 4 8 5-5 5 15H2L8 3z"/>
                </svg>
              </div>
              <div class="cf-wp-identity-col">
                <div class="cf-wp-name-row">
                  <span class="cf-wp-app-title">{{ app.name || app.app_id }}</span>
                </div>
                <div class="cf-wp-domain-subrow">
                  <a
                    v-if="app.domain"
                    :href="`http://${app.domain}`"
                    target="_blank"
                    class="cf-wp-domain-url"
                    @click.stop
                  >
                    {{ app.domain }}
                  </a>
                  <span v-else class="cf-wp-domain-url" style="color: #71717a;">
                    127.0.0.1:{{ app.active_port }}
                  </span>
                  <span class="cf-wp-badge-extra">Port :{{ app.active_port }}</span>
                </div>
              </div>
            </div>

            <div class="cf-wp-card-right" @click.prevent.stop>
              <div class="cf-wp-repo-pill">
                <CfIcon name="server" size="12" />
                <span>{{ app.app_id }}</span>
              </div>
              <span class="cf-wp-time-text">{{ formatRelative(app.updated_at) }}</span>
              <CfDropdownMenu :items="appMenuItems" align="right" @select="handleAppMenu($event, app)">
                <template #trigger>
                  <button type="button" class="cf-wp-card-more-btn" title="More options">⋯</button>
                </template>
              </CfDropdownMenu>
            </div>
          </div>

          <!-- Tier 2: Sub-footer Bar (40px) -->
          <div class="cf-wp-card-tier-footer">
            <div class="cf-wp-commit-msg">
              <CfIcon name="cis-page-text" size="13" style="opacity: 0.6;" />
              <span>{{ app.binary_entrypoint || './server' }}</span>
              <span v-if="app.current_version" style="opacity: 0.7;">({{ app.current_version }})</span>
            </div>

            <div class="cf-wp-card-telemetry">
              <div class="cf-wp-telem-item">
                <span>{{ app.ram_limit_mb || 256 }} MB RAM</span>
              </div>
              <div class="cf-wp-telem-item">
                <span>{{ app.cpu_quota_pct || 100 }}% CPU</span>
              </div>
              <div class="cf-wp-telem-item">
                <span>{{ app.auto_restart ? 'Auto-restart' : 'Manual' }}</span>
              </div>
            </div>
          </div>
        </router-link>
      </div>

      <!-- 5. PAGINATION ROW -->
      <div v-if="filteredApps.length" class="cf-wp-pagination-row">
        <span>Showing 1-{{ filteredApps.length }} of {{ filteredApps.length }}</span>
        <div class="cf-wp-page-controls">
          <button class="cf-wp-page-btn" disabled title="First page">«</button>
          <button class="cf-wp-page-btn" disabled title="Previous page">‹</button>
          <span class="cf-wp-page-num">1</span>
          <button class="cf-wp-page-btn" disabled title="Next page">›</button>
          <button class="cf-wp-page-btn" disabled title="Last page">»</button>
        </div>
      </div>

      <!-- 6. BOTTOM STREAM CARDS -->
      <div class="cf-wp-bottom-cards-stream">
        <div class="cf-wp-stream-card">
          <h2 class="cf-wp-stream-card-title">Node details</h2>
          <div class="cf-wp-stream-row">
            <span class="cf-wp-stream-label">Node ID</span>
            <div class="cf-wp-stream-val">
              <span>{{ nodeInfo.node_id || 'local-node' }}</span>
              <button class="cf-mini-icon-btn" title="Copy Node ID" @click="copyText(nodeInfo.node_id || 'local-node', 'node')">
                <CfIcon v-if="copiedKey !== 'node'" name="cis-clipboard" size="13" />
                <CfIcon v-else name="cis-checkmark-1-medium" size="13" style="color: #10b981;" />
              </button>
            </div>
          </div>
          <div class="cf-wp-stream-row">
            <span class="cf-wp-stream-label">Hostname</span>
            <div class="cf-wp-stream-val">
              <span>{{ nodeInfo.hostname || 'localhost' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Deploy App Modal -->
    <CfModal v-model="showCreateModal" title="Deploy an App" size="md">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <CfTabs
          v-model="deployMode"
          :tabs="[
            { label: 'Create App', value: 'create' },
            { label: 'Deploy from Git', value: 'git' }
          ]"
        />

        <template v-if="deployMode === 'create'">
          <CfInput v-model="createForm.app_id" label="App ID" placeholder="e.g. my-api" required />
          <CfInput v-model="createForm.name" label="Display Name" placeholder="e.g. My API Service" required />
          <CfInput v-model="createForm.domain" label="Domain (optional)" placeholder="e.g. api.northlab.dev" />
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem;">
            <CfInput v-model="createForm.ram_limit_mb" label="RAM Limit (MB)" type="number" />
            <CfInput v-model="createForm.cpu_quota_pct" label="CPU Quota (%)" type="number" />
          </div>
        </template>

        <template v-else>
          <CfInput v-model="gitForm.repo_url" label="Git Repository URL" placeholder="https://github.com/user/repo.git" required />
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem;">
            <CfInput v-model="gitForm.branch" label="Branch" placeholder="main" />
            <CfInput v-model="gitForm.app_id" label="App ID" placeholder="my-git-app" required />
          </div>
          <CfInput v-model="gitForm.domain" label="Domain (optional)" placeholder="app.northlab.dev" />
        </template>

        <CfAlert v-if="createError" variant="danger" title="Error">{{ createError }}</CfAlert>
      </div>

      <template #footer>
        <CfButton variant="outline" @click="showCreateModal = false">Cancel</CfButton>
        <CfButton variant="primary" :loading="creating" @click="submitCreate">
          {{ deployMode === 'git' ? 'Clone & Deploy' : 'Create App' }}
        </CfButton>
      </template>
    </CfModal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const apps = ref([])
const nodeInfo = ref({})
const loading = ref(true)
const searchQuery = ref('')
const filterStatus = ref('all')
const sortBy = ref('modified')

const filterMenuItems = [
  { label: 'Show all', value: 'all' },
  { label: 'Running', value: 'running' },
  { label: 'Stopped', value: 'stopped' }
]
const sortMenuItems = [
  { label: 'Last modified', value: 'modified' },
  { label: 'Name', value: 'name' },
  { label: 'Port', value: 'port' }
]
const appMenuItems = [
  { label: 'Open Console', value: 'open' },
  { label: 'Copy App ID', value: 'copy' }
]

const filterLabel = computed(() => filterMenuItems.find(i => i.value === filterStatus.value)?.label || 'Show all')
const sortLabel = computed(() => sortMenuItems.find(i => i.value === sortBy.value)?.label || 'Last modified')

function handleAppMenu(item, app) {
  const act = item?.value || item
  if (act === 'open') router.push(`/apps/${app.app_id}`)
  else if (act === 'copy') copyText(app.app_id, app.app_id)
}
const showCreateModal = ref(false)
const deployMode = ref('create')
const creating = ref(false)
const createError = ref('')
const copiedKey = ref('')

const createForm = reactive({
  app_id: '',
  name: '',
  domain: '',
  ram_limit_mb: '256',
  cpu_quota_pct: '100'
})

const gitForm = reactive({
  repo_url: '',
  branch: 'main',
  app_id: '',
  domain: ''
})

function authHeaders() {
  const token = localStorage.getItem('orbit_token') || ''
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function fetchApps() {
  loading.value = true
  try {
    const [appsRes, nodeRes] = await Promise.all([
      fetch('/api/v1/apps', { headers: authHeaders() }),
      fetch('/api/v1/node', { headers: authHeaders() })
    ])
    if (appsRes.ok) {
      const data = await appsRes.json()
      apps.value = Array.isArray(data) ? data : []
    }
    if (nodeRes.ok) {
      nodeInfo.value = await nodeRes.json()
    }
  } catch (e) {
    console.warn('Could not load apps from Orbit backend', e)
  } finally {
    loading.value = false
  }
}

const runningCount = computed(() => apps.value.filter(a => a.status === 'running').length)
const totalRamMB = computed(() => apps.value.reduce((sum, a) => sum + (Number(a.ram_limit_mb) || 0), 0))

const filteredApps = computed(() => {
  let list = [...apps.value]
  if (filterStatus.value === 'running') {
    list = list.filter(a => a.status === 'running')
  } else if (filterStatus.value === 'stopped') {
    list = list.filter(a => a.status !== 'running')
  }
  const q = searchQuery.value.toLowerCase().trim()
  if (q) {
    list = list.filter(a =>
      (a.name || '').toLowerCase().includes(q) ||
      (a.app_id || '').toLowerCase().includes(q) ||
      (a.domain || '').toLowerCase().includes(q)
    )
  }
  if (sortBy.value === 'name') {
    list.sort((a, b) => (a.name || a.app_id).localeCompare(b.name || b.app_id))
  } else if (sortBy.value === 'port') {
    list.sort((a, b) => (a.active_port || 0) - (b.active_port || 0))
  }
  return list
})

async function submitCreate() {
  createError.value = ''
  creating.value = true
  try {
    if (deployMode.value === 'create') {
      const res = await fetch('/api/v1/apps', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', ...authHeaders() },
        body: JSON.stringify({
          app_id: createForm.app_id.trim(),
          name: createForm.name.trim(),
          domain: createForm.domain.trim() || undefined,
          ram_limit_mb: Number(createForm.ram_limit_mb) || 256,
          cpu_quota_pct: Number(createForm.cpu_quota_pct) || 100
        })
      })
      const body = await res.json().catch(() => ({}))
      if (!res.ok) throw new Error(body.error || `HTTP ${res.status}`)
    } else {
      const res = await fetch('/api/v1/apps/install-git', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', ...authHeaders() },
        body: JSON.stringify({
          repo_url: gitForm.repo_url.trim(),
          branch: gitForm.branch.trim() || 'main',
          app_id: gitForm.app_id.trim(),
          domain: gitForm.domain.trim() || undefined
        })
      })
      const body = await res.json().catch(() => ({}))
      if (!res.ok) throw new Error(body.error || `HTTP ${res.status}`)
    }
    showCreateModal.value = false
    createForm.app_id = ''
    createForm.name = ''
    createForm.domain = ''
    gitForm.repo_url = ''
    gitForm.app_id = ''
    await fetchApps()
  } catch (err) {
    createError.value = err.message
  } finally {
    creating.value = false
  }
}

function formatRelative(ts) {
  if (!ts) return 'Just now'
  const diff = Math.floor((Date.now() - new Date(ts).getTime()) / 1000)
  if (diff < 60) return 'Just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

async function copyText(text, key) {
  try {
    await navigator.clipboard.writeText(text)
    copiedKey.value = key
    setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = ''
    }, 2000)
  } catch (e) {}
}

onMounted(fetchApps)
</script>
