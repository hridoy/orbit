<template>
  <div class="kumo-table-container">
    <!-- Optional Built-in Toolbar -->
    <div v-if="searchable || densityToggle || columnPicker || $slots.toolbar" class="kumo-table-toolbar">
      <div v-if="searchable" class="kumo-table-search">
        <CfInput v-model="searchQuery" :placeholder="searchPlaceholder" size="sm" />
      </div>
      <div class="kumo-table-toolbar-right">
        <slot name="toolbar" />
        <CfPopover v-if="columnPicker" title="Visible Columns" width="210px" align="right">
          <template #trigger>
            <button type="button" class="kumo-density-btn">⚙ Columns</button>
          </template>
          <div style="display: flex; flex-direction: column; gap: 6px;">
            <label v-for="col in columns" :key="col.key" style="display: flex; align-items: center; gap: 8px; font-size: 12.5px; cursor: pointer;">
              <input
                type="checkbox"
                class="kumo-table-checkbox"
                :checked="!hiddenCols.includes(col.key)"
                @change="toggleCol(col.key)"
              />
              <span>{{ col.label }}</span>
            </label>
          </div>
        </CfPopover>
        <button
          v-if="densityToggle"
          type="button"
          class="kumo-density-btn"
          @click="isCompact = !isCompact"
        >
          {{ isCompact ? 'Default Rows' : 'Compact Rows' }}
        </button>
      </div>
    </div>

    <div class="kumo-table-wrapper">
      <table class="kumo-table" :class="{ 'is-compact': isCompact }">
        <thead>
          <tr>
            <th v-if="selectable" style="width: 40px; text-align: center;">
              <input
                type="checkbox"
                class="kumo-table-checkbox"
                :checked="isAllSelected"
                @change="$emit('toggle-select-all')"
              />
            </th>
            <th
              v-for="col in visibleColumns"
              :key="col.key"
              :style="{ width: col.width, textAlign: col.align || 'left' }"
              :class="{ 'is-sortable': col.sortable }"
              @click="col.sortable && toggleSort(col.key)"
            >
              <span class="kumo-th-inner" :style="{ justifyContent: col.align === 'right' ? 'flex-end' : 'flex-start' }">
                <span>{{ col.label }}</span>
                <span v-if="col.sortable" class="kumo-sort-indicator">
                  {{ sortKey === col.key ? (sortDir === 'asc' ? '↑' : '↓') : '↕' }}
                </span>
              </span>
            </th>
            <th v-if="$slots.actions" style="width: 48px;"></th>
          </tr>
        </thead>
        <tbody>
          <!-- Loading State Skeletons -->
          <template v-if="loading">
            <tr v-for="i in skeletonRows" :key="i">
              <td v-if="selectable" style="text-align: center;">
                <CfSkeleton width="16px" height="16px" />
              </td>
              <td v-for="col in visibleColumns" :key="col.key">
                <CfSkeleton :width="col.skeletonWidth || '65%'" height="12px" />
              </td>
              <td v-if="$slots.actions">
                <CfSkeleton width="20px" height="20px" />
              </td>
            </tr>
          </template>

          <!-- Empty State -->
          <tr v-else-if="processedData.length === 0">
            <td :colspan="totalColumns" class="kumo-table-empty">
              <slot name="empty">
                <CfIcon name="cis-circle-info" size="24" style="color: var(--text-color-kumo-subtle); margin-bottom: 0.5rem;" />
                <div>No records found</div>
              </slot>
            </td>
          </tr>

          <!-- Data Rows -->
          <tr
            v-for="(row, idx) in processedData"
            v-else
            :key="rowKey ? row[rowKey] : idx"
            :class="{ 'is-selected': selectedKeys && selectedKeys.includes(row[rowKey]) }"
            @click="$emit('row-click', row)"
          >
            <td v-if="selectable" style="text-align: center;" @click.stop>
              <input
                type="checkbox"
                class="kumo-table-checkbox"
                :checked="selectedKeys && selectedKeys.includes(row[rowKey])"
                @change="$emit('toggle-select', row[rowKey])"
              />
            </td>
            <td
              v-for="col in visibleColumns"
              :key="col.key"
              :style="{ textAlign: col.align || 'left' }"
            >
              <slot :name="`col-${col.key}`" :row="row" :val="row[col.key]">
                {{ row[col.key] !== undefined && row[col.key] !== null ? row[col.key] : '—' }}
              </slot>
            </td>
            <td v-if="$slots.actions" class="kumo-table-actions-cell" @click.stop>
              <slot name="actions" :row="row" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  columns: {
    type: Array,
    required: true // [{ key, label, width, align, skeletonWidth, sortable }]
  },
  data: {
    type: Array,
    default: () => []
  },
  rowKey: {
    type: String,
    default: 'id'
  },
  loading: {
    type: Boolean,
    default: false
  },
  skeletonRows: {
    type: Number,
    default: 3
  },
  selectable: {
    type: Boolean,
    default: false
  },
  selectedKeys: {
    type: Array,
    default: () => []
  },
  searchable: {
    type: Boolean,
    default: false
  },
  searchPlaceholder: {
    type: String,
    default: 'Filter table rows...'
  },
  densityToggle: {
    type: Boolean,
    default: false
  },
  columnPicker: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-select-all', 'toggle-select', 'row-click'])

const searchQuery = ref('')
const sortKey = ref('')
const sortDir = ref('asc')
const isCompact = ref(false)
const hiddenCols = ref([])

function toggleCol(key) {
  const idx = hiddenCols.value.indexOf(key)
  if (idx > -1) hiddenCols.value.splice(idx, 1)
  else if (hiddenCols.value.length < props.columns.length - 1) hiddenCols.value.push(key)
}

const visibleColumns = computed(() =>
  props.columns.filter(c => !hiddenCols.value.includes(c.key))
)

function toggleSort(key) {
  if (sortKey.value === key) {
    if (sortDir.value === 'asc') sortDir.value = 'desc'
    else {
      sortKey.value = ''
      sortDir.value = 'asc'
    }
  } else {
    sortKey.value = key
    sortDir.value = 'asc'
  }
}

const processedData = computed(() => {
  let rows = [...(props.data || [])]
  if (props.searchable && searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    rows = rows.filter(r =>
      props.columns.some(c => String(r[c.key] ?? '').toLowerCase().includes(q))
    )
  }
  if (sortKey.value) {
    const k = sortKey.value
    const dir = sortDir.value === 'asc' ? 1 : -1
    rows.sort((a, b) => {
      const av = String(a[k] ?? '')
      const bv = String(b[k] ?? '')
      return av.localeCompare(bv, undefined, { numeric: true }) * dir
    })
  }
  return rows
})

const isAllSelected = computed(() => {
  return props.data && props.data.length > 0 && props.selectedKeys.length === props.data.length
})

const totalColumns = computed(() => {
  let count = visibleColumns.value.length
  if (props.selectable) count++
  return count + 1
})
</script>

<style scoped>
.kumo-table-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.kumo-table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.kumo-table-search {
  width: 260px;
  max-width: 100%;
}

.kumo-table-toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.kumo-density-btn {
  height: 28px;
  padding: 0 10px;
  border-radius: 6px;
  border: 1px solid var(--color-kumo-line);
  background: var(--color-kumo-base);
  color: var(--text-color-kumo-subtle);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
}

.kumo-density-btn:hover {
  background: var(--color-kumo-tint);
  color: var(--text-color-kumo-strong);
}

.kumo-table-wrapper {
  width: 100%;
  background-color: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  border-radius: 8px;
  overflow-x: auto;
  box-shadow: var(--shadow-kumo-card, 0 1px 2px rgba(0, 0, 0, 0.04));
  box-sizing: border-box;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

.kumo-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.kumo-table thead th {
  height: 40px;
  padding: 8px 14px;
  background-color: var(--color-kumo-canvas);
  color: var(--text-color-kumo-subtle);
  font-size: 13px;
  font-weight: 500;
  border-bottom: 1px solid var(--color-kumo-line);
  box-sizing: border-box;
  user-select: none;
}

.kumo-table thead th.is-sortable {
  cursor: pointer;
}

.kumo-table thead th.is-sortable:hover {
  color: var(--text-color-kumo-strong);
}

.kumo-th-inner {
  display: flex;
  align-items: center;
  gap: 6px;
}

.kumo-sort-indicator {
  font-size: 11px;
  opacity: 0.65;
}

.kumo-table-checkbox {
  accent-color: var(--color-kumo-brand);
  cursor: pointer;
}

.kumo-table tbody tr {
  border-top: 1px solid var(--color-kumo-hairline);
  height: 44px;
  transition: background-color 0.12s ease;
}

.kumo-table.is-compact tbody tr {
  height: 34px;
}

.kumo-table tbody tr:hover {
  background-color: var(--color-kumo-overlay);
}

.kumo-table tbody tr.is-selected {
  background-color: rgba(0, 81, 195, 0.06);
}

.kumo-table tbody td {
  height: 44px;
  padding: 8px 16px;
  font-size: 13px;
  color: var(--text-color-kumo-default);
  vertical-align: middle;
  font-variant-numeric: tabular-nums;
  box-sizing: border-box;
}

.kumo-table.is-compact tbody td {
  height: 34px;
  padding: 4px 16px;
  font-size: 12.5px;
}

.kumo-table-empty {
  text-align: center;
  padding: 3rem 1rem !important;
  color: var(--text-color-kumo-subtle);
}
</style>
