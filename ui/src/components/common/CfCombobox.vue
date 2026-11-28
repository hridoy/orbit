<template>
  <div ref="rootRef" class="cf-combobox">
    <label v-if="label" class="cf-combobox-label">{{ label }}</label>
    <button
      type="button"
      class="cf-combobox-trigger"
      :class="{ 'is-open': open }"
      @click="open = !open"
    >
      <span class="cf-combobox-value">{{ selectedLabel || placeholder }}</span>
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="cf-combobox-caret">
        <path d="m7 15 5 5 5-5"/><path d="m7 9 5-5 5 5"/>
      </svg>
    </button>

    <div v-if="open" class="cf-combobox-popover">
      <div class="cf-combobox-search-row">
        <CfIcon name="cis-magnifying-glass" size="14" />
        <input
          v-model="query"
          type="text"
          :placeholder="searchPlaceholder"
          class="cf-combobox-search-input"
        />
      </div>

      <div class="cf-combobox-list" role="listbox">
        <button
          v-for="opt in filteredOptions"
          :key="opt.value"
          type="button"
          role="option"
          :aria-selected="modelValue === opt.value"
          class="cf-combobox-option"
          :class="{ 'is-selected': modelValue === opt.value }"
          @click="selectOption(opt)"
        >
          <div class="cf-combobox-opt-text">
            <span class="cf-combobox-opt-label">{{ opt.label }}</span>
            <span v-if="opt.meta" class="cf-combobox-opt-meta">{{ opt.meta }}</span>
          </div>
          <span v-if="modelValue === opt.value" class="cf-combobox-check">✓</span>
        </button>
        <div v-if="filteredOptions.length === 0" class="cf-combobox-empty">
          No matches for "{{ query }}"
        </div>
      </div>

      <div v-if="$slots.footer" class="cf-combobox-footer">
        <slot name="footer" :close="() => (open = false)" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  options: { type: Array, default: () => [] }, // [{ value, label, meta? }]
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Select...' },
  searchPlaceholder: { type: String, default: 'Search...' }
})

const emit = defineEmits(['update:modelValue', 'change'])
const open = ref(false)
const query = ref('')
const rootRef = ref(null)

const normalizedOptions = computed(() =>
  props.options.map(o => (typeof o === 'object' ? o : { value: o, label: String(o) }))
)

const selectedLabel = computed(() => {
  const found = normalizedOptions.value.find(o => o.value === props.modelValue)
  return found ? found.label : ''
})

const filteredOptions = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return normalizedOptions.value
  return normalizedOptions.value.filter(
    o => o.label.toLowerCase().includes(q) || (o.meta && o.meta.toLowerCase().includes(q))
  )
})

function selectOption(opt) {
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
  open.value = false
  query.value = ''
}

function onDocClick(e) {
  if (rootRef.value && !rootRef.value.contains(e.target)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))
</script>

<style scoped>
.cf-combobox {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-combobox-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-combobox-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  height: 36px;
  padding: 0 12px;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  color: var(--text-color-kumo-default);
  font-family: var(--font-sans);
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s ease;
}

.cf-combobox-trigger:hover,
.cf-combobox-trigger.is-open {
  border-color: var(--color-kumo-brand);
}

.cf-combobox-caret {
  color: var(--text-color-kumo-subtle);
  flex-shrink: 0;
}

.cf-combobox-popover {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 50;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-elevated);
  box-shadow: var(--shadow-kumo-dialog);
  overflow: hidden;
}

.cf-combobox-search-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  color: var(--text-color-kumo-subtle);
}

.cf-combobox-search-input {
  flex: 1;
  border: none;
  background: transparent;
  outline: none;
  font-size: 13px;
  color: var(--text-color-kumo-default);
}

.cf-combobox-list {
  max-height: 210px;
  overflow-y: auto;
  padding: 4px;
}

.cf-combobox-option {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 7px 10px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--text-color-kumo-default);
  text-align: left;
  cursor: pointer;
  font-size: 13px;
}

.cf-combobox-option:hover,
.cf-combobox-option.is-selected {
  background-color: var(--color-kumo-tint);
}

.cf-combobox-opt-text {
  display: flex;
  flex-direction: column;
}

.cf-combobox-opt-meta {
  font-size: 11px;
  font-family: var(--font-mono);
  color: var(--text-color-kumo-subtle);
}

.cf-combobox-check {
  color: var(--color-kumo-brand);
  font-weight: 700;
}

.cf-combobox-empty {
  padding: 12px;
  text-align: center;
  font-size: 12.5px;
  color: var(--text-color-kumo-subtle);
}

.cf-combobox-footer {
  padding: 6px 8px;
  border-top: 1px solid var(--color-kumo-hairline);
  background-color: var(--color-kumo-canvas);
}
</style>
