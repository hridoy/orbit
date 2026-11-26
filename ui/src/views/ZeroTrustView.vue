<template>
  <div>
    <div class="breadcrumbs">
      <span>Security</span>
      <span class="divider">/</span>
      <span class="current">Zero Trust Center</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Cloudflare Zero Trust & Tunnels</h1>
        <p class="page-subtitle">Identity-aware access proxy, sovereign Cloudflare Tunnels (cloudflared), and private subnet CIDR routing without open inbound ports.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showCreateModal = true">+ Create Tunnel</button>
    </div>

    <!-- Quick Stats Cards -->
    <div class="stats-overview-grid" style="margin-bottom: 1.5rem;">
      <CfStatCard
        label="Active Tunnels"
        value="2"
        unit="Connected"
        delta-type="up"
        delta="Online"
        tooltip="Sovereign outbound daemon"
      />
      <CfStatCard
        label="Access Applications"
        value="4"
        unit="Apps"
        delta-type="up"
        delta="Secured"
        tooltip="Identity-aware edge gates"
      />
      <CfStatCard
        label="DNS Security Gateway"
        value="100%"
        delta-type="up"
        delta="Protected"
        tooltip="Malware & C2 filtering"
      />
      <CfStatCard
        label="Active WARP Clients"
        value="18"
        unit="Enrolled"
        delta-type="up"
        delta="WireGuard"
        tooltip="WireGuard kernel protocol"
      />
    </div>

    <!-- Tunnels List Table -->
    <div class="kumo-table-container" style="margin-bottom: 2rem;">
      <div class="table-header-title">
        <span>Cloudflare Tunnels (cloudflared)</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Outbound-only encrypted tunnels</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Tunnel Name</th>
            <th>Tunnel ID</th>
            <th>Connector Status</th>
            <th>Ingress Routes</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tunnels" :key="t.id">
            <td>
              <div style="display: flex; align-items: center; gap: 0.5rem;">
                <span style="font-size: 1.1rem;">🚇</span>
                <strong>{{ t.name }}</strong>
              </div>
            </td>
            <td><code>{{ t.id }}</code></td>
            <td>
              <span class="cf-pill" :class="t.status === 'Active' ? 'active' : 'stopped'">
                <span class="status-dot"></span>
                <span>{{ t.status }}</span>
              </span>
            </td>
            <td>
              <span class="badge-tag-beta">{{ t.route }}</span>
            </td>
            <td>{{ t.created }}</td>
            <td>
              <div style="display: flex; gap: 0.35rem;">
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="copyToken(t.token)">Copy Token</button>
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="configureTunnel(t)">Configure</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Access Applications Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Protected Access Applications</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">SAML 2.0 / OIDC Protected</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Application</th>
            <th>Internal Destination</th>
            <th>Identity Policy</th>
            <th>Session Duration</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="app in accessApps" :key="app.name">
            <td><strong>{{ app.name }}</strong></td>
            <td><code style="color: var(--cf-brand-orange);">{{ app.domain }}</code></td>
            <td>{{ app.policy }}</td>
            <td>{{ app.session }}</td>
            <td>
              <span class="cf-pill active">
                <span class="status-dot"></span>
                <span>Active</span>
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Tunnel Modal -->
    <div v-if="showCreateModal" class="cf-modal-backdrop" @click.self="showCreateModal = false">
      <div class="cf-modal-card" style="max-width: 540px;">
        <div class="cf-modal-header">
          <div class="cf-modal-title">Create Cloudflare Tunnel</div>
          <button class="cf-modal-close-btn" @click="showCreateModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>Tunnel Name</label>
            <input v-model="newTunnelName" type="text" class="cf-input" placeholder="e.g. k8s-production-edge" />
          </div>
          <p style="font-size: 0.82rem; color: var(--text-color-kumo-subtle);">
            After creating the tunnel, you will receive a single <code>cloudflared service install</code> command to run on your private server or container.
          </p>
        </div>
        <div class="cf-modal-footer">
          <button class="btn-cf btn-cf-secondary" @click="showCreateModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="createTunnel">Create & Generate Token</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const showCreateModal = ref(false)
const newTunnelName = ref('')

const tunnels = ref([
  {
    name: 'orbit-cluster-tunnel',
    id: 'f941029c-8821-4f10-9112-aa00119944bb',
    status: 'Active',
    route: '10.0.0.0/16 (Private VPC)',
    created: 'Oct 24, 2026',
    token: 'eyJhIjoiMTEwMjk0YWZkYiIsInQiOiJmOTQxMDI5Yy04ODIxLTRmMTAtOTExMi1hYTAwMTE5OTQ0YmIifQ=='
  },
  {
    name: 'northlab-internal-services',
    id: 'a8819200-11ff-4848-bc01-998822331100',
    status: 'Active',
    route: 'internal.northlab.dev',
    created: 'Nov 01, 2026',
    token: 'eyJhIjoiMjIwMzgxM2ZhYiIsInQiOiJhODgxOTIwMC0xMWZmLTQ4NDgtYmMwMS05OTg4MjIzMzExMDAifQ=='
  }
])

const accessApps = [
  { name: 'Internal Admin Dashboard', domain: 'admin.northlab.dev', policy: 'Email ending in @northlab.dev + MFA', session: '24 hours' },
  { name: 'Grafana Observability', domain: 'grafana.sovereign-app.dev', policy: 'Engineers & DevOps Group', session: '12 hours' },
  { name: 'K8s Cluster API Ingress', domain: 'k8s.orbit-cluster.internal', policy: 'Service Token + Certificate pinning', session: '1 hour' },
  { name: 'Git Repository Mirror', domain: 'git.northlab.dev', policy: 'Hardware Security Key (FIDO2)', session: '4 hours' }
]

function copyToken(token) {
  navigator.clipboard.writeText(`cloudflared service install ${token}`)
  alert('Connector token copied to clipboard! Run on your host:\ncloudflared service install <token>')
}

function configureTunnel(t) {
  alert(`Opening ingress rules for tunnel: ${t.name}`)
}

function createTunnel() {
  if (!newTunnelName.value) return
  tunnels.value.push({
    name: newTunnelName.value,
    id: 'cf_tun_' + Math.random().toString(36).substring(2, 10),
    status: 'Active',
    route: '0.0.0.0/0',
    created: 'Just now',
    token: 'eyJhIjoi' + Math.random().toString(36).substring(2, 15)
  })
  newTunnelName.value = ''
  showCreateModal.value = false
}
</script>
