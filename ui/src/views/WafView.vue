<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Security</router-link>
      <span class="divider">/</span>
      <span class="current">WAF & Rules</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Web Application Firewall (WAF)</h1>
        <p class="page-subtitle">Inspect incoming HTTP traffic with OWASP Core Rulesets, bot mitigation, and rate limiting.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="alert('Create Rule wizard opened.')">+ Create Firewall Rule</button>
    </div>

    <!-- Stats Grid -->
    <div class="stats-overview-grid">
      <MetricCard label="Threats Blocked (24h)" value="12" delta="↘ 41%" strokeColor="#dc2626" footer="SQLi & XSS mitigated" />
      <MetricCard label="Managed Ruleset" value="OWASP v3.3.4" strokeColor="#059669" footer="Sensitivity: Normal" />
      <MetricCard label="Bot Fight Mode" value="Active" strokeColor="#059669" footer="Automated bot challenges" />
      <MetricCard label="DDoS Mitigation" value="Autonomous" strokeColor="#059669" footer="Layer 7 mitigation" />
    </div>

    <!-- Managed Rulesets Table -->
    <div class="kumo-table-container" style="margin-top: 1.5rem;">
      <div class="table-header-title">
        <span>Active Firewall Rulesets</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Ruleset Name</th>
            <th>Type</th>
            <th>Action</th>
            <th>Sensitivity</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><strong>Cloudflare Managed Ruleset</strong></td>
            <td><span class="badge-tag-beta">Managed</span></td>
            <td>Block</td>
            <td>Normal (Default)</td>
            <td><span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span></td>
          </tr>
          <tr>
            <td><strong>OWASP Core Ruleset v3.3.4</strong></td>
            <td><span class="badge-tag-beta">Managed</span></td>
            <td>Score Anomaly (>25)</td>
            <td>Paranoia Level 1</td>
            <td><span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span></td>
          </tr>
          <tr>
            <td><strong>Block Malicious User-Agents (sqlmap, nikto)</strong></td>
            <td><span class="badge-tag-beta">Custom Expression</span></td>
            <td>Immediate Drop</td>
            <td>Exact Match</td>
            <td><span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import MetricCard from '@/components/MetricCard.vue'
</script>
