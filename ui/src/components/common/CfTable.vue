<template>
  <div class="kumo-table-wrapper">
    <table class="kumo-table">
      <thead>
        <tr>
          <th v-if="selectable" style="width: 40px; text-align: center;">
            <input
              type="checkbox"
              :checked="isAllSelected"
              @change="$emit('toggle-select-all')"
            />
          </th>
          <th
            v-for="col in columns"
            :key="col.key"
            :style="{ width: col.width, textAlign: col.align || 'left' }"
          >
            {{ col.label }}
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
            <td v-for="col in columns" :key="col.key">
              <CfSkeleton :width="col.skeletonWidth || '65%'" height="12px" />
            </td>
            <td v-if="$slots.actions">
              <CfSkeleton width="20px" height="20px" />
            </td>
          </tr>
        </template>

        <!-- Empty State -->
        <tr v-else-if="!data || data.length === 0">
          <td :colspan="totalColumns" class="kumo-table-empty">
            <slot name="empty">
              <CfIcon name="cis-circle-info" size="24" style="color: #a1a1aa; margin-bottom: 0.5rem;" />
              <div>No records found</div>
            </slot>
          </td>
        </tr>

        <!-- Data Rows -->
        <tr
          v-for="(row, idx) in data"
          v-else
          :key="rowKey ? row[rowKey] : idx"
          :class="{ 'is-selected': selectedKeys && selectedKeys.includes(row[rowKey]) }"
          @click="$emit('row-click', row)"
        >
          <td v-if="selectable" style="text-align: center;" @click.stop>
            <input
              type="checkbox"
              :checked="selectedKeys && selectedKeys.includes(row[rowKey])"
              @change="$emit('toggle-select', row[rowKey])"
            />
          </td>
          <td
            v-for="col in columns"
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
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  columns: {
    type: Array,
    required: true // [{ key, label, width, align, skeletonWidth }]
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
  }
})

defineEmits(['toggle-select-all', 'toggle-select', 'row-click'])

const isAllSelected = computed(() => {
  return props.data && props.data.length > 0 && props.selectedKeys.length === props.data.length
})

const totalColumns = computed(() => {
  let count = props.columns.length
  if (props.selectable) count++
  return count + 1 // including optional actions slot
})
</script>

<style scoped>
.kumo-table-wrapper {
  width: 100%;
  background-color: var(--color-kumo-base, #ffffff);
  border: 1px solid var(--color-kumo-line, #e2e8f0);
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
  height: 44px;
  padding: 8px 14px;
  background-color: var(--color-kumo-canvas, #fafafa);
  color: var(--text-color-kumo-subtle, #71717a);
  font-size: 13px;
  font-weight: 500;
  border-bottom: 1px solid var(--color-kumo-line, #e2e8f0);
  box-sizing: border-box;
}

[data-theme="dark"] .kumo-table thead th {
  background-color: #0c0f17;
  border-bottom-color: #1e293b;
  color: #a1a1aa;
}

.kumo-table tbody tr {
  border-top: 1px solid var(--color-kumo-line, #f4f4f5);
  height: 44px;
  transition: background-color 0.12s ease;
}

[data-theme="dark"] .kumo-table tbody tr {
  border-top-color: #1e293b;
}

.kumo-table tbody tr:hover {
  background-color: #f8fafc;
}

[data-theme="dark"] .kumo-table tbody tr:hover {
  background-color: rgba(255, 255, 255, 0.02);
}

.kumo-table tbody tr.is-selected {
  background-color: rgba(0, 81, 195, 0.03);
}

.kumo-table tbody td {
  height: 44px;
  padding: 8px 16px;
  font-size: 13px;
  color: var(--text-color-kumo-default, #18181b);
  vertical-align: middle;
  font-variant-numeric: tabular-nums;
  box-sizing: border-box;
}

[data-theme="dark"] .kumo-table tbody td {
  color: #f8fafc;
}

.kumo-table-empty {
  text-align: center;
  padding: 3rem 1rem !important;
  color: #71717a;
}
</style>
