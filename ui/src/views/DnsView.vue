<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Account</router-link>
      <span class="divider">/</span>
      <router-link to="/zones/northlab.dev">{{ zone.name }}</router-link>
      <span class="divider">/</span>
      <span class="current">DNS Records</span>
    </div>

    <div class="page-title-row">
      <div>
        <div style="display: flex; align-items: center; gap: 0.65rem; flex-wrap: wrap;">
          <h1 class="page-title">DNS records for {{ zone.name }}</h1>
          <a href="https://developers.cloudflare.com/dns/" target="_blank" rel="noopener noreferrer" class="cf-view-docs-pill">
            <span>View docs</span>
          </a>
        </div>
        <p class="page-subtitle">Configure authoritative DNS records with global Anycast routing and automatic DDoS protection.</p>
      </div>
      <div style="display: flex; align-items: center; gap: 0.5rem;">
        <CfBadge variant="neutral" size="md">DNS Setup: Full</CfBadge>
        <CfButton variant="secondary" icon="info" @click="exportBindZone">DNS documentation</CfButton>
      </div>
    </div>

    <!-- Domain Sub-Navigation Tabs Bar -->
    <div class="cf-tab-nav">
      <router-link to="/zones/northlab.dev" class="cf-tab-btn">Overview</router-link>
      <router-link to="/dns" class="cf-tab-btn active">DNS Records</router-link>
      <router-link to="/ssl" class="cf-tab-btn">SSL/TLS</router-link>
      <router-link to="/waf" class="cf-tab-btn">Security (WAF)</router-link>
      <router-link to="/caching" class="cf-tab-btn">Caching</router-link>
      <router-link to="/rules" class="cf-tab-btn">Rules</router-link>
      <router-link to="/network" class="cf-tab-btn">Network</router-link>
    </div>

    <!-- Collapsible Recommendations Banner -->
    <CfAccordion
      :items="[{ value: 'rec', title: 'Recommendations', badge: '1', content: 'DNSSEC is not enabled for northlab.dev. Enable DNSSEC to protect your domain from DNS spoofing and cache poisoning attacks.' }]"
      style="margin-bottom: 1.25rem;"
    />

    <!-- Authentic Cloudflare 36px Toolbar -->
    <div class="cf-dns-toolbar">
      <div class="cf-dns-toolbar-left">
        <div style="width: 260px;">
          <CfInput v-model="searchQuery" placeholder="Search DNS Records" prefix-icon="search" clearable />
        </div>
        <CfButton variant="secondary" size="md">Filters</CfButton>
        <CfButton variant="secondary" size="md">Display options</CfButton>
      </div>
      <div class="cf-dns-toolbar-right">
        <CfButton variant="secondary" size="md" @click="exportBindZone">Import</CfButton>
        <CfButton variant="secondary" size="md" icon="download" @click="exportBindZone">Export</CfButton>
        <CfButton variant="primary" size="md" icon="plus" @click="showAddModal = true">Add record</CfButton>
      </div>
    </div>

    <!-- DNS Records Table Card with Quota Strip -->
    <CfCard body-class="dns-table-body">
      <div class="cf-dns-quota-bar">
        <span>You have used <strong>{{ records.length }}</strong> of <strong>200</strong> available DNS records in this domain.</span>
      </div>

      <CfTable :columns="dnsColumns" :data="filteredRecords" row-key="id">
        <template #col-type="{ val }">
          <span class="dns-type-badge">{{ val }}</span>
        </template>
        <template #col-name="{ val }">
          <span class="cf-fqdn">
            <strong style="color: var(--text-color-kumo-strong);">{{ val === '@' ? zone.name : val }}</strong>
            <span v-if="val !== '@' && !val.includes('.')" style="color: var(--text-color-kumo-subtle);">.{{ zone.name }}</span>
          </span>
        </template>
        <template #col-content="{ val }">
          <code style="font-family: var(--font-mono); font-size: 0.82rem;">{{ val }}</code>
        </template>
        <template #col-proxied="{ row }">
          <span
            class="dns-status-cloud"
            :class="row.proxied ? 'proxied' : 'dns-only'"
            style="cursor: pointer;"
            title="Click to toggle proxy status"
            @click="row.proxied = !row.proxied"
          >
            <span>☁️</span>
            <span>{{ row.proxied ? 'Proxied' : 'DNS only' }}</span>
          </span>
        </template>
        <template #actions="{ row }">
          <div style="display: flex; gap: 0.35rem; justify-content: flex-end;">
            <CfButton variant="secondary" size="sm" @click="alert(`Editing record: ${row.name}`)">Edit</CfButton>
            <CfButton variant="danger-outline" size="sm" icon="trash" icon-only title="Delete record" @click="deleteRecord(row.id)" />
          </div>
        </template>
      </CfTable>

      <template #footer>
        <CfPagination v-model:current-page="currentPage" :total="filteredRecords.length" :page-size="25" />
      </template>
    </CfCard>

    <!-- Add Record Modal Dialog -->
    <CfModal v-model="showAddModal" :title="`Add DNS Record to ${zone.name}`" size="md">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <CfSelect
          v-model="newRecord.type"
          label="Record Type"
          :options="recordTypeOptions"
        />
        <CfInput
          v-model="newRecord.name"
          label="Name (Use @ for root)"
          placeholder="subdomain or @"
        />
        <CfInput
          v-model="newRecord.content"
          label="Target / IPv4 Address"
          placeholder="198.51.100.24"
        />
        <div style="display: flex; flex-direction: column; gap: 0.5rem;">
          <span style="font-size: 13px; font-weight: 500; color: var(--text-color-kumo-default);">Proxy Status</span>
          <CfSwitch v-model="newRecord.proxied" :label="newRecord.proxied ? 'Proxied (Orange Cloud — CDN & DDoS Shield Active)' : 'DNS Only (Bypass Cloudflare Proxy)'" />
        </div>
      </div>
      <template #footer>
        <CfButton variant="outline" @click="showAddModal = false">Cancel</CfButton>
        <CfButton variant="primary" @click="saveNewRecord">Save Record</CfButton>
      </template>
    </CfModal>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'
import { mockData } from '@/data/mockData.js'

const zone = mockData.activeZone
const records = reactive([...zone.dnsRecords])
const searchQuery = ref('')
const showAddModal = ref(false)
const currentPage = ref(1)

const dnsColumns = [
  { key: 'type', label: 'Type', width: '76px' },
  { key: 'name', label: 'Name', width: '175px' },
  { key: 'content', label: 'Content' },
  { key: 'proxied', label: 'Proxy status', width: '125px' },
  { key: 'ttl', label: 'TTL', width: '72px' }
]

const recordTypeOptions = [
  { value: 'A', label: 'A (IPv4 Address)' },
  { value: 'AAAA', label: 'AAAA (IPv6 Address)' },
  { value: 'CNAME', label: 'CNAME (Domain Alias)' },
  { value: 'TXT', label: 'TXT (Text Record)' },
  { value: 'MX', label: 'MX (Mail Exchanger)' }
]

const newRecord = reactive({
  type: 'A',
  name: '',
  content: '',
  proxied: true,
  ttl: 'Auto'
})

const filteredRecords = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return records
  return records.filter(r => r.name.toLowerCase().includes(q) || r.content.toLowerCase().includes(q) || r.type.toLowerCase().includes(q))
})

function saveNewRecord() {
  if (!newRecord.name || !newRecord.content) {
    alert('Please enter record name and content.')
    return
  }
  records.unshift({
    id: 'dns_' + Date.now(),
    type: newRecord.type,
    name: newRecord.name,
    content: newRecord.content,
    proxied: newRecord.proxied,
    ttl: newRecord.ttl,
    modified: 'Just now'
  })
  showAddModal.value = false
  newRecord.name = ''
  newRecord.content = ''
}

function deleteRecord(id) {
  const idx = records.findIndex(r => r.id === id)
  if (idx !== -1) records.splice(idx, 1)
}

function exportBindZone() {
  alert('BIND Zone file generated and downloaded.')
}
</script>

<style scoped>
:deep(.dns-table-body) {
  padding: 0;
}

.cf-view-docs-pill {
  display: inline-flex;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 9999px;
  background: var(--color-kumo-recessed);
  border: 1px solid var(--color-kumo-line);
  font-size: 11.5px;
  font-weight: 500;
  color: var(--text-color-kumo-subtle);
  text-decoration: none;
}

.cf-view-docs-pill:hover {
  color: var(--text-color-kumo-default);
}

.cf-dns-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.cf-dns-toolbar-left,
.cf-dns-toolbar-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.cf-dns-quota-bar {
  padding: 0.65rem 1.25rem;
  font-size: 12.5px;
  color: var(--text-color-kumo-subtle);
  background: var(--color-kumo-canvas);
  border-bottom: 1px solid var(--color-kumo-line);
}
</style>
