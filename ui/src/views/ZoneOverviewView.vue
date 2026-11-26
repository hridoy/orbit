<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Account</router-link>
      <span class="divider">/</span>
      <span class="current">{{ zone.name }}</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="font-size: 12px; font-weight: 500; color: var(--text-color-kumo-subtle); margin-bottom: 0.2rem;">Overview</div>
        <div style="display: flex; align-items: center; gap: 0.65rem; flex-wrap: wrap;">
          <h1 class="page-title">{{ zone.name }}</h1>
          <CfBadge variant="success" dot size="sm">{{ zone.status }}</CfBadge>
          <CfBadge variant="neutral" size="sm">{{ zone.plan }}</CfBadge>
        </div>
        <p class="page-subtitle">Zone-scoped edge configurations, Anycast DNS, SSL/TLS certificates, and edge caching.</p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <router-link to="/dns" style="text-decoration: none;">
          <CfButton variant="primary">Manage DNS Records</CfButton>
        </router-link>
        <CfButton variant="secondary" @click="alert('Purge Everything executed.')">Purge Cache</CfButton>
      </div>
    </div>

    <!-- Domain Sub-Navigation Tabs Bar -->
    <div class="cf-tab-nav">
      <router-link to="/zones/northlab.dev" class="cf-tab-btn active">Overview</router-link>
      <router-link to="/dns" class="cf-tab-btn">DNS Records</router-link>
      <router-link to="/ssl" class="cf-tab-btn">SSL/TLS</router-link>
      <router-link to="/waf" class="cf-tab-btn">Security (WAF)</router-link>
      <router-link to="/caching" class="cf-tab-btn">Caching</router-link>
      <router-link to="/rules" class="cf-tab-btn">Rules</router-link>
      <router-link to="/network" class="cf-tab-btn">Network</router-link>
    </div>

    <!-- Domain Overview Content Cards -->
    <div class="stats-overview-grid">
      <CfStatCard
        label="SSL Encryption Mode"
        :value="zone.ssl.mode"
        delta="Universal SSL"
        delta-type="up"
        tooltip="Global edge SSL certificate"
      />
      <CfStatCard
        label="DNS Records Count"
        :value="zone.dnsRecords.length"
        unit="Records"
        delta="Anycast"
        delta-type="neutral"
        tooltip="Anycast global DNS zone"
      />
      <CfStatCard
        label="Cache Level"
        :value="zone.caching.cachingLevel"
        delta="Standard"
        delta-type="up"
        :tooltip="`Browser TTL: ${zone.caching.browserCacheTtl}`"
      />
      <CfStatCard
        label="HTTP/3 (QUIC)"
        value="Active"
        delta="0-RTT"
        delta-type="up"
        tooltip="0-RTT Resumption Enabled"
      />
    </div>

    <!-- 2-Column Zone Overview: DNS Preview + Quick Actions & API Card -->
    <div class="cf-zone-grid">
      <!-- Left: Quick DNS Records Preview Table -->
      <CfCard :title="`Configured DNS Records for ${zone.name}`" body-class="zone-table-body">
        <template #actions>
          <router-link to="/dns" style="text-decoration: none;">
            <CfButton variant="secondary" size="sm">Open Full DNS Manager →</CfButton>
          </router-link>
        </template>

        <CfTable :columns="dnsPreviewCols" :data="zone.dnsRecords.slice(0, 4)" row-key="id">
          <template #col-type="{ val }">
            <span class="dns-type-badge">{{ val }}</span>
          </template>
          <template #col-name="{ val }">
            <strong>{{ val }}</strong>
          </template>
          <template #col-content="{ val }">
            <code style="font-family: var(--font-mono); font-size: 0.82rem;">{{ val }}</code>
          </template>
          <template #col-proxied="{ row }">
            <span class="dns-status-cloud" :class="row.proxied ? 'proxied' : 'dns-only'">
              <span>☁️</span>
              <span>{{ row.proxied ? 'Proxied' : 'DNS only' }}</span>
            </span>
          </template>
        </CfTable>
      </CfCard>

      <!-- Right: Quick Actions & API Credentials (Real Cloudflare Zone Sidebar) -->
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <CfCard title="Quick Actions">
          <div style="display: flex; flex-direction: column; gap: 0.85rem;">
            <div class="cf-qa-row">
              <span class="cf-qa-label">Under Attack Mode</span>
              <CfSwitch v-model="underAttack" />
            </div>
            <div class="cf-qa-row">
              <span class="cf-qa-label">Development Mode</span>
              <CfSwitch v-model="devMode" />
            </div>
          </div>
        </CfCard>

        <CfCard title="API">
          <div style="display: flex; flex-direction: column; gap: 0.75rem;">
            <div class="cf-api-row">
              <span class="cf-api-label">Zone ID</span>
              <div class="cf-api-val">
                <code>8a4f2c9b1e7d...309a</code>
                <CfButton variant="ghost" size="sm" icon="clipboard" icon-only title="Copy Zone ID" @click="copyId('8a4f2c9b1e7d401c92a8f0129384309a')" />
              </div>
            </div>
            <div class="cf-api-row">
              <span class="cf-api-label">Account ID</span>
              <div class="cf-api-val">
                <code>7f91c2b8a34e...81d2</code>
                <CfButton variant="ghost" size="sm" icon="clipboard" icon-only title="Copy Account ID" @click="copyId('7f91c2b8a34e9012f1a8c4e5b67881d2')" />
              </div>
            </div>
          </div>
        </CfCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const zone = mockData.activeZone
const underAttack = ref(false)
const devMode = ref(false)

const dnsPreviewCols = [
  { key: 'type', label: 'Type', width: '80px' },
  { key: 'name', label: 'Name', width: '140px' },
  { key: 'content', label: 'Content' },
  { key: 'proxied', label: 'Proxy', width: '115px' }
]

function copyId(val) {
  navigator.clipboard?.writeText(val)
}
</script>

<style scoped>
:deep(.zone-table-body) {
  padding: 0;
}

.cf-zone-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.25rem;
  margin-top: 1.5rem;
  align-items: start;
}

@media (max-width: 1024px) {
  .cf-zone-grid {
    grid-template-columns: 1fr;
  }
}

.cf-qa-row,
.cf-api-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
}

.cf-qa-label,
.cf-api-label {
  color: var(--text-color-kumo-subtle);
  font-weight: 500;
}

.cf-api-val {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-color-kumo-default);
}
</style>
