<template>
  <div class="cf-bucket-page">
    <!-- Breadcrumbs -->
    <div class="breadcrumbs">
      <router-link to="/storage">Object Storage</router-link>
      <span class="divider">/</span>
      <span class="current">{{ bucketName }}</span>
    </div>

    <!-- Header -->
    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.6rem; flex-wrap: wrap;">
          <CfIcon name="r2" size="22" style="color: #f6821f;" />
          <h1 class="page-title">{{ bucketName }}</h1>
          <CfBadge :variant="bucket.is_public ? 'warning' : 'neutral'" size="sm">
            {{ bucket.is_public ? 'Public Read' : 'Private' }}
          </CfBadge>
          <CfBadge v-if="bucket.versioning_enabled" variant="info" size="sm">Versioning</CfBadge>
          <CfBadge v-if="bucket.erasure_enabled" variant="success" size="sm">RS 4+2 Erasure</CfBadge>
          <CfBadge v-if="bucket.sse_enabled" variant="info" size="sm">SSE-AES256</CfBadge>
          <CfBadge v-if="bucket.object_lock_enabled" variant="danger" size="sm">
            WORM Lock ({{ bucket.retention_days }}d)
          </CfBadge>
          <CfBadge v-if="bucket.lifecycle_days > 0" variant="neutral" size="sm">
            TTL {{ bucket.lifecycle_days }}d
          </CfBadge>
          <CfBadge v-if="bucket.webhook_url" variant="success" size="sm">Webhook</CfBadge>
        </div>
        <p class="page-subtitle">
          S3-compatible bucket on local Orbit node · Created {{ formatDate(bucket.created_at) }}
        </p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <CfButton variant="secondary" icon="refresh" @click="loadAll">Refresh</CfButton>
        <CfButton variant="primary" icon="upload" @click="showUploadModal = true">Upload Object</CfButton>
      </div>
    </div>

    <!-- Segmented Tabs (Versions tab unlocks automatically when Versioning is enabled) -->
    <div class="cf-tabs-bar" role="tablist">
      <button
        role="tab"
        class="cf-tab-btn"
        :class="{ active: activeTab === 'objects' }"
        @click="activeTab = 'objects'"
      >
        Objects ({{ objects.length }})
      </button>
      <button
        v-if="bucket.versioning_enabled"
        role="tab"
        class="cf-tab-btn"
        :class="{ active: activeTab === 'versions' }"
        @click="activeTab = 'versions'; fetchVersions()"
      >
        Versions ({{ versionCount }})
      </button>
      <button
        role="tab"
        class="cf-tab-btn"
        :class="{ active: activeTab === 'metrics' }"
        @click="activeTab = 'metrics'; if (bucket.erasure_enabled) fetchErasureHealth()"
      >
        Metrics & Quota
      </button>
      <button
        role="tab"
        class="cf-tab-btn"
        :class="{ active: activeTab === 'settings' }"
        @click="activeTab = 'settings'"
      >
        Settings & Advanced
      </button>
    </div>

    <CfAlert v-if="actionError" variant="danger" title="Operation Failed" style="margin-bottom: 0.25rem;">
      {{ actionError }}
    </CfAlert>
    <CfAlert v-if="actionSuccess" variant="success" title="Success" style="margin-bottom: 0.25rem;">
      {{ actionSuccess }}
    </CfAlert>

    <!-- TAB 1: OBJECTS -->
    <div v-if="activeTab === 'objects'" class="cf-tab-panel">
      <div class="cf-meta-strip">
        <div class="cf-meta-card">
          <span class="cf-meta-label">Storage Used</span>
          <span class="cf-meta-value">{{ formatBytes(bucket.used_bytes) }} / {{ bucket.quota_bytes > 0 ? formatBytes(bucket.quota_bytes) : 'Unlimited' }}</span>
        </div>
        <div class="cf-meta-card">
          <span class="cf-meta-label">Objects</span>
          <span class="cf-meta-value">{{ objectCount }}</span>
        </div>
        <div class="cf-meta-card">
          <span class="cf-meta-label">Region</span>
          <span class="cf-meta-value">{{ region }}</span>
        </div>
        <div class="cf-meta-card">
          <span class="cf-meta-label">S3 Endpoint</span>
          <span class="cf-meta-value" style="font-family: monospace; font-size: 0.8rem;">{{ s3BucketUrl }}</span>
        </div>
      </div>

      <CfCard title="Objects" body-class="table-body">
        <template #actions>
          <div style="width: 260px;">
            <CfInput v-model="searchPrefix" placeholder="Filter by key prefix..." prefix-icon="search" clearable />
          </div>
        </template>

        <div v-if="loading" style="padding: 2rem; display: flex; justify-content: center;">
          <CfLoader size="26" />
        </div>

        <CfEmptyState
          v-else-if="!filteredObjects.length"
          title="Bucket is empty"
          description="Upload an object via the dashboard, AWS CLI, or any S3 SDK."
          action-text="Upload Object"
          action-icon="upload"
          @action="showUploadModal = true"
        />

        <CfTable v-else :columns="objectColumns" :data="filteredObjects" row-key="key">
          <template #col-key="{ row }">
            <div style="display: flex; align-items: center; gap: 0.45rem;">
              <CfIcon name="file" size="14" style="color: #71717a;" />
              <span style="font-weight: 500; font-family: monospace;">{{ row.key }}</span>
            </div>
          </template>
          <template #col-size="{ val }">
            <span>{{ formatBytes(val) }}</span>
          </template>
          <template #actions="{ row }">
            <div style="display: flex; gap: 0.35rem; justify-content: flex-end;">
              <CfButton variant="secondary" size="sm" @click="openPresign(row.key)">Presign</CfButton>
              <CfButton variant="secondary" size="sm" @click="downloadObject(row.key)">Download</CfButton>
              <CfButton
                v-if="bucket.versioning_enabled"
                variant="secondary"
                size="sm"
                @click="viewKeyVersions(row.key)"
              >
                History
              </CfButton>
              <CfButton
                variant="danger-outline"
                size="sm"
                icon="trash"
                icon-only
                title="Delete object"
                @click="promptDeleteObject(row.key)"
              />
            </div>
          </template>
        </CfTable>
      </CfCard>
    </div>

    <!-- TAB 2: VERSIONS (Unlocked when Versioning is enabled) -->
    <div v-if="activeTab === 'versions' && bucket.versioning_enabled" class="cf-tab-panel">
      <CfCard title="Object Version History" subtitle="Restore previous versions or purge delete markers">
        <template #actions>
          <div style="display: flex; gap: 0.5rem; align-items: center;">
            <CfInput v-model="versionKeyFilter" placeholder="Filter by key..." clearable style="width: 220px;" />
            <CfButton variant="secondary" size="sm" @click="fetchVersions">Reload</CfButton>
          </div>
        </template>

        <CfEmptyState
          v-if="!filteredVersions.length"
          title="No archived versions yet"
          description="Overwriting or deleting objects while Versioning is enabled automatically archives previous versions here."
        />

        <CfTable v-else :columns="versionColumns" :data="filteredVersions" row-key="version_id">
          <template #col-object_key="{ val }">
            <code style="font-size: 0.85rem;">{{ val }}</code>
          </template>
          <template #col-type="{ row }">
            <CfBadge :variant="row.is_delete_marker ? 'danger' : 'info'" size="sm">
              {{ row.is_delete_marker ? 'Delete Marker' : 'Version' }}
            </CfBadge>
          </template>
          <template #col-size_bytes="{ val }">
            <span>{{ formatBytes(val) }}</span>
          </template>
          <template #col-created_at="{ val }">
            <span>{{ formatDate(val) }}</span>
          </template>
          <template #actions="{ row }">
            <div style="display: flex; gap: 0.35rem; justify-content: flex-end;">
              <CfButton
                v-if="!row.is_delete_marker"
                variant="primary"
                size="sm"
                @click="restoreVersion(row.version_id)"
              >
                Restore
              </CfButton>
              <CfButton variant="danger-outline" size="sm" @click="purgeVersion(row.version_id)">Purge</CfButton>
            </div>
          </template>
        </CfTable>
      </CfCard>
    </div>

    <!-- TAB 3: METRICS & QUOTA -->
    <div v-if="activeTab === 'metrics'" class="cf-tab-panel">
      <div class="cf-metrics-grid">
        <CfCard title="Quota Utilization">
          <div style="display: flex; flex-direction: column; gap: 0.75rem;">
            <div style="display: flex; justify-content: space-between; font-size: 0.875rem;">
              <span>Used: <strong>{{ formatBytes(bucket.used_bytes) }}</strong></span>
              <span>Limit: <strong>{{ bucket.quota_bytes > 0 ? formatBytes(bucket.quota_bytes) : 'Unlimited' }}</strong></span>
            </div>
            <div class="cf-progress-track">
              <div class="cf-progress-fill" :style="{ width: quotaPercent + '%' }"></div>
            </div>
            <span style="font-size: 0.8rem; color: #71717a;">{{ quotaPercent }}% of configured bucket quota</span>
          </div>
        </CfCard>

        <CfCard title="Bucket Inventory">
          <div class="cf-kv-list">
            <div class="cf-kv-row"><span>Active Objects</span><strong>{{ objectCount }}</strong></div>
            <div class="cf-kv-row"><span>Archived Versions</span><strong>{{ versionCount }}</strong></div>
            <div class="cf-kv-row"><span>Versioning</span><strong>{{ bucket.versioning_enabled ? 'Enabled' : 'Disabled' }}</strong></div>
            <div class="cf-kv-row"><span>WORM Object Lock</span><strong>{{ bucket.object_lock_enabled ? `Enabled (${bucket.retention_days}d)` : 'Disabled' }}</strong></div>
            <div class="cf-kv-row"><span>Erasure Coding</span><strong>{{ bucket.erasure_enabled ? 'RS (4 Data + 2 Parity)' : 'Single Copy' }}</strong></div>
          </div>
        </CfCard>
      </div>

      <CfCard
        v-if="bucket.erasure_enabled && erasureHealth"
        title="Multi-Drive Erasure Coding Health (Reed-Solomon 4+2)"
        subtitle="Tolerates up to 2 simultaneous drive failures with automatic shard healing"
      >
        <template #actions>
          <CfButton variant="secondary" size="sm" @click="healErasureShards">
            Verify & Heal Shards
          </CfButton>
        </template>
        <div class="cf-drives-grid">
          <div v-for="d in erasureHealth.drives" :key="d.drive_index" class="cf-drive-card">
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <strong>Drive {{ d.drive_index }}</strong>
              <CfBadge :variant="d.missing_shards + d.corrupt_shards === 0 ? 'success' : 'warning'" size="sm">
                {{ d.missing_shards + d.corrupt_shards === 0 ? 'Healthy' : 'Degraded' }}
              </CfBadge>
            </div>
            <div style="font-size: 0.78rem; color: #71717a; margin-top: 0.35rem;">
              Shards: {{ d.healthy_shards }} OK · {{ d.missing_shards }} missing · {{ d.corrupt_shards }} corrupt
            </div>
          </div>
        </div>
      </CfCard>
    </div>

    <!-- TAB 4: SETTINGS & ADVANCED FEATURES -->
    <div v-if="activeTab === 'settings'" class="cf-tab-panel">
      <CfCard
        title="Advanced Storage Features"
        subtitle="Enable enterprise S3 capabilities on this bucket. Features unlock automatically once enabled."
      >
        <template #actions>
          <CfButton variant="secondary" size="sm" @click="enableAllAdvanced">
            Enable All Advanced
          </CfButton>
        </template>

        <div class="cf-settings-stack">
          <label class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">Reed-Solomon Erasure Coding (4 Data + 2 Parity)</div>
              <div class="cf-toggle-desc">Stripe new objects across 6 storage drives with GF(2^8) parity. Recovers & auto-heals from up to 2 drive failures.</div>
            </div>
            <input type="checkbox" v-model="settingsForm.erasure_enabled" class="cf-switch" />
          </label>

          <label class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">SSE-KMS Envelope Encryption (AES-256-GCM)</div>
              <div class="cf-toggle-desc">Encrypt every object at rest with a per-object 256-bit DEK wrapped by Orbit's MasterKey.</div>
            </div>
            <input type="checkbox" v-model="settingsForm.sse_enabled" class="cf-switch" />
          </label>

          <label class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">Object Versioning</div>
              <div class="cf-toggle-desc">Automatically archive previous versions on overwrite or delete and unlock the Versions history tab.</div>
            </div>
            <input type="checkbox" v-model="settingsForm.versioning_enabled" class="cf-switch" />
          </label>

          <div class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">Object Lock (WORM Retention)</div>
              <div class="cf-toggle-desc">Block object deletion or overwrite until the retention period expires.</div>
            </div>
            <div style="display: flex; align-items: center; gap: 0.75rem;">
              <CfInput
                v-if="settingsForm.object_lock_enabled"
                v-model.number="settingsForm.retention_days"
                type="number"
                min="1"
                placeholder="Days"
                style="width: 110px;"
              />
              <input type="checkbox" v-model="settingsForm.object_lock_enabled" class="cf-switch" />
            </div>
          </div>

          <label class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">Public Read Access</div>
              <div class="cf-toggle-desc">Allow unauthenticated GET and HEAD requests to objects in this bucket.</div>
            </div>
            <input type="checkbox" v-model="settingsForm.is_public" class="cf-switch" />
          </label>

          <div class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">Lifecycle Auto-Expiry (TTL Days)</div>
              <div class="cf-toggle-desc">Expire objects older than N days (0 disables lifecycle expiry). Locked objects are preserved.</div>
            </div>
            <div style="display: flex; align-items: center; gap: 0.5rem;">
              <CfInput v-model.number="settingsForm.lifecycle_days" type="number" min="0" style="width: 110px;" />
              <CfButton
                v-if="bucket.lifecycle_days > 0"
                variant="secondary"
                size="sm"
                @click="runLifecycleNow"
              >
                Run Now
              </CfButton>
            </div>
          </div>

          <div class="cf-toggle-row">
            <div>
              <div class="cf-toggle-title">Storage Quota</div>
              <div class="cf-toggle-desc">Enforce hard storage limit on S3 PUT requests (0 = Unlimited).</div>
            </div>
            <CfSelect
              :model-value="settingsForm.quota_bytes"
              :options="[
                { value: 0, label: 'Unlimited' },
                { value: 1073741824, label: '1 GB' },
                { value: 10737418240, label: '10 GB' },
                { value: 53687091200, label: '50 GB' },
                { value: 107374182400, label: '100 GB' },
              ]"
              @update:model-value="settingsForm.quota_bytes = Number($event)"
            />
          </div>

          <div class="cf-toggle-row" style="align-items: flex-start;">
            <div>
              <div class="cf-toggle-title">Event Webhook URL</div>
              <div class="cf-toggle-desc">Receive async JSON POST notifications on ObjectCreated and ObjectRemoved events.</div>
            </div>
            <CfInput
              v-model="settingsForm.webhook_url"
              placeholder="https://example.com/webhooks/s3"
              style="width: 320px;"
            />
          </div>

          <div style="display: flex; justify-content: flex-end; padding-top: 0.5rem;">
            <CfButton variant="primary" :loading="savingSettings" @click="saveSettings">
              Save Bucket Settings
            </CfButton>
          </div>
        </div>
      </CfCard>
    </div>

    <!-- Upload Modal -->
    <CfModal v-model="showUploadModal" title="Upload Object" size="sm">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <CfInput v-model="uploadKey" label="Object Key" placeholder="e.g. backups/db.sql or hello.txt" required />
        <div>
          <label style="display: block; font-size: 0.85rem; font-weight: 500; margin-bottom: 0.35rem;">Select File</label>
          <input type="file" @change="onFilePicked" />
        </div>
        <CfAlert v-if="uploadError" variant="danger" title="Upload Error">{{ uploadError }}</CfAlert>
      </div>
      <template #footer>
        <CfButton variant="outline" @click="showUploadModal = false">Cancel</CfButton>
        <CfButton variant="primary" :loading="uploading" @click="uploadObject">Upload</CfButton>
      </template>
    </CfModal>

    <!-- Presigned URL Modal -->
    <CfModal v-model="showPresignModal" title="Generate Presigned S3 URL" size="md">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <div style="font-size: 0.875rem;">
          Object: <code>{{ presignKey }}</code>
        </div>
        <div style="display: flex; gap: 0.75rem; align-items: center;">
          <label style="font-size: 0.85rem; font-weight: 500;">Expiry:</label>
          <CfSelect
            :model-value="presignTTL"
            :options="[
              { value: 900, label: '15 Minutes' },
              { value: 3600, label: '1 Hour' },
              { value: 86400, label: '24 Hours' },
              { value: 604800, label: '7 Days' },
            ]"
            @update:model-value="presignTTL = Number($event); generatePresign()"
          />
        </div>
        <CfCopyField v-if="presignUrl" label="Signed URL" :value="presignUrl" />
      </div>
      <template #footer>
        <CfButton variant="primary" @click="showPresignModal = false">Done</CfButton>
      </template>
    </CfModal>

    <!-- Delete Object Confirmation -->
    <CfConfirmDialog
      v-model="showDeleteDialog"
      title="Delete Object"
      :message="`Delete '${objectToDelete}' from bucket '${bucketName}'?`"
      confirm-label="Delete Object"
      danger
      @confirm="confirmDeleteObject"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { apiFetch } from '@/api/client.js'

const route = useRoute()
const bucketName = computed(() => String(route.params.id || ''))
const activeTab = ref('objects')

const loading = ref(true)
const bucket = ref({
  bucket_id: '',
  is_public: false,
  quota_bytes: 10737418240,
  used_bytes: 0,
  versioning_enabled: false,
  object_lock_enabled: false,
  erasure_enabled: false,
  sse_enabled: false,
  retention_days: 0,
  lifecycle_days: 0,
  webhook_url: ''
})
const erasureHealth = ref(null)
const region = ref('us-east-1')
const objectCount = ref(0)
const versionCount = ref(0)
const objects = ref([])
const versions = ref([])
const searchPrefix = ref('')
const versionKeyFilter = ref('')

const actionError = ref('')
const actionSuccess = ref('')

const settingsForm = ref({
  is_public: false,
  versioning_enabled: false,
  object_lock_enabled: false,
  erasure_enabled: false,
  sse_enabled: false,
  retention_days: 7,
  lifecycle_days: 0,
  quota_bytes: 10737418240,
  webhook_url: ''
})
const savingSettings = ref(false)

const showUploadModal = ref(false)
const uploadKey = ref('')
const uploadFile = ref(null)
const uploading = ref(false)
const uploadError = ref('')

const showPresignModal = ref(false)
const presignKey = ref('')
const presignTTL = ref(3600)
const presignUrl = ref('')

const showDeleteDialog = ref(false)
const objectToDelete = ref('')

const s3BucketUrl = computed(() => `${window.location.origin}/s3/${bucketName.value}`)

const objectColumns = [
  { key: 'key', label: 'Key' },
  { key: 'size', label: 'Size', width: '110px' },
  { key: 'modified', label: 'Last Modified', width: '170px' }
]

const versionColumns = [
  { key: 'object_key', label: 'Object Key' },
  { key: 'type', label: 'Type', width: '125px' },
  { key: 'size_bytes', label: 'Size', width: '100px' },
  { key: 'created_at', label: 'Archived At', width: '165px' }
]

const filteredObjects = computed(() => {
  const q = searchPrefix.value.trim().toLowerCase()
  if (!q) return objects.value
  return objects.value.filter(o => o.key.toLowerCase().includes(q))
})

const filteredVersions = computed(() => {
  const q = versionKeyFilter.value.trim().toLowerCase()
  if (!q) return versions.value
  return versions.value.filter(v => v.object_key.toLowerCase().includes(q))
})

const quotaPercent = computed(() => {
  if (!bucket.value.quota_bytes || bucket.value.quota_bytes <= 0) return 0
  return Math.min(100, Math.round((bucket.value.used_bytes / bucket.value.quota_bytes) * 100))
})

function formatBytes(bytes) {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

function formatDate(raw) {
  if (!raw) return '—'
  return new Date(raw).toLocaleString()
}

function parseListObjectsXml(xmlText) {
  const doc = new DOMParser().parseFromString(xmlText, 'application/xml')
  return Array.from(doc.querySelectorAll('Contents')).map(n => ({
    key: n.querySelector('Key')?.textContent || '',
    size: Number(n.querySelector('Size')?.textContent || 0),
    modified: formatDate(n.querySelector('LastModified')?.textContent || '')
  })).filter(o => o.key)
}

async function loadAll() {
  loading.value = true
  actionError.value = ''
  try {
    const [detailRes, listRes] = await Promise.all([
      apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}`),
      apiFetch(`/s3/${encodeURIComponent(bucketName.value)}?list-type=2`)
    ])
    if (detailRes.ok) {
      const data = await detailRes.json()
      bucket.value = data.bucket
      objectCount.value = data.object_count || 0
      versionCount.value = data.version_count || 0
      region.value = data.region || 'us-east-1'
      settingsForm.value = {
        is_public: !!data.bucket.is_public,
        versioning_enabled: !!data.bucket.versioning_enabled,
        object_lock_enabled: !!data.bucket.object_lock_enabled,
        erasure_enabled: !!data.bucket.erasure_enabled,
        sse_enabled: !!data.bucket.sse_enabled,
        retention_days: data.bucket.retention_days || 7,
        lifecycle_days: data.bucket.lifecycle_days || 0,
        quota_bytes: data.bucket.quota_bytes ?? 10737418240,
        webhook_url: data.bucket.webhook_url || ''
      }
      if (data.bucket.erasure_enabled) {
        fetchErasureHealth()
      }
    }
    if (listRes.ok) {
      objects.value = parseListObjectsXml(await listRes.text())
    }
  } catch (e) {
    actionError.value = e.message
  } finally {
    loading.value = false
  }
}

async function fetchVersions() {
  const res = await apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/versions`)
  if (res.ok) {
    const data = await res.json()
    versions.value = data.versions || []
    versionCount.value = versions.value.length
  }
}

function viewKeyVersions(key) {
  versionKeyFilter.value = key
  activeTab.value = 'versions'
  fetchVersions()
}

function enableAllAdvanced() {
  settingsForm.value.erasure_enabled = true
  settingsForm.value.sse_enabled = true
  settingsForm.value.versioning_enabled = true
  settingsForm.value.object_lock_enabled = true
  if (!settingsForm.value.retention_days || settingsForm.value.retention_days < 1) {
    settingsForm.value.retention_days = 7
  }
  if (!settingsForm.value.lifecycle_days || settingsForm.value.lifecycle_days < 1) {
    settingsForm.value.lifecycle_days = 30
  }
}

async function fetchErasureHealth() {
  const res = await apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/erasure/health`)
  if (res.ok) {
    erasureHealth.value = await res.json()
  }
}

async function healErasureShards() {
  actionError.value = ''
  actionSuccess.value = ''
  const res = await apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/erasure/heal`, {
    method: 'POST'
  })
  if (res.ok) {
    const data = await res.json()
    actionSuccess.value = `Erasure heal complete: ${data.healed_shards} shard(s) repaired.`
    await fetchErasureHealth()
  }
}

async function saveSettings() {
  savingSettings.value = true
  actionError.value = ''
  actionSuccess.value = ''
  try {
    const res = await apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/settings`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(settingsForm.value)
    })
    if (!res.ok) throw new Error(await res.text())
    actionSuccess.value = 'Bucket settings saved and active.'
    await loadAll()
  } catch (e) {
    actionError.value = e.message
  } finally {
    savingSettings.value = false
  }
}

async function runLifecycleNow() {
  actionError.value = ''
  actionSuccess.value = ''
  const res = await apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/lifecycle/run`, {
    method: 'POST'
  })
  if (res.ok) {
    const data = await res.json()
    actionSuccess.value = `Lifecycle run complete: ${data.expired_count} expired object(s) removed.`
    await loadAll()
  }
}

function onFilePicked(e) {
  const file = e.target.files?.[0]
  if (!file) return
  uploadFile.value = file
  if (!uploadKey.value) uploadKey.value = file.name
}

async function uploadObject() {
  uploadError.value = ''
  if (!uploadKey.value.trim() || !uploadFile.value) {
    uploadError.value = 'Both object key and file are required.'
    return
  }
  uploading.value = true
  try {
    const res = await apiFetch(`/s3/${encodeURIComponent(bucketName.value)}/${uploadKey.value.trim()}`, {
      method: 'PUT',
      headers: { 'Content-Type': uploadFile.value.type || 'application/octet-stream' },
      body: uploadFile.value
    })
    if (!res.ok) throw new Error(await res.text())
    showUploadModal.value = false
    uploadKey.value = ''
    uploadFile.value = null
    actionSuccess.value = 'Object uploaded successfully.'
    await loadAll()
  } catch (e) {
    uploadError.value = e.message
  } finally {
    uploading.value = false
  }
}

async function downloadObject(key) {
  const res = await apiFetch(`/s3/${encodeURIComponent(bucketName.value)}/${key}`)
  if (!res.ok) {
    actionError.value = `Download failed (HTTP ${res.status})`
    return
  }
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = key.split('/').pop() || 'object'
  a.click()
  URL.revokeObjectURL(url)
}

function openPresign(key) {
  presignKey.value = key
  presignUrl.value = ''
  showPresignModal.value = true
  generatePresign()
}

async function generatePresign() {
  const res = await apiFetch(`/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/presign`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ key: presignKey.value, method: 'GET', expires_seconds: presignTTL.value })
  })
  if (res.ok) {
    const data = await res.json()
    presignUrl.value = data.url
  }
}

function promptDeleteObject(key) {
  objectToDelete.value = key
  showDeleteDialog.value = true
}

async function confirmDeleteObject() {
  actionError.value = ''
  actionSuccess.value = ''
  const res = await apiFetch(`/s3/${encodeURIComponent(bucketName.value)}/${objectToDelete.value}`, {
    method: 'DELETE'
  })
  if (!res.ok) {
    actionError.value = await res.text()
    return
  }
  actionSuccess.value = `Deleted '${objectToDelete.value}'.`
  await loadAll()
  if (bucket.value.versioning_enabled) await fetchVersions()
}

async function restoreVersion(versionId) {
  actionError.value = ''
  const res = await apiFetch(
    `/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/versions/${encodeURIComponent(versionId)}/restore`,
    { method: 'POST' }
  )
  if (!res.ok) {
    actionError.value = await res.text()
    return
  }
  actionSuccess.value = `Restored version ${versionId}.`
  await loadAll()
  await fetchVersions()
}

async function purgeVersion(versionId) {
  const res = await apiFetch(
    `/api/v1/storage/buckets/${encodeURIComponent(bucketName.value)}/versions/${encodeURIComponent(versionId)}`,
    { method: 'DELETE' }
  )
  if (res.ok) {
    await fetchVersions()
  }
}

onMounted(loadAll)
</script>

<style scoped>
.cf-bucket-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.cf-tabs-bar {
  display: flex;
  gap: 0.4rem;
  border-bottom: 1px solid #e4e4e7;
  padding-bottom: 0.5rem;
}

.cf-tab-btn {
  padding: 0.45rem 0.9rem;
  border-radius: 6px;
  border: 1px solid transparent;
  background: transparent;
  color: #52525b;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
}

.cf-tab-btn.active {
  background: #ffffff;
  color: #18181b;
  border-color: #d4d4d8;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.cf-tab-panel {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.cf-meta-strip {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}

.cf-meta-card {
  background: #ffffff;
  border: 1px solid #e4e4e7;
  border-radius: 8px;
  padding: 0.85rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.cf-meta-label {
  font-size: 0.75rem;
  color: #71717a;
}

.cf-meta-value {
  font-size: 0.95rem;
  font-weight: 600;
  color: #18181b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cf-metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
}

.cf-progress-track {
  height: 10px;
  background: #f4f4f5;
  border-radius: 999px;
  overflow: hidden;
}

.cf-progress-fill {
  height: 100%;
  background: #0051c3;
  transition: width 0.25s ease;
}

.cf-kv-list {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.cf-kv-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
  color: #52525b;
}

.cf-settings-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.cf-toggle-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding-bottom: 0.9rem;
  border-bottom: 1px solid #f4f4f5;
}

.cf-toggle-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #18181b;
}

.cf-toggle-desc {
  font-size: 0.8rem;
  color: #71717a;
  margin-top: 0.15rem;
}

.cf-switch {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.cf-select-input {
  padding: 0.45rem 0.65rem;
  border: 1px solid #d4d4d8;
  border-radius: 6px;
  background: #ffffff;
  font-size: 0.875rem;
}

.cf-drives-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.85rem;
}

.cf-drive-card {
  border: 1px solid #e4e4e7;
  border-radius: 8px;
  padding: 0.75rem 0.9rem;
  background: #fafafa;
}
</style>
