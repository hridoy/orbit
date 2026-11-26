<template>
  <div>
    <div class="breadcrumbs">
      <span>Rules</span>
      <span class="divider">/</span>
      <span class="current">Transform & Page Rules</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Edge Transform Rules</h1>
        <p class="page-subtitle">Customize HTTP request and response headers, URL rewriting, and canonical redirect rules directly at the edge.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showCreateModal = true">+ Create Rule</button>
    </div>

    <!-- Rule Type Filters -->
    <div style="display: flex; gap: 0.5rem; margin-bottom: 1.25rem;">
      <button 
        v-for="cat in categories" 
        :key="cat"
        class="btn-cf btn-cf-sm"
        :class="activeCategory === cat ? 'btn-cf-primary' : 'btn-cf-secondary'"
        @click="activeCategory = cat"
      >
        {{ cat }}
      </button>
    </div>

    <!-- Rules Table -->
    <div class="kumo-table-container">
      <div class="table-header-title">
        <span>Active Rules ({{ filteredRules.length }})</span>
        <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Rule quota: 4 / 20 used</span>
      </div>
      <table class="kumo-table">
        <thead>
          <tr>
            <th style="width: 50px;">Order</th>
            <th>Rule Description</th>
            <th>Expression / Match Criteria</th>
            <th>Action Applied</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(rule, idx) in filteredRules" :key="rule.id">
            <td style="font-weight: 600; color: var(--text-color-kumo-subtle);">#{{ idx + 1 }}</td>
            <td>
              <strong>{{ rule.name }}</strong>
              <div style="font-size: 0.78rem; color: var(--text-color-kumo-subtle);">Type: {{ rule.type }}</div>
            </td>
            <td>
              <code style="font-size: 0.78rem; background: var(--color-kumo-canvas); padding: 0.2rem 0.4rem; border-radius: var(--radius-sm);">
                {{ rule.expression }}
              </code>
            </td>
            <td>
              <span class="badge-tag-beta">{{ rule.action }}</span>
            </td>
            <td>
              <span 
                class="cf-pill cursor-pointer" 
                :class="rule.active ? 'active' : 'stopped'"
                @click="rule.active = !rule.active"
              >
                <span class="status-dot"></span>
                <span>{{ rule.active ? 'Active' : 'Paused' }}</span>
              </span>
            </td>
            <td>
              <div style="display: flex; gap: 0.35rem;">
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="editRule(rule)">Edit</button>
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="deleteRule(rule.id)" style="color: #ef4444;">Delete</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Rule Modal -->
    <div v-if="showCreateModal" class="cf-modal-backdrop" @click.self="showCreateModal = false">
      <div class="cf-modal-card" style="max-width: 580px;">
        <div class="cf-modal-header">
          <div class="cf-modal-title">Create Transform Rule</div>
          <button class="cf-modal-close-btn" @click="showCreateModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>Rule Name</label>
            <input v-model="newRule.name" type="text" class="cf-input" placeholder="e.g. Redirect www to root" />
          </div>

          <div class="cf-form-group">
            <label>Rule Category</label>
            <select v-model="newRule.type" class="cf-select">
              <option value="Redirect Rule">Dynamic URL Redirect</option>
              <option value="Transform Rule">HTTP Response Header Modification</option>
              <option value="Page Rule">Cache Bypass & Edge Settings</option>
              <option value="Rewrite Rule">URL Path Rewriting</option>
            </select>
          </div>

          <div class="cf-form-group">
            <label>Match Expression (Wireshark-syntax)</label>
            <input v-model="newRule.expression" type="text" class="cf-input" placeholder="http.request.uri.path matches &quot;^/api&quot;" />
          </div>

          <div class="cf-form-group">
            <label>Target Action</label>
            <input v-model="newRule.action" type="text" class="cf-input" placeholder="e.g. Set header Strict-Transport-Security" />
          </div>
        </div>
        <div class="cf-modal-footer">
          <button class="btn-cf btn-cf-secondary" @click="showCreateModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="saveNewRule">Save & Deploy Rule</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const categories = ['All Rules', 'Redirect Rule', 'Transform Rule', 'Page Rule', 'Rewrite Rule']
const activeCategory = ref('All Rules')
const showCreateModal = ref(false)

const rules = ref([
  {
    id: 1,
    name: 'Enforce Canonical API Redirection',
    type: 'Redirect Rule',
    expression: 'http.request.uri.path matches "^/v1/api"',
    action: '301 Permanent Redirect to /api/v1',
    active: true
  },
  {
    id: 2,
    name: 'Add Security Response Headers',
    type: 'Transform Rule',
    expression: 'ssl eq true',
    action: 'Set HSTS & X-Content-Type-Options: nosniff',
    active: true
  },
  {
    id: 3,
    name: 'Strip Marketing UTM Query Params',
    type: 'Rewrite Rule',
    expression: 'http.request.uri.query contains "utm_"',
    action: 'Preserve URI path & strip query parameters',
    active: true
  },
  {
    id: 4,
    name: 'Bypass Edge Cache on Admin Routes',
    type: 'Page Rule',
    expression: 'http.request.uri.path matches "^/(admin|api/auth)"',
    action: 'Cache Level: Bypass, Disable Apps',
    active: true
  }
])

const filteredRules = computed(() => {
  if (activeCategory.value === 'All Rules') return rules.value
  return rules.value.filter(r => r.type === activeCategory.value)
})

const newRule = ref({
  name: '',
  type: 'Redirect Rule',
  expression: '',
  action: ''
})

function saveNewRule() {
  if (!newRule.value.name) return
  rules.value.push({
    id: Date.now(),
    name: newRule.value.name,
    type: newRule.value.type,
    expression: newRule.value.expression || 'true',
    action: newRule.value.action || 'Default Action',
    active: true
  })
  newRule.value = { name: '', type: 'Redirect Rule', expression: '', action: '' }
  showCreateModal.value = false
}

function deleteRule(id) {
  rules.value = rules.value.filter(r => r.id !== id)
}

function editRule(rule) {
  alert(`Editing rule: ${rule.name}`)
}
</script>
