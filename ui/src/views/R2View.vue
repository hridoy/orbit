<template>
  <div class="cf-r2-page-wrapper">
    <!-- 1. STICKY SUBHEADER & SEGMENTED PILL TABLIST -->
    <div class="cf-r2-subnav-bar">
      <div class="cf-r2-pill-tablist" role="tablist">
        <button
          role="tab"
          class="cf-r2-pill-tab"
          :class="{ active: currentTab === 'overview' }"
          @click="currentTab = 'overview'"
        >
          Overview
        </button>
        <button
          role="tab"
          class="cf-r2-pill-tab"
          :class="{ active: currentTab === 'metrics' }"
          @click="currentTab = 'metrics'"
        >
          Metrics
        </button>
        <button
          role="tab"
          class="cf-r2-pill-tab"
          :class="{ active: currentTab === 'slurper' }"
          @click="currentTab = 'slurper'"
        >
          Data migration
        </button>
      </div>
    </div>

    <!-- 2. TWO-COLUMN SPLIT LAYOUT (WITH CONTINUOUS VERTICAL DIVIDER) -->
    <div v-if="currentTab === 'overview'" class="cf-r2-split-layout">
      <!-- Left Main Content Column -->
      <main class="cf-r2-main-content">
        <div class="cf-r2-inner-bounded">
          <!-- Hero Title Row -->
          <div class="cf-r2-hero-row">
            <div class="cf-r2-heading-group">
              <div class="cf-r2-title-line">
                <h1 class="cf-r2-page-title">R2 Object Storage</h1>
                <a
                  href="https://developers.cloudflare.com/r2/"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="cf-r2-view-docs-badge"
                >
                  <span>View docs</span>
                  <svg width="11" height="11" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M200,64V168a8,8,0,0,1-16,0V83.31L69.66,197.66a8,8,0,0,1-11.32-11.32L172.69,72H88a8,8,0,0,1,0-16H192A8,8,0,0,1,200,64Z"></path>
                  </svg>
                </a>
              </div>
              <p class="cf-r2-subtitle-text">High-performance storage for files and objects with zero egress charges.</p>
            </div>
            <a href="javascript:void(0)" class="cf-r2-primary-btn" @click="showCreateModal = true">
              <svg width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                <path d="M224,128a8,8,0,0,1-8,8H136v80a8,8,0,0,1-16,0V136H40a8,8,0,0,1,0-16h80V40a8,8,0,0,1,16,0v80h80A8,8,0,0,1,224,128Z"></path>
              </svg>
              <span>Create bucket</span>
            </a>
          </div>

          <!-- Integrated Search Toolbar (722px) -->
          <div class="cf-r2-integrated-toolbar">
            <div class="cf-r2-toolbar-search">
              <CfIcon name="cis-magnifying-glass" size="14" customClass="cf-r2-search-svg" />
              <input
                v-model="searchQuery"
                type="text"
                class="cf-r2-toolbar-input"
                placeholder="Search buckets..."
                aria-label="Search buckets..."
              />
            </div>
            <button class="cf-r2-toolbar-refresh" aria-label="Refresh" @click="handleRefresh">
              <CfIcon name="cis-arrow-rotate-clockwise" size="14" />
            </button>
          </div>

          <!-- Unified Table Card (722px) -->
          <div class="cf-r2-table-wrapper">
            <table class="cf-r2-pure-table">
              <thead>
                <tr>
                  <th class="col-name" @click="toggleSort">
                    <div style="display: flex; align-items: center; gap: 4px;">
                      <span>Buckets</span>
                      <svg width="12" height="12" viewBox="0 0 256 256" fill="currentColor" style="opacity: 0.5;">
                        <path d="M181.66,133.66l-48,48a8,8,0,0,1-11.32,0l-48-48a8,8,0,0,1,11.32-11.32L120,156.69V48a8,8,0,0,1,16,0V156.69l34.34-34.35a8,8,0,0,1,11.32,11.32Zm-96-35.32a8,8,0,0,0,11.32,0L131.31,64V172a8,8,0,0,0,16,0V64l34.35,34.34a8,8,0,0,0,11.32-11.32l-48-48a8,8,0,0,0-11.32,0l-48,48A8,8,0,0,0,85.66,98.34Z"></path>
                      </svg>
                    </div>
                  </th>
                  <th class="col-right">Objects</th>
                  <th class="col-right">Size</th>
                </tr>
              </thead>
              <tbody v-if="filteredBuckets.length > 0">
                <tr v-for="bucket in filteredBuckets" :key="bucket.id">
                  <td class="col-name">
                    <router-link :to="`/r2/${bucket.id}`" class="cf-r2-link">
                      {{ bucket.name }}
                    </router-link>
                  </td>
                  <td class="col-right">{{ bucket.objectsCount.toLocaleString() }}</td>
                  <td class="col-right">{{ bucket.sizeFormatted }}</td>
                </tr>
              </tbody>
              <tbody v-else>
                <tr>
                  <td colspan="3" style="text-align: center; padding: 2.5rem 1rem; color: #71717a;">
                    No buckets found
                  </td>
                </tr>
              </tbody>
            </table>

            <!-- Embedded Pagination inside Card -->
            <div class="cf-r2-table-footer-pagination">
              <div class="cf-r2-pagination-btn-group">
                <button class="cf-r2-page-control-btn" disabled aria-label="Go to previous page">
                  <svg width="12" height="12" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M165.66,202.34a8,8,0,0,1-11.32,11.32l-80-80a8,8,0,0,1,0-11.32l80-80a8,8,0,0,1,11.32,11.32L91.31,128Z"></path>
                  </svg>
                </button>
                <button class="cf-r2-page-control-btn" disabled aria-label="Go to next page">
                  <svg width="12" height="12" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M181.66,133.66l-80,80a8,8,0,0,1-11.32-11.32L164.69,128,90.34,53.66a8,8,0,0,1,11.32-11.32l80,80A8,8,0,0,1,181.66,133.66Z"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </main>

      <!-- Right Continuous Integrated Sidebar Column (400px) -->
      <aside class="cf-r2-sidebar-column" role="complementary">
        <!-- Section 1: Usage -->
        <div class="cf-r2-aside-section">
          <div class="cf-r2-usage-headline">
            <h2 class="cf-r2-usage-title">Usage</h2>
            <span class="cf-r2-usage-date-range">November 11 - December 11</span>
          </div>

          <div style="display: flex; flex-direction: column; gap: 0.25rem;">
            <div class="cf-r2-usage-row">
              <div class="cf-r2-usage-tag">
                <span>Class A Operations</span>
                <button type="button" class="cf-r2-info-icon-btn" title="Mutating actions like PutObject, CreateBucket.">
                  <svg width="13" height="13" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M128,24A104,104,0,1,0,232,128,104.11,104.11,0,0,0,128,24Zm0,192a88,88,0,1,1,88-88A88.1,88.1,0,0,1,128,216Zm16-40a8,8,0,0,1-8,8,16,16,0,0,1-16-16V128a8,8,0,0,1,0-16,16,16,0,0,1,16,16v40A8,8,0,0,1,144,176ZM112,84a12,12,0,1,1,12,12A12,12,0,0,1,112,84Z"></path>
                  </svg>
                </button>
              </div>
              <div class="cf-r2-metric-digits">54</div>
            </div>

            <div class="cf-r2-usage-row">
              <div class="cf-r2-usage-tag">
                <span>Class B Operations</span>
                <button type="button" class="cf-r2-info-icon-btn" title="Read actions like GetObject.">
                  <svg width="13" height="13" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M128,24A104,104,0,1,0,232,128,104.11,104.11,0,0,0,128,24Zm0,192a88,88,0,1,1,88-88A88.1,88.1,0,0,1,128,216Zm16-40a8,8,0,0,1-8,8,16,16,0,0,1-16-16V128a8,8,0,0,1,0-16,16,16,0,0,1,16,16v40A8,8,0,0,1,144,176ZM112,84a12,12,0,1,1,12,12A12,12,0,0,1,112,84Z"></path>
                  </svg>
                </button>
              </div>
              <div class="cf-r2-metric-digits">16</div>
            </div>

            <div class="cf-r2-usage-row">
              <div class="cf-r2-usage-tag">
                <span>Total storage</span>
                <button type="button" class="cf-r2-info-icon-btn" title="Total payload bytes stored in R2.">
                  <svg width="13" height="13" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M128,24A104,104,0,1,0,232,128,104.11,104.11,0,0,0,128,24Zm0,192a88,88,0,1,1,88-88A88.1,88.1,0,0,1,128,216Zm16-40a8,8,0,0,1-8,8,16,16,0,0,1-16-16V128a8,8,0,0,1,0-16,16,16,0,0,1,16,16v40A8,8,0,0,1,144,176ZM112,84a12,12,0,1,1,12,12A12,12,0,0,1,112,84Z"></path>
                  </svg>
                </button>
              </div>
              <div class="cf-r2-metric-digits">49.98 MB</div>
            </div>
          </div>
        </div>

        <!-- Section 2: Billing Dashboard -->
        <div class="cf-r2-aside-section">
          <div>
            <router-link to="/billing" class="cf-r2-billing-title-link">
              <span>Billing Dashboard</span>
              <svg width="12" height="12" viewBox="0 0 256 256" fill="currentColor">
                <path d="M181.66,133.66l-80,80a8,8,0,0,1-11.32-11.32L164.69,128,90.34,53.66a8,8,0,0,1,11.32-11.32l80,80A8,8,0,0,1,181.66,133.66Z"></path>
              </svg>
            </router-link>
          </div>

          <div class="cf-r2-billing-inner-box">
            <div class="cf-r2-billing-box-top">
              <span class="cf-r2-box-top-label">
                <span>Billable usage (current period)</span>
                <svg width="11" height="11" viewBox="0 0 256 256" fill="currentColor">
                  <path d="M128,24A104,104,0,1,0,232,128,104.11,104.11,0,0,0,128,24Zm0,192a88,88,0,1,1,88-88A88.1,88.1,0,0,1,128,216Zm16-40a8,8,0,0,1-8,8,16,16,0,0,1-16-16V128a8,8,0,0,1,0-16,16,16,0,0,1,16,16v40A8,8,0,0,1,144,176ZM112,84a12,12,0,1,1,12,12A12,12,0,0,1,112,84Z"></path>
                </svg>
              </span>
              <button class="cf-r2-budget-text-btn" @click="handleBudgetAlert">
                + Add Budget Alert
              </button>
            </div>

            <div class="cf-r2-gauge-row">
              <div class="cf-r2-gauge-circle-wrap">
                <svg viewBox="0 0 48 48" width="48" height="48">
                  <circle
                    cx="24"
                    cy="24"
                    r="20"
                    fill="none"
                    stroke="#e2e8f0"
                    stroke-width="4"
                  ></circle>
                </svg>
                <span class="cf-r2-gauge-center-price">$0.00</span>
              </div>
              <div class="cf-r2-gauge-desc">
                No billable usage incurred yet
              </div>
            </div>
          </div>
        </div>

        <!-- Section 3: Account Details -->
        <div class="cf-r2-aside-section">
          <h3 class="cf-r2-details-header">Account Details</h3>

          <div style="display: flex; flex-direction: column; gap: 0.5rem;">
            <!-- Account ID Field Box -->
            <div class="cf-r2-kv-box">
              <span class="cf-r2-kv-label">Account ID</span>
              <span class="cf-r2-kv-val" :title="accountId">{{ accountId }}</span>
              <button
                class="cf-r2-mini-copy-btn"
                :title="copiedKey === 'acc' ? 'Copied!' : 'Copy Account ID'"
                @click="copyText(accountId, 'acc')"
              >
                <CfIcon v-if="copiedKey !== 'acc'" name="cis-clipboard" size="13" />
                <CfIcon v-else name="cis-checkmark-1-medium" size="13" style="color: #10b981;" />
              </button>
            </div>

            <!-- S3 API Field Box -->
            <div class="cf-r2-kv-box">
              <span class="cf-r2-kv-label">S3 API</span>
              <span class="cf-r2-kv-val" :title="s3ApiEndpoint">{{ s3ApiEndpoint }}</span>
              <button
                class="cf-r2-mini-copy-btn"
                :title="copiedKey === 's3' ? 'Copied!' : 'Copy S3 API'"
                @click="copyText(s3ApiEndpoint, 's3')"
              >
                <CfIcon v-if="copiedKey !== 's3'" name="cis-clipboard" size="13" />
                <CfIcon v-else name="cis-checkmark-1-medium" size="13" style="color: #10b981;" />
              </button>
            </div>
          </div>

          <!-- Manage API Tokens points directly to /r2/api-tokens -->
          <router-link to="/r2/api-tokens" class="cf-r2-manage-tokens-full-btn">
            <CfIcon name="cis-key-1" size="13" />
            <span>Manage API Tokens</span>
          </router-link>
        </div>
      </aside>
    </div>

    <!-- METRICS TAB -->
    <div v-else-if="currentTab === 'metrics'" style="padding: 2.25rem 2rem;">
      <div class="cf-r2-hero-row" style="margin-bottom: 1.5rem;">
        <div>
          <h2 class="cf-r2-page-title">R2 Metrics & Analytics</h2>
          <p class="cf-r2-subtitle-text">Real-time storage bandwidth, operations, and egress metrics across all buckets.</p>
        </div>
      </div>

      <div class="stats-overview-grid" style="margin-bottom: 1.5rem;">
        <CfStatCard
          label="Total Stored Data"
          value="49.98 MB"
          delta="Active"
          delta-type="up"
          tooltip="Total object storage volume across all buckets"
        />
        <CfStatCard
          label="Total Buckets"
          value="1"
          unit="Bucket"
          delta="APAC"
          delta-type="neutral"
          tooltip="Active buckets in account"
        />
        <CfStatCard
          label="Class A Requests"
          value="54"
          delta="Free Tier"
          delta-type="up"
          tooltip="1,000,000 requests/month included free"
        />
        <CfStatCard
          label="Class B Requests"
          value="16"
          delta="Free Tier"
          delta-type="up"
          tooltip="10,000,000 requests/month included free"
        />
      </div>
    </div>

    <!-- DATA MIGRATION TAB -->
    <div v-else-if="currentTab === 'slurper'" style="padding: 2.25rem 2rem;">
      <div class="cf-r2-hero-row" style="margin-bottom: 1.5rem;">
        <div>
          <h2 class="cf-r2-page-title">Super Slurper (Data Migration)</h2>
          <p class="cf-r2-subtitle-text">Migrate object storage buckets from AWS S3, Google Cloud Storage, or Azure Blob to Cloudflare R2 seamlessly with zero downtime.</p>
        </div>
        <button class="cf-r2-primary-btn" @click="handleNewMigration">
          + Start new migration
        </button>
      </div>

      <div class="cf-r2-table-wrapper" style="max-width: 722px;">
        <table class="cf-r2-pure-table">
          <thead>
            <tr>
              <th class="col-name">Job Name</th>
              <th>Source</th>
              <th>Destination Bucket</th>
              <th>Status</th>
              <th class="col-right">Objects Transferred</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td class="col-name"><strong>aws-s3-production-sync</strong></td>
              <td>Amazon S3 (us-east-1)</td>
              <td>code-collection</td>
              <td><span class="status-badge status-active">Complete</span></td>
              <td class="col-right">397 / 397</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- CREATE BUCKET MODAL -->
    <div v-if="showCreateModal" class="cf-r2-modal-overlay" @click.self="showCreateModal = false">
      <div class="cf-r2-modal-card">
        <div class="cf-r2-modal-header">
          <h2 class="cf-r2-modal-title">Create bucket</h2>
          <button class="cf-r2-mini-copy-btn" @click="showCreateModal = false">✕</button>
        </div>
        <form @submit.prevent="createBucket">
          <div class="cf-r2-modal-body">
            <div class="cf-r2-form-group">
              <label class="cf-r2-form-label" for="modal-bucket-name">Bucket name</label>
              <input
                id="modal-bucket-name"
                v-model="newBucket.name"
                type="text"
                class="cf-r2-form-input"
                placeholder="e.g. my-production-assets"
                required
              />
              <span class="cf-r2-form-hint">
                Must be between 3 and 63 characters long and only contain lowercase letters, numbers, and hyphens.
              </span>
            </div>

            <div class="cf-r2-form-group">
              <label class="cf-r2-form-label" for="modal-bucket-loc">Location hint</label>
              <select id="modal-bucket-loc" v-model="newBucket.location" class="cf-r2-form-select">
                <option value="auto">Automatic (Optimized for low latency)</option>
                <option value="apac">Asia-Pacific (APAC)</option>
                <option value="wnam">Western North America (WNAM)</option>
                <option value="enam">Eastern North America (ENAM)</option>
                <option value="weur">Western Europe (WEUR)</option>
                <option value="eeur">Eastern Europe (EEUR)</option>
              </select>
            </div>
          </div>
          <div class="cf-r2-modal-footer">
            <button type="button" class="cf-r2-btn-secondary" @click="showCreateModal = false">
              Cancel
            </button>
            <button type="submit" class="cf-r2-primary-btn">
              Create bucket
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const currentTab = ref('overview')
const searchQuery = ref('')
const copiedKey = ref('')
const showCreateModal = ref(false)

const accountId = '330b08634cd698ef2ff8587d66ac4b2c'
const s3ApiEndpoint = `https://${accountId}.r2.cloudflarestorage.com`

// Authentic Cloudflare single bucket row
const buckets = ref([
  {
    id: 'code-collection',
    name: 'code-collection',
    region: 'APAC',
    objectsCount: 397,
    sizeFormatted: '49.98 MB'
  }
])

const newBucket = ref({
  name: '',
  location: 'auto',
  storageClass: 'Standard'
})

const filteredBuckets = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return buckets.value
  return buckets.value.filter(b => b.name.toLowerCase().includes(q))
})

const toggleSort = () => {
  buckets.value.reverse()
}

const handleRefresh = () => {
  const temp = [...buckets.value]
  buckets.value = []
  setTimeout(() => {
    buckets.value = temp
  }, 100)
}

const copyText = async (text, key) => {
  try {
    await navigator.clipboard.writeText(text)
    copiedKey.value = key
    setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = ''
    }, 2000)
  } catch (err) {
    console.warn('Copy failed:', err)
  }
}

const handleBudgetAlert = () => {
  alert('Redirecting to Cloudflare Billing Alerts configuration.')
}

const handleNewMigration = () => {
  alert('Starting Super Slurper automated migration wizard.')
}

const createBucket = () => {
  if (!newBucket.value.name) return
  const cleanName = newBucket.value.name.toLowerCase().replace(/[^a-z0-9-]/g, '-')
  buckets.value.unshift({
    id: cleanName,
    name: cleanName,
    region: newBucket.value.location.toUpperCase(),
    objectsCount: 0,
    sizeFormatted: '0 B'
  })
  showCreateModal.value = false
  newBucket.value.name = ''
}
</script>
