<template>
  <div style="display: none;">
    <!-- Floating Customizer Trigger Button (hidden in live mirror copy) -->
    <button class="cf-customizer-trigger" @click="toggleCustomizer(true)" title="Customize Theme & Layout">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/>
        <circle cx="12" cy="12" r="3"/>
      </svg>
    </button>

    <!-- Customizer Backdrop & Drawer Panel -->
    <div class="cf-customizer-backdrop" :class="{ open: customizerOpen }" @click="toggleCustomizer(false)"></div>
    <aside class="cf-customizer-drawer" :class="{ open: customizerOpen }">
      <div class="cf-customizer-head">
        <div>
          <h3 style="margin: 0; font-size: 1.05rem;">Theme Customizer</h3>
          <span style="font-size: 0.75rem; color: var(--text-color-kumo-subtle);">Envato Commercial Options</span>
        </div>
        <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="toggleCustomizer(false)" style="padding: 0.25rem 0.5rem;">✕</button>
      </div>

      <div class="cf-customizer-body">
        <!-- Color Scheme -->
        <div>
          <div class="cf-customizer-section-title">Color Scheme</div>
          <div style="display: flex; gap: 0.5rem;">
            <button
              class="btn-cf"
              :class="theme === 'light' ? 'btn-cf-primary' : 'btn-cf-secondary'"
              style="flex: 1; justify-content: center;"
              @click="theme = 'light'; toggleTheme()"
            >
              ☀️ Light
            </button>
            <button
              class="btn-cf"
              :class="theme === 'dark' ? 'btn-cf-primary' : 'btn-cf-secondary'"
              style="flex: 1; justify-content: center;"
              @click="theme = 'dark'; toggleTheme()"
            >
              🌙 Dark
            </button>
          </div>
        </div>

        <!-- Brand Accent Colors -->
        <div>
          <div class="cf-customizer-section-title">Brand Accent Color</div>
          <div class="color-swatches-grid">
            <button
              v-for="color in swatches"
              :key="color.name"
              class="color-swatch-btn"
              :class="[`swatch-${color.name}`, { active: accent === color.name }]"
              :title="color.label"
              @click="setAccent(color.name)"
            >
              {{ accent === color.name ? '✓' : '' }}
            </button>
          </div>
        </div>

        <!-- Sidebar Layout Style -->
        <div>
          <div class="cf-customizer-section-title">Sidebar Style</div>
          <div style="display: flex; gap: 0.5rem;">
            <button
              class="btn-cf"
              :class="!sidebarMini ? 'btn-cf-primary' : 'btn-cf-secondary'"
              style="flex: 1; justify-content: center;"
              @click="setSidebarMini(false)"
            >
              Expanded (256px)
            </button>
            <button
              class="btn-cf"
              :class="sidebarMini ? 'btn-cf-primary' : 'btn-cf-secondary'"
              style="flex: 1; justify-content: center;"
              @click="setSidebarMini(true)"
            >
              Mini (68px)
            </button>
          </div>
        </div>

        <!-- Standalone Auth & Error Links -->
        <div>
          <div class="cf-customizer-section-title">Auth & System Pages</div>
          <div style="display: flex; flex-direction: column; gap: 0.35rem;">
            <router-link to="/login" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>🔐 Sign In / Login</span> ↗</router-link>
            <router-link to="/register" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>📝 Sign Up / Register</span> ↗</router-link>
            <router-link to="/2fa" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>🛡️ Two-Factor Auth (2FA)</span> ↗</router-link>
            <router-link to="/forgot-password" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>🔑 Forgot Password</span> ↗</router-link>
            <router-link to="/404" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>⚠️ 404 Not Found</span> ↗</router-link>
            <router-link to="/500" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>💥 500 Server Error</span> ↗</router-link>
            <router-link to="/maintenance" class="kumo-nav-sub-item" @click="toggleCustomizer(false)"><span>🚧 Maintenance Mode</span> ↗</router-link>
          </div>
        </div>
      </div>

      <div class="cf-customizer-foot">
        <button class="btn-cf btn-cf-secondary" @click="resetDefaults" style="flex: 1; justify-content: center;">Reset Defaults</button>
        <button class="btn-cf btn-cf-primary" @click="toggleCustomizer(false)" style="flex: 1; justify-content: center;">Done</button>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { useTheme } from '@/composables/useTheme.js'

const { theme, accent, sidebarMini, customizerOpen, toggleTheme, setAccent, setSidebarMini, toggleCustomizer, resetDefaults } = useTheme()

const swatches = [
  { name: 'orange', label: 'Cloudflare Orange' },
  { name: 'blue', label: 'Deep Blue' },
  { name: 'emerald', label: 'Emerald Green' },
  { name: 'purple', label: 'Royal Purple' },
  { name: 'red', label: 'Ruby Red' }
]
</script>
