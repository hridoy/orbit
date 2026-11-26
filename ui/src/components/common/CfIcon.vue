<template>
  <svg
    :width="size || width || 16"
    :height="size || height || 16"
    :viewBox="resolvedIcon.viewBox"
    fill="currentColor"
    :class="['cf-icon', customClass]"
    role="presentation"
    aria-hidden="true"
    v-html="resolvedIcon.content"
  ></svg>
</template>

<script setup>
import { computed } from 'vue'
import iconsCatalog from '@/assets/icons/icons-catalog.json'

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: {
    type: [Number, String],
    default: null
  },
  width: {
    type: [Number, String],
    default: 16
  },
  height: {
    type: [Number, String],
    default: 16
  },
  customClass: {
    type: String,
    default: ''
  }
})

// Semantic aliases mapping common names to official Cloudflare Kumo sprite IDs
const ICON_ALIASES = {
  'check': 'icon-cis-checkmark-1-medium',
  'check-circle': 'icon-cis-check-circle',
  'x': 'icon-cis-cross-medium',
  'close': 'icon-cis-cross-medium',
  'minus': 'icon-cis-minus-medium',
  'chevron-down': 'icon-cis-chevron-down-medium',
  'chevron-up': 'icon-cis-chevron-top-medium',
  'chevron-left': 'icon-cis-chevron-left-medium',
  'chevron-right': 'icon-cis-chevron-right-medium',
  'plus': 'icon-cis-plus-medium',
  'refresh': 'icon-cis-arrow-rotate-clockwise',
  'search': 'icon-cis-magnifying-glass',
  'trash': 'icon-cis-trash-can',
  'clipboard': 'icon-cis-clipboard',
  'copy': 'icon-cis-clipboard',
  'download': 'icon-cis-arrow-wall-down',
  'more-horizontal': 'icon-cis-dot-grid-1x3-horizontal',
  'more-vertical': 'icon-cis-dot-grid-1x3-vertical',
  'server': 'icon-cis-server',
  'database': 'icon-d1',
  'cpu': 'icon-edgeworker',
  'layers': 'icon-cis-layers-three',
  'folder': 'icon-cis-folder-1',
  'bar-chart-2': 'icon-chart',
  'settings': 'icon-cis-settings-slider-hor',
  'alert-triangle': 'icon-cis-exclamation-triangle',
  'info': 'icon-cis-circle-info',
  'lock': 'icon-cis-lock',
  'key': 'icon-cis-key-1',
  'user': 'icon-cis-user',
  'globe': 'icon-cis-globe',
  'shield': 'icon-cis-shield-check',
  'clock': 'icon-cis-clock'
}

// Inline fallbacks for external brand icons not in Cloudflare's sprite
const INLINE_FALLBACKS = {
  'github': {
    viewBox: '0 0 256 256',
    content: '<path d="M128,24A104,104,0,0,0,24,128c0,46,29.8,85,71.1,98.8,5.2,1,7.1-2.3,7.1-5v-18c-28.9,6.3-35-13.9-35-13.9-4.7-12-11.5-15.2-11.5-15.2-9.4-6.4.7-6.3.7-6.3,10.4.7,15.9,10.7,15.9,10.7,9.3,15.9,24.3,11.3,30.2,8.6a22.2,22.2,0,0,1,6.6-13.9c-23.1-2.6-47.4-11.5-47.4-51.4a40.2,40.2,0,0,1,10.7-27.9,37.3,37.3,0,0,1,1-27.5s8.7-2.8,28.6,10.7a98.8,98.8,0,0,1,52,0c19.9-13.5,28.6-10.7,28.6-10.7a37.3,37.3,0,0,1,1,27.5,40.2,40.2,0,0,1,10.7,27.9c0,40-24.4,48.7-47.6,51.3a24.9,24.9,0,0,1,7.1,19.3v28.7c0,2.8,1.9,6.1,7.2,5C202.2,213,232,174,232,128A104,104,0,0,0,128,24Z"/>'
  }
}

const resolvedIcon = computed(() => {
  const raw = props.name || ''
  const aliased = ICON_ALIASES[raw]
  if (aliased && iconsCatalog[aliased]) {
    return iconsCatalog[aliased]
  }
  const key = raw.startsWith('icon-') ? raw : `icon-${raw}`
  if (iconsCatalog[key]) return iconsCatalog[key]
  if (iconsCatalog[raw]) return iconsCatalog[raw]
  if (INLINE_FALLBACKS[raw]) return INLINE_FALLBACKS[raw]
  return { viewBox: '0 0 24 24', content: '' }
})
</script>

<style scoped>
.cf-icon {
  display: inline-block;
  vertical-align: middle;
  flex-shrink: 0;
}
</style>
