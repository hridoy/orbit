<template>
  <div>
    <div class="breadcrumbs">
      <span>Observability</span>
      <span class="divider">/</span>
      <span class="current">Account Analytics</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Account Telemetry & Edge Analytics</h1>
        <p class="page-subtitle">Real-time telemetry, HTTP response status distribution, and resource quotas across 330+ Anycast edge data centers.</p>
      </div>
      <div style="display: flex; gap: 0.5rem; align-items: center;">
        <select v-model="timeRange" class="cf-select" style="font-size: 0.82rem; height: 36px;">
          <option value="24h">Past 24 Hours</option>
          <option value="7d">Past 7 Days</option>
          <option value="30d">Past 30 Days</option>
        </select>
        <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="refreshData">↻ Refresh</button>
      </div>
    </div>

    <!-- 4 Sparkline Metric Cards -->
    <div class="stats-overview-grid" style="margin-bottom: 1.5rem;">
      <MetricCard
        label="Total Edge Requests"
        :value="analytics.totalRequests"
        :delta="analytics.totalRequestsDelta"
        sparkline="M0 16 Q 25 4, 50 14 T 100 4"
        strokeColor="var(--cf-brand-orange)"
        footer="Worldwide Anycast CDN traffic"
      />
      <MetricCard
        label="Cache Hit Ratio"
        :value="analytics.cacheHitRate"
        :delta="analytics.cacheHitRateDelta"
        sparkline="M0 18 Q 30 8, 60 12 T 100 2"
        strokeColor="#059669"
        footer="Static asset edge hits"
      />
      <MetricCard
        label="Worker Invocations"
        :value="analytics.workerInvocations"
        :delta="analytics.workerInvocationsDelta"
        sparkline="M0 14 Q 35 18, 70 6 T 100 10"
        strokeColor="#3b82f6"
        footer="V8 serverless compute executions"
      />
      <MetricCard
        label="Median CPU Time"
        :value="analytics.cpuTimeP90"
        :delta="analytics.cpuTimeDelta"
        sparkline="M0 10 Q 25 15, 60 8 T 100 12"
        strokeColor="#8b5cf6"
        footer="P90 compute execution duration"
      />
    </div>

    <!-- Telemetry Details Grid -->
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(360px, 1fr)); gap: 1.5rem; margin-bottom: 1.5rem;">
      <!-- HTTP Status Code Breakdown -->
      <div class="ui-kit-section">
        <h3 class="ui-kit-title">HTTP Response Code Breakdown</h3>
        <p class="ui-kit-desc">Edge response code distribution across all proxied zones.</p>

        <div style="display: flex; flex-direction: column; gap: 1rem; margin-top: 1.25rem;">
          <div v-for="item in analytics.statusBreakdown" :key="item.code">
            <div style="display: flex; justify-content: space-between; font-size: 0.84rem; margin-bottom: 0.35rem;">
              <strong>{{ item.code }}</strong>
              <span style="color: var(--text-color-kumo-subtle);">{{ item.percent }}% ({{ item.count }})</span>
            </div>
            <div style="height: 8px; background: var(--color-kumo-canvas); border-radius: var(--radius-full); overflow: hidden;">
              <div 
                :style="{ width: item.percent + '%', height: '100%' }"
                :class="item.color === 'fill-blue' ? 'bg-blue' : item.color === 'fill-orange' ? 'bg-orange' : 'bg-red'"
                style="border-radius: var(--radius-full); transition: width 0.5s ease;"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Resource Quotas & Limits -->
      <div class="ui-kit-section">
        <h3 class="ui-kit-title">Resource Quotas & Allocation</h3>
        <p class="ui-kit-desc">Account-level compute, storage, and database limits.</p>

        <div style="display: flex; flex-direction: column; gap: 1rem; margin-top: 1.25rem;">
          <div v-for="quota in analytics.resourceQuotas" :key="quota.name">
            <div style="display: flex; justify-content: space-between; font-size: 0.84rem; margin-bottom: 0.35rem;">
              <strong>{{ quota.name }}</strong>
              <span style="color: var(--text-color-kumo-subtle);">{{ quota.label }}</span>
            </div>
            <div style="height: 8px; background: var(--color-kumo-canvas); border-radius: var(--radius-full); overflow: hidden;">
              <div 
                :style="{ width: quota.percent + '%', height: '100%' }"
                :class="quota.color === 'fill-green' ? 'bg-green' : quota.color === 'fill-purple' ? 'bg-purple' : 'bg-blue'"
                style="border-radius: var(--radius-full); transition: width 0.5s ease;"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Global Edge Data Centers (POPs) Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Top Edge Data Centers (POPs)</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Active Anycast mesh</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>IATA Code</th>
            <th>Metro Location</th>
            <th>Region</th>
            <th>Requests Handled</th>
            <th>Cache Hit Rate</th>
            <th>Median Latency</th>
            <th>Health</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pop in popLocations" :key="pop.code">
            <td><code style="font-weight: 700;">{{ pop.code }}</code></td>
            <td><strong>{{ pop.city }}</strong></td>
            <td>{{ pop.region }}</td>
            <td>{{ pop.requests }}</td>
            <td><span style="color: #059669; font-weight: 600;">{{ pop.hitRate }}</span></td>
            <td>{{ pop.latency }}</td>
            <td>
              <span class="cf-pill active">
                <span class="status-dot"></span>
                <span>Operational</span>
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'
import MetricCard from '@/components/MetricCard.vue'

const analytics = mockData.analytics
const timeRange = ref('24h')

const popLocations = [
  { code: 'IAD', city: 'Ashburn / Washington DC', region: 'North America (ENAM)', requests: '3,840,120', hitRate: '94.2%', latency: '4 ms' },
  { code: 'NRT', city: 'Tokyo (Narita)', region: 'Asia-Pacific (APAC)', requests: '2,410,900', hitRate: '91.8%', latency: '7 ms' },
  { code: 'FRA', city: 'Frankfurt am Main', region: 'Europe (EU)', requests: '1,980,450', hitRate: '93.5%', latency: '6 ms' },
  { code: 'LHR', city: 'London Heathrow', region: 'Europe (UK)', requests: '1,420,000', hitRate: '95.1%', latency: '5 ms' },
  { code: 'SIN', city: 'Singapore', region: 'Asia-Pacific (APAC)', requests: '984,200', hitRate: '89.4%', latency: '8 ms' }
]

function refreshData() {
  alert('Analytics telemetry updated from Anycast edge.')
}
</script>

<style scoped>
.bg-blue { background-color: #3b82f6; }
.bg-orange { background-color: var(--cf-brand-orange); }
.bg-red { background-color: #ef4444; }
.bg-green { background-color: #10b981; }
.bg-purple { background-color: #8b5cf6; }
</style>
