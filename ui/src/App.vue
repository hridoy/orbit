<template>
  <!-- Orbit Sovereign First-Run Setup / Login Gate -->
  <AuthGate v-if="authChecked && !signedIn" :setup-required="setupRequired" @authed="onAuthed" />

  <!-- Blank Layout for Auth & Error Pages -->
  <div v-else-if="route.meta.layout === 'blank'">
    <router-view />
  </div>

  <!-- Standard Dashboard Layout -->
  <div v-else>
    <!-- Mobile Sidebar Backdrop -->
    <div 
      class="mobile-sidebar-backdrop" 
      :class="{ active: mobileSidebarOpen }" 
      @click="toggleMobileSidebar(false)"
    ></div>

    <div class="kumo-shell">
      <TheSidebar />
      <div class="kumo-main">
        <TheHeader />
        <main class="kumo-content">
          <router-view />
        </main>
        <TheFooter />
      </div>
    </div>

    <!-- Theme Customizer Drawer -->
    <ThemeCustomizer />

    <!-- Spotlight Command Palette -->
    <CommandPalette />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AuthGate from '@/AuthGate.vue'
import TheSidebar from '@/components/TheSidebar.vue'
import TheHeader from '@/components/TheHeader.vue'
import TheFooter from '@/components/TheFooter.vue'
import ThemeCustomizer from '@/components/ThemeCustomizer.vue'
import CommandPalette from '@/components/CommandPalette.vue'
import { useTheme } from '@/composables/useTheme.js'
import { loadServer } from '@/composables/serverState.js'

const route = useRoute()
const { mobileSidebarOpen, toggleMobileSidebar } = useTheme()

const authToken = ref(localStorage.getItem('orbit_token') || '')
const urlParams = new URLSearchParams(window.location.search)
if (urlParams.get('token')) {
  authToken.value = urlParams.get('token')
  try { localStorage.setItem('orbit_token', authToken.value) } catch (e) {}
}

const signedIn = ref(false)
const authChecked = ref(false)
const setupRequired = ref(false)

async function hydrateNodeState() {
  try {
    await loadServer()
  } catch (e) {}
}

async function checkAuth() {
  try {
    if (authToken.value) {
      const res = await fetch('/api/v1/node', { headers: { Authorization: 'Bearer ' + authToken.value } })
      if (res.ok) {
        signedIn.value = true
        authChecked.value = true
        hydrateNodeState()
        return
      }
    }
    const statusRes = await fetch('/api/v1/auth/status')
    if (statusRes.ok) {
      const body = await statusRes.json()
      setupRequired.value = !!body.setup_required
      signedIn.value = false
      authChecked.value = true
      return
    }
  } catch (e) {
    // Standalone preview without Go backend
  }
  signedIn.value = true
  authChecked.value = true
}

function onAuthed(token) {
  authToken.value = token
  try { localStorage.setItem('orbit_token', token) } catch (e) {}
  signedIn.value = true
  hydrateNodeState()
}

onMounted(() => {
  checkAuth()
})
</script>
