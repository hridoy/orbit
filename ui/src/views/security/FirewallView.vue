<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Security</router-link>
      <span class="divider">/</span>
      <span class="current">Firewall</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Firewall & Security Rules</h1>
        <p class="page-subtitle">Inspect incoming HTTP/TCP traffic, rate limiting, and active firewall rules on this VPS.</p>
      </div>
      <CfButton variant="primary" icon="plus" @click="alert('Create Firewall Rule wizard.')">Create Rule</CfButton>
    </div>

    <!-- Stats Grid -->
    <div class="stats-overview-grid">
      <CfStatCard label="Threats Blocked (24h)" value="12" delta="↘ 41%" delta-type="down" tooltip="SQLi & XSS mitigated" />
      <CfStatCard label="Managed Ruleset" value="OWASP" delta="v3.3.4" delta-type="up" tooltip="Sensitivity: Normal" />
      <CfStatCard label="Rate Limiter" value="Active" delta="100 req/s" delta-type="up" />
      <CfStatCard label="DDoS Mitigation" value="Enabled" delta="Layer 7" delta-type="up" />
    </div>

    <!-- Active Rulesets Table -->
    <CfCard title="Active Firewall Rulesets" body-class="fw-table-body" style="margin-top: 1.5rem;">
      <CfTable :columns="fwCols" :data="rules" row-key="name">
        <template #col-type="{ val }">
          <CfBadge variant="info" size="sm">{{ val }}</CfBadge>
        </template>
        <template #col-status="{ val }">
          <CfBadge variant="success" size="sm" dot>{{ val }}</CfBadge>
        </template>
      </CfTable>
    </CfCard>
  </div>
</template>

<script setup>
const fwCols = [
  { key: 'name', label: 'Ruleset Name' },
  { key: 'type', label: 'Type', width: '160px' },
  { key: 'action', label: 'Action', width: '180px' },
  { key: 'sensitivity', label: 'Sensitivity', width: '160px' },
  { key: 'status', label: 'Status', width: '110px' }
]

const rules = [
  { name: 'Orbit Managed Ruleset', type: 'Managed', action: 'Block', sensitivity: 'Normal (Default)', status: 'Active' },
  { name: 'OWASP Core Ruleset v3.3.4', type: 'Managed', action: 'Score Anomaly (>25)', sensitivity: 'Paranoia Level 1', status: 'Active' },
  { name: 'Block Malicious User-Agents (sqlmap, nikto)', type: 'Custom Expression', action: 'Immediate Drop', sensitivity: 'Exact Match', status: 'Active' }
]
</script>

<style scoped>
:deep(.fw-table-body) {
  padding: 0;
}
</style>
