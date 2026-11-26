<template>
  <div class="cf-workers-page-real">
    <div class="cf-workers-container">
      <!-- 1. PAGE HEADER ROW -->
      <div class="cf-wp-header-row">
        <div class="cf-wp-title-block">
          <div class="cf-wp-h1-row">
            <h1 class="cf-wp-h1">Apps & Services</h1>
            <a 
              href="https://developers.cloudflare.com/workers/" 
              target="_blank" 
              rel="noopener noreferrer" 
              class="cf-wp-view-docs-pill"
            >
              <span>Docs</span>
              <svg width="10" height="10" viewBox="0 0 256 256" fill="currentColor">
                <path d="M200,64V168a8,8,0,0,1-16,0V83.31L69.66,197.66a8,8,0,0,1-11.32-11.32L172.69,72H88a8,8,0,0,1,0-16H192A8,8,0,0,1,200,64Z"></path>
              </svg>
            </a>
          </div>
          <p class="cf-wp-subtitle">
            Deploy, restart, and monitor containers, background services, and web apps on this VPS node.
          </p>
        </div>

        <button class="cf-wp-create-btn" @click="showCreateModal = true">
          <span>Deploy service</span>
        </button>
      </div>

      <!-- 2. NODE RESOURCE USAGE PANEL -->
      <div class="cf-wp-usage-panel">
        <div class="cf-wp-usage-header">
          <span class="cf-wp-usage-title">
            <span>Node Resources</span>
            <CfIcon name="cis-circle-info" size="13" />
          </span>
          <div class="cf-wp-usage-actions">
            <button class="cf-wp-time-pill" type="button">
              <span>Live</span>
              <CfIcon name="cis-chevron-down-medium" size="11" />
            </button>
          </div>
        </div>

        <div class="cf-wp-usage-body">
          <!-- Row 1: CPU Allocation -->
          <div class="cf-wp-usage-row">
            <div class="cf-wp-usage-left">
              <div class="cf-wp-donut-mini"></div>
              <span>CPU cores active</span>
            </div>
            <span class="cf-wp-usage-right">2.0 / 8 vCPU</span>
          </div>

          <!-- Row 2: Memory Allocation -->
          <div class="cf-wp-usage-row">
            <div class="cf-wp-usage-left">
              <div class="cf-wp-donut-mini" style="border-top-color: #f6821f;"></div>
              <span>Memory allocated</span>
            </div>
            <span class="cf-wp-usage-right">6.4 / 16 GB</span>
          </div>
        </div>
      </div>

      <!-- 3. UNIFIED JOINED TOOLBAR -->
      <div class="cf-wp-toolbar-joined">
        <div class="cf-wp-tb-search">
          <CfIcon name="cis-magnifying-glass" size="14" style="color: #71717a;" />
          <input
            v-model="searchQuery"
            type="text"
            class="cf-wp-tb-input"
            placeholder="Search applications"
          />
        </div>

        <div class="cf-wp-tb-divider"></div>

        <select v-model="filterType" class="cf-wp-tb-select">
          <option value="all">Show all</option>
          <option value="worker">Workers</option>
          <option value="pages">Pages</option>
        </select>

        <div class="cf-wp-tb-divider"></div>

        <select v-model="sortBy" class="cf-wp-tb-select">
          <option value="modified">Sort by: Last modified</option>
          <option value="name">Sort by: Name</option>
          <option value="requests">Sort by: Requests</option>
        </select>

        <div class="cf-wp-tb-divider"></div>

        <button class="cf-wp-tb-btn" title="Refresh" @click="handleRefresh">
          <CfIcon name="cis-arrow-rotate-clockwise" size="14" />
        </button>
      </div>

      <!-- 4. TWO-TIER APPLICATION CARD (code-collection) -->
      <div class="cf-wp-card-stream">
        <router-link
          v-for="app in filteredApps"
          :key="app.id"
          :to="`/workers/${app.id}`"
          class="cf-wp-app-card-two-tier"
        >
          <!-- Tier 1: Main Header (74px) -->
          <div class="cf-wp-card-tier-main">
            <div class="cf-wp-card-left">
              <div class="cf-wp-logo-badge">
                <CfIcon name="edgeworker" size="20" />
              </div>
              <div class="cf-wp-identity-col">
                <div class="cf-wp-name-row">
                  <span class="cf-wp-app-title">{{ app.name }}</span>
                </div>
                <div class="cf-wp-domain-subrow">
                  <a :href="`https://${app.domain}`" target="_blank" class="cf-wp-domain-url" @click.stop>
                    {{ app.domain }}
                  </a>
                  <span class="cf-wp-badge-extra">+ 1 other route</span>
                </div>
              </div>
            </div>

            <div class="cf-wp-card-right">
              <div class="cf-wp-repo-pill">
                <svg width="13" height="13" viewBox="0 0 256 256" fill="currentColor">
                  <path d="M128,24A104,104,0,0,0,24,128c0,46,29.8,85,71.1,98.8,5.2,1,7.1-2.3,7.1-5v-18c-28.9,6.3-35-13.9-35-13.9-4.7-12-11.5-15.2-11.5-15.2-9.4-6.4.7-6.3.7-6.3,10.4.7,15.9,10.7,15.9,10.7,9.3,15.9,24.3,11.3,30.2,8.6a22.2,22.2,0,0,1,6.6-13.9c-23.1-2.6-47.4-11.5-47.4-51.4a40.2,40.2,0,0,1,10.7-27.9,37.3,37.3,0,0,1,1-27.5s8.7-2.8,28.6,10.7a98.8,98.8,0,0,1,52,0c19.9-13.5,28.6-10.7,28.6-10.7a37.3,37.3,0,0,1,1,27.5,40.2,40.2,0,0,1,10.7,27.9c0,40-24.4,48.7-47.6,51.3a24.9,24.9,0,0,1,7.1,19.3v28.7c0,2.8,1.9,6.1,7.2,5C202.2,213,232,174,232,128A104,104,0,0,0,128,24Z"/>
                </svg>
                <span>{{ app.repo }}</span>
              </div>
              <span class="cf-wp-time-text">{{ app.lastModified }}</span>
              <div class="cf-mini-icon-btn" title="More options" @click.prevent.stop="alertMenu(app.name)">
                <svg width="16" height="16" viewBox="0 0 256 256" fill="currentColor">
                  <circle cx="128" cy="128" r="16"></circle>
                  <circle cx="64" cy="128" r="16"></circle>
                  <circle cx="192" cy="128" r="16"></circle>
                </svg>
              </div>
            </div>
          </div>

          <!-- Tier 2: Sub-footer Bar (40px) -->
          <div class="cf-wp-card-tier-footer">
            <div class="cf-wp-commit-msg">
              <svg width="13" height="13" viewBox="0 0 256 256" fill="currentColor">
                <path d="M128,80a24,24,0,1,0-24,24A24,24,0,0,0,128,80Zm-24-8a8,8,0,1,1-8,8A8,8,0,0,1,104,72ZM216,40H40A16,16,0,0,0,24,56V200a16,16,0,0,0,16,16H216a16,16,0,0,0,16-16V56A16,16,0,0,0,216,40Zm0,160H40V56H216V200Z" opacity="0.4"/>
              </svg>
              <span>{{ app.commitMsg }}</span>
            </div>

            <div class="cf-wp-card-telemetry">
              <div class="cf-wp-telem-item">
                <CfIcon name="chart" size="13" style="color: #71717a;" />
                <span>{{ app.requests }} requests</span>
              </div>
              <div class="cf-wp-telem-item">
                <CfIcon name="cis-clock-alert" size="13" style="color: #71717a;" />
                <span>{{ app.errors }} errors</span>
              </div>
            </div>
          </div>
        </router-link>
      </div>

      <!-- 5. PAGINATION ROW -->
      <div class="cf-wp-pagination-row">
        <span>Showing 1-{{ filteredApps.length }} of {{ filteredApps.length }}</span>
        <div class="cf-wp-page-controls">
          <button class="cf-wp-page-btn" disabled title="First page">«</button>
          <button class="cf-wp-page-btn" disabled title="Previous page">‹</button>
          <span class="cf-wp-page-num">1</span>
          <button class="cf-wp-page-btn" disabled title="Next page">›</button>
          <button class="cf-wp-page-btn" disabled title="Last page">»</button>
        </div>
      </div>

      <!-- 6. BOTTOM STREAM CARDS (Account details & Cloudflare Access) -->
      <div class="cf-wp-bottom-cards-stream">
        <!-- Card 1: Account details -->
        <div class="cf-wp-stream-card">
          <h2 class="cf-wp-stream-card-title">Account details</h2>
          <div class="cf-wp-stream-row">
            <span class="cf-wp-stream-label">Account ID</span>
            <div class="cf-wp-stream-val">
              <span>{{ accountId }}</span>
              <button class="cf-mini-icon-btn" title="Copy account ID" @click="copyText(accountId, 'acc')">
                <CfIcon v-if="copiedKey !== 'acc'" name="cis-clipboard" size="13" />
                <CfIcon v-else name="cis-checkmark-1-medium" size="13" style="color: #10b981;" />
              </button>
            </div>
          </div>
          <div class="cf-wp-stream-row">
            <span class="cf-wp-stream-label">Subdomain</span>
            <div class="cf-wp-stream-val">
              <span>{{ subdomain }}</span>
              <button class="cf-mini-icon-btn" title="Copy subdomain" @click="copyText(subdomain, 'subdomain')">
                <CfIcon v-if="copiedKey !== 'subdomain'" name="cis-clipboard" size="13" />
                <CfIcon v-else name="cis-checkmark-1-medium" size="13" style="color: #10b981;" />
              </button>
            </div>
          </div>
        </div>

        <!-- Card 2: Cloudflare Access -->
        <div class="cf-wp-stream-card">
          <h2 class="cf-wp-stream-card-title">Cloudflare Access</h2>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <p style="font-size: 13px; color: #71717a; margin: 0; max-width: 600px;">
              Apply one account-wide Access policy across your Workers.
            </p>
            <label class="cf-wp-switch">
              <input type="checkbox" v-model="accessEnabled" />
              <span class="cf-wp-slider"></span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Application Modal -->
    <div v-if="showCreateModal" class="cf-r2-modal-overlay" @click.self="showCreateModal = false">
      <div class="cf-r2-modal-card" style="max-width: 640px;">
        <div class="cf-r2-modal-header">
          <h2 class="cf-r2-modal-title">Create an application</h2>
          <button class="cf-mini-icon-btn" @click="showCreateModal = false">✕</button>
        </div>
        <div class="cf-r2-modal-body" style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
          <div class="cf-app-card" style="cursor: pointer;" @click="createAppType('worker')">
            <div class="cf-app-card-title-group">
              <div class="cf-app-icon-badge">
                <CfIcon name="edgeworker" size="18" />
              </div>
              <span class="cf-app-name">Workers</span>
            </div>
            <p style="font-size: 12px; color: #71717a; margin: 0;">Deploy serverless code with instant V8 isolates.</p>
          </div>

          <div class="cf-app-card" style="cursor: pointer;" @click="createAppType('pages')">
            <div class="cf-app-card-title-group">
              <div class="cf-app-icon-badge pages-badge">
                <CfIcon name="workers-pages" size="18" />
              </div>
              <span class="cf-app-name">Pages</span>
            </div>
            <p style="font-size: 12px; color: #71717a; margin: 0;">Connect a GitHub repository for fullstack Jamstack apps.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')
const filterType = ref('all')
const sortBy = ref('modified')
const showCreateModal = ref(false)
const copiedKey = ref('')
const accessEnabled = ref(false)

const subdomain = ref('northlabs.workers.dev')
const accountId = ref('330b08634cd698ef2ff8587d66ac4b2c')

const apps = ref([
  {
    id: 'code-collection',
    name: 'code-collection',
    type: 'Worker',
    domain: 'codecollection.northlab.dev',
    repo: 'hridoy/code-collection',
    branch: 'main',
    commitMsg: 'feat(ui): inline author with title and clean typography',
    author: 'hridoy',
    lastModified: '9h ago',
    requests: '13',
    errors: '13'
  }
])

const filteredApps = computed(() => {
  let list = apps.value
  if (filterType.value === 'worker') {
    list = list.filter(a => a.type === 'Worker')
  } else if (filterType.value === 'pages') {
    list = list.filter(a => a.type === 'Pages')
  }
  const q = searchQuery.value.toLowerCase().trim()
  if (q) {
    list = list.filter(a => a.name.toLowerCase().includes(q) || a.domain.toLowerCase().includes(q))
  }
  return list
})

const handleRefresh = () => {
  alert('Refreshed application telemetry from Cloudflare edge.')
}

const copyText = async (text, key) => {
  try {
    await navigator.clipboard.writeText(text)
    copiedKey.value = key
    setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = ''
    }, 2000)
  } catch (err) {
    console.warn(err)
  }
}

const alertMenu = (name) => {
  alert(`Options for "${name}":\n- Manage worker\n- Edit code\n- Visit site\n- View deployments\n- Settings`)
}

const createAppType = (type) => {
  showCreateModal.value = false
  alert(`Creating new ${type === 'worker' ? 'Cloudflare Worker' : 'Cloudflare Pages'} project...`)
}
</script>
