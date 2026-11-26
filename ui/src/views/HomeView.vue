<template>
  <div class="cf-account-home">
    <!-- 1. Centered Hero Search Area (VPS Controller Command Center) -->
    <div class="cf-home-hero">
      <div class="cf-hero-pill">
        <span class="cf-dir-status-dot"></span>
        <span>{{ mockData.account.hostname }} · {{ mockData.account.plan }}</span>
      </div>

      <h1 class="cf-hero-title">VPS Controller</h1>

      <button type="button" class="cf-hero-search-bar" @click="openCommandPalette">
        <span class="cf-hero-search-left">
          <CfIcon name="search" size="16" style="color: #71717a;" />
          <span class="cf-hero-search-placeholder">Search domains, apps, databases, or storage buckets...</span>
        </span>
        <span class="cf-hero-kbd-group">
          <kbd class="cf-kbd">⌘</kbd>
          <kbd class="cf-kbd">K</kbd>
        </span>
      </button>
    </div>

    <!-- 2. 3-Column Borderless Quick Directory (Domains | Apps | Storage & DB) -->
    <div class="cf-quick-dir-grid">
      <!-- Column 1: Domains & Reverse Proxy -->
      <div class="cf-dir-col">
        <div class="cf-dir-col-head">
          <router-link to="/zones/northlab.dev" class="cf-dir-col-title">
            <span>Domains</span>
            <CfIcon name="chevron-right" size="13" />
          </router-link>
          <CfButton variant="ghost" size="sm" icon="more-horizontal" icon-only title="Domain options" />
        </div>
        <div class="cf-dir-list">
          <router-link
            v-for="d in mockData.domains.slice(0, 4)"
            :key="d.id"
            to="/zones/northlab.dev"
            class="cf-dir-row"
          >
            <div class="cf-dir-row-left">
              <span class="cf-dir-status-dot"></span>
              <span class="cf-dir-row-name">{{ d.name }}</span>
            </div>
            <div class="cf-dir-row-right">
              <CfBadge variant="neutral" size="sm">{{ d.plan }}</CfBadge>
              <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
            </div>
          </router-link>
        </div>
      </div>

      <!-- Column 2: Apps & Services -->
      <div class="cf-dir-col">
        <div class="cf-dir-col-head">
          <router-link to="/workers" class="cf-dir-col-title">
            <span>Apps & Services</span>
            <CfIcon name="chevron-right" size="13" />
          </router-link>
          <CfButton variant="ghost" size="sm" icon="more-horizontal" icon-only title="Service options" />
        </div>
        <div class="cf-dir-list">
          <router-link
            v-for="w in mockData.workers"
            :key="w.app_id"
            :to="`/workers/${w.app_id}`"
            class="cf-dir-row"
          >
            <div class="cf-dir-row-left">
              <CfIcon name="edgeworker" size="14" style="color: #f6821f;" />
              <span class="cf-dir-row-name">{{ w.name }}</span>
            </div>
            <div class="cf-dir-row-right">
              <span class="cf-dir-meta">{{ w.invocations_24h }} req</span>
              <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
            </div>
          </router-link>
          <router-link to="/waf" class="cf-dir-row">
            <div class="cf-dir-row-left">
              <CfIcon name="security-waf" size="14" style="color: #10b981;" />
              <span class="cf-dir-row-name">UFW Firewall & Proxy</span>
            </div>
            <div class="cf-dir-row-right">
              <span class="cf-dir-meta">Active</span>
              <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
            </div>
          </router-link>
        </div>
      </div>

      <!-- Column 3: Storage & Database -->
      <div class="cf-dir-col">
        <div class="cf-dir-col-head">
          <span class="cf-dir-col-title" style="cursor: default;">
            <span>Storage & Database</span>
          </span>
          <CfButton variant="ghost" size="sm" icon="more-horizontal" icon-only title="Storage & DB options" />
        </div>
        <div class="cf-dir-list">
          <router-link to="/r2" class="cf-dir-row">
            <div class="cf-dir-row-left">
              <span class="cf-dir-prefix">Storage /</span>
              <span class="cf-dir-row-name">Object Storage (S3)</span>
            </div>
            <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
          </router-link>
          <router-link to="/d1" class="cf-dir-row">
            <div class="cf-dir-row-left">
              <span class="cf-dir-prefix">Database /</span>
              <span class="cf-dir-row-name">SQL Databases</span>
            </div>
            <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
          </router-link>
          <router-link to="/ssl" class="cf-dir-row">
            <div class="cf-dir-row-left">
              <span class="cf-dir-prefix">Domain /</span>
              <span class="cf-dir-row-name">SSL Certificates</span>
            </div>
            <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
          </router-link>
          <router-link to="/zero-trust" class="cf-dir-row">
            <div class="cf-dir-row-left">
              <span class="cf-dir-prefix">Security /</span>
              <span class="cf-dir-row-name">SSH & Tunnels</span>
            </div>
            <CfIcon name="chevron-right" size="13" class="cf-dir-chevron" />
          </router-link>
        </div>
      </div>
    </div>

    <!-- 3. VPS Node Status Banner -->
    <div class="cf-announcement-banner">
      <div class="cf-announce-left">
        <div class="cf-announce-thumb">
          <CfIcon name="server" size="22" style="color: #ffffff;" />
        </div>
        <div>
          <div class="cf-announce-title-row">
            <span class="cf-announce-title">{{ mockData.account.name }}</span>
            <CfBadge variant="success" size="sm" dot>{{ mockData.account.status }}</CfBadge>
          </div>
          <p class="cf-announce-sub">
            {{ mockData.account.region }} · Uptime: <strong>{{ mockData.account.uptime }}</strong> · Reverse Proxy & Auto-SSL Active
          </p>
        </div>
      </div>
      <div class="cf-announce-right" style="display: flex; gap: 0.5rem;">
        <router-link to="/analytics" style="text-decoration: none;">
          <CfButton variant="secondary" size="md">System Metrics</CfButton>
        </router-link>
        <router-link to="/workers" style="text-decoration: none;">
          <button type="button" class="cf-orange-cta-btn">Deploy App</button>
        </router-link>
      </div>
    </div>

    <!-- 4. Node Telemetry Header & Sparkline Metric Grid -->
    <div class="cf-analytics-section">
      <div class="cf-analytics-head">
        <h2 class="cf-section-heading">Node Telemetry</h2>
        <div class="cf-analytics-controls">
          <div style="width: 160px;">
            <CfSelect v-model="analyticsRange" :options="['Last 24 hours', 'Last 1 hour', 'Last 7 days']" size="md" />
          </div>
          <CfButton variant="secondary" size="md" icon="refresh" icon-only title="Refresh telemetry" />
        </div>
      </div>

      <div class="stats-overview-grid">
        <CfStatCard
          label="CPU Load"
          :value="mockData.account.cpuUsage"
          :delta="mockData.account.cpuDelta"
          delta-type="down"
          tooltip="Average CPU utilization across 8 vCPU cores"
          :sparkline-data="[32, 28, 30, 26, 29, 25, 24]"
        />
        <CfStatCard
          label="Memory (RAM)"
          :value="mockData.account.ramUsage"
          :delta="mockData.account.ramDelta"
          delta-type="neutral"
          tooltip="Active physical memory consumption"
          :sparkline-data="[5.8, 6.0, 6.2, 6.1, 6.3, 6.4]"
        />
        <CfStatCard
          label="NVMe Disk"
          :value="mockData.account.diskUsage"
          :delta="mockData.account.diskDelta"
          delta-type="up"
          tooltip="Root volume + S3 Object Storage usage"
          :sparkline-data="[62, 63, 64, 65, 67, 68]"
        />
        <CfStatCard
          label="Network I/O"
          :value="mockData.account.networkIo"
          :delta="mockData.account.networkDelta"
          delta-type="up"
          tooltip="Combined eth0 ingress + egress throughput"
          :sparkline-data="[95, 110, 105, 128, 134, 142]"
        />
      </div>
    </div>

    <!-- 5. Running Apps & Services -->
    <div>
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.85rem;">
        <h2 class="cf-section-heading">Running Apps & Services ({{ mockData.workers.length }})</h2>
        <router-link to="/workers" style="text-decoration: none;">
          <CfButton variant="secondary" size="sm">Manage All Services →</CfButton>
        </router-link>
      </div>

      <div class="cf-item-card-list">
        <router-link
          v-for="w in mockData.workers"
          :key="w.app_id"
          :to="`/workers/${w.app_id}`"
          class="cf-item-card"
          style="text-decoration: none;"
        >
          <div class="cf-item-card-top">
            <div class="cf-item-card-info">
              <div class="cf-item-card-header-row">
                <span class="cf-item-card-title">{{ w.name }}</span>
                <CfBadge variant="success" dot size="sm">{{ w.status }}</CfBadge>
              </div>
              <div style="display: flex; align-items: center; gap: 0.5rem; margin-top: 0.2rem;">
                <span style="font-family: var(--font-mono); font-size: 0.8rem; color: #0051c3;">{{ w.custom_domain }}</span>
                <span style="color: #a1a1aa;">•</span>
                <span class="cf-badge-repo">
                  <CfIcon name="github" size="12" />
                  <span>{{ w.repo }} ({{ w.branch }})</span>
                </span>
              </div>
            </div>
            <div class="cf-item-card-actions">
              <span style="font-size: 0.78rem; color: #71717a;">Deployed {{ w.last_deployed }}</span>
              <CfButton variant="ghost" size="sm" icon="more-horizontal" icon-only />
            </div>
          </div>

          <div class="cf-item-card-bottom">
            <div class="cf-item-card-meta">
              <span class="cf-stat-chip">Release: <strong>{{ w.current_version }}</strong></span>
            </div>
            <div class="cf-item-card-stats">
              <span class="cf-stat-chip">
                <CfIcon name="bar-chart-2" size="13" />
                <span>Requests: <strong>{{ w.invocations_24h }}</strong></span>
              </span>
              <span class="cf-stat-chip">
                <CfIcon name="cpu" size="13" />
                <span>Latency: <strong>{{ w.median_cpu_time }}</strong></span>
              </span>
              <span class="cf-stat-chip">
                <CfIcon name="alert-triangle" size="13" />
                <span>Errors: <strong>{{ w.errors_24h }}</strong></span>
              </span>
            </div>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const analyticsRange = ref('Last 24 hours')

function openCommandPalette() {
  window.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', metaKey: true }))
}
</script>

<style scoped>
.cf-account-home {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

/* Hero Search Section */
.cf-home-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1.5rem 0 0.5rem 0;
  gap: 1rem;
}

.cf-hero-pill {
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
  box-shadow: var(--shadow-kumo-xs);
}

.cf-hero-title {
  font-size: 28px;
  font-weight: 600;
  letter-spacing: -0.025em;
  color: var(--text-color-kumo-strong);
  margin: 0;
}

.cf-hero-search-bar {
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

.cf-hero-search-bar:hover {
  border-color: #0051c3;
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.12);
}

.cf-hero-search-left {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.cf-hero-search-placeholder {
  font-size: 13.5px;
  color: var(--text-color-kumo-subtle);
}

.cf-hero-kbd-group {
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
  font-family: var(--font-sans);
  font-weight: 500;
  color: var(--text-color-kumo-subtle);
  background: var(--color-kumo-recessed);
  border: 1px solid var(--color-kumo-line);
  border-radius: 4px;
}

/* 3-Column Borderless Quick Directory */
.cf-quick-dir-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 2rem;
}

@media (max-width: 900px) {
  .cf-quick-dir-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
}

.cf-dir-col {
  display: flex;
  flex-direction: column;
}

.cf-dir-col-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--color-kumo-line);
}

.cf-dir-col-title {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
  text-decoration: none;
}

.cf-dir-col-title:hover {
  color: #0051c3;
}

.cf-dir-list {
  display: flex;
  flex-direction: column;
}

.cf-dir-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.65rem 0.25rem;
  border-bottom: 1px solid var(--color-kumo-hairline);
  text-decoration: none;
  color: var(--text-color-kumo-default);
  transition: background-color 0.12s ease;
}

.cf-dir-row:hover {
  background: var(--color-kumo-fill);
}

.cf-dir-row-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.cf-dir-status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background-color: #10b981;
  flex-shrink: 0;
}

.cf-dir-prefix {
  font-size: 13px;
  color: var(--text-color-kumo-subtle);
  white-space: nowrap;
}

.cf-dir-row-name {
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text-color-kumo-strong);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cf-dir-row-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}

.cf-dir-meta {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-dir-chevron {
  color: var(--text-color-kumo-inactive);
}

/* Announcement Banner */
.cf-announcement-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  padding: 1rem 1.25rem;
  border-radius: 10px;
  background: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  box-shadow: var(--shadow-kumo-card);
}

@media (max-width: 720px) {
  .cf-announcement-banner {
    flex-direction: column;
    align-items: flex-start;
  }
}

.cf-announce-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.cf-announce-thumb {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f6821f, #e05300);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.cf-announce-title-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.cf-announce-title {
  font-size: 14.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-announce-sub {
  font-size: 13px;
  color: var(--text-color-kumo-subtle);
  margin: 0.2rem 0 0 0;
}

.cf-orange-cta-btn {
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
  transition: background-color 0.15s ease;
}

.cf-orange-cta-btn:hover {
  background: #ee730a;
}

/* Analytics Section */
.cf-analytics-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.cf-section-heading {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
  margin: 0;
  letter-spacing: -0.015em;
}

.cf-analytics-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
</style>
