<template>
  <CfModal :model-value="modelValue" :title="title" size="sm" @update:model-value="$emit('update:modelValue', $event)">
    <div class="cf-confirm-body">
      <p class="cf-confirm-msg">{{ message }}</p>
      <CfInput
        v-if="confirmText"
        v-model="typed"
        :label="`Type '${confirmText}' to confirm`"
        :placeholder="confirmText"
      />
    </div>
    <template #footer>
      <CfButton variant="outline" @click="$emit('update:modelValue', false)">Cancel</CfButton>
      <CfButton
        :variant="danger ? 'danger' : 'primary'"
        :disabled="!!confirmText && typed !== confirmText"
        @click="onConfirm"
      >
        {{ confirmLabel }}
      </CfButton>
    </template>
  </CfModal>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: 'Confirm Action' },
  message: { type: String, default: 'Are you sure you want to proceed?' },
  confirmLabel: { type: String, default: 'Confirm' },
  confirmText: { type: String, default: '' },
  danger: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue', 'confirm'])
const typed = ref('')

watch(() => props.modelValue, (val) => {
  if (val) typed.value = ''
})

function onConfirm() {
  emit('confirm')
  emit('update:modelValue', false)
}
</script>

<style scoped>
.cf-confirm-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cf-confirm-msg {
  margin: 0;
  font-size: 13.5px;
  line-height: 1.5;
  color: var(--text-color-kumo-default);
}
</style>
