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
    path: '/apps/create',
    alias: '/workers/create',
    name: 'CreateApp',
    component: () => import('@/views/apps/CreateAppView.vue'),
    meta: { title: 'Create Application — Orbit VPS Controller' }
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
    path: '/storage/:id',
    alias: '/r2/:id',
    name: 'BucketDetail',
    component: () => import('@/views/storage/BucketDetailView.vue'),
    meta: { title: 'Storage Bucket — Orbit VPS Controller' }
  },

  // Domains & SSL
  {
    path: '/ssl',
    name: 'SslDomains',
    component: () => import('@/views/domains/SslDomainsView.vue'),
    meta: { title: 'SSL & Domains — Orbit VPS Controller' }
  },

  // Security & Network

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
