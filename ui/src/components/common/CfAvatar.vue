<script setup>
import { computed } from 'vue'

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  src: {
    type: String,
    default: ''
  },
  size: {
    type: String,
    default: 'md', // 'xs' | 'sm' | 'md' | 'lg' | 'xl'
    validator: (v) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(v)
  },
  status: {
    type: String,
    default: '', // 'online' | 'busy' | 'away' | 'offline' | ''
    validator: (v) => ['', 'online', 'busy', 'away', 'offline'].includes(v)
  }
})

const initials = computed(() => {
  if (!props.name) return ''
  const parts = props.name.trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return parts[0].slice(0, 2).toUpperCase()
})

const bgHue = computed(() => {
  if (!props.name) return 210
  let hash = 0
  for (let i = 0; i < props.name.length; i++) {
    hash = props.name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return Math.abs(hash % 360)
})
</script>

<template>
  <div class="cf-avatar-wrap" :class="[`size-${size}`]">
    <img v-if="src" :src="src" :alt="name" class="cf-avatar-img" />
    <span 
      v-else 
      class="cf-avatar-fallback" 
      :style="{ backgroundColor: `hsl(${bgHue}, 65%, 45%)` }"
    >
      {{ initials || '?' }}
    </span>

    <span 
      v-if="status" 
      class="cf-avatar-status" 
      :class="[`status-${status}`]"
    ></span>
  </div>
</template>

<style scoped>
.cf-avatar-wrap {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  flex-shrink: 0;
  user-select: none;
}

.size-xs {
  width: 20px;
  height: 20px;
  font-size: 9px;
}

.size-sm {
  width: 24px;
  height: 24px;
  font-size: 10px;
}

.size-md {
  width: 32px;
  height: 32px;
  font-size: 12px;
}

.size-lg {
  width: 40px;
  height: 40px;
  font-size: 14px;
}

.size-xl {
  width: 56px;
  height: 56px;
  font-size: 18px;
}

.cf-avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.cf-avatar-fallback {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.cf-avatar-status {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 28%;
  height: 28%;
  min-width: 6px;
  min-height: 6px;
  border-radius: 50%;
  border: 2px solid var(--kumo-surface-default, #ffffff);
}

.status-online {
  background-color: #10b981;
}

.status-busy {
  background-color: #ef4444;
}

.status-away {
  background-color: #f59e0b;
}

.status-offline {
  background-color: #9ca3af;
}
</style>
