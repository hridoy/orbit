<template>
  <div class="cf-tag-input-wrap">
    <label v-if="label" class="cf-tag-label">{{ label }}</label>
    <div class="cf-tag-box" @click="inputRef?.focus()">
      <span v-for="(tag, idx) in modelValue" :key="tag + idx" class="cf-tag-chip">
        <span>{{ tag }}</span>
        <button
          type="button"
          class="cf-tag-remove"
          :aria-label="`Remove ${tag}`"
          @click.stop="removeTag(idx)"
        >
          ×
        </button>
      </span>
      <input
        ref="inputRef"
        v-model="draft"
        type="text"
        :placeholder="modelValue.length ? '' : placeholder"
        class="cf-tag-native"
        @keydown.enter.prevent="addTag"
        @keydown.backspace="onBackspace"
      />
    </div>
    <span v-if="hint" class="cf-tag-hint">{{ hint }}</span>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Add tag and press Enter...'
  },
  hint: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])
const draft = ref('')
const inputRef = ref(null)

function addTag() {
  const val = draft.value.trim()
  if (!val || props.modelValue.includes(val)) {
    draft.value = ''
    return
  }
  emit('update:modelValue', [...props.modelValue, val])
  draft.value = ''
}

function removeTag(idx) {
  const next = [...props.modelValue]
  next.splice(idx, 1)
  emit('update:modelValue', next)
}

function onBackspace() {
  if (!draft.value && props.modelValue.length > 0) {
    removeTag(props.modelValue.length - 1)
  }
}
</script>

<style scoped>
.cf-tag-input-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.cf-tag-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color-kumo-default);
}

.cf-tag-box {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  min-height: 36px;
  padding: 4px 8px;
  border-radius: var(--radius-md, 6px);
  border: 1px solid var(--kumo-control-border-default);
  background-color: var(--kumo-control-bg-default);
  cursor: text;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.cf-tag-box:focus-within {
  border-color: var(--color-kumo-brand);
  box-shadow: 0 0 0 3px rgba(0, 81, 195, 0.15);
}

.cf-tag-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 24px;
  padding: 0 8px;
  border-radius: 9999px;
  background-color: var(--color-kumo-fill);
  color: var(--text-color-kumo-strong);
  font-size: 12px;
  font-weight: 500;
}

.cf-tag-remove {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-subtle);
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
  padding: 0 2px;
}

.cf-tag-remove:hover {
  color: var(--text-color-kumo-danger);
}

.cf-tag-native {
  flex: 1;
  min-width: 110px;
  border: none;
  background: transparent;
  outline: none;
  font-family: var(--font-sans);
  font-size: 13px;
  color: var(--text-color-kumo-default);
  height: 26px;
}

.cf-tag-hint {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}
</style>
