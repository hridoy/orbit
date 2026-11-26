<template>
  <div class="cf-tokens-page">
    <!-- Breadcrumbs -->
    <nav class="cf-tokens-breadcrumb-bar" aria-label="breadcrumb">
      <router-link to="/r2" class="cf-tokens-breadcrumb-link">
        <CfIcon name="r2" size="18" />
        <span>R2 Object Storage</span>
      </router-link>
      <svg class="cf-tokens-breadcrumb-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
      </svg>
      <span v-if="!showCreateForm" class="cf-tokens-breadcrumb-current">API Tokens</span>
      <template v-else>
        <a href="javascript:void(0)" class="cf-tokens-breadcrumb-link" @click="showCreateForm = false">API Tokens</a>
        <svg class="cf-tokens-breadcrumb-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
        </svg>
        <span class="cf-tokens-breadcrumb-current">{{ createType === 'account' ? 'Create Account API Token' : 'Create User API Token' }}</span>
      </template>
    </nav>

    <!-- Main Heading -->
    <h1 class="cf-tokens-main-title">API Tokens</h1>

    <!-- MAIN VIEW: ACCOUNT & USER TOKENS LIST -->
    <template v-if="!showCreateForm">
      <!-- Section 1: Account API Tokens -->
      <section class="cf-tokens-section">
        <h2 class="cf-tokens-section-title">Account API Tokens</h2>
        <div class="cf-tokens-desc-row">
          <p class="cf-tokens-desc-text">
            Create API tokens for service authentication and bucket access that are tied to the entire account. These remain active even when you leave the organization, making them ideal for production systems. (recommended)
          </p>
          <a href="javascript:void(0)" class="cf-tokens-btn-kumo" @click="openCreate('account')">
            <span aria-hidden="true" class="cf-btn-gloss"></span>
            <span class="cf-btn-label">Create Account API token</span>
          </a>
        </div>

        <!-- Account Tokens Table or Empty State -->
        <div v-if="accountTokens.length === 0" class="cf-tokens-empty-card" data-testid="no-tokens">
          <p style="margin: 0;">There are currently no Account API tokens for your buckets</p>
        </div>
        <div v-else class="cf-tokens-table-card">
          <table class="cf-tokens-table">
            <thead>
              <tr>
                <th>
                  <div style="display: flex; align-items: center; gap: 4px;">
                    <span>Token name</span>
                    <svg width="12" height="12" viewBox="0 0 256 256" fill="currentColor" style="opacity: 0.5;">
                      <path d="M181.66,133.66l-48,48a8,8,0,0,1-11.32,0l-48-48a8,8,0,0,1,11.32-11.32L120,156.69V48a8,8,0,0,1,16,0V156.69l34.34-34.35a8,8,0,0,1,11.32,11.32Zm-96-35.32a8,8,0,0,0,11.32,0L131.31,64V172a8,8,0,0,0,16,0V64l34.35,34.34a8,8,0,0,0,11.32-11.32l-48-48a8,8,0,0,0-11.32,0l-48,48A8,8,0,0,0,85.66,98.34Z"></path>
                    </svg>
                  </div>
                </th>
                <th>Applied to</th>
                <th>Permission</th>
                <th>Issued on</th>
                <th>Status</th>
                <th style="width: 48px;"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="token in accountTokens" :key="token.id">
                <td>{{ token.name }}</td>
                <td>{{ token.appliedTo }}</td>
                <td>{{ token.permission }}</td>
                <td>{{ token.issuedOn }}</td>
                <td><span class="cf-tokens-badge-active">Active</span></td>
                <td style="position: relative;">
                  <button class="cf-tokens-action-btn" aria-label="Actions" @click="toggleMenu(token.id)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 256 256">
                      <path d="M144,128a16,16,0,1,1-16-16A16,16,0,0,1,144,128ZM60,112a16,16,0,1,0,16,16A16,16,0,0,0,60,112Zm136,0a16,16,0,1,0,16,16A16,16,0,0,0,196,112Z"/>
                    </svg>
                  </button>
                  <div v-if="activeMenuId === token.id" class="cf-tokens-menu-popover">
                    <button class="cf-tokens-menu-item" @click="editToken(token)">Edit</button>
                    <button class="cf-tokens-menu-item" @click="rollToken(token)">Roll</button>
                    <button class="cf-tokens-menu-item danger" @click="deleteToken(token, 'account')">Delete</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- Section 2: User API Tokens -->
      <section class="cf-tokens-section" style="margin-top: 1.5rem;">
        <h2 class="cf-tokens-section-title">User API Tokens</h2>
        <div class="cf-tokens-desc-row">
          <p class="cf-tokens-desc-text">
            Create API tokens for service authentication and bucket access specific to your user account. These become inactive if you leave the organization, making them ideal for personal access or development work.
          </p>
          <a href="javascript:void(0)" class="cf-tokens-btn-kumo" @click="openCreate('user')">
            <span aria-hidden="true" class="cf-btn-gloss"></span>
            <span class="cf-btn-label">Create User API token</span>
          </a>
        </div>

        <!-- User Tokens Table -->
        <div class="cf-tokens-table-card">
          <table class="cf-tokens-table">
            <thead>
              <tr>
                <th>
                  <div style="display: flex; align-items: center; gap: 4px;">
                    <span>Token name</span>
                    <svg width="12" height="12" viewBox="0 0 256 256" fill="currentColor" style="opacity: 0.5;">
                      <path d="M181.66,133.66l-48,48a8,8,0,0,1-11.32,0l-48-48a8,8,0,0,1,11.32-11.32L120,156.69V48a8,8,0,0,1,16,0V156.69l34.34-34.35a8,8,0,0,1,11.32,11.32Zm-96-35.32a8,8,0,0,0,11.32,0L131.31,64V172a8,8,0,0,0,16,0V64l34.35,34.34a8,8,0,0,0,11.32-11.32l-48-48a8,8,0,0,0-11.32,0l-48,48A8,8,0,0,0,85.66,98.34Z"></path>
                    </svg>
                  </div>
                </th>
                <th>Applied to</th>
                <th>Permission</th>
                <th>Issued on</th>
                <th>Status</th>
                <th style="width: 48px;"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="token in userTokens" :key="token.id">
                <td>{{ token.name }}</td>
                <td>{{ token.appliedTo }}</td>
                <td>{{ token.permission }}</td>
                <td>{{ token.issuedOn }}</td>
                <td><span class="cf-tokens-badge-active">Active</span></td>
                <td style="position: relative;">
                  <button class="cf-tokens-action-btn" aria-label="Actions" @click="toggleMenu(token.id)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 256 256">
                      <path d="M144,128a16,16,0,1,1-16-16A16,16,0,0,1,144,128ZM60,112a16,16,0,1,0,16,16A16,16,0,0,0,60,112Zm136,0a16,16,0,1,0,16,16A16,16,0,0,0,196,112Z"/>
                    </svg>
                  </button>
                  <div v-if="activeMenuId === token.id" class="cf-tokens-menu-popover">
                    <button class="cf-tokens-menu-item" @click="editToken(token)">Edit</button>
                    <button class="cf-tokens-menu-item" @click="rollToken(token)">Roll</button>
                    <button class="cf-tokens-menu-item danger" @click="deleteToken(token, 'user')">Delete</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </template>

    <!-- CREATE TOKEN FORM VIEW -->
    <div v-else class="cf-token-create-card">
      <h2 class="cf-tokens-section-title" style="margin-bottom: 0.5rem;">
        {{ createType === 'account' ? 'Create Account API Token' : 'Create User API Token' }}
      </h2>

      <!-- Token name -->
      <div class="cf-form-section">
        <label class="cf-form-title" for="modal-token-name">Token name</label>
        <input
          id="modal-token-name"
          v-model="tokenForm.name"
          type="text"
          class="cf-input-control"
          placeholder="e.g. R2 CI/CD Deployment Token"
          required
        />
      </div>

      <!-- Permissions -->
      <div class="cf-form-section">
        <label class="cf-form-title">Permissions</label>
        <p class="cf-form-subtext">Specify the R2 Storage permission type for this token.</p>
        <div class="cf-perm-grid">
          <div
            v-for="p in permissionOptions"
            :key="p.id"
            class="cf-perm-option"
            :class="{ active: tokenForm.permission === p.id }"
            @click="tokenForm.permission = p.id"
          >
            <input
              type="radio"
              :name="'perm-radio'"
              :value="p.id"
              :checked="tokenForm.permission === p.id"
              style="margin-top: 3px;"
            />
            <div class="cf-perm-text-group">
              <span class="cf-perm-name">{{ p.name }}</span>
              <span class="cf-perm-desc">{{ p.desc }}</span>
            </div>
          </div>
        </div>
        <p class="cf-form-subtext" style="margin-top: 0.5rem;">
          Permissions can be further customized or created for multiple accounts from your Profile.
        </p>
      </div>

      <!-- Specify bucket(s) -->
      <div class="cf-form-section">
        <label class="cf-form-title">Specify bucket(s)</label>
        <div style="display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.35rem;">
          <label style="display: flex; align-items: center; gap: 0.5rem; font-size: 14px; cursor: pointer;">
            <input v-model="tokenForm.bucketScope" type="radio" value="all" />
            <span>Apply to all buckets in this account (including newly created buckets)</span>
          </label>
          <label style="display: flex; align-items: center; gap: 0.5rem; font-size: 14px; cursor: pointer;">
            <input v-model="tokenForm.bucketScope" type="radio" value="specific" />
            <span>Apply to specific buckets only</span>
          </label>
        </div>

        <div v-if="tokenForm.bucketScope === 'specific'" style="margin-top: 0.75rem;">
          <select v-model="tokenForm.selectedBucket" class="cf-input-control" style="max-width: 380px;">
            <option value="">Select a bucket...</option>
            <option value="code-collection">code-collection</option>
            <option value="media-assets-cdn">media-assets-cdn</option>
            <option value="backups-archive">backups-archive</option>
          </select>
          <span v-if="!tokenForm.selectedBucket" style="display: block; font-size: 12px; color: #dc2626; margin-top: 0.35rem;">
            At least one bucket is required
          </span>
        </div>
      </div>

      <!-- TTL -->
      <div class="cf-form-section">
        <label class="cf-form-title">TTL (Time to Live)</label>
        <p class="cf-form-subtext">Define how long this token will stay active:</p>
        <select v-model="tokenForm.ttl" class="cf-input-control" style="max-width: 240px; margin-top: 0.35rem;">
          <option value="forever">Forever</option>
          <option value="24h">24 hours</option>
          <option value="7d">1 week</option>
          <option value="30d">30 days</option>
          <option value="6m">6 months</option>
          <option value="1y">1 year</option>
        </select>
      </div>

      <!-- Client IP Address Filtering -->
      <div class="cf-form-section">
        <label class="cf-form-title">Client IP Address Filtering</label>
        <p class="cf-form-subtext">
          Select IP addresses or ranges of IP addresses to filter. This filter limits the client IP addresses that can use the API token with Cloudflare. By default, this token will apply to all addresses.
        </p>
        <div style="display: flex; gap: 1rem; margin-top: 0.5rem; flex-wrap: wrap;">
          <div style="flex: 1; min-width: 200px;">
            <label style="font-size: 12px; color: #737373; display: block; margin-bottom: 0.25rem;">Include:</label>
            <input v-model="tokenForm.ipInclude" type="text" class="cf-input-control" placeholder="e.g. 192.168.1.88" />
          </div>
          <div style="flex: 1; min-width: 200px;">
            <label style="font-size: 12px; color: #737373; display: block; margin-bottom: 0.25rem;">Exclude:</label>
            <input v-model="tokenForm.ipExclude" type="text" class="cf-input-control" placeholder="e.g. 192.168.1.88" />
          </div>
        </div>
      </div>

      <!-- Actions Bar -->
      <div class="cf-form-actions-bar">
        <button type="button" class="cf-btn-outline" @click="showCreateForm = false">
          Cancel
        </button>
        <button
          type="button"
          class="cf-tokens-btn-kumo"
          :disabled="isSubmitDisabled"
          @click="submitToken"
        >
          <span aria-hidden="true" class="cf-btn-gloss"></span>
          <span class="cf-btn-label">{{ createType === 'account' ? 'Create Account API Token' : 'Create User API Token' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const showCreateForm = ref(false)
const createType = ref('account')
const activeMenuId = ref(null)

// Authentic initial data matching real Cloudflare
const accountTokens = ref([])

const userTokens = ref([
  {
    id: 'tok-1',
    name: 'code-collection build token',
    appliedTo: 'All buckets',
    permission: 'Admin Read & Write',
    issuedOn: 'Nov 3, 2026',
    status: 'Active'
  },
  {
    id: 'tok-2',
    name: 'github-actions-sync',
    appliedTo: 'All buckets',
    permission: 'Object Read & Write',
    issuedOn: 'Oct 15, 2026',
    status: 'Active'
  }
])

const permissionOptions = [
  {
    id: 'Admin Read & Write',
    name: 'Admin Read & Write',
    desc: 'Allows the ability to create, list, and delete buckets, edit bucket configuration, read, write, and list objects, and read and write access to data catalog tables and associated metadata.'
  },
  {
    id: 'Admin Read only',
    name: 'Admin Read only',
    desc: 'Allows the ability to list buckets and view bucket configuration, read and list objects, and read access to data catalog tables and associated metadata.'
  },
  {
    id: 'Object Read & Write',
    name: 'Object Read & Write',
    desc: 'Allows the ability to read, write, and list objects in specific buckets.'
  },
  {
    id: 'Object Read only',
    name: 'Object Read only',
    desc: 'Allows the ability to read and list objects in specific buckets.'
  }
]

const tokenForm = ref({
  name: '',
  permission: 'Object Read only',
  bucketScope: 'all',
  selectedBucket: '',
  ttl: 'forever',
  ipInclude: '',
  ipExclude: ''
})

const isSubmitDisabled = computed(() => {
  if (!tokenForm.value.name.trim()) return true
  if (tokenForm.value.bucketScope === 'specific' && !tokenForm.value.selectedBucket) return true
  return false
})

const openCreate = (type) => {
  createType.value = type
  tokenForm.value.name = type === 'account' ? 'R2 Account Token' : 'R2 User Token'
  tokenForm.value.permission = 'Object Read only'
  tokenForm.value.bucketScope = 'all'
  tokenForm.value.selectedBucket = ''
  tokenForm.value.ttl = 'forever'
  tokenForm.value.ipInclude = ''
  tokenForm.value.ipExclude = ''
  showCreateForm.value = true
}

const toggleMenu = (id) => {
  activeMenuId.value = activeMenuId.value === id ? null : id
}

const editToken = (token) => {
  alert(`Editing token: ${token.name}`)
  activeMenuId.value = null
}

const rollToken = (token) => {
  alert(`Rolling access keys for ${token.name}. Previous secret key is revoked.`)
  activeMenuId.value = null
}

const deleteToken = (token, type) => {
  if (confirm(`Permanently delete token "${token.name}"?`)) {
    if (type === 'account') {
      accountTokens.value = accountTokens.value.filter(t => t.id !== token.id)
    } else {
      userTokens.value = userTokens.value.filter(t => t.id !== token.id)
    }
  }
  activeMenuId.value = null
}

const submitToken = () => {
  const applied = tokenForm.value.bucketScope === 'all' ? 'All buckets' : tokenForm.value.selectedBucket
  const newToken = {
    id: `tok-${Date.now()}`,
    name: tokenForm.value.name,
    appliedTo: applied,
    permission: tokenForm.value.permission,
    issuedOn: 'Just now',
    status: 'Active'
  }

  if (createType.value === 'account') {
    accountTokens.value.unshift(newToken)
  } else {
    userTokens.value.unshift(newToken)
  }

  showCreateForm.value = false
  alert(`Generated token "${newToken.name}" successfully! Access Key ID & Secret Access Key generated.`)
}
</script>
