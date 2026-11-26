<template>
  <div>
    <div class="breadcrumbs">
      <span>Manage Account</span>
      <span class="divider">/</span>
      <span class="current">Profile & Preferences</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">User Profile & Account Preferences</h1>
        <p class="page-subtitle">Manage personal account credentials, sessions, security enforcement, and alert subscriptions.</p>
      </div>
    </div>

    <!-- Alert Banner -->
    <CfAlert v-if="savedNotice" variant="success" title="Profile Updated" dismissible style="margin-bottom: 1.5rem;" @dismiss="savedNotice = false">
      Your account preferences and profile details have been saved successfully.
    </CfAlert>

    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(360px, 1fr)); gap: 1.5rem;">
      <!-- Profile Information -->
      <CfCard title="User Profile Information" subtitle="Update your display name and regional localization preferences.">
        <form @submit.prevent="saveProfile" style="display: flex; flex-direction: column; gap: 1rem;">
          <div style="display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); gap: 1rem;">
            <CfInput v-model="profile.firstName" label="First Name" />
            <CfInput v-model="profile.lastName" label="Last Name" />
          </div>

          <CfInput
            :model-value="profile.email"
            label="Primary Email Address"
            type="email"
            disabled
            hint="Contact enterprise support to change primary login identity."
          />

          <CfInput v-model="profile.timezone" label="Timezone" />

          <CfSelect
            v-model="profile.language"
            label="Dashboard Language"
            :options="languageOptions"
          />

          <div>
            <CfButton type="submit" variant="primary" icon="check">Save Profile</CfButton>
          </div>
        </form>
      </CfCard>

      <!-- Security & Active Sessions -->
      <div style="display: flex; flex-direction: column; gap: 1.5rem;">
        <!-- Multi-Factor Authentication -->
        <CfCard
          title="Two-Factor Authentication (2FA)"
          subtitle="Hardware security keys (FIDO2/WebAuthn) and TOTP authenticator app enforcement."
        >
          <template #actions>
            <CfBadge variant="success" dot>Enforced</CfBadge>
          </template>
          <router-link to="/auth/2fa" style="text-decoration: none;">
            <CfButton variant="secondary" size="sm" icon="shield">Re-verify 2FA Session</CfButton>
          </router-link>
        </CfCard>

        <!-- Active Security Sessions -->
        <CfCard
          title="Active Security Sessions"
          subtitle="Review connected devices and revoke unauthorized browser sessions."
        >
          <div style="display: flex; flex-direction: column; gap: 0.75rem;">
            <div
              v-for="(sess, idx) in sessions"
              :key="idx"
              style="display: flex; justify-content: space-between; align-items: center; padding: 0.75rem; border: 1px solid var(--color-kumo-line); border-radius: var(--radius-md);"
            >
              <div>
                <div style="font-weight: 600; font-size: 0.88rem; color: var(--text-color-kumo-default);">{{ sess.device }}</div>
                <div style="font-size: 0.78rem; color: var(--text-color-kumo-subtle);">
                  {{ sess.browser }} • {{ sess.location }} (<code>{{ sess.ip }}</code>)
                </div>
              </div>
              <div>
                <CfBadge v-if="sess.current" variant="success" dot size="sm">Current Session</CfBadge>
                <CfButton v-else variant="secondary" size="sm" @click="revokeSession(idx)">Revoke</CfButton>
              </div>
            </div>
          </div>
        </CfCard>

        <!-- Notifications & Alerts -->
        <CfCard
          title="Email Notifications"
          subtitle="Manage system broadcasts and real-time threat alerts."
        >
          <div style="display: flex; flex-direction: column; gap: 0.75rem;">
            <CfCheckbox v-model="notifications.ddosAlerts" label="Real-time DDoS threat mitigation alerts" />
            <CfCheckbox v-model="notifications.billingAlerts" label="Monthly billing statements and quota thresholds" />
            <CfCheckbox v-model="notifications.weeklyDigests" label="Weekly edge analytics digest" />
          </div>
        </CfCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { mockData } from '@/data/mockData.js'

const profile = ref({ ...mockData.settings.profile })
const sessions = ref([...mockData.settings.security.activeSessions])
const notifications = ref({ ...mockData.settings.notifications })
const savedNotice = ref(false)

const languageOptions = [
  'English (United States)',
  'German (Deutsch)',
  'Japanese (日本語)',
  'French (Français)'
]

function saveProfile() {
  savedNotice.value = true
  setTimeout(() => {
    savedNotice.value = false
  }, 4000)
}

function revokeSession(index) {
  sessions.value.splice(index, 1)
  alert('Session revoked immediately.')
}
</script>
