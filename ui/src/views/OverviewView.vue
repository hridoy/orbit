<template>
  <div class="cf-overview-page">
    <!-- 1. Centered Hero Search Area -->
    <div class="cf-ov-hero">
      <div class="cf-ov-pill">
        <span class="cf-ov-dot"></span>
        <span>{{ server.hostname || 'localhost' }}<template v-if="server.os"> · {{ server.os }} {{ server.arch }}</template></span>
      </div>

      <h1 class="cf-ov-title">VPS Controller</h1>

      <button type="button" class="cf-ov-search-bar" @click="openCommandPalette">
        <span class="cf-ov-search-left">
          <CfIcon name="search" size="16" style="color: #71717a;" />
          <span class="cf-ov-search-placeholder">Search domains, apps, databases, or storage...</span>
        </span>
        <span class="cf-ov-kbd-group">
          <kbd class="cf-kbd">⌘</kbd>
          <kbd class="cf-kbd">K</kbd>
        </span>
      </button>
    </div>

    <!-- 2. 3-Column Borderless Quick Directory (Domains | Apps | Databases) -->
    <div class="cf-ov-dir-grid">
      <!-- Column 1: Domains -->
      <div class="cf-ov-col">
        <div class="cf-ov-col-head">
          <router-link to="/ssl" class="cf-ov-col-title">
            <span>Domains</span>
            <CfIcon name="chevron-right" size="13" />
          </router-link>
          <router-link to="/settings" class="cf-ov-col-link">Settings</router-link>
        </div>
        <div class="cf-ov-list">
          <p v-if="!server.domains.length && !server.panelDomain" class="cf-ov-empty">
            No domains yet. Add one in SSL & Domains.
          </p>
          <router-link v-if="server.panelDomain" to="/settings" class="cf-ov-row">
            <div class="cf-ov-row-left">
              <span class="cf-ov-dot"></span>
              <span class="cf-ov-row-name">{{ server.panelDomain }}</span>
            </div>
            <div class="cf-ov-row-right">
              <CfBadge variant="neutral" size="sm">panel</CfBadge>
              <CfIcon name="chevron-right" size="13" class="cf-ov-chevron" />
            </div>
          </router-link>
          <router-link v-for="d in server.domains" :key="d.domain" to="/ssl" class="cf-ov-row">
            <div class="cf-ov-row-left">
              <span class="cf-ov-dot"></span>
              <span class="cf-ov-row-name">{{ d.domain }}</span>
            </div>
            <div class="cf-ov-row-right">
              <CfBadge :variant="d.status === 'active' ? 'success' : 'neutral'" size="sm">{{ d.status }}</CfBadge>
              <CfIcon name="chevron-right" size="13" class="cf-ov-chevron" />
            </div>
          </router-link>
        </div>
      </div>

      <!-- Column 2: Apps -->
      <div class="cf-ov-col">
        <div class="cf-ov-col-head">
          <router-link to="/apps" class="cf-ov-col-title">
            <span>Apps</span>
            <CfIcon name="chevron-right" size="13" />
          </router-link>
          <span class="cf-ov-count">{{ server.apps.length }}</span>
        </div>
        <div class="cf-ov-list">
          <p v-if="!server.apps.length" class="cf-ov-empty">No apps deployed yet.</p>
          <router-link
            v-for="a in server.apps"
            :key="a.app_id"
            :to="`/apps/${a.app_id}`"
            class="cf-ov-row"
          >
            <div class="cf-ov-row-left">
              <CfIcon name="edgeworker" size="14" style="color: #f6821f;" />
              <span class="cf-ov-row-name">{{ a.name || a.app_id }}</span>
            </div>
            <div class="cf-ov-row-right">
              <CfBadge :variant="a.status === 'running' ? 'success' : 'neutral'" size="sm" dot>
                {{ a.status || 'stopped' }}
              </CfBadge>
              <CfIcon name="chevron-right" size="13" class="cf-ov-chevron" />
            </div>
          </router-link>
        </div>
      </div>

      <!-- Column 3: Databases -->
      <div class="cf-ov-col">
        <div class="cf-ov-col-head">
          <router-link to="/databases" class="cf-ov-col-title">
            <span>Databases</span>
            <CfIcon name="chevron-right" size="13" />
          </router-link>
          <span class="cf-ov-count">{{ server.databases.length }}</span>
        </div>
        <div class="cf-ov-list">
          <p v-if="!server.databases.length" class="cf-ov-empty">No databases provisioned yet.</p>
          <router-link
            v-for="d in server.databases"
            :key="d.db_id"
            :to="`/databases/${d.db_id}`"
            class="cf-ov-row"
          >
            <div class="cf-ov-row-left">
              <CfIcon name="d1" size="14" style="color: #0051c3;" />
              <span class="cf-ov-row-name">{{ d.name }}</span>
            </div>
            <div class="cf-ov-row-right">
              <CfBadge variant="info" size="sm">{{ d.engine }}</CfBadge>
              <CfIcon name="chevron-right" size="13" class="cf-ov-chevron" />
            </div>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 3. VPS Node Status Banner -->
    <div class="cf-ov-banner">
      <div class="cf-ov-banner-left">
        <div class="cf-ov-banner-thumb">
          <CfIcon name="server" size="22" style="color: #ffffff;" />
        </div>
        <div>
          <div style="display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
            <span class="cf-ov-banner-title">{{ server.hostname || 'Orbit VPS Node' }}</span>
            <CfBadge variant="success" size="sm" dot>Node Healthy</CfBadge>
          </div>
          <p class="cf-ov-banner-sub">
            {{ server.os || 'Linux' }} {{ server.arch || 'amd64' }} · {{ server.apps.length }} apps · {{ server.databases.length }} databases · Reverse Proxy & Auto-SSL Ready
          </p>
        </div>
      </div>
      <div style="display: flex; gap: 0.5rem; align-items: center;">
        <router-link to="/metrics" style="text-decoration: none;">
          <CfButton variant="secondary" size="md">System Metrics</CfButton>
        </router-link>
        <router-link to="/apps" style="text-decoration: none;">
          <button type="button" class="cf-ov-orange-cta">Deploy App</button>
        </router-link>
      </div>
    </div>

    <!-- 4. Live Node Telemetry Strip -->
    <div>
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.85rem;">
        <h2 class="cf-ov-section-heading">Node Telemetry</h2>
        <CfButton variant="secondary" size="sm" icon="refresh" @click="refreshAll">Refresh</CfButton>
      </div>

      <div class="stats-overview-grid">
        <CfStatCard
          label="CPU Usage"
          :value="`${telemetry.cpu_usage_pct.toFixed(1)}%`"
          delta="Host CPU"
          delta-type="up"
        />
        <CfStatCard
          label="Memory Used"
          :value="formatBytes(telemetry.ram_used_bytes)"
          :delta="`of ${formatBytes(telemetry.ram_total_bytes)}`"
          delta-type="neutral"
        />
        <CfStatCard
          label="Apps"
          :value="`${runningAppsCount} / ${server.apps.length}`"
          delta="Running"
          delta-type="up"
        />
        <CfStatCard
          label="Databases"
          :value="server.databases.length"
          delta="Provisioned"
          delta-type="neutral"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted, onUnmounted } from 'vue'
import { server, loadServer } from '@/composables/serverState.js'
import { apiFetch } from '@/api/client.js'

const telemetry = reactive({
  cpu_usage_pct: 0,
  ram_used_bytes: 0,
  ram_total_bytes: 0
})

const runningAppsCount = computed(() => server.apps.filter(a => a.status === 'running').length)

async function loadTelemetry() {
  try {
    const res = await apiFetch('/api/v1/telemetry/metrics')
    if (res.ok) {
      const m = await res.json()
      telemetry.cpu_usage_pct = Number(m.cpu_usage_pct || 0)
      telemetry.ram_used_bytes = Number(m.ram_used_bytes || 0)
      telemetry.ram_total_bytes = Number(m.ram_total_bytes || 0)
    }
  } catch {}
}

function refreshAll() {
  loadServer()
  loadTelemetry()
}

function openCommandPalette() {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }))
}

function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

let timer
onMounted(() => {
  refreshAll()
  timer = setInterval(refreshAll, 10000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.cf-overview-page {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.cf-ov-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1.25rem 0 0.25rem 0;
  gap: 0.9rem;
}

.cf-ov-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.3rem 0.85rem;
  border-radius: 9999px;
  border: 1px solid var(--color-kumo-line);
  background: var(--color-kumo-base);
  font-size: 12.5px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-ov-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #10b981;
  flex-shrink: 0;
}

.cf-ov-title {
  font-size: 28px;
  font-weight: 600;
  letter-spacing: -0.025em;
  color: var(--text-color-kumo-strong);
  margin: 0;
}

.cf-ov-search-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  max-width: 560px;
  height: 40px;
  padding: 0 0.875rem;
  border-radius: 8px;
  background: var(--color-kumo-base);
  border: 1px solid var(--kumo-control-border-default, rgba(0, 0, 0, 0.14));
  box-shadow: 0 0 0 3px var(--color-kumo-fill, rgba(0, 0, 0, 0.03));
  cursor: pointer;
  transition: all 0.15s ease;
}

.cf-ov-search-bar:hover {
  border-color: #0051c3;
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.12);
}

.cf-ov-search-left {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.cf-ov-search-placeholder {
  font-size: 13.5px;
  color: var(--text-color-kumo-subtle);
}

.cf-ov-kbd-group {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.cf-kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 5px;
  font-size: 11px;
  font-weight: 500;
  color: var(--text-color-kumo-subtle);
  background: var(--color-kumo-recessed);
  border: 1px solid var(--color-kumo-line);
  border-radius: 4px;
}

/* 3-Column Borderless Quick Directory */
.cf-ov-dir-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 2rem;
}

@media (max-width: 900px) {
  .cf-ov-dir-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
}

.cf-ov-col {
  display: flex;
  flex-direction: column;
}

.cf-ov-col-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--color-kumo-line);
}

.cf-ov-col-title {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
  text-decoration: none;
}

.cf-ov-col-title:hover {
  color: #0051c3;
}

.cf-ov-col-link,
.cf-ov-count {
  font-size: 12.5px;
  color: var(--text-color-kumo-subtle);
  text-decoration: none;
}

.cf-ov-list {
  display: flex;
  flex-direction: column;
}

.cf-ov-empty {
  font-size: 13px;
  color: var(--text-color-kumo-subtle);
  padding: 0.75rem 0;
  margin: 0;
}

.cf-ov-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.65rem 0.25rem;
  border-bottom: 1px solid var(--color-kumo-hairline);
  text-decoration: none;
  color: var(--text-color-kumo-default);
}

.cf-ov-row:hover {
  background: var(--color-kumo-fill);
}

.cf-ov-row-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.cf-ov-row-name {
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text-color-kumo-strong);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cf-ov-row-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}

.cf-ov-chevron {
  color: var(--text-color-kumo-inactive);
}

/* Node Banner */
.cf-ov-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  padding: 1rem 1.25rem;
  border-radius: 10px;
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  box-shadow: var(--shadow-kumo-card);
  flex-wrap: wrap;
}

.cf-ov-banner-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.cf-ov-banner-thumb {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f6821f, #e05300);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.cf-ov-banner-title {
  font-size: 14.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-ov-banner-sub {
  font-size: 13px;
  color: var(--text-color-kumo-subtle);
  margin: 0.2rem 0 0 0;
}

.cf-ov-orange-cta {
  height: 36px;
  padding: 0 16px;
  border-radius: 6px;
  background: #f6821f;
  color: #ffffff;
  font-size: 13.5px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  white-space: nowrap;
}

.cf-ov-orange-cta:hover {
  background: #ee730a;
}

.cf-ov-section-heading {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
  margin: 0;
}
</style>
