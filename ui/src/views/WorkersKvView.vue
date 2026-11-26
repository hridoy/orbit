<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Storage & Databases</router-link>
      <span class="divider">/</span>
      <span class="current">Workers KV Namespaces</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Workers KV Key-Value Store</h1>
        <p class="page-subtitle">Low-latency distributed global key-value store with read replication to 330+ Cloudflare edge locations.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="alert('Create KV Namespace wizard opened.')">+ Create Namespace</button>
    </div>

    <!-- Namespaces Table -->
    <div class="kumo-table-container" style="margin-bottom: 2rem;">
      <div class="table-header-title">
        <span>Active KV Namespaces ({{ kv.namespaces.length }})</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Namespace Name</th>
            <th>Namespace ID</th>
            <th>Total Keys</th>
            <th>Storage Size</th>
            <th>24h Reads / Writes</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ns in kv.namespaces" :key="ns.id">
            <td><strong style="color: var(--text-color-kumo-strong);">🔑 {{ ns.name }}</strong></td>
            <td><code>{{ ns.id }}</code></td>
            <td>{{ ns.keys_count.toLocaleString() }} keys</td>
            <td>{{ ns.size }}</td>
            <td>{{ ns.read_ops_24h }} / {{ ns.write_ops_24h }}</td>
            <td>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="alert(`Viewing keys in ${ns.name}`)">Browse Keys</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Sample Key-Value Browser -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Keys in EDGE_CONFIG_STORE</span>
        <div style="display: flex; gap: 0.5rem;">
          <input v-model="keySearch" type="text" class="cf-input" placeholder="Search keys by prefix..." style="width: 220px; font-size: 0.8rem; padding: 0.3rem 0.6rem;" />
          <button class="btn-cf btn-cf-primary btn-cf-sm" @click="showAddKey = true">+ Add Key</button>
        </div>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Key Name</th>
            <th>Value (JSON / String)</th>
            <th>Expiration TTL</th>
            <th>Metadata</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="k in filteredKeys" :key="k.key">
            <td><strong>{{ k.key }}</strong></td>
            <td><code>{{ k.value }}</code></td>
            <td>{{ k.ttl }}</td>
            <td><span class="badge-tag-beta">{{ k.metadata }}</span></td>
            <td>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="copyValue(k.value, k.key)">
                {{ copiedKey === k.key ? '✓ Copied' : 'Copy' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { mockData } from '@/data/mockData.js'

const kv = mockData.workersKv
const keySearch = ref('')
const copiedKey = ref(null)
const showAddKey = ref(false)

const filteredKeys = computed(() => {
  const q = keySearch.value.toLowerCase().trim()
  if (!q) return kv.sampleKeys
  return kv.sampleKeys.filter(k => k.key.toLowerCase().includes(q) || k.value.toLowerCase().includes(q))
})

function copyValue(val, key) {
  if (navigator.clipboard) navigator.clipboard.writeText(val)
  copiedKey.value = key
  setTimeout(() => { copiedKey.value = null }, 2000)
}
</script>
