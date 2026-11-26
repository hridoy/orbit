<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Protect & Connect</router-link>
      <span class="divider">/</span>
      <span class="current">Turnstile CAPTCHA</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.75rem;">
          <h1 class="page-title">🛡️ Cloudflare Turnstile</h1>
          <span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span>
          <span class="badge-tag-beta">Privacy First</span>
        </div>
        <p class="page-subtitle">Smart CAPTCHA alternative protecting websites from bad bots with seamless zero-friction verification.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="alert('Add Turnstile Widget wizard.')">+ Add Widget</button>
    </div>

    <!-- Stats Strip -->
    <div class="stats-overview-grid">
      <MetricCard label="Total Solves (24h)" value="145,430" delta="↗ 8.4%" strokeColor="var(--cf-brand-orange)" footer="Across 3 active domains" />
      <MetricCard label="Visitor Pass Rate" value="99.8%" delta="↗ 0.2%" strokeColor="#059669" footer="Zero friction human pass" />
      <MetricCard label="Average Solve Time" value="0.4 sec" delta="↘ 0.1s" strokeColor="#0051c3" footer="Non-interactive challenge" />
      <MetricCard label="Bots Mitigated" value="290" delta="↘ 14" strokeColor="#dc2626" footer="Automated scrapers blocked" />
    </div>

    <!-- Interactive Turnstile Widget Simulator -->
    <div class="ui-kit-section" style="margin-top: 1.5rem;">
      <h3 class="ui-kit-title">Interactive Turnstile Widget Simulator</h3>
      <p class="ui-kit-desc">Experience the exact visitor verification challenge embedded in authentication forms.</p>
      <TurnstileWidget />
    </div>

    <!-- Configured Widgets Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Configured Turnstile Widgets ({{ turnstile.widgets.length }})</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Widget Name</th>
            <th>Site Key</th>
            <th>Domain Scope</th>
            <th>Challenge Mode</th>
            <th>24h Solves</th>
            <th>Pass Rate</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in turnstile.widgets" :key="w.id">
            <td><strong>{{ w.name }}</strong></td>
            <td><code>{{ w.id }}</code></td>
            <td>{{ w.domain }}</td>
            <td><span class="badge-tag-beta">{{ w.mode }}</span></td>
            <td>{{ w.solves_24h }}</td>
            <td><strong style="color: #059669;">{{ w.passRate }}</strong></td>
            <td><span class="cf-pill active"><span class="status-dot"></span><span>{{ w.status }}</span></span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { mockData } from '@/data/mockData.js'
import MetricCard from '@/components/MetricCard.vue'
import TurnstileWidget from '@/components/TurnstileWidget.vue'

const turnstile = mockData.turnstile
</script>
