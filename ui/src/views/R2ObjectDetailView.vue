<template>
  <div class="cf-bucket-page cf-object-details-page">
    <!-- 1. BREADCRUMBS & TOP BAR -->
    <div class="cf-bucket-breadcrumb-bar">
      <router-link to="/r2" class="cf-bucket-breadcrumb-link">
        <CfIcon name="r2" size="18" />
        <span>R2 Object Storage</span>
      </router-link>
      <span class="cf-breadcrumb-sep">›</span>
      <router-link :to="{ path: `/r2/${bucketId}`, query: prefix ? { prefix } : {} }" class="cf-bucket-breadcrumb-link">
        <span>{{ bucketId }}</span>
      </router-link>
      <span class="cf-breadcrumb-sep">›</span>
      <div class="cf-bucket-breadcrumb-name">
        <span class="truncate">{{ objectKey }}</span>
      </div>
    </div>

    <!-- 2. SUBHEADER ACTION BAR -->
    <div class="cf-obj-detail-header-bar">
      <div class="cf-obj-header-left">
        <span class="cf-obj-key-title font-mono">{{ objectKey }}</span>
        <button class="cf-mini-icon-btn" :title="copied ? 'Copied!' : 'Copy key'" @click="copyKey">
          <CfIcon v-if="!copied" name="cis-clipboard" size="14" />
          <CfIcon v-else name="cis-checkmark-1-medium" size="14" style="color: #10b981;" />
        </button>
      </div>
      <div class="cf-obj-header-right">
        <button class="cf-btn-outline" @click="handleDownload">
          <CfIcon name="cis-arrow-wall-down" size="14" />
          <span>Download</span>
        </button>
        <button class="cf-btn-danger" @click="handleDelete">
          <CfIcon name="cis-trash-can" size="14" />
          <span>Delete</span>
        </button>
      </div>
    </div>

    <!-- 3. MAIN BODY SECTION: 3 LAYERCARDS -->
    <div class="cf-bucket-content-body">
      <!-- CARD 1: Object Details -->
      <div class="cf-layer-card">
        <div class="cf-layer-card-header">
          <span>Object Details</span>
        </div>
        <div class="cf-layer-card-body">
          <div class="cf-obj-details-grid">
            <div class="cf-obj-detail-item">
              <span class="cf-detail-label">Date Created</span>
              <span class="cf-detail-val">{{ objectInfo.dateCreated }}</span>
            </div>
            <div class="cf-obj-detail-item">
              <span class="cf-detail-label">Type</span>
              <span class="cf-detail-val">{{ objectInfo.contentType }}</span>
            </div>
            <div class="cf-obj-detail-item">
              <span class="cf-detail-label">Storage Class</span>
              <span class="cf-detail-val">{{ objectInfo.storageClass }}</span>
            </div>
            <div class="cf-obj-detail-item">
              <span class="cf-detail-label">Size</span>
              <span class="cf-detail-val">{{ objectInfo.size }}</span>
            </div>
            <div class="cf-obj-detail-item">
              <span class="cf-detail-label">ETag</span>
              <span class="cf-detail-val font-mono">{{ objectInfo.etag }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- CARD 2: Custom Metadata -->
      <div class="cf-layer-card">
        <div class="cf-layer-card-header">
          <span>Custom Metadata</span>
        </div>
        <div class="cf-layer-card-body">
          <p v-if="!objectInfo.customMetadata || Object.keys(objectInfo.customMetadata).length === 0" class="cf-no-metadata-text">
            No custom metadata set
          </p>
          <div v-else class="cf-metadata-table-wrapper">
            <table class="cf-metadata-table">
              <thead>
                <tr>
                  <th>Key</th>
                  <th>Value</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(val, key) in objectInfo.customMetadata" :key="key">
                  <td class="font-mono">{{ key }}</td>
                  <td>{{ val }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- CARD 3: Object Preview -->
      <div class="cf-layer-card">
        <div class="cf-layer-card-header">
          <span>Object Preview</span>
        </div>
        <div class="cf-layer-card-body">
          <div v-if="isImage" class="cf-preview-container">
            <div class="cf-preview-box">
              <img :src="previewSrc" :alt="objectKey" class="cf-preview-img" />
            </div>
            <span class="cf-preview-caption">{{ objectKey }} · {{ objectInfo.size }}</span>
          </div>
          <div v-else class="cf-preview-unsupported">
            <CfIcon name="cis-page-text" size="36" style="color: #71717a;" />
            <p>Preview is not available for this file type.</p>
            <button class="cf-btn-outline" @click="handleDownload">Download to view</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const bucketId = computed(() => route.params.bucket || route.params.id || 'code-collection')
const objectKey = computed(() => {
  return route.params.key ? decodeURIComponent(route.params.key) : 'thumbnails/6208/thumb.jpg'
})
const prefix = computed(() => route.query.prefix || '')

const copied = ref(false)

const copyKey = async () => {
  try {
    await navigator.clipboard.writeText(objectKey.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch (err) {
    console.warn(err)
  }
}

// Mock object metadata
const objectInfo = ref({
  dateCreated: '04 Nov 2026 15:42:34 GMT+6',
  contentType: objectKey.value.endsWith('.jpg') || objectKey.value.endsWith('.jpeg') ? 'image/jpeg' : (objectKey.value.endsWith('.png') ? 'image/png' : 'application/octet-stream'),
  storageClass: 'Standard',
  size: objectKey.value.endsWith('.jpg') ? '106.31 kB' : (objectKey.value.endsWith('.tar.gz') ? '14.24 MB' : '35.70 MB'),
  etag: '"9b24e68e21a42f5341cdb375b14197e8"',
  customMetadata: {}
})

const isImage = computed(() => {
  const k = objectKey.value.toLowerCase()
  return k.endsWith('.jpg') || k.endsWith('.jpeg') || k.endsWith('.png') || k.endsWith('.webp') || k.endsWith('.svg')
})

// Clean gradient placeholder preview image if local mock
const previewSrc = computed(() => {
  return 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="600" height="400" viewBox="0 0 600 400"><rect width="600" height="400" fill="%23f6821f"/><circle cx="300" cy="200" r="80" fill="%23ffffff" opacity="0.3"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-family="sans-serif" font-size="22" font-weight="bold" fill="%23ffffff">Cloudflare R2 Object</text></svg>'
})

const handleDownload = () => {
  alert(`Downloading "${objectKey.value}" from Cloudflare R2 edge...`)
}

const handleDelete = () => {
  if (confirm(`Are you sure you want to permanently delete "${objectKey.value}"?`)) {
    alert(`Object "${objectKey.value}" deleted.`)
    router.push({
      path: `/r2/${bucketId.value}`,
      query: prefix.value ? { prefix: prefix.value } : {}
    })
  }
}
</script>

<style scoped>
.cf-object-details-page {
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

.cf-breadcrumb-sep {
  color: #a1a1aa;
  font-size: 13px;
}

.cf-obj-detail-header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 2rem;
  background-color: var(--kumo-base, #ffffff);
  border-bottom: 1px solid var(--color-kumo-line, #e2e8f0);
}

[data-theme="dark"] .cf-obj-detail-header-bar {
  background-color: #0b0f19;
  border-bottom-color: #1e293b;
}

.cf-obj-header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.cf-obj-key-title {
  font-size: 14px;
  font-weight: 500;
  color: #18181b;
}

[data-theme="dark"] .cf-obj-key-title {
  color: #f8fafc;
}

.cf-obj-header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

/* LayerCard Styling matching Cloudflare Kumo LayerCard */
.cf-layer-card {
  display: flex;
  flex-direction: column;
  width: 100%;
  border-radius: 8px;
  background-color: #ffffff;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

[data-theme="dark"] .cf-layer-card {
  background-color: #0c0f17;
  border-color: #1e293b;
}

.cf-layer-card-header {
  padding: 0.875rem 1.25rem;
  font-size: 14px;
  font-weight: 500;
  color: #71717a;
  background-color: #ffffff;
  border-bottom: 1px solid #f4f4f5;
}

[data-theme="dark"] .cf-layer-card-header {
  background-color: #0c0f17;
  border-bottom-color: #1e293b;
  color: #a1a1aa;
}

.cf-layer-card-body {
  padding: 1.25rem;
  background-color: #ffffff;
}

[data-theme="dark"] .cf-layer-card-body {
  background-color: #070a12;
}

.cf-obj-details-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20px 80px;
}

.cf-obj-detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 140px;
}

.cf-detail-label {
  font-size: 13px;
  font-weight: 500;
  color: #71717a;
}

.cf-detail-val {
  font-size: 14px;
  font-weight: 500;
  color: #18181b;
}

[data-theme="dark"] .cf-detail-val {
  color: #f8fafc;
}

.cf-no-metadata-text {
  font-size: 13px;
  color: #71717a;
  margin: 0;
}

.cf-preview-container {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  max-width: 480px;
}

.cf-preview-box {
  width: 100%;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  background-color: #fafafa;
}

[data-theme="dark"] .cf-preview-box {
  border-color: #1e293b;
  background-color: #151a26;
}

.cf-preview-img {
  display: block;
  width: 100%;
  height: auto;
  object-fit: contain;
}

.cf-preview-caption {
  font-size: 12px;
  color: #71717a;
  font-family: var(--font-mono);
}

.cf-preview-unsupported {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2.5rem;
  gap: 0.75rem;
  color: #71717a;
  font-size: 13px;
}

/* Button styles */
.cf-btn-outline {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  height: 32px;
  padding: 0 12px;
  font-size: 13px;
  font-weight: 500;
  color: #18181b;
  background-color: #ffffff;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

[data-theme="dark"] .cf-btn-outline {
  background-color: #070a12;
  border-color: #27272a;
  color: #f8fafc;
}

.cf-btn-outline:hover {
  background-color: #f4f4f5;
}

.cf-btn-danger {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  height: 32px;
  padding: 0 12px;
  font-size: 13px;
  font-weight: 500;
  color: #ffffff;
  background-color: #dc2626;
  border: 1px solid #b91c1c;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.cf-btn-danger:hover {
  background-color: #b91c1c;
}
</style>
