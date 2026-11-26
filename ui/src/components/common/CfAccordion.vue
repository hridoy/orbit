<template>
  <div class="kumo-accordion">
    <!-- Multi-item mode via `items` prop -->
    <template v-if="items && items.length">
      <div
        v-for="(item, idx) in items"
        :key="item.value || idx"
        :class="['kumo-accordion-item', { 'is-expanded': isExpanded(item.value || idx) }]"
      >
        <button
          type="button"
          class="kumo-accordion-header"
          :aria-expanded="isExpanded(item.value || idx) ? 'true' : 'false'"
          @click="toggle(item.value || idx)"
        >
          <div class="kumo-accordion-title-row">
            <CfIcon v-if="item.icon" :name="item.icon" size="16" class="kumo-accordion-icon" />
            <span class="kumo-accordion-title">{{ item.title }}</span>
            <span v-if="item.badge" class="kumo-accordion-badge">{{ item.badge }}</span>
          </div>
          <CfIcon name="cis-chevron-down-medium" size="14" class="kumo-accordion-chevron" />
        </button>

        <div v-show="isExpanded(item.value || idx)" class="kumo-accordion-content">
          <slot :name="`item-${item.value || idx}`" :item="item">
            {{ item.content }}
          </slot>
        </div>
      </div>
    </template>

    <!-- Single-item standalone mode via `title` prop + default slot -->
    <div
      v-else
      :class="['kumo-accordion-item', { 'is-expanded': singleOpen }]"
    >
      <button
        type="button"
        class="kumo-accordion-header"
        :aria-expanded="singleOpen ? 'true' : 'false'"
        @click="singleOpen = !singleOpen"
      >
        <div class="kumo-accordion-title-row">
          <CfIcon v-if="icon" :name="icon" size="16" class="kumo-accordion-icon" />
          <span class="kumo-accordion-title">{{ title }}</span>
          <span v-if="badge" class="kumo-accordion-badge">{{ badge }}</span>
        </div>
        <CfIcon name="cis-chevron-down-medium" size="14" class="kumo-accordion-chevron" />
      </button>

      <div v-show="singleOpen" class="kumo-accordion-content">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => [] // [{ title, content, value, icon, badge }]
  },
  title: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: ''
  },
  badge: {
    type: String,
    default: ''
  },
  defaultOpen: {
    type: Boolean,
    default: false
  },
  multiple: {
    type: Boolean,
    default: false
  },
  defaultExpanded: {
    type: [Array, String, Number],
    default: () => []
  }
})

const singleOpen = ref(props.defaultOpen)

const expandedKeys = ref(
  Array.isArray(props.defaultExpanded)
    ? [...props.defaultExpanded]
    : [props.defaultExpanded].filter(Boolean)
)

const isExpanded = (key) => expandedKeys.value.includes(key)

const toggle = (key) => {
  if (isExpanded(key)) {
    expandedKeys.value = expandedKeys.value.filter(k => k !== key)
  } else {
    if (props.multiple) {
      expandedKeys.value.push(key)
    } else {
      expandedKeys.value = [key]
    }
  }
}
</script>

<style scoped>
.kumo-accordion {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 100%;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  box-sizing: border-box;
}

.kumo-accordion-item {
  border: 1px solid var(--color-kumo-line, #e2e8f0);
  border-radius: 8px;
  background-color: var(--color-kumo-base, #ffffff);
  overflow: hidden;
  transition: border-color 0.15s ease;
}

.kumo-accordion-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0.875rem 1.25rem;
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: left;
  user-select: none;
}

.kumo-accordion-header:hover {
  background-color: var(--color-kumo-tint, #f8fafc);
}

.kumo-accordion-title-row {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.kumo-accordion-icon {
  color: var(--text-color-kumo-subtle, #71717a);
}

.kumo-accordion-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-kumo-default, #18181b);
}

.kumo-accordion-badge {
  font-size: 11px;
  font-weight: 500;
  padding: 1px 7px;
  border-radius: 9999px;
  background-color: var(--color-kumo-recessed, #f4f4f5);
  color: var(--text-color-kumo-subtle, #52525b);
}

.kumo-accordion-chevron {
  color: var(--text-color-kumo-subtle, #71717a);
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.kumo-accordion-item.is-expanded .kumo-accordion-chevron {
  transform: rotate(180deg);
}

.kumo-accordion-content {
  padding: 0 1.25rem 1rem 1.25rem;
  font-size: 13.5px;
  line-height: 1.6;
  color: var(--text-color-kumo-subtle, #52525b);
  border-top: 1px solid var(--color-kumo-hairline, #f4f4f5);
  padding-top: 0.875rem;
}
</style>
