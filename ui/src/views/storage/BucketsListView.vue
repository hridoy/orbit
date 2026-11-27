<template>
  <div class="cf-storage-page">
    <div class="breadcrumbs">
      <router-link to="/">Storage</router-link>
      <span class="divider">/</span>
      <span class="current">Object Storage</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.65rem; flex-wrap: wrap;">
          <h1 class="page-title">Object Storage (S3)</h1>
          <CfBadge variant="success" size="sm" dot>S3 Compatible</CfBadge>
        </div>
        <p class="page-subtitle">
          Zero-egress local S3-compatible object storage buckets and SigV4 access credentials on this VPS node.
        </p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <CfButton variant="secondary" icon="refresh" @click="fetchAll">Refresh</CfButton>
        <CfButton variant="primary" icon="plus" @click="showCreateModal = true">Create bucket</CfButton>
      </div>
    </div>

    <!-- 2-Column Asymmetric Layout (Buckets Table + S3 Credentials Sidebar) -->
    <div class="cf-storage-grid">
      <!-- Left: Buckets List -->
      <div class="cf-storage-main">
        <CfCard :title="`Buckets (${filteredBuckets.length})`" body-class="storage-table-body">
          <template #actions>
            <div style="width: 240px;">
              <CfInput v-model="searchQuery" placeholder="Search buckets..." prefix-icon="search" clearable />
            </div>
          </template>

          <div v-if="loading" style="padding: 2rem; display: flex; justify-content: center;">
            <CfLoader size="26" />
          </div>

          <CfEmptyState
            v-else-if="!filteredBuckets.length"
            title="No storage buckets found"
            description="Create an S3 bucket to store backups, media assets, or application artifacts."
            action-text="Create bucket"
            action-icon="plus"
            @action="showCreateModal = true"
          />

          <CfTable v-else :columns="bucketColumns" :data="filteredBuckets" row-key="name">
            <template #col-name="{ row }">
              <router-link
                :to="`/storage/${row.name}`"
                style="font-weight: 600; color: #0051c3; text-decoration: none; display: inline-flex; align-items: center; gap: 0.45rem;"
              >
                <CfIcon name="r2" size="15" style="color: #f6821f;" />
                <span>{{ row.name }}</span>
              </router-link>
            </template>
            <template #col-region="{ val }">
              <CfBadge variant="neutral" size="sm">{{ val }}</CfBadge>
            </template>
            <template #actions="{ row }">
              <div style="display: flex; gap: 0.35rem; justify-content: flex-end;">
                <router-link :to="`/storage/${row.name}`" style="text-decoration: none;">
                  <CfButton variant="secondary" size="sm">Browse</CfButton>
                </router-link>
                <CfButton variant="danger-outline" size="sm" icon="trash" icon-only title="Delete bucket" @click="promptDelete(row.name)" />
              </div>
            </template>
          </CfTable>
        </CfCard>

        <!-- Multi-Tenant IAM Access Keys & Scoped Policies -->
        <CfCard
          title="Tenant IAM Access Keys & Policies"
          subtitle="Issue scoped S3 credentials restricted by bucket, key prefix, and read/write/delete permissions"
          style="margin-top: 1.25rem;"
        >
          <template #actions>
            <CfButton variant="secondary" size="sm" icon="plus" @click="showCreateIamModal = true">
              Create IAM Key
            </CfButton>
          </template>

          <CfEmptyState
            v-if="!iamKeys.length"
            title="No tenant IAM keys issued"
            description="Create scoped S3 keys for individual apps or CI pipelines instead of sharing the root S3 keypair."
            action-text="Create IAM Key"
            action-icon="plus"
            @action="showCreateIamModal = true"
          />

          <CfTable v-else :columns="iamColumns" :data="iamKeys" row-key="key_id">
            <template #col-access_key="{ val }">
              <code style="font-size: 0.82rem;">{{ val }}</code>
            </template>
            <template #col-bucket_scope="{ row }">
              <span><code>{{ row.bucket_scope }}</code><span v-if="row.prefix_scope">/{{ row.prefix_scope }}*</span></span>
            </template>
            <template #col-permissions="{ val }">
              <CfBadge variant="info" size="sm">{{ val }}</CfBadge>
            </template>
            <template #actions="{ row }">
              <CfButton variant="danger-outline" size="sm" @click="revokeIamKey(row.key_id)">Revoke</CfButton>
            </template>
          </CfTable>
        </CfCard>
      </div>

      <!-- Right: S3 API Credentials & Endpoint (Using CfCopyField from cloudflare-ui) -->
      <div class="cf-storage-sidebar">
        <CfCard title="S3 API Credentials" subtitle="Use with AWS CLI, rclone, or any S3 SDK">
          <div style="display: flex; flex-direction: column; gap: 0.85rem;">
            <CfCopyField label="S3 Endpoint" :value="fullEndpoint" />
            <CfCopyField label="Region" :value="creds.region || 'auto'" />
            <CfCopyField label="Access Key ID" :value="creds.access_key || 'Loading...'" />
            <CfCopyField label="Secret Access Key" :value="creds.secret_key || 'Loading...'" secret />
          </div>
        </CfCard>

        <CfCard title="Storage Summary">
          <div class="cf-storage-kv">
            <div class="cf-storage-kv-row">
              <span>Total Buckets</span>
              <strong>{{ buckets.length }}</strong>
            </div>
            <div class="cf-storage-kv-row">
              <span>Protocol</span>
              <code>AWS SigV4 / S3</code>
            </div>
          </div>
        </CfCard>
      </div>
    </div>

    <!-- Create Bucket Modal -->
    <CfModal v-model="showCreateModal" title="Create S3 Bucket" size="sm">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <CfInput
          v-model="newBucketName"
          label="Bucket Name"
          placeholder="e.g. app-backups"
          hint="Lowercase letters, numbers, and hyphens only."
          required
        />
        <CfAlert v-if="createError" variant="danger" title="Error">{{ createError }}</CfAlert>
      </div>
      <template #footer>
        <CfButton variant="outline" @click="showCreateModal = false">Cancel</CfButton>
        <CfButton variant="primary" :loading="creating" @click="createBucket">Create bucket</CfButton>
      </template>
    </CfModal>

    <!-- Delete Confirmation Dialog (Copied from cloudflare-ui) -->
    <CfConfirmDialog
      v-model="showDeleteDialog"
      title="Delete Bucket"
      :message="`This will permanently delete bucket '${bucketToDelete}' and all objects inside it.`"
      :confirm-text="bucketToDelete"
      confirm-label="Delete Bucket"
      danger
      @confirm="confirmDelete"
    />

    <!-- Create Tenant IAM Key Modal -->
    <CfModal v-model="showCreateIamModal" title="Create Scoped IAM Key" size="sm">
      <div style="display: flex; flex-direction: column; gap: 0.9rem;">
        <CfInput v-model="iamForm.label" label="Key Label" placeholder="e.g. backups-worker" required />
        <div>
          <label style="display: block; font-size: 0.85rem; font-weight: 500; margin-bottom: 0.3rem;">Bucket Scope</label>
          <select v-model="iamForm.bucket_scope" class="cf-iam-select">
            <option value="*">All Buckets (*)</option>
            <option v-for="b in buckets" :key="b.name" :value="b.name">{{ b.name }}</option>
          </select>
        </div>
        <CfInput v-model="iamForm.prefix_scope" label="Prefix Scope (optional)" placeholder="e.g. prod/" />
        <div>
          <label style="display: block; font-size: 0.85rem; font-weight: 500; margin-bottom: 0.3rem;">Permissions</label>
          <select v-model="iamForm.permissions" class="cf-iam-select">
            <option value="read">Read Only (GET / HEAD)</option>
            <option value="read,write">Read + Write (GET / PUT)</option>
            <option value="read,write,delete">Read + Write + Delete</option>
            <option value="admin">Full Admin</option>
          </select>
        </div>
        <div v-if="createdIamSecret" style="display: flex; flex-direction: column; gap: 0.65rem; padding-top: 0.5rem; border-top: 1px solid #e4e4e7;">
          <CfAlert variant="success" title="IAM Key Created">Copy the Secret Key now — it will not be shown again.</CfAlert>
          <CfCopyField label="Access Key ID" :value="createdIamSecret.access_key" />
          <CfCopyField label="Secret Key" :value="createdIamSecret.secret_key" secret />
        </div>
      </div>
      <template #footer>
        <CfButton variant="outline" @click="closeIamModal">Close</CfButton>
        <CfButton v-if="!createdIamSecret" variant="primary" :loading="creatingIam" @click="createIamKey">Issue Key</CfButton>
      </template>
    </CfModal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { apiFetch } from '@/api/client.js'

const buckets = ref([])
const creds = ref({ access_key: '', secret_key: '', endpoint: '/s3', region: 'auto' })
const loading = ref(true)
const searchQuery = ref('')
const showCreateModal = ref(false)
const newBucketName = ref('')
const creating = ref(false)
const createError = ref('')

const showDeleteDialog = ref(false)
const bucketToDelete = ref('')

const iamKeys = ref([])
const showCreateIamModal = ref(false)
const creatingIam = ref(false)
const createdIamSecret = ref(null)
const iamForm = ref({ label: '', bucket_scope: '*', prefix_scope: '', permissions: 'read' })

const iamColumns = [
  { key: 'label', label: 'Label' },
  { key: 'access_key', label: 'Access Key' },
  { key: 'bucket_scope', label: 'Scope' },
  { key: 'permissions', label: 'Permissions', width: '140px' }
]

const bucketColumns = [
  { key: 'name', label: 'Bucket Name' },
  { key: 'region', label: 'Region', width: '110px' },
  { key: 'created', label: 'Created', width: '160px' }
]

const fullEndpoint = computed(() => {
  const ep = creds.value.endpoint || '/s3'
  if (ep.startsWith('http')) return ep
  return `${window.location.origin}${ep}`
})

const filteredBuckets = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return buckets.value
  return buckets.value.filter(b => b.name.toLowerCase().includes(q))
})

function parseS3BucketsXml(xmlText) {
  const parser = new DOMParser()
  const doc = parser.parseFromString(xmlText, 'application/xml')
  const nodes = Array.from(doc.querySelectorAll('Bucket'))
  return nodes.map(n => {
    const name = n.querySelector('Name')?.textContent || ''
    const rawDate = n.querySelector('CreationDate')?.textContent || ''
    const created = rawDate ? new Date(rawDate).toLocaleDateString() : 'Active'
    return { name, region: creds.value.region || 'auto', created }
  }).filter(b => b.name)
}

async function fetchAll() {
  loading.value = true
  try {
    const [s3Res, credRes, iamRes] = await Promise.all([
      apiFetch('/s3/'),
      apiFetch('/api/v1/storage/credentials'),
      apiFetch('/api/v1/storage/iam-keys')
    ])
    if (credRes.ok) {
      creds.value = await credRes.json()
    }
    if (s3Res.ok) {
      const xml = await s3Res.text()
      buckets.value = parseS3BucketsXml(xml)
    }
    if (iamRes.ok) {
      const data = await iamRes.json()
      iamKeys.value = data.keys || []
    }
  } catch (e) {
    console.warn('Failed to fetch S3 storage:', e)
  } finally {
    loading.value = false
  }
}

async function createBucket() {
  createError.value = ''
  const name = newBucketName.value.trim().toLowerCase()
  if (!name) {
    createError.value = 'Bucket name is required.'
    return
  }
  creating.value = true
  try {
    const res = await apiFetch(`/s3/${encodeURIComponent(name)}`, { method: 'PUT' })
    if (!res.ok) {
      const txt = await res.text()
      throw new Error(txt || `HTTP ${res.status}`)
    }
    showCreateModal.value = false
    newBucketName.value = ''
    await fetchAll()
  } catch (e) {
    createError.value = e.message
  } finally {
    creating.value = false
  }
}

function promptDelete(name) {
  bucketToDelete.value = name
  showDeleteDialog.value = true
}

async function confirmDelete() {
  if (!bucketToDelete.value) return
  try {
    await apiFetch(`/s3/${encodeURIComponent(bucketToDelete.value)}`, { method: 'DELETE' })
    await fetchAll()
  } catch (e) {}
}

async function createIamKey() {
  creatingIam.value = true
  try {
    const res = await apiFetch('/api/v1/storage/iam-keys', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(iamForm.value)
    })
    if (res.ok) {
      const data = await res.json()
      createdIamSecret.value = { access_key: data.key.access_key, secret_key: data.secret_key }
      await fetchAll()
    }
  } finally {
    creatingIam.value = false
  }
}

async function revokeIamKey(keyId) {
  const res = await apiFetch(`/api/v1/storage/iam-keys/${encodeURIComponent(keyId)}`, { method: 'DELETE' })
  if (res.ok) await fetchAll()
}

function closeIamModal() {
  showCreateIamModal.value = false
  createdIamSecret.value = null
  iamForm.value = { label: '', bucket_scope: '*', prefix_scope: '', permissions: 'read' }
}

onMounted(fetchAll)
</script>

<style scoped>
.cf-storage-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.cf-storage-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
  align-items: start;
}

@media (max-width: 1024px) {
  .cf-storage-grid {
    grid-template-columns: 1fr;
  }
}

.cf-storage-sidebar {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

:deep(.storage-table-body) {
  padding: 0;
}

.cf-storage-kv {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.cf-iam-select {
  width: 100%;
  padding: 0.45rem 0.65rem;
  border: 1px solid #d4d4d8;
  border-radius: 6px;
  background: #ffffff;
  font-size: 0.875rem;
}

.cf-storage-kv-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  color: var(--text-color-kumo-subtle);
}

.cf-storage-kv-row strong {
  color: var(--text-color-kumo-strong);
}
</style>
