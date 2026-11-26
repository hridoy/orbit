<template>
  <div class="kumo-pagination">
    <span class="kumo-pagination-counter">
      Showing {{ startItem }}–{{ endItem }} of {{ total }} records
    </span>

    <div class="kumo-pagination-group">
      <button
        type="button"
        class="kumo-page-nav-btn"
        :disabled="currentPage <= 1"
        title="First page"
        @click="goToPage(1)"
      >
        «
      </button>
      <button
        type="button"
        class="kumo-page-nav-btn"
        :disabled="currentPage <= 1"
        title="Previous page"
        @click="goToPage(currentPage - 1)"
      >
        ‹
      </button>

      <span class="kumo-page-current-pill">
        {{ currentPage }} / {{ totalPages || 1 }}
      </span>

      <button
        type="button"
        class="kumo-page-nav-btn"
        :disabled="currentPage >= totalPages"
        title="Next page"
        @click="goToPage(currentPage + 1)"
      >
        ›
      </button>
      <button
        type="button"
        class="kumo-page-nav-btn"
        :disabled="currentPage >= totalPages"
        title="Last page"
        @click="goToPage(totalPages)"
      >
        »
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  currentPage: {
    type: Number,
    default: 1
  },
  total: {
    type: Number,
    required: true
  },
  pageSize: {
    type: Number,
    default: 10
  }
})

const emit = defineEmits(['update:currentPage', 'page-change'])

const totalPages = computed(() => Math.ceil(props.total / props.pageSize) || 1)

const startItem = computed(() => {
  if (props.total === 0) return 0
  return (props.currentPage - 1) * props.pageSize + 1
})

const endItem = computed(() => {
  return Math.min(props.currentPage * props.pageSize, props.total)
})

const goToPage = (p) => {
  if (p < 1 || p > totalPages.value || p === props.currentPage) return
  emit('update:currentPage', p)
  emit('page-change', p)
}
</script>

<style scoped>
.kumo-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.25rem;
  flex-wrap: wrap;
  padding: 0.75rem 0;
  font-family: var(--font-sans);
  font-size: 13px;
  color: var(--text-color-kumo-subtle, #71717a);
  letter-spacing: -0.01em;
  box-sizing: border-box;
}

.kumo-pagination-group {
  display: inline-flex;
  align-items: center;
  border: 1px solid var(--color-kumo-line, #d1d5db);
  border-radius: 6px;
  background-color: var(--color-kumo-base, #ffffff);
  overflow: hidden;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.kumo-page-nav-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-color-kumo-default, #18181b);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
}

[data-theme="dark"] .kumo-page-nav-btn {
  color: #f8fafc;
}

.kumo-page-nav-btn:hover:not(:disabled) {
  background-color: #f4f4f5;
}

[data-theme="dark"] .kumo-page-nav-btn:hover:not(:disabled) {
  background-color: #151a26;
}

.kumo-page-nav-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.kumo-page-current-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 10px;
  height: 32px;
  font-size: 12.5px;
  font-weight: 500;
  color: var(--text-color-kumo-default, #18181b);
  border-left: 1px solid var(--color-kumo-line, #e2e8f0);
  border-right: 1px solid var(--color-kumo-line, #e2e8f0);
  background-color: #fafafa;
  font-variant-numeric: tabular-nums;
}

[data-theme="dark"] .kumo-page-current-pill {
  background-color: #151a26;
  border-color: #27272a;
  color: #f8fafc;
}
</style>
