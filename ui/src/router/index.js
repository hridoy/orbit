import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  // Dashboard Core Views
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: { title: 'Account Home — Cloudflare' }
  },
  {
    path: '/zones/northlab.dev',
    name: 'ZoneOverview',
    component: () => import('@/views/ZoneOverviewView.vue'),
    meta: { title: 'Zone Overview (northlab.dev) — Cloudflare' }
  },
  {
    path: '/dns',
    name: 'Dns',
    component: () => import('@/views/DnsView.vue'),
    meta: { title: 'Anycast DNS Management — Cloudflare' }
  },
  {
    path: '/ssl',
    name: 'Ssl',
    component: () => import('@/views/SslView.vue'),
    meta: { title: 'SSL/TLS Edge Certificates — Cloudflare' }
  },
  {
    path: '/caching',
    name: 'Caching',
    component: () => import('@/views/CachingView.vue'),
    meta: { title: 'Caching & Tiered CDN — Cloudflare' }
  },
  {
    path: '/rules',
    name: 'Rules',
    component: () => import('@/views/RulesView.vue'),
    meta: { title: 'Edge Transform & Page Rules — Cloudflare' }
  },
  {
    path: '/network',
    name: 'Network',
    component: () => import('@/views/NetworkView.vue'),
    meta: { title: 'Network Protocols (HTTP/3 & gRPC) — Cloudflare' }
  },
  {
    path: '/waf',
    name: 'Waf',
    component: () => import('@/views/WafView.vue'),
    meta: { title: 'WAF & Security Rules — Cloudflare' }
  },
  {
    path: '/apps',
    alias: '/workers',
    name: 'Apps',
    component: () => import('@/views/WorkersView.vue'),
    meta: { title: 'Apps — Orbit VPS Controller' }
  },
  {
    path: '/apps/:id',
    alias: '/workers/:id',
    name: 'AppDetail',
    component: () => import('@/views/WorkerDetailView.vue'),
    meta: { title: 'App Console — Orbit VPS Controller' }
  },
  {
    path: '/ai',
    name: 'WorkersAi',
    component: () => import('@/views/WorkersAiView.vue'),
    meta: { title: 'Workers AI & AI Gateway — Cloudflare' }
  },
  {
    path: '/d1',
    name: 'D1',
    component: () => import('@/views/D1View.vue'),
    meta: { title: 'D1 SQLite Databases — Cloudflare' }
  },
  {
    path: '/d1/:id',
    name: 'D1Detail',
    component: () => import('@/views/D1DetailView.vue'),
    meta: { title: 'D1 SQL Query Console — Cloudflare' }
  },
  {
    path: '/r2',
    name: 'R2',
    component: () => import('@/views/R2View.vue'),
    meta: { title: 'R2 Zero-Egress Storage — Cloudflare' }
  },
  {
    path: '/r2/api-tokens',
    name: 'R2ApiTokens',
    component: () => import('@/views/TokensView.vue'),
    meta: { title: 'R2 API Tokens — Cloudflare' }
  },
  {
    path: '/r2/:id',
    name: 'R2Detail',
    component: () => import('@/views/R2DetailView.vue'),
    meta: { title: 'R2 Object Bucket — Cloudflare' }
  },
  {
    path: '/r2/:bucket/objects/:key(.*)/details',
    name: 'R2ObjectDetail',
    component: () => import('@/views/R2ObjectDetailView.vue'),
    meta: { title: 'R2 Object Details — Cloudflare' }
  },
  {
    path: '/kv',
    name: 'WorkersKv',
    component: () => import('@/views/WorkersKvView.vue'),
    meta: { title: 'Workers KV Global Store — Cloudflare' }
  },
  {
    path: '/turnstile',
    name: 'Turnstile',
    component: () => import('@/views/TurnstileView.vue'),
    meta: { title: 'Turnstile CAPTCHA & Security — Cloudflare' }
  },
  {
    path: '/analytics',
    name: 'Analytics',
    component: () => import('@/views/AnalyticsView.vue'),
    meta: { title: 'Account Telemetry & Analytics — Cloudflare' }
  },
  {
    path: '/zero-trust',
    name: 'ZeroTrust',
    component: () => import('@/views/ZeroTrustView.vue'),
    meta: { title: 'Zero Trust Access & Tunnels — Cloudflare' }
  },
  {
    path: '/team',
    name: 'Team',
    component: () => import('@/views/TeamView.vue'),
    meta: { title: 'Members & Roles (RBAC) — Cloudflare' }
  },
  {
    path: '/billing',
    name: 'Billing',
    component: () => import('@/views/BillingView.vue'),
    meta: { title: 'Billing & Subscriptions — Cloudflare' }
  },
  {
    path: '/tokens',
    name: 'Tokens',
    component: () => import('@/views/TokensView.vue'),
    meta: { title: 'API Tokens & Access Credentials — Cloudflare' }
  },
  {
    path: '/audit-logs',
    name: 'AuditLogs',
    component: () => import('@/views/AuditLogsView.vue'),
    meta: { title: 'Audit Logs Activity Stream — Cloudflare' }
  },
  {
    path: '/integrations',
    name: 'Integrations',
    component: () => import('@/views/IntegrationsView.vue'),
    meta: { title: 'Credentials & Integrations' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue'),
    meta: { title: 'Profile & Account Preferences — Cloudflare' }
  },
  {
    path: '/components',
    name: 'Components',
    component: () => import('@/views/ComponentsView.vue'),
    meta: { title: 'UI Kit Showcase (Kumo Design System) — Cloudflare' }
  },

  // Authentication Views (Blank Layout)
  {
    path: '/auth/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { layout: 'blank', title: 'Log in — Cloudflare' }
  },
  {
    path: '/auth/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { layout: 'blank', title: 'Create an account — Cloudflare' }
  },
  {
    path: '/auth/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPasswordView.vue'),
    meta: { layout: 'blank', title: 'Reset password — Cloudflare' }
  },
  {
    path: '/auth/2fa',
    name: 'TwoFactor',
    component: () => import('@/views/auth/TwoFactorView.vue'),
    meta: { layout: 'blank', title: 'Two-Factor Authentication — Cloudflare' }
  },

  // Error & Maintenance Views (Blank Layout)
  {
    path: '/errors/500',
    name: 'Error500',
    component: () => import('@/views/errors/Error500View.vue'),
    meta: { layout: 'blank', title: '500 Server Error — Cloudflare' }
  },
  {
    path: '/errors/maintenance',
    name: 'Maintenance',
    component: () => import('@/views/errors/MaintenanceView.vue'),
    meta: { layout: 'blank', title: 'Scheduled Maintenance — Cloudflare' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/errors/Error404View.vue'),
    meta: { layout: 'blank', title: '404 Page Not Found — Cloudflare' }
  }
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.afterEach((to) => {
  if (to.meta && to.meta.title) {
    document.title = to.meta.title
  }
})
