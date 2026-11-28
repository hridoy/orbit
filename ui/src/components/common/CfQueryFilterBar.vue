<template>
  <div class="cf-filter-bar">
    <div v-for="(chip, idx) in modelValue" :key="idx" class="cf-filter-chip">
      <span class="cf-fc-field">{{ chip.field }}</span>
      <span class="cf-fc-op">{{ chip.op }}</span>
      <span class="cf-fc-val">{{ chip.value }}</span>
      <button type="button" class="cf-fc-rm" @click="removeChip(idx)">×</button>
    </div>

    <div class="cf-filter-adder">
      <select v-model="draftField" class="cf-fa-select">
        <option v-for="f in fields" :key="f" :value="f">{{ f }}</option>
      </select>
      <select v-model="draftOp" class="cf-fa-select">
        <option value="=">=</option>
        <option value="!=">!=</option>
        <option value="contains">contains</option>
      </select>
      <input
        v-model="draftVal"
        type="text"
        placeholder="Value..."
        class="cf-fa-input"
        @keydown.enter.prevent="addChip"
      />
      <CfButton variant="secondary" size="sm" @click="addChip">+ Add</CfButton>
      <button v-if="modelValue.length" type="button" class="cf-fa-clear" @click="$emit('update:modelValue', [])">
        Clear
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  modelValue: { type: Array, default: () => [] }, // [{ field, op, value }]
  fields: { type: Array, default: () => ['status', 'actor', 'resource', 'action'] }
})

const emit = defineEmits(['update:modelValue'])
const draftField = ref(props.fields[0] || 'status')
const draftOp = ref('=')
const draftVal = ref('')

function addChip() {
  const val = draftVal.value.trim()
  if (!val) return
  emit('update:modelValue', [...props.modelValue, { field: draftField.value, op: draftOp.value, value: val }])
  draftVal.value = ''
}

function removeChip(idx) {
  const next = [...props.modelValue]
  next.splice(idx, 1)
  emit('update:modelValue', next)
}
</script>

<style scoped>
.cf-filter-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
}

.cf-filter-chip {
  display: inline-flex;
  align-items: center;
  height: 26px;
  border-radius: 6px;
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-recessed);
  font-size: 12px;
  overflow: hidden;
}

.cf-fc-field {
  padding: 0 7px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-fc-op {
  padding: 0 5px;
  border-left: 1px solid var(--color-kumo-line);
  border-right: 1px solid var(--color-kumo-line);
  color: var(--text-color-kumo-subtle);
  font-family: var(--font-mono);
}

.cf-fc-val {
  padding: 0 7px;
  font-family: var(--font-mono);
  color: var(--text-color-kumo-link);
}

.cf-fc-rm {
  border: none;
  border-left: 1px solid var(--color-kumo-line);
  background: transparent;
  color: var(--text-color-kumo-subtle);
  padding: 0 6px;
  height: 100%;
  cursor: pointer;
}

.cf-fc-rm:hover {
  color: var(--text-color-kumo-danger);
}

.cf-filter-adder {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.cf-fa-select,
.cf-fa-input {
  height: 28px;
  padding: 0 8px;
  border-radius: 6px;
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  color: var(--text-color-kumo-default);
  font-size: 12px;
  outline: none;
}

.cf-fa-clear {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-subtle);
  font-size: 12px;
  cursor: pointer;
  text-decoration: underline;
}
</style>
