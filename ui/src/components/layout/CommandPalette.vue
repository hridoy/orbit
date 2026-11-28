<template>
  <div
    v-if="commandPaletteOpen"
    class="cf-command-backdrop"
    @click.self="toggleCommandPalette(false)"
  >
    <div class="cf-command-palette">
      <div class="cf-command-search-head">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
        </svg>
        <input
          ref="inputRef"
          v-model="searchQuery"
          type="text"
          placeholder="Quick search workloads, databases, buckets, zones... (Type to filter)"
        />
        <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="toggleCommandPalette(false)">Esc</button>
      </div>

      <div class="cf-command-results">
        <div class="cf-command-section-title">Navigation & Quick Switcher ({{ filteredItems.length }} results)</div>
        <div
          v-for="item in filteredItems"
          :key="item.path"
          class="cf-command-item"
          @click="selectItem(item)"
        >
          <div class="item-left">
            <span>{{ item.icon }}</span>
            <strong>{{ item.title }}</strong>
          </div>
          <span class="cf-pill stopped">{{ item.badge }}</span>
        </div>
      </div>

      <div class="cf-command-foot">
        <span>Press <kbd>↑</kbd> <kbd>↓</kbd> to navigate, <kbd>Enter</kbd> to select</span>
        <kbd>ESC to close</kbd>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from '@/composables/useTheme.js'

const router = useRouter()
const { commandPaletteOpen, toggleCommandPalette } = useTheme()
const searchQuery = ref('')
const inputRef = ref(null)

const allItems = [
  { icon: '🏠', title: 'Account Home', path: '/', badge: 'Overview' },
  { icon: '🛡️', title: 'SSL/TLS Certificates', path: '/ssl', badge: 'Strict' },
  { icon: '⚙️', title: 'Profile & Preferences', path: '/settings', badge: 'Account' },
]

const filteredItems = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return allItems
  return allItems.filter(item => item.title.toLowerCase().includes(q) || item.badge.toLowerCase().includes(q))
})

function selectItem(item) {
  router.push(item.path)
  toggleCommandPalette(false)
}

watch(commandPaletteOpen, (open) => {
  if (open) {
    searchQuery.value = ''
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
})

function handleGlobalKeydown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    toggleCommandPalette()
  } else if (e.key === 'Escape' && commandPaletteOpen.value) {
    toggleCommandPalette(false)
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})
</script>
