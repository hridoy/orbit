<template>
  <ul class="cf-tree" role="tree">
    <li v-for="node in nodes" :key="node.id" class="cf-tree-node">
      <div
        class="cf-tree-row"
        :class="{ 'is-selected': selectedId === node.id }"
        @click="onNodeClick(node)"
      >
        <span class="cf-tree-chevron" :class="{ invisible: !node.children?.length }">
          {{ expanded[node.id] ? '▾' : '▸' }}
        </span>
        <CfIcon v-if="node.icon" :name="node.icon" size="14" class="cf-tree-icon" />
        <span class="cf-tree-label">{{ node.label }}</span>
        <span v-if="node.badge" class="cf-tree-badge">{{ node.badge }}</span>
      </div>

      <ul v-if="node.children?.length && expanded[node.id]" class="cf-tree-children">
        <li
          v-for="child in node.children"
          :key="child.id"
          class="cf-tree-row is-leaf"
          :class="{ 'is-selected': selectedId === child.id }"
          @click.stop="$emit('select', child)"
        >
          <CfIcon v-if="child.icon" :name="child.icon" size="13" class="cf-tree-icon" />
          <span class="cf-tree-label">{{ child.label }}</span>
          <span v-if="child.badge" class="cf-tree-badge">{{ child.badge }}</span>
        </li>
      </ul>
    </li>
  </ul>
</template>

<script setup>
import { reactive } from 'vue'

const props = defineProps({
  nodes: { type: Array, required: true }, // [{ id, label, icon?, badge?, children? }]
  selectedId: { type: String, default: '' },
  defaultExpanded: { type: Array, default: () => [] }
})

const emit = defineEmits(['select'])
const expanded = reactive(
  Object.fromEntries(props.defaultExpanded.map(id => [id, true]))
)

function onNodeClick(node) {
  if (node.children?.length) {
    expanded[node.id] = !expanded[node.id]
  }
  emit('select', node)
}
</script>

<style scoped>
.cf-tree,
.cf-tree-children {
  list-style: none;
  margin: 0;
  padding: 0;
}

.cf-tree-children {
  padding-left: 18px;
  border-left: 1px solid var(--color-kumo-hairline);
  margin-left: 10px;
}

.cf-tree-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: 6px;
  font-size: 12.5px;
  color: var(--text-color-kumo-default);
  cursor: pointer;
  user-select: none;
}

.cf-tree-row:hover {
  background-color: var(--color-kumo-tint);
}

.cf-tree-row.is-selected {
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-strong);
  font-weight: 600;
}

.cf-tree-chevron {
  width: 12px;
  font-size: 11px;
  color: var(--text-color-kumo-subtle);
}

.cf-tree-chevron.invisible {
  visibility: hidden;
}

.cf-tree-icon {
  color: var(--text-color-kumo-subtle);
  flex-shrink: 0;
}

.cf-tree-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cf-tree-badge {
  font-family: var(--font-mono);
  font-size: 10.5px;
  color: var(--text-color-kumo-subtle);
  background: var(--color-kumo-recessed);
  padding: 1px 5px;
  border-radius: 4px;
}
</style>
