/* ==========================================================================
   Cloudflare Kumo — useTheme Composable (Theme & Layout Engine)
   ========================================================================== */

import { ref, watch, onMounted } from 'vue'

const theme = ref(localStorage.getItem('cf_theme') || 'light')
const accent = ref(localStorage.getItem('cf_accent') || 'orange')
const sidebarMini = ref(localStorage.getItem('cf_sidebar_mini') === 'true')
const mobileSidebarOpen = ref(false)
const customizerOpen = ref(false)
const commandPaletteOpen = ref(false)

export function useTheme() {
  function applyTheme() {
    if (theme.value === 'dark') {
      document.documentElement.classList.add('dark')
      document.body.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
      document.body.classList.remove('dark')
    }
    localStorage.setItem('cf_theme', theme.value)
  }

  function applyAccent() {
    document.documentElement.setAttribute('data-accent', accent.value)
    localStorage.setItem('cf_accent', accent.value)
  }

  function applySidebarMini() {
    if (sidebarMini.value) {
      document.body.classList.add('sidebar-mini')
    } else {
      document.body.classList.remove('sidebar-mini')
    }
    localStorage.setItem('cf_sidebar_mini', sidebarMini.value)
  }

  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    applyTheme()
  }

  function setAccent(color) {
    accent.value = color
    applyAccent()
  }

  function setSidebarMini(val) {
    sidebarMini.value = val
    applySidebarMini()
  }

  function toggleCustomizer(val) {
    customizerOpen.value = val !== undefined ? val : !customizerOpen.value
  }

  function toggleMobileSidebar(val) {
    mobileSidebarOpen.value = val !== undefined ? val : !mobileSidebarOpen.value
  }

  function toggleCommandPalette(val) {
    commandPaletteOpen.value = val !== undefined ? val : !commandPaletteOpen.value
  }

  function resetDefaults() {
    theme.value = 'light'
    accent.value = 'orange'
    sidebarMini.value = false
    applyTheme()
    applyAccent()
    applySidebarMini()
    customizerOpen.value = false
  }

  onMounted(() => {
    applyTheme()
    applyAccent()
    applySidebarMini()
  })

  return {
    theme,
    accent,
    sidebarMini,
    mobileSidebarOpen,
    customizerOpen,
    commandPaletteOpen,
    toggleTheme,
    setAccent,
    setSidebarMini,
    toggleCustomizer,
    toggleMobileSidebar,
    toggleCommandPalette,
    resetDefaults
  }
}
