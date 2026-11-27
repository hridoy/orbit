import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  // Overview
  {
    path: '/',
    name: 'Overview',
    component: () => import('@/views/OverviewView.vue'),
    meta: { title: 'Overview — Orbit VPS Controller' }
  },

  // Apps
  {
    path: '/apps',
    alias: '/workers',
    name: 'Apps',
    component: () => import('@/views/apps/AppsListView.vue'),
    meta: { title: 'Apps — Orbit VPS Controller' }
  },
  {
    path: '/apps/:id',
    alias: '/workers/:id',
    name: 'AppDetail',
    component: () => import('@/views/apps/AppDetailView.vue'),
    meta: { title: 'App Console — Orbit VPS Controller' }
  },

  // Databases
  {
    path: '/databases',
    alias: '/d1',
    name: 'Databases',
    component: () => import('@/views/databases/DatabasesListView.vue'),
    meta: { title: 'Databases — Orbit VPS Controller' }
  },
  {
    path: '/databases/:id',
    alias: '/d1/:id',
    name: 'DatabaseDetail',
    component: () => import('@/views/databases/DatabaseDetailView.vue'),
    meta: { title: 'Database Console — Orbit VPS Controller' }
  },

  // Storage (S3 Object Storage)
  {
    path: '/storage',
    alias: '/r2',
    name: 'Storage',
    component: () => import('@/views/storage/BucketsListView.vue'),
    meta: { title: 'Object Storage — Orbit VPS Controller' }
  },
  {
    path: '/storage/api-keys',
    alias: '/r2/api-tokens',
    name: 'StorageApiKeys',
    component: () => import('@/views/system/ApiKeysView.vue'),
    meta: { title: 'Storage API Keys — Orbit VPS Controller' }
  },
  {
    path: '/storage/:id',
    alias: '/r2/:id',
    name: 'BucketDetail',
    component: () => import('@/views/storage/BucketDetailView.vue'),
    meta: { title: 'Storage Bucket — Orbit VPS Controller' }
  },
  {
    path: '/storage/:bucket/objects/:key(.*)/details',
    alias: '/r2/:bucket/objects/:key(.*)/details',
    name: 'ObjectDetail',
    component: () => import('@/views/storage/ObjectDetailView.vue'),
    meta: { title: 'Object Details — Orbit VPS Controller' }
  },

  // Domains & SSL
  {
    path: '/domains',
    alias: '/zones/northlab.dev',
    name: 'DomainOverview',
    component: () => import('@/views/domains/DomainOverviewView.vue'),
    meta: { title: 'Domains — Orbit VPS Controller' }
  },
  {
    path: '/dns',
    name: 'DnsRecords',
    component: () => import('@/views/domains/DnsRecordsView.vue'),
    meta: { title: 'DNS Records — Orbit VPS Controller' }
  },
  {
    path: '/ssl',
    name: 'SslDomains',
    component: () => import('@/views/domains/SslDomainsView.vue'),
    meta: { title: 'SSL & Domains — Orbit VPS Controller' }
  },
  {
    path: '/proxy-cache',
    alias: '/caching',
    name: 'ProxyCache',
    component: () => import('@/views/domains/ProxyCacheView.vue'),
    meta: { title: 'Proxy & Cache — Orbit VPS Controller' }
  },

  // Security & Network
  {
    path: '/firewall',
    alias: '/waf',
    name: 'Firewall',
    component: () => import('@/views/security/FirewallView.vue'),
    meta: { title: 'Firewall — Orbit VPS Controller' }
  },
  {
    path: '/tunnels',
    alias: '/zero-trust',
    name: 'Tunnels',
    component: () => import('@/views/security/TunnelsView.vue'),
    meta: { title: 'SSH & Tunnels — Orbit VPS Controller' }
  },
  {
    path: '/network',
    name: 'Network',
    component: () => import('@/views/security/NetworkView.vue'),
    meta: { title: 'Network — Orbit VPS Controller' }
  },

  // System
  {
    path: '/metrics',
    alias: '/analytics',
    name: 'Metrics',
    component: () => import('@/views/system/MetricsView.vue'),
    meta: { title: 'System Metrics — Orbit VPS Controller' }
  },
  {
    path: '/users',
    alias: '/team',
    name: 'Users',
    component: () => import('@/views/system/UsersView.vue'),
    meta: { title: 'Users — Orbit VPS Controller' }
  },
  {
    path: '/api-keys',
    alias: '/tokens',
    name: 'ApiKeys',
    component: () => import('@/views/system/ApiKeysView.vue'),
    meta: { title: 'API Keys — Orbit VPS Controller' }
  },
  {
    path: '/audit-logs',
    name: 'AuditLogs',
    component: () => import('@/views/system/AuditLogsView.vue'),
    meta: { title: 'Audit Logs — Orbit VPS Controller' }
  },
  {
    path: '/integrations',
    name: 'Integrations',
    component: () => import('@/views/system/IntegrationsView.vue'),
    meta: { title: 'Credentials — Orbit VPS Controller' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/system/SettingsView.vue'),
    meta: { title: 'Settings — Orbit VPS Controller' }
  },

  // Auth & Errors
  {
    path: '/auth/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { layout: 'blank', title: 'Sign In — Orbit' }
  },
  {
    path: '/auth/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { layout: 'blank', title: 'Create Account — Orbit' }
  },
  {
    path: '/auth/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPasswordView.vue'),
    meta: { layout: 'blank', title: 'Reset Password — Orbit' }
  },
  {
    path: '/auth/2fa',
    name: 'TwoFactor',
    component: () => import('@/views/auth/TwoFactorView.vue'),
    meta: { layout: 'blank', title: 'Two-Factor Verification — Orbit' }
  },
  {
    path: '/error/500',
    name: 'Error500',
    component: () => import('@/views/errors/Error500View.vue'),
    meta: { layout: 'blank', title: '500 Internal Server Error — Orbit' }
  },
  {
    path: '/maintenance',
    name: 'Maintenance',
    component: () => import('@/views/errors/MaintenanceView.vue'),
    meta: { layout: 'blank', title: 'Scheduled Maintenance — Orbit' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/errors/Error404View.vue'),
    meta: { layout: 'blank', title: '404 Not Found — Orbit' }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }
})

export { router }
