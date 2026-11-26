<template>
  <div>
    <div class="breadcrumbs">
      <span>Optimization</span>
      <span class="divider">/</span>
      <span class="current">Caching & CDN</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Configuration</h1>
        <p class="page-subtitle">Manage caching behavior, purge edge assets, and accelerate static content delivery across 330+ POPs.</p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <CfButton variant="secondary" icon="info">Documentation</CfButton>
      </div>
    </div>

    <!-- Notification Banner -->
    <CfAlert v-if="purgedNotice" variant="success" title="Purge Dispatched" dismissible style="margin-bottom: 1.5rem;" @dismiss="purgedNotice = false">
      Edge cache purge instruction broadcast to all global POPs. Propagation completes in ~5s.
    </CfAlert>

    <!-- Authentic Cloudflare Stacked Single-Setting Cards -->
    <div class="cf-caching-stack">
      <!-- 1. Purge Cache Card (with Red-Outlined Purge Everything) -->
      <CfCard>
        <div class="cf-setting-row">
          <div class="cf-setting-info">
            <h3 class="cf-setting-title">Purge Cache</h3>
            <p class="cf-setting-desc">
              Clear cached files to force Cloudflare to fetch a fresh version of those files from your web server.
            </p>
          </div>
          <div class="cf-setting-actions">
            <CfButton variant="secondary" @click="openPurgeModal('custom')">Custom Purge</CfButton>
            <CfButton variant="danger-outline" @click="openPurgeModal('everything')">Purge Everything</CfButton>
          </div>
        </div>
      </CfCard>

      <!-- 2. Caching Level -->
      <CfCard>
        <div class="cf-setting-row">
          <div class="cf-setting-info">
            <h3 class="cf-setting-title">Caching Level</h3>
            <p class="cf-setting-desc">
              Determine how much of your website's static content you want Cloudflare to cache. Increased caching can speed up page load time.
            </p>
          </div>
          <div class="cf-setting-actions" style="width: 260px;">
            <CfSelect v-model="cachingLevel" :options="cachingLevelOptions" />
          </div>
        </div>
      </CfCard>

      <!-- 3. Browser Cache TTL -->
      <CfCard>
        <div class="cf-setting-row">
          <div class="cf-setting-info">
            <h3 class="cf-setting-title">Browser Cache TTL</h3>
            <p class="cf-setting-desc">
              Determine the length of time Cloudflare instructs a visitor's browser to cache files. During this period, the browser loads the files from its local cache.
            </p>
          </div>
          <div class="cf-setting-actions" style="width: 260px;">
            <CfSelect v-model="browserTtl" :options="browserTtlOptions" />
          </div>
        </div>
      </CfCard>

      <!-- 4. Always Online -->
      <CfCard>
        <div class="cf-setting-row">
          <div class="cf-setting-info">
            <div style="display: flex; align-items: center; gap: 0.5rem;">
              <h3 class="cf-setting-title">Always Online™</h3>
              <CfBadge :variant="alwaysOnline ? 'success' : 'neutral'" size="sm" :dot="alwaysOnline">
                {{ alwaysOnline ? 'Enabled' : 'Disabled' }}
              </CfBadge>
            </div>
            <p class="cf-setting-desc">
              Keep visitors browsing your site even if your origin server goes down by serving cached copies from the Internet Archive.
            </p>
          </div>
          <div class="cf-setting-actions">
            <CfSwitch v-model="alwaysOnline" />
          </div>
        </div>
      </CfCard>

      <!-- 5. Development Mode -->
      <CfCard>
        <div class="cf-setting-row">
          <div class="cf-setting-info">
            <div style="display: flex; align-items: center; gap: 0.5rem;">
              <h3 class="cf-setting-title">Development Mode</h3>
              <CfBadge :variant="devMode ? 'warning' : 'neutral'" size="sm" :dot="devMode">
                {{ devMode ? 'Active (Bypassing)' : 'Off' }}
              </CfBadge>
            </div>
            <p class="cf-setting-desc">
              Temporarily bypass Cloudflare's cache allowing you to see changes to your origin server in real time. Automatically turns off after 3 hours.
            </p>
          </div>
          <div class="cf-setting-actions">
            <CfSwitch v-model="devMode" />
          </div>
        </div>
      </CfCard>
    </div>

    <!-- Purge Cache Modal -->
    <CfModal v-model="showPurgeModal" :title="`Purge Cache — ${mockData.activeZone.name}`" size="md">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <p style="font-size: 0.88rem; color: var(--text-color-kumo-subtle); margin: 0;">
          Clearing cached assets forces Cloudflare to re-fetch the latest versions from your origin server.
        </p>

        <CfTabs
          v-model="purgeType"
          :tabs="[
            { label: 'Purge Everything', value: 'everything' },
            { label: 'Purge by URL / Tag', value: 'custom' }
          ]"
        />

        <CfAlert v-if="purgeType === 'everything'" variant="warning" title="Notice">
          Purging everything will temporarily increase traffic load on your origin server while assets are re-cached.
        </CfAlert>

        <CfTextarea
          v-else
          v-model="customPurgeUrls"
          label="Specific URLs or Cache-Tags (one per line)"
          placeholder="https://northlab.dev/styles.css&#10;https://northlab.dev/app.js"
          :rows="3"
        />
      </div>

      <template #footer>
        <CfButton variant="outline" @click="showPurgeModal = false">Cancel</CfButton>
        <CfButton variant="primary" @click="handlePurge">Confirm Purge</CfButton>
      </template>
    </CfModal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const cachingLevel = ref(mockData.activeZone.caching.cachingLevel || 'Standard')
const browserTtl = ref(mockData.activeZone.caching.browserCacheTtl || '4 hours')
const devMode = ref(false)
const alwaysOnline = ref(true)
const showPurgeModal = ref(false)
const purgeType = ref('everything')
const customPurgeUrls = ref('')
const purgedNotice = ref(false)

const cachingLevelOptions = [
  { value: 'No Query String', label: 'No Query String' },
  { value: 'Ignore Query String', label: 'Ignore Query String' },
  { value: 'Standard', label: 'Standard' }
]

const browserTtlOptions = [
  'Respect Existing Headers',
  '2 hours',
  '4 hours',
  '8 hours',
  '1 day',
  '8 days'
]

function openPurgeModal(type) {
  purgeType.value = type
  showPurgeModal.value = true
}

function handlePurge() {
  showPurgeModal.value = false
  purgedNotice.value = true
  setTimeout(() => {
    purgedNotice.value = false
  }, 5000)
}
</script>

<style scoped>
.cf-caching-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.cf-setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
}

@media (max-width: 768px) {
  .cf-setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
}

.cf-setting-info {
  flex: 1;
  min-width: 0;
}

.cf-setting-title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-setting-desc {
  margin: 0.35rem 0 0 0;
  font-size: 13px;
  line-height: 1.5;
  color: var(--text-color-kumo-subtle);
}

.cf-setting-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}
</style>
