<template>
  <div class="kumo-code-box">
    <div class="kumo-code-header">
      <div class="kumo-code-lang">
        <span v-if="filename" class="kumo-code-filename">{{ filename }}</span>
        <span v-else class="kumo-code-lang-badge">{{ language }}</span>
      </div>
      <button
        type="button"
        class="kumo-code-copy-btn"
        :title="copied ? 'Copied!' : 'Copy snippet'"
        @click="copySnippet"
      >
        <CfIcon v-if="!copied" name="cis-clipboard" size="13" />
        <CfIcon v-else name="cis-checkmark-1-medium" size="13" style="color: #10b981;" />
        <span>{{ copied ? 'Copied' : 'Copy' }}</span>
      </button>
    </div>

    <pre class="kumo-code-pre"><code><slot>{{ code }}</slot></code></pre>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  code: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'bash'
  },
  filename: {
    type: String,
    default: null
  }
})

const copied = ref(false)

const copySnippet = async () => {
  try {
    await navigator.clipboard.writeText(props.code)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch (err) {
    console.warn(err)
  }
}
</script>

<style scoped>
.kumo-code-box {
  display: flex;
  flex-direction: column;
  width: 100%;
  border-radius: 8px;
  border: 1px solid var(--color-kumo-line, #e2e8f0);
  background-color: #fafafa;
  overflow: hidden;
  font-family: var(--font-mono);
  font-size: 12.5px;
  box-sizing: border-box;
}

[data-theme="dark"] .kumo-code-box {
  background-color: #0d1117;
  border-color: #1e293b;
}

.kumo-code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.85rem;
  background-color: #f4f4f5;
  border-bottom: 1px solid var(--color-kumo-line, #e2e8f0);
}

[data-theme="dark"] .kumo-code-header {
  background-color: #161b22;
  border-bottom-color: #21262d;
}

.kumo-code-lang-badge {
  font-size: 10.5px;
  font-weight: 600;
  text-transform: uppercase;
  color: #71717a;
  letter-spacing: 0.05em;
}

.kumo-code-filename {
  font-size: 12px;
  font-weight: 500;
  color: #18181b;
}

[data-theme="dark"] .kumo-code-filename {
  color: #f8fafc;
}

.kumo-code-copy-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 11.5px;
  font-weight: 500;
  color: #71717a;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.15s ease;
}

.kumo-code-copy-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: #18181b;
}

[data-theme="dark"] .kumo-code-copy-btn:hover {
  background-color: rgba(255, 255, 255, 0.08);
  color: #f8fafc;
}

.kumo-code-pre {
  margin: 0;
  padding: 1rem;
  overflow-x: auto;
  line-height: 1.6;
  color: #24292f;
}

[data-theme="dark"] .kumo-code-pre {
  color: #c9d1d9;
}
</style>
