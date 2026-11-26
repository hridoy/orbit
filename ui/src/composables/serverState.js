import { reactive } from 'vue'
import { api } from '@/composables/api.js'

// Live server state from the agent. Nothing here is invented: a list stays empty until the
// agent reports it.
export const server = reactive({
  loaded: false,
  hostname: '',
  nodeId: '',
  os: '',
  arch: '',
  panelDomain: '',
  apps: [],
  databases: [],
  domains: [],
})

export async function loadServer() {
  const results = await Promise.allSettled([
    api('/api/v1/node'),
    api('/api/v1/apps'),
    api('/api/v1/databases'),
    api('/api/v1/tls/domains'),
    api('/api/v1/server/settings'),
  ])
  const [node, apps, dbs, domains, settings] = results.map((r) => (r.status === 'fulfilled' ? r.value : null))
  if (node) {
    server.hostname = node.hostname || server.hostname
    server.nodeId = node.node_id || server.nodeId
    server.os = node.os || ''
    server.arch = node.arch || ''
  }
  if (Array.isArray(apps)) server.apps = apps
  if (Array.isArray(dbs)) server.databases = dbs
  if (Array.isArray(domains)) server.domains = domains
  if (settings) server.panelDomain = settings.panel_domain || ''
  server.loaded = true
}
