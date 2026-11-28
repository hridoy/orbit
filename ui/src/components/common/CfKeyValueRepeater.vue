<template>
  <div class="cf-kv-repeater">
    <div v-for="(row, idx) in modelValue" :key="idx" class="cf-kv-card">
      <div class="cf-kv-row">
        <span class="cf-kv-tag">Key</span>
        <input
          :value="row.key"
          type="text"
          placeholder="VARIABLE_NAME"
          class="cf-kv-input"
          @input="updateField(idx, 'key', $event.target.value)"
        />
        <button
          v-if="modelValue.length > 1"
          type="button"
          class="cf-kv-rm"
          title="Remove row"
          @click="removeRow(idx)"
        >
          ×
        </button>
      </div>
      <div class="cf-kv-row">
        <span class="cf-kv-tag">Value</span>
        <input
          :value="row.value"
          :type="row.secret ? 'password' : 'text'"
          placeholder="value..."
          class="cf-kv-input"
          @input="updateField(idx, 'value', $event.target.value)"
        />
        <label class="cf-kv-secret">
          <input
            type="checkbox"
            :checked="row.secret"
            @change="updateField(idx, 'secret', $event.target.checked)"
          />
          <span>Secret</span>
        </label>
      </div>
    </div>

    <div class="cf-kv-footer">
      <CfButton variant="secondary" size="sm" @click="addRow">+ Add</CfButton>
      <span class="cf-kv-hint">{{ hint }}</span>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [{ key: '', value: '', secret: false }]
  },
  hint: { type: String, default: 'Or paste key-value pairs directly.' }
})

const emit = defineEmits(['update:modelValue'])

function updateField(idx, field, val) {
  const next = props.modelValue.map((r, i) => (i === idx ? { ...r, [field]: val } : r))
  emit('update:modelValue', next)
}

function addRow() {
  emit('update:modelValue', [...props.modelValue, { key: '', value: '', secret: false }])
}

function removeRow(idx) {
  const next = [...props.modelValue]
  next.splice(idx, 1)
  emit('update:modelValue', next)
}
</script>

<style scoped>
.cf-kv-repeater {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.cf-kv-card {
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  overflow: hidden;
}

.cf-kv-row {
  display: flex;
  align-items: center;
  height: 36px;
  border-bottom: 1px solid var(--color-kumo-hairline);
}

.cf-kv-row:last-child {
  border-bottom: none;
}

.cf-kv-tag {
  width: 64px;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 10px;
  background-color: var(--color-kumo-recessed);
  border-right: 1px solid var(--color-kumo-hairline);
  font-size: 12px;
  font-weight: 500;
  color: var(--text-color-kumo-subtle);
  flex-shrink: 0;
}

.cf-kv-input {
  flex: 1;
  min-width: 0;
  border: none;
  background: transparent;
  padding: 0 10px;
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: var(--text-color-kumo-strong);
  outline: none;
}

.cf-kv-secret {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 0 10px;
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
  cursor: pointer;
  user-select: none;
}

.cf-kv-rm {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-subtle);
  padding: 0 10px;
  cursor: pointer;
  font-size: 15px;
}

.cf-kv-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.cf-kv-hint {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}
</style>
