<template>
  <div class="cf-notif-card">
    <div class="cf-notif-head">
      <div style="display: flex; align-items: center; gap: 8px;">
        <strong>Notifications</strong>
        <CfBadge v-if="unreadCount" variant="brand" size="sm">{{ unreadCount }} new</CfBadge>
      </div>
      <button v-if="unreadCount" type="button" class="cf-notif-mark-btn" @click="$emit('mark-all-read')">
        Mark all read
      </button>
    </div>

    <div class="cf-notif-list">
      <div
        v-for="(n, idx) in items"
        :key="idx"
        class="cf-notif-row"
        :class="{ 'is-unread': !n.read }"
        @click="$emit('select', n)"
      >
        <span class="cf-notif-dot" :class="{ active: !n.read }"></span>
        <div class="cf-notif-body">
          <div class="cf-notif-title">{{ n.title }}</div>
          <div class="cf-notif-msg">{{ n.message }}</div>
          <div class="cf-notif-time">{{ n.time }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    required: true // [{ title, message, time, read? }]
  }
})

defineEmits(['mark-all-read', 'select'])
const unreadCount = computed(() => props.items.filter(i => !i.read).length)
</script>

<style scoped>
.cf-notif-card {
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-notif-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-notif-mark-btn {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-link);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
}

.cf-notif-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.cf-notif-row:last-child {
  border-bottom: none;
}

.cf-notif-row:hover {
  background-color: var(--color-kumo-tint);
}

.cf-notif-row.is-unread {
  background-color: var(--color-kumo-info-tint);
}

.cf-notif-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  margin-top: 6px;
  background: transparent;
  flex-shrink: 0;
}

.cf-notif-dot.active {
  background-color: var(--color-kumo-brand);
}

.cf-notif-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cf-notif-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-notif-msg {
  font-size: 12.5px;
  color: var(--text-color-kumo-default);
}

.cf-notif-time {
  font-size: 11px;
  color: var(--text-color-kumo-subtle);
}
</style>
