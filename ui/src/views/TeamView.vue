<template>
  <div>
    <div class="breadcrumbs">
      <span>Manage Account</span>
      <span class="divider">/</span>
      <span class="current">Members & Permissions</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">Account Members & Roles ({{ filteredMembers.length }})</h1>
        <p class="page-subtitle">Manage member access, role-based access control (RBAC), and mandatory two-factor authentication.</p>
      </div>
      <button class="btn-cf btn-cf-primary" @click="showInviteModal = true">+ Invite Members</button>
    </div>

    <!-- Filter Bar -->
    <div style="margin-bottom: 1.25rem;">
      <input 
        v-model="searchQuery" 
        type="text" 
        class="cf-input" 
        placeholder="Filter members by name, email, or role..." 
        style="max-width: 360px;"
      />
    </div>

    <!-- Team Members Table -->
    <div class="kumo-table-container">
      <table class="kumo-table">
        <thead>
          <tr>
            <th>Member Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Permission Groups</th>
            <th>2FA Status</th>
            <th>Last Active</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in filteredMembers" :key="m.id">
            <td>
              <div style="display: flex; align-items: center; gap: 0.65rem;">
                <div class="user-avatar-initials" :style="{ background: m.color }">{{ m.initials }}</div>
                <strong>{{ m.name }}</strong>
              </div>
            </td>
            <td>{{ m.email }}</td>
            <td><span class="badge-tag-beta">{{ m.role }}</span></td>
            <td>{{ m.groups }}</td>
            <td>
              <span class="cf-pill" :class="m.twoFactor ? 'active' : 'stopped'">
                <span class="status-dot"></span>
                <span>{{ m.twoFactor ? 'Enforced' : 'Disabled' }}</span>
              </span>
            </td>
            <td>{{ m.lastActive }}</td>
            <td>
              <div style="display: flex; gap: 0.35rem;">
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="manageMember(m)">Manage</button>
                <button 
                  v-if="m.role !== 'Super Administrator'" 
                  class="btn-cf btn-cf-secondary btn-cf-sm" 
                  @click="removeMember(m.id)"
                  style="color: #ef4444;"
                >
                  Remove
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Invite Member Modal -->
    <div v-if="showInviteModal" class="cf-modal-backdrop" @click.self="showInviteModal = false">
      <div class="cf-modal-card" style="max-width: 520px;">
        <div class="cf-modal-header">
          <div class="cf-modal-title">Invite New Team Member</div>
          <button class="cf-modal-close-btn" @click="showInviteModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>Member Email Address</label>
            <input v-model="inviteForm.email" type="email" class="cf-input" placeholder="colleague@company.com" />
          </div>

          <div class="cf-form-group">
            <label>Account Role</label>
            <select v-model="inviteForm.role" class="cf-select">
              <option value="Administrator">Administrator (Manage domains & workers)</option>
              <option value="Developer">Developer (Deploy code & read logs)</option>
              <option value="Billing Only">Billing Only (View & pay invoices)</option>
              <option value="Read Only">Read Only (Telemetry and auditing)</option>
            </select>
          </div>

          <div class="cf-form-group">
            <label>Initial Permission Groups</label>
            <div style="display: flex; flex-direction: column; gap: 0.5rem; margin-top: 0.25rem;">
              <label style="font-size: 0.82rem; display: flex; align-items: center; gap: 0.5rem;">
                <input type="checkbox" checked /> Workers Deployers
              </label>
              <label style="font-size: 0.82rem; display: flex; align-items: center; gap: 0.5rem;">
                <input type="checkbox" checked /> DNS Zone Managers
              </label>
              <label style="font-size: 0.82rem; display: flex; align-items: center; gap: 0.5rem;">
                <input type="checkbox" /> D1 Database Editors
              </label>
            </div>
          </div>
        </div>
        <div class="cf-modal-footer">
          <button class="btn-cf btn-cf-secondary" @click="showInviteModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="sendInvite">Send Invitation</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { mockData } from '@/data/mockData.js'

const members = ref([...mockData.teamMembers])
const searchQuery = ref('')
const showInviteModal = ref(false)

const inviteForm = ref({
  email: '',
  role: 'Developer'
})

const filteredMembers = computed(() => {
  if (!searchQuery.value) return members.value
  const q = searchQuery.value.toLowerCase()
  return members.value.filter(m => 
    m.name.toLowerCase().includes(q) ||
    m.email.toLowerCase().includes(q) ||
    m.role.toLowerCase().includes(q)
  )
})

function sendInvite() {
  if (!inviteForm.value.email) return
  const initials = inviteForm.value.email.substring(0, 2).toUpperCase()
  members.value.push({
    id: 'usr_' + Date.now(),
    name: inviteForm.value.email.split('@')[0],
    email: inviteForm.value.email,
    role: inviteForm.value.role,
    groups: 'Workers Deployers, DNS Managers',
    twoFactor: true,
    lastActive: 'Invited (Pending)',
    initials,
    color: '#059669'
  })
  inviteForm.value.email = ''
  showInviteModal.value = false
  alert('Invitation sent successfully!')
}

function removeMember(id) {
  if (confirm('Are you sure you want to revoke this user access?')) {
    members.value = members.value.filter(m => m.id !== id)
  }
}

function manageMember(m) {
  alert(`Opening permissions drawer for ${m.name} (${m.email})`)
}
</script>
