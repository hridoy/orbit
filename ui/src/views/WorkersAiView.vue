<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Build</router-link>
      <span class="divider">/</span>
      <span>AI</span>
      <span class="divider">/</span>
      <span class="current">Workers AI & Gateway</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.75rem;">
          <h1 class="page-title">🤖 Workers AI & AI Gateway</h1>
          <span class="cf-pill active"><span class="status-dot"></span><span>GPU Edge Active</span></span>
          <span class="badge-tag-beta">Tensor Core</span>
        </div>
        <p class="page-subtitle">Serverless GPU inference running Meta Llama 3.3, Mistral, Whisper, and SDXL across Cloudflare's Anycast network.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="runInference">▶ Run Prompt Inference</button>
    </div>

    <!-- AI Gateway Telemetry Cards -->
    <div class="stats-overview-grid">
      <MetricCard
        label="24h AI Inferences"
        :value="ai.gateway.totalRequests"
        delta="↗ 24.8%"
        sparkline="M0 22 L 20 18 L 45 20 L 70 8 L 100 4"
        strokeColor="var(--cf-brand-orange)"
        :footer="`Tokens: ${ai.gateway.tokensGenerated}`"
      />
      <MetricCard
        label="Gateway Cache Hit"
        :value="ai.gateway.cacheHitRate"
        delta="↗ 12.1%"
        sparkline="M0 16 Q 30 22, 60 10 T 100 2"
        strokeColor="#059669"
        :footer="`Saved: ${ai.gateway.costSaved}`"
      />
      <MetricCard
        label="Median Latency"
        :value="ai.gateway.avgLatency"
        delta="↘ 14ms"
        sparkline="M0 8 L 25 15 L 50 12 L 75 20 L 100 16"
        strokeColor="#0051c3"
        footer="Streaming first-token"
      />
      <MetricCard
        label="Active GPU Cluster"
        value="NVIDIA H100 / L40S"
        sparkline="M0 14 L 30 18 L 60 10 L 100 8"
        strokeColor="#7c3aed"
        footer="Serverless edge execution"
      />
    </div>

    <!-- Interactive AI Prompt Playground -->
    <div class="ui-kit-section" style="margin-top: 1.5rem;">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; flex-wrap: wrap; gap: 0.75rem;">
        <div>
          <h3 class="ui-kit-title" style="margin: 0;">Interactive AI Prompt Playground</h3>
          <p class="ui-kit-desc" style="margin: 0.25rem 0 0 0;">Test real-time model inference and response streaming directly from Cloudflare V8 isolates.</p>
        </div>
        <div style="display: flex; gap: 0.5rem; align-items: center;">
          <select v-model="selectedModel" class="cf-select" style="width: 240px; font-size: 0.82rem;">
            <option v-for="m in ai.models" :key="m.id" :value="m.id">{{ m.name }}</option>
          </select>
          <button class="btn-cf btn-cf-primary btn-cf-sm" @click="runInference">▶ Run (Ctrl+Enter)</button>
        </div>
      </div>

      <div class="cf-form-group">
        <label>Prompt Input</label>
        <textarea
          v-model="promptInput"
          class="cf-input"
          style="height: 70px; resize: vertical; font-family: var(--font-mono); font-size: 0.84rem;"
          @keydown.ctrl.enter="runInference"
        ></textarea>
      </div>

      <div>
        <label style="font-size: 0.82rem; font-weight: 600; color: var(--text-color-kumo-strong); margin-bottom: 0.35rem; display: block;">Streaming Inference Output</label>
        <div style="background: #0f172a; border: 1px solid #1e293b; border-radius: var(--radius-md); padding: 1.15rem; color: #38bdf8; font-family: var(--font-mono); font-size: 0.84rem; line-height: 1.6; min-height: 110px;">
          <span v-if="isStreaming" style="color: #94a3b8;">Streaming from Cloudflare Anycast GPU (IAD Edge)... </span>
          <span>{{ streamedResponse }}</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-top: 0.5rem; font-size: 0.76rem; color: var(--text-color-kumo-subtle);">
          <span>Generation speed: <strong>142 tokens/sec</strong></span>
          <span>Edge latency: <strong>114 ms</strong> • Tokens used: <strong>89 tokens</strong></span>
        </div>
      </div>
    </div>

    <!-- Supported Models Catalog Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Supported Edge Models Catalog ({{ ai.models.length }})</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Available via <code>@cloudflare/ai</code> SDK</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Model Name</th>
            <th>Author</th>
            <th>Task / Category</th>
            <th>Inference Speed</th>
            <th>Context Window</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in ai.models" :key="m.id">
            <td><strong style="color: var(--text-color-kumo-strong);">🧠 {{ m.name }}</strong></td>
            <td>{{ m.author }}</td>
            <td><span class="badge-tag-beta">{{ m.category }}</span></td>
            <td>{{ m.speed }}</td>
            <td><code>{{ m.context }}</code></td>
            <td><span class="cf-pill active"><span class="status-dot"></span><span>{{ m.status }}</span></span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'
import MetricCard from '@/components/MetricCard.vue'

const ai = mockData.workersAi
const selectedModel = ref('llama-3.3-70b')
const promptInput = ref(ai.samplePrompt)
const streamedResponse = ref(ai.sampleResponse)
const isStreaming = ref(false)

function runInference() {
  isStreaming.value = true
  streamedResponse.value = ''
  const fullText = ai.sampleResponse
  let idx = 0
  const interval = setInterval(() => {
    if (idx < fullText.length) {
      streamedResponse.value += fullText.charAt(idx)
      idx += 3
    } else {
      clearInterval(interval)
      streamedResponse.value = fullText
      isStreaming.value = false
    }
  }, 15)
}
</script>
