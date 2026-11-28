<template>
  <div class="cf-joined-toolbar">
    <div class="cf-jt-search">
      <CfIcon name="cis-magnifying-glass" size="14" class="cf-jt-icon" />
      <input
        :value="search"
        type="text"
        :placeholder="placeholder"
        class="cf-jt-input"
        @input="$emit('update:search', $event.target.value)"
      />
    </div>

    <select
      :value="filter"
      class="cf-jt-select"
      @change="$emit('update:filter', $event.target.value)"
    >
      <option v-for="opt in filterOptions" :key="opt" :value="opt">{{ opt }}</option>
    </select>

    <select
      :value="sort"
      class="cf-jt-select"
      @change="$emit('update:sort', $event.target.value)"
    >
      <option v-for="opt in sortOptions" :key="opt" :value="opt">{{ opt }}</option>
    </select>

    <button type="button" class="cf-jt-refresh" title="Refresh" @click="$emit('refresh')">
      ↻
    </button>
  </div>
</template>

<script setup>
defineProps({
  search: { type: String, default: '' },
  filter: { type: String, default: 'Show all' },
  sort: { type: String, default: 'Last modified' },
  placeholder: { type: String, default: 'Search applications...' },
  filterOptions: { type: Array, default: () => ['Show all', 'Workers', 'Pages'] },
  sortOptions: { type: Array, default: () => ['Last modified', 'Name', 'Requests', 'Errors'] }
})

defineEmits(['update:search', 'update:filter', 'update:sort', 'refresh'])
</script>

<style scoped>
.cf-joined-toolbar {
  display: flex;
  align-items: stretch;
  height: 36px;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  overflow: hidden;
  width: 100%;
}

.cf-jt-search {
  flex: 1;
  min-width: 160px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
}

.cf-jt-icon {
  color: var(--text-color-kumo-subtle);
  flex-shrink: 0;
}

.cf-jt-input {
  flex: 1;
  border: none;
  background: transparent;
  outline: none;
  font-size: 13px;
  color: var(--text-color-kumo-default);
}

.cf-jt-select,
.cf-jt-refresh {
  border: none;
  border-left: 1px solid var(--kumo-control-border-default);
  background: transparent;
  color: var(--text-color-kumo-default);
  font-family: var(--font-sans);
  font-size: 13px;
  padding: 0 12px;
  cursor: pointer;
  outline: none;
}

.cf-jt-select:hover,
.cf-jt-refresh:hover {
  background-color: var(--color-kumo-tint);
}

.cf-jt-refresh {
  padding: 0 12px;
  font-size: 15px;
}
</style>
