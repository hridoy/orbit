/* ==========================================================================
   Cloudflare Dashboard (Kumo Design System) — Vue 3 Commercial Application
   Envato Standard / ThemeForest / CodeCanyon
   Vue 3.5 + Vite 6 + Vue Router 4 + Composition API
   ========================================================================== */

import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router'

// Import Kumo Design System Stylesheets
import '@/assets/styles/kumo-tokens.css'
import '@/assets/styles/kumo-base.css'
import '@/assets/styles/kumo-sidebar.css'
import '@/assets/styles/kumo-header.css'
import '@/assets/styles/kumo-cards.css'
import '@/assets/styles/kumo-tables.css'
import '@/assets/styles/kumo-components.css'
import '@/assets/styles/kumo-modals.css'
import '@/assets/styles/kumo-r2.css'
import '@/assets/styles/kumo-bucket-details.css'
import '@/assets/styles/kumo-tokens-page.css'
import '@/assets/styles/kumo-workers.css'
import '@/assets/styles/kumo-pages.css'
import '@/assets/styles/kumo-customizer.css'
import '@/assets/styles/kumo-responsive.css'

// Global Cloudflare Kumo Design System Components
import * as KumoComponents from '@/components/common'

const app = createApp(App)

// Register all Kumo Design System components globally
for (const [name, component] of Object.entries(KumoComponents)) {
  app.component(name, component)
}
app.use(router)
app.mount('#app')
