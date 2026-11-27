<template>
  <div class="cf-log-viewer">
    <div class="cf-log-toolbar">
      <div class="cf-log-title">
        <span class="cf-log-live-dot" :class="{ paused }"></span>
        <span>{{ title }}</span>
      </div>
      <div class="cf-log-controls">
        <input v-model="filter" type="text" placeholder="Filter logs..." class="cf-log-search" />
        <button type="button" class="cf-log-btn" @click="paused = !paused">
          {{ paused ? 'Resume' : 'Pause' }}
        </button>
        <button type="button" class="cf-log-btn" @click="$emit('clear')">Clear</button>
      </div>
    </div>

    <div class="cf-log-body">
      <div v-for="(entry, idx) in filteredLogs" :key="idx" class="cf-log-line">
        <span class="cf-log-ts">{{ entry.time }}</span>
        <span class="cf-log-lvl" :class="`lvl-${(entry.level || 'info').toLowerCase()}`">
          {{ (entry.level || 'INFO').toUpperCase() }}
        </span>
        <span class="cf-log-msg">{{ entry.message }}</span>
      </div>
      <div v-if="filteredLogs.length === 0" class="cf-log-empty">
        No log entries match filter.
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: { type: String, default: 'Live Tail Logs' },
  logs: {
    type: Array,
    default: () => [] // [{ time, level, message }]
  }
})

defineEmits(['clear'])
const filter = ref('')
const paused = ref(false)

const filteredLogs = computed(() => {
  const q = filter.value.trim().toLowerCase()
  if (!q) return props.logs
  return props.logs.filter(
    l => l.message.toLowerCase().includes(q) || (l.level && l.level.toLowerCase().includes(q))
  )
})
</script>

<style scoped>
.cf-log-viewer {
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: #0d1117;
  color: #e6edf3;
  overflow: hidden;
  font-family: var(--font-mono);
}

.cf-log-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 8px 12px;
  background-color: #161b22;
  border-bottom: 1px solid #30363d;
  flex-wrap: wrap;
}

.cf-log-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: var(--font-sans);
  font-size: 12.5px;
  font-weight: 600;
}

.cf-log-live-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #3fb950;
}

.cf-log-live-dot.paused {
  background-color: #d29922;
}

.cf-log-controls {
  display: flex;
  align-items: center;
  gap: 6px;
}

.cf-log-search {
  height: 26px;
  padding: 0 8px;
  border-radius: 4px;
  border: 1px solid #30363d;
  background-color: #0d1117;
  color: #e6edf3;
  font-size: 12px;
  outline: none;
}

.cf-log-btn {
  height: 26px;
  padding: 0 10px;
  border-radius: 4px;
  border: 1px solid #30363d;
  background-color: #21262d;
  color: #e6edf3;
  font-size: 11.5px;
  cursor: pointer;
}

.cf-log-btn:hover {
  background-color: #30363d;
}

.cf-log-body {
  max-height: 210px;
  overflow-y: auto;
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  line-height: 1.5;
}

.cf-log-line {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.cf-log-ts {
  color: #8b949e;
  flex-shrink: 0;
}

.cf-log-lvl {
  font-size: 10.5px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
  flex-shrink: 0;
}

.lvl-info { background: rgba(56, 139, 253, 0.15); color: #58a6ff; }
.lvl-warn { background: rgba(210, 153, 34, 0.18); color: #e3b341; }
.lvl-error { background: rgba(248, 81, 73, 0.18); color: #ff7b72; }

.cf-log-msg {
  word-break: break-word;
}

.cf-log-empty {
  padding: 16px;
  text-align: center;
  color: #8b949e;
}
</style>
