<template>
  <div
    class="cf-dropzone"
    :class="{ 'is-dragging': dragging }"
    @dragover.prevent="dragging = true"
    @dragleave.prevent="dragging = false"
    @drop.prevent="onDrop"
    @click="fileInput?.click()"
  >
    <input
      ref="fileInput"
      type="file"
      :multiple="multiple"
      :accept="accept"
      style="display: none;"
      @change="onFileSelect"
    />
    <CfIcon name="cis-cloud-upload" size="28" class="cf-dropzone-icon" />
    <div class="cf-dropzone-title">{{ title }}</div>
    <div class="cf-dropzone-hint">{{ hint }}</div>

    <div v-if="files.length" class="cf-dropzone-files" @click.stop>
      <div v-for="(f, idx) in files" :key="f.name + idx" class="cf-dropzone-file">
        <span class="cf-file-name">{{ f.name }}</span>
        <span class="cf-file-size">{{ formatSize(f.size) }}</span>
        <button type="button" class="cf-file-rm" @click="removeFile(idx)">×</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: { type: String, default: 'Drag and drop files here, or click to browse' },
  hint: { type: String, default: 'Supports JSON, TOML, JS/TS, or binary assets up to 50 MB' },
  multiple: { type: Boolean, default: true },
  accept: { type: String, default: '' }
})

const emit = defineEmits(['files-change'])
const dragging = ref(false)
const fileInput = ref(null)
const files = ref([])

function handleFiles(list) {
  const incoming = Array.from(list || [])
  files.value = props.multiple ? [...files.value, ...incoming] : incoming.slice(0, 1)
  emit('files-change', files.value)
}

function onDrop(e) {
  dragging.value = false
  handleFiles(e.dataTransfer?.files)
}

function onFileSelect(e) {
  handleFiles(e.target?.files)
}

function removeFile(idx) {
  files.value.splice(idx, 1)
  emit('files-change', files.value)
}

function formatSize(bytes) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<style scoped>
.cf-dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 28px 20px;
  border-radius: var(--radius-lg, 8px);
  border: 1.5px dashed var(--color-kumo-line);
  background-color: var(--color-kumo-canvas);
  text-align: center;
  cursor: pointer;
  transition: border-color 0.15s ease, background-color 0.15s ease;
}

.cf-dropzone:hover,
.cf-dropzone.is-dragging {
  border-color: var(--color-kumo-brand);
  background-color: var(--color-kumo-tint);
}

.cf-dropzone-icon {
  color: var(--text-color-kumo-subtle);
}

.cf-dropzone-title {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--text-color-kumo-strong);
}

.cf-dropzone-hint {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-dropzone-files {
  width: 100%;
  max-width: 420px;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.cf-dropzone-file {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 6px;
  background-color: var(--color-kumo-base);
  border: 1px solid var(--color-kumo-line);
  font-size: 12px;
}

.cf-file-name {
  flex: 1;
  text-align: left;
  color: var(--text-color-kumo-strong);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cf-file-size {
  font-family: var(--font-mono);
  color: var(--text-color-kumo-subtle);
}

.cf-file-rm {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-subtle);
  cursor: pointer;
  font-size: 14px;
}
</style>
