<template>
  <header class="kumo-header">
    <!-- Header Left: Zone Switcher or Clean Space -->
    <div class="header-left">
      <!-- Mobile Menu Button -->
      <button class="mobile-menu-btn" @click="toggleMobileSidebar(true)" aria-label="Toggle Mobile Sidebar">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/>
        </svg>
      </button>

      <!-- Zone Switcher Component (Appears when in domain scope) -->
      <div v-if="isZoneScoped" class="zone-switcher-wrapper">
        <div class="zone-switcher-trigger-group">
          <!-- Star Domain Button -->
          <button 
            class="star-domain-btn" 
            :title="isStarred ? 'Unstar domain' : 'Star domain'" 
            type="button" 
            @click="isStarred = !isStarred"
          >
            <CfIcon name="cis-star-filled" size="14" :style="{ color: isStarred ? '#f59e0b' : 'currentColor' }" customClass="star-icon" />
          </button>

          <!-- Zone Selector Trigger Button -->
          <button 
            class="zone-trigger-btn" 
            :class="{ active: zoneDropdownOpen }"
            type="button" 
            @click="zoneDropdownOpen = !zoneDropdownOpen"
          >
            <p class="zone-name">{{ currentZoneName }}</p>
            <span class="cf-plan-badge">free</span>
            <CfIcon name="cis-chevron-grabber-vertical" size="14" customClass="caret-icon" />
          </button>
        </div>

        <!-- 448px Zone Popover -->
        <div v-if="zoneDropdownOpen" class="cf-zone-popover">
          <!-- Popover Pointer Arrow -->
          <div class="cf-popover-arrow">
            <svg width="20" height="10" viewBox="0 0 20 10" fill="none">
              <path d="M9.66437 2.60207L4.80758 6.97318C4.07308 7.63423 3.11989 8 2.13172 8H0V10H20V8H18.5349C17.5468 8 16.5936 7.63423 15.8591 6.97318L11.0023 2.60207C10.622 2.2598 10.0447 2.25979 9.66437 2.60207Z" class="arrow-fill"/>
            </svg>
          </div>

          <!-- Search & Star Filter Header -->
          <div class="zone-search-row">
            <button 
              class="star-filter-btn" 
              :class="{ active: starredOnly }" 
              title="Show starred zones only" 
              type="button"
              @click="starredOnly = !starredOnly"
            >
              <CfIcon name="cis-star-filled" size="14" :style="{ color: starredOnly ? '#f59e0b' : 'currentColor' }" />
            </button>
            <input 
              v-model="zoneSearch" 
              placeholder="Search zone names..." 
              class="zone-search-input" 
              autofocus
            />
            <CfIcon name="cis-magnifying-glass" size="14" customClass="search-icon" />
          </div>

          <!-- Zone List Items -->
          <div class="zone-items-list">
            <router-link 
              v-for="zone in filteredZones" 
              :key="zone.id" 
              :to="'/zones/' + zone.id" 
              class="zone-item-row" 
              :class="{ 'is-current': zone.name === currentZoneName }"
              @click="zoneDropdownOpen = false"
            >
              <div class="zone-item-left">
                <span class="star-spacer"></span>
                <span class="zone-item-name">{{ zone.name }}</span>
                <CfIcon v-if="zone.name === currentZoneName" name="cis-checkmark-1-medium" size="14" customClass="active-check" />
              </div>
              <div class="zone-item-right">
                <span class="cf-plan-badge">{{ (zone.plan || 'free').toLowerCase().replace(' plan', '').replace(' website', '') }}</span>
              </div>
            </router-link>
            <div v-if="filteredZones.length === 0" style="padding: 16px; text-align: center; font-size: 13px; color: var(--text-color-kumo-subtle);">
              No zones found matching "{{ zoneSearch }}"
            </div>
          </div>

          <!-- Zone Popover Footer -->
          <div class="zone-popover-footer">
            <span class="domain-count">{{ filteredZones.length }} domain{{ filteredZones.length === 1 ? '' : 's' }}</span>
            <router-link to="/" class="go-all-link" @click="zoneDropdownOpen = false">Go to all domains →</router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Actions Group -->
    <div class="header-right">
      <!-- 1. Ask AI Button -->
      <button 
        class="header-btn" 
        :class="{ active: askAiDrawerOpen }" 
        type="button" 
        @click="askAiDrawerOpen = !askAiDrawerOpen"
      >
        <CfIcon name="flare-ask-ai" size="17" customClass="btn-icon" />
        <span>Ask AI</span>
      </button>

      <!-- 2. Support Link Button -->
      <a href="https://support.cloudflare.com" target="_blank" class="header-btn header-link">
        <CfIcon name="cis-circle-questionmark" size="17" customClass="btn-icon" />
        <span>Support</span>
      </a>

      <!-- 3. Notifications Bell Button -->
      <button 
        class="header-icon-btn" 
        type="button" 
        aria-label="Notifications" 
        @click="alertNotifications"
      >
        <CfIcon name="cis-bell" size="17" customClass="btn-icon" />
      </button>

      <!-- 4. User Menu Trigger Button & Popover -->
      <div class="user-menu-wrapper">
        <button 
          class="header-icon-btn user-btn" 
          :class="{ active: userMenuOpen }"
          type="button" 
          aria-label="User menu" 
          @click="userMenuOpen = !userMenuOpen"
        >
          <CfIcon name="cis-user" size="17" customClass="btn-icon" />
        </button>

        <!-- 280px User Menu Dropdown Popover -->
        <div v-if="userMenuOpen" class="cf-user-popover">
          <!-- Header: Email & Account Info -->
          <router-link to="/settings" class="user-info-section" style="text-decoration: none;" @click="userMenuOpen = false">
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <div style="min-width: 0;">
                <div class="user-email">Dashboard administrator</div>
                <div class="user-account-name">{{ server.hostname || 'This server' }}</div>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 256 256" class="text-kumo-subtle" style="flex-shrink: 0;">
                <circle cx="128" cy="128" r="16"></circle>
                <circle cx="64" cy="128" r="16"></circle>
                <circle cx="192" cy="128" r="16"></circle>
              </svg>
            </div>
          </router-link>

          <div class="cf-popover-divider"></div>

          <!-- Main Menu Items -->
          <div class="cf-menu-list">
            <!-- Account Home -->
            <router-link to="/" class="cf-menu-item" @click="userMenuOpen = false">
              <span>Account Home</span>
            </router-link>

            <!-- Profile -->
            <router-link to="/settings" class="cf-menu-item" @click="userMenuOpen = false">
              <span>Profile</span>
            </router-link>

            <!-- Appearance Submenu Trigger -->
            <div 
              class="cf-menu-item has-submenu" 
              :class="{ open: activeSubmenu === 'appearance' }"
              @mouseenter="activeSubmenu = 'appearance'" 
              @mouseleave="activeSubmenu = null"
            >
              <div class="item-inner">
                <span>Appearance</span>
                <svg class="chevron-right" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                  <path d="M181.66,133.66l-80,80a8,8,0,0,1-11.32-11.32L164.69,128,90.34,53.66a8,8,0,0,1,11.32-11.32l80,80A8,8,0,0,1,181.66,133.66Z"/>
                </svg>
              </div>

              <!-- Appearance Submenu Flyout -->
              <div v-if="activeSubmenu === 'appearance'" class="cf-submenu">
                <button 
                  class="cf-submenu-item" 
                  :class="{ active: currentAppearance === 'light' }" 
                  @click.stop="setAppearance('light')"
                >
                  <span>Light mode</span>
                  <svg v-if="currentAppearance === 'light'" class="check-icon" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"/>
                  </svg>
                </button>
                <button 
                  class="cf-submenu-item" 
                  :class="{ active: currentAppearance === 'dark' }" 
                  @click.stop="setAppearance('dark')"
                >
                  <span>Dark mode</span>
                  <svg v-if="currentAppearance === 'dark'" class="check-icon" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"/>
                  </svg>
                </button>
                <button 
                  class="cf-submenu-item" 
                  :class="{ active: currentAppearance === 'system' }" 
                  @click.stop="setAppearance('system')"
                >
                  <span>System</span>
                  <svg v-if="currentAppearance === 'system'" class="check-icon" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"/>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Language Submenu Trigger -->
            <div 
              class="cf-menu-item has-submenu" 
              :class="{ open: activeSubmenu === 'language' }"
              @mouseenter="activeSubmenu = 'language'" 
              @mouseleave="activeSubmenu = null"
            >
              <div class="item-inner">
                <span>Language</span>
                <svg class="chevron-right" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                  <path d="M181.66,133.66l-80,80a8,8,0,0,1-11.32-11.32L164.69,128,90.34,53.66a8,8,0,0,1,11.32-11.32l80,80A8,8,0,0,1,181.66,133.66Z"/>
                </svg>
              </div>

              <!-- Language Submenu Flyout -->
              <div v-if="activeSubmenu === 'language'" class="cf-submenu">
                <button 
                  v-for="lang in ['English (US)', 'Español', 'Français', 'Deutsch', '日本語']" 
                  :key="lang"
                  class="cf-submenu-item"
                  :class="{ active: currentLanguage === lang }"
                  @click.stop="currentLanguage = lang"
                >
                  <span>{{ lang }}</span>
                  <svg v-if="currentLanguage === lang" class="check-icon" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"/>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Timezone Submenu Trigger -->
            <div 
              class="cf-menu-item has-submenu" 
              :class="{ open: activeSubmenu === 'timezone' }"
              @mouseenter="activeSubmenu = 'timezone'" 
              @mouseleave="activeSubmenu = null"
            >
              <div class="item-inner">
                <span>Timezone</span>
                <svg class="chevron-right" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                  <path d="M181.66,133.66l-80,80a8,8,0,0,1-11.32-11.32L164.69,128,90.34,53.66a8,8,0,0,1,11.32-11.32l80,80A8,8,0,0,1,181.66,133.66Z"/>
                </svg>
              </div>

              <!-- Timezone Submenu Flyout -->
              <div v-if="activeSubmenu === 'timezone'" class="cf-submenu">
                <button 
                  class="cf-submenu-item"
                  :class="{ active: currentTimezone === 'UTC' }"
                  @click.stop="currentTimezone = 'UTC'"
                >
                  <span>UTC</span>
                  <svg v-if="currentTimezone === 'UTC'" class="check-icon" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"/>
                  </svg>
                </button>
                <button 
                  class="cf-submenu-item"
                  :class="{ active: currentTimezone === 'Local' }"
                  @click.stop="currentTimezone = 'Local'"
                >
                  <span>Local (Asia/Dhaka)</span>
                  <svg v-if="currentTimezone === 'Local'" class="check-icon" width="14" height="14" viewBox="0 0 256 256" fill="currentColor">
                    <path d="M229.66,77.66l-128,128a8,8,0,0,1-11.32,0l-56-56a8,8,0,0,1,11.32-11.32L96,188.69,218.34,66.34a8,8,0,0,1,11.32,11.32Z"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <div class="cf-popover-divider"></div>

          <!-- Copy Account ID -->
          <button class="cf-menu-item" style="width: 100%; border: none; background: transparent;" @click="copyAccountId">
            <span>{{ copiedId ? '✓ Account ID Copied' : 'Copy account ID' }}</span>
          </button>

          <div class="cf-popover-divider"></div>

          <!-- Log out in Real Cloudflare Red -->
          <router-link to="/auth/login" class="cf-menu-item logout-item" @click="userMenuOpen = false">
            <span style="color: #ef4444; font-weight: 500;">Log out</span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Ask AI Slide-out Drawer -->
    <div v-if="askAiDrawerOpen" class="ask-ai-drawer">
      <div class="drawer-header">
        <div style="display: flex; align-items: center; gap: 0.5rem;">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="#F6821F" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9.23 12C12.73 11.84 16.32 9.17 17.57 6.04C17.1 7.78 17.28 9.26 18.07 10.36C18.78 11.37 19.99 11.96 21.25 12.02C17.73 12.18 14.16 14.85 12.91 17.99C13.38 16.24 13.2 14.76 12.41 13.66C11.7 12.65 10.49 12.07 9.23 12ZM9.23 12H2.75"></path>
          </svg>
          <strong style="font-size: 0.95rem;">Ask AI</strong>
        </div>
        <button class="drawer-close-btn" @click="askAiDrawerOpen = false">✕</button>
      </div>
      <div class="drawer-body">
        <p style="font-size: 0.85rem; color: var(--text-color-kumo-subtle); margin-bottom: 1rem;">
          Ask questions about your DNS zones, Workers, or get recommended WAF rules for your application.
        </p>
        <div class="ai-sample-prompts">
          <button class="sample-prompt-btn" @click="queryAi('How do I configure Full (Strict) SSL for northlab.dev?')">
            "How do I configure Full (Strict) SSL?"
          </button>
          <button class="sample-prompt-btn" @click="queryAi('Why are my Workers requests returning 11ms latency?')">
            "Why are my Workers requests returning 11ms latency?"
          </button>
        </div>
        <div v-if="aiAnswer" class="ai-response-box">
          {{ aiAnswer }}
        </div>
      </div>
      <div class="drawer-footer">
        <input 
          v-model="aiQuery" 
          type="text" 
          placeholder="Ask a question about your account..." 
          class="cf-input" 
          @keyup.enter="runAiQuery"
        />
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { server } from '@/composables/serverState.js'
import { useTheme } from '@/composables/useTheme.js'

const route = useRoute()
const { theme, toggleTheme, toggleMobileSidebar } = useTheme()

const isStarred = ref(false)
const zoneDropdownOpen = ref(false)
const userMenuOpen = ref(false)
const askAiDrawerOpen = ref(false)
const activeSubmenu = ref(null)

const currentAppearance = ref(theme.value || 'light')
const currentLanguage = ref('English (US)')
const currentTimezone = ref('Local')

const zoneSearch = ref('')
const starredOnly = ref(false)
const copiedId = ref(false)

const aiQuery = ref('')
const aiAnswer = ref('')

const isZoneScoped = computed(() => {
  return route.path.startsWith('/zones') || 
         route.path === '/dns' || 
         route.path === '/ssl' || 
         route.path === '/caching' || 
         route.path === '/rules' || 
         route.path === '/network' || 
         route.path === '/waf'
})

const currentZoneName = computed(() => {
  return server.panelDomain || server.hostname || 'This server'
})

const filteredZones = computed(() => {
  let list = server.domains.map((d) => ({ name: d.domain }))
  if (starredOnly.value) {
    list = list.filter(z => z.name === server.panelDomain)
  }
  if (zoneSearch.value.trim()) {
    const q = zoneSearch.value.trim().toLowerCase()
    list = list.filter(z => z.name.toLowerCase().includes(q))
  }
  return list
})

function alertNotifications() {
  alert('Notifications: All systems operational. 0 unread security alerts.')
}

function setAppearance(mode) {
  currentAppearance.value = mode
  if (mode === 'dark' && theme.value !== 'dark') {
    toggleTheme()
  } else if (mode === 'light' && theme.value !== 'light') {
    toggleTheme()
  } else if (mode === 'system') {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    if (prefersDark && theme.value !== 'dark') toggleTheme()
    else if (!prefersDark && theme.value !== 'light') toggleTheme()
  }
}

function copyAccountId() {
  navigator.clipboard.writeText(server.nodeId)
  copiedId.value = true
  setTimeout(() => {
    copiedId.value = false
  }, 3000)
}

function queryAi(prompt) {
  aiQuery.value = prompt
  runAiQuery()
}

function runAiQuery() {
  if (!aiQuery.value) return
  aiAnswer.value = 'Analyzing your Cloudflare Anycast edge configuration... ' + aiQuery.value + ' - Recommendation: All 8 DNS records are proxied with TLS 1.3 and 0-RTT resumption active worldwide.'
}
</script>
