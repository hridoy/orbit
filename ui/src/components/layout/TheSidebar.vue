<template>
  <aside 
    class="kumo-sidebar" 
    :class="{ 
      'sidebar-collapsed': sidebarMini, 
      'mobile-open': mobileSidebarOpen 
    }"
  >
    <!-- Top Section: Logo + Account Switcher -->
    <div class="sidebar-top-section">
      <!-- Logo Row -->
      <div class="sidebar-logo-row">
        <router-link to="/" class="sidebar-logo-link" title="Cloudflare Dashboard">
          <svg viewBox="0 0 49 22" fill="none" xmlns="http://www.w3.org/2000/svg" class="cf-logo-flame" width="34" height="20">
            <path d="M33.204 20.4C33.3649 19.9741 33.4217 19.5159 33.3695 19.0636C33.3173 18.6113 33.1577 18.1781 32.904 17.8C32.6435 17.4876 32.3239 17.2297 31.9636 17.0409C31.6032 16.8522 31.2092 16.7363 30.804 16.7L13.404 16.5C13.304 16.5 13.204 16.4 13.104 16.4C13.0808 16.3825 13.0618 16.3599 13.0488 16.3339C13.0358 16.3078 13.029 16.2791 13.029 16.25C13.029 16.2209 13.0358 16.1922 13.0488 16.1662C13.0618 16.1401 13.0808 16.1175 13.104 16.1C13.204 15.9 13.304 15.8 13.504 15.8L31.004 15.6C32.115 15.4767 33.1731 15.0597 34.0695 14.3918C34.9659 13.7239 35.6681 12.8293 36.104 11.8L37.104 9.20002C37.104 9.10002 37.204 9.00001 37.104 8.90001C36.5604 6.47843 35.2411 4.30052 33.3466 2.69721C31.4521 1.09391 29.086 0.152865 26.6079 0.0170769C24.1298 -0.118712 21.675 0.558179 19.6167 1.94489C17.5584 3.33161 16.009 5.35233 15.204 7.70002C14.159 6.95365 12.8843 6.59957 11.604 6.70002C10.4291 6.83102 9.33369 7.35777 8.49774 8.19372C7.66179 9.02966 7.13505 10.1251 7.00404 11.3C6.93745 11.9014 6.97125 12.5097 7.10404 13.1C5.20298 13.1526 3.39743 13.9448 2.07147 15.3081C0.745511 16.6714 0.00377461 18.4982 0.00403983 20.4C-0.0123708 20.7695 0.0212659 21.1395 0.104038 21.5C0.10863 21.5781 0.141713 21.6517 0.19701 21.707C0.252307 21.7623 0.325975 21.7954 0.404041 21.8H32.504C32.704 21.8 32.904 21.7 32.904 21.5L33.204 20.4Z" fill="#FF5E1F"></path>
            <path d="M38.704 9.20002H38.204C38.104 9.20002 38.004 9.30001 37.904 9.40001L37.204 11.8C37.0431 12.2259 36.9864 12.6841 37.0386 13.1364C37.0908 13.5887 37.2504 14.0219 37.504 14.4C37.7646 14.7124 38.0842 14.9704 38.4445 15.1591C38.8049 15.3479 39.1989 15.4637 39.604 15.5L43.304 15.7C43.404 15.7 43.504 15.8 43.604 15.8C43.6273 15.8175 43.6462 15.8401 43.6592 15.8662C43.6723 15.8922 43.679 15.9209 43.679 15.95C43.679 15.9791 43.6723 16.0078 43.6592 16.0339C43.6462 16.0599 43.6273 16.0826 43.604 16.1C43.504 16.3 43.404 16.4 43.204 16.4L39.404 16.6C38.293 16.7233 37.2349 17.1403 36.3386 17.8082C35.4422 18.4761 34.74 19.3707 34.304 20.4L34.104 21.3C34.004 21.4 34.104 21.6 34.304 21.6H47.504C47.5448 21.6058 47.5863 21.6021 47.6254 21.5891C47.6644 21.5761 47.6999 21.5541 47.729 21.525C47.7581 21.4959 47.7801 21.4604 47.7931 21.4214C47.8061 21.3823 47.8099 21.3408 47.804 21.3C48.0421 20.4527 48.1764 19.5797 48.204 18.7C48.1882 16.1854 47.1822 13.7782 45.404 12C43.6259 10.2218 41.2187 9.21587 38.704 9.20002Z" fill="#FF9911"></path>
          </svg>
        </router-link>
      </div>

      <!-- Account Switcher Button (with Dropdown) -->
      <div v-if="!sidebarMini" class="account-switcher-container">
        <button 
          class="account-switcher-btn" 
          :class="{ active: accountDropdownOpen }"
          @click="accountDropdownOpen = !accountDropdownOpen"
        >
          <span class="account-switcher-text">{{ server.hostname || 'This server' }}</span>
          <CfIcon name="cis-chevron-grabber-vertical" size="14" class="caret-icon" />
        </button>

        <!-- Account Selector Dropdown Menu -->
        <div v-if="accountDropdownOpen" class="account-dropdown-menu">
          <div class="dropdown-search-row">
            <CfIcon name="cis-magnifying-glass" size="14" />
            <input type="text" placeholder="Search accounts..." />
          </div>
          <div class="account-item selected">
            <div class="account-avatar-dot"></div>
            <div style="flex: 1; min-width: 0;">
              <div class="account-item-title">{{ server.hostname || 'This server' }}</div>
              <div class="account-item-id">{{ server.nodeId }}</div>
            </div>
            <CfIcon name="cis-checkmark-1-medium" size="14" customClass="text-kumo-brand" />
          </div>
          <div class="dropdown-divider"></div>
          <button class="create-account-btn" @click="alertCreateAccount">+ Create Account</button>
        </div>
      </div>

      <!-- Quick Search Button -->
      <button v-if="!sidebarMini" class="sidebar-search-btn" @click="toggleCommandPalette(true)">
        <CfIcon name="cis-magnifying-glass" size="15" />
        <span>Quick search...</span>
        <kbd>⌘ K</kbd>
      </button>
    </div>

    <!-- Sliding Views Container (Context Switcher between Account & Zone Level) -->
    <div class="sidebar-sliding-container" :class="{ 'in-zone-view': isZoneScoped }">
      <div class="sliding-views-track">
        
        <!-- ============================================================== -->
        <!-- VIEW 1: SIMPLE & MINIMAL VPS CONTROLLER NAVIGATION             -->
        <!-- ============================================================== -->
        <div class="sliding-view account-view">
          <!-- Overview -->
          <router-link to="/" class="sidebar-nav-item" active-class="active" exact>
            <CfIcon name="home" size="18" />
            <span v-if="!sidebarMini">Overview</span>
          </router-link>

          <!-- Section: DOMAIN -->
          <div v-if="!sidebarMini" class="sidebar-section-header">Domain</div>

          <router-link to="/ssl" class="sidebar-nav-item" active-class="active">
            <CfIcon name="cis-lock" size="18" />
            <span v-if="!sidebarMini">SSL & Domains</span>
          </router-link>

          <!-- Section: APPS -->
          <div v-if="!sidebarMini" class="sidebar-section-header">Apps</div>

          <router-link to="/apps" class="sidebar-nav-item" active-class="active">
            <CfIcon name="edgeworker" size="18" />
            <span v-if="!sidebarMini">Apps</span>
          </router-link>

          <!-- Section: STORAGE (Separate) -->
          <div v-if="!sidebarMini" class="sidebar-section-header">Storage</div>

          <router-link to="/storage" class="sidebar-nav-item" active-class="active">
            <CfIcon name="r2" size="18" />
            <span v-if="!sidebarMini">Object Storage</span>
          </router-link>

          <!-- Section: DATABASE (Separate) -->
          <div v-if="!sidebarMini" class="sidebar-section-header">Database</div>

          <router-link to="/databases" class="sidebar-nav-item" active-class="active">
            <CfIcon name="d1" size="18" />
            <span v-if="!sidebarMini">Databases</span>
          </router-link>

          <!-- Section: SECURITY -->
          <div v-if="!sidebarMini" class="sidebar-section-header">Security</div>



          <!-- Section: SYSTEM -->
          <div v-if="!sidebarMini" class="sidebar-section-header">System</div>

          <router-link to="/metrics" class="sidebar-nav-item" active-class="active">
            <CfIcon name="web-analytics-logo" size="18" />
            <span v-if="!sidebarMini">Metrics</span>
          </router-link>

          <router-link to="/users" class="sidebar-nav-item" active-class="active">
            <CfIcon name="cis-user-duo" size="18" />
            <span v-if="!sidebarMini">Users</span>
          </router-link>


          <router-link to="/integrations" class="sidebar-nav-item" active-class="active">
            <CfIcon name="cis-settings-slider-hor" size="18" />
            <span v-if="!sidebarMini">Credentials</span>
          </router-link>

          <router-link to="/settings" class="sidebar-nav-item" active-class="active">
            <CfIcon name="cis-settings-gear-2" size="18" />
            <span v-if="!sidebarMini">Settings</span>
          </router-link>
        </div>

        <!-- ============================================================== -->
        <!-- VIEW 2: ZONE LEVEL NAVIGATION (e.g. northlab.dev)              -->
        <!-- ============================================================== -->
        <div class="sliding-view zone-view">
          <!-- Back to Domains Button -->
          <router-link to="/" class="back-to-domains-btn">
            <CfIcon name="cis-arrow-left" size="16" />
            <span>Back to Domains</span>
          </router-link>

          <!-- Zone Overview -->
          <router-link to="/ssl" class="sidebar-nav-item" active-class="active" exact>
            <CfIcon name="cis-globe" size="18" />
            <span>Overview</span>
          </router-link>

          <!-- DNS Records -->

          <!-- SSL/TLS -->
          <router-link to="/ssl" class="sidebar-nav-item" active-class="active">
            <CfIcon name="cis-lock" size="18" />
            <span>SSL / TLS</span>
          </router-link>

          <!-- Caching -->

          <!-- Rules -->

          <!-- Network -->

          <!-- WAF -->
        </div>

      </div>
    </div>

    <!-- Bottom Fixed Trigger (Collapse Sidebar Button) -->
    <div class="sidebar-bottom-trigger">
      <button 
        class="collapse-trigger-btn" 
        :title="sidebarMini ? 'Expand sidebar' : 'Collapse sidebar'"
        @click="toggleSidebarMini"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
          <path d="M21.25 6.72v10.56a2.97 2.97 0 0 1-2.97 2.97H5.72a2.97 2.97 0 0 1-2.97-2.97V6.72a2.97 2.97 0 0 1 2.97-2.97h12.56a2.97 2.97 0 0 1 2.97 2.97" />
          <path d="M6.25 7.25v9.5" :style="{ transform: sidebarMini ? 'translateX(7px)' : 'translateX(0)' }" class="sidebar-split-line" />
        </svg>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRoute } from 'vue-router'
import { server } from '@/composables/serverState.js'
import { useTheme } from '@/composables/useTheme.js'

const route = useRoute()
const { sidebarMini, setSidebarMini, mobileSidebarOpen, toggleCommandPalette } = useTheme()

const accountDropdownOpen = ref(false)

const openGroups = reactive({
  recents: false,
  compute: false,
  storage: true,
  security: false
})

const isZoneScoped = computed(() => false)

function toggleGroup(key) {
  openGroups[key] = !openGroups[key]
}

function toggleSidebarMini() {
  setSidebarMini(!sidebarMini.value)
}

function alertCreateAccount() {
  accountDropdownOpen.value = false
  alert('Redirecting to Cloudflare Account Creation...')
}
</script>
