<template>
  <div class="kumo-components-showcase">
    <!-- Header -->
    <div class="showcase-header">
      <div>
        <div style="display: flex; align-items: center; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 0.35rem;">
          <h1 class="showcase-title">Cloudflare Kumo Design System</h1>
          <CfBadge variant="brand" size="sm">Commercial Kit</CfBadge>
          <CfBadge variant="success" size="sm" dot>26 Core Components</CfBadge>
          <CfBadge variant="info" size="sm">Vue 3.5</CfBadge>
        </div>
        <p class="showcase-subtitle">
          Exhaustive suite of shared, reusable, 1:1 authentic Cloudflare UI components with full dark mode and accessibility support.
        </p>
      </div>
      <div style="display: flex; gap: 0.5rem;">
        <CfButton variant="outline" icon="refresh" @click="resetStates">Reset</CfButton>
        <CfButton variant="primary" icon="plus" @click="showSampleModal = true">Preview Modal</CfButton>
        <CfButton variant="secondary" icon="layers" @click="showDrawer = true">Open Drawer</CfButton>
      </div>
    </div>

    <!-- Section 1: Buttons -->
    <CfCard title="1. Buttons (CfButton)" subtitle="Primary, secondary, outline, danger, ghost, link with loading states & icons" class="showcase-section">
      <div class="component-row">
        <CfButton variant="primary">Primary Button</CfButton>
        <CfButton variant="primary" icon="server">Deploy Worker</CfButton>
        <CfButton variant="secondary">Secondary</CfButton>
        <CfButton variant="outline">Outline</CfButton>
        <CfButton variant="danger" icon="trash">Delete Object</CfButton>
        <CfButton variant="ghost">Ghost</CfButton>
        <CfButton variant="link">Inline Link</CfButton>
        <CfButton variant="primary" :loading="btnLoading" @click="simulateBtnLoad">Click to Load</CfButton>
        <CfButton variant="secondary" icon="search" icon-only />
      </div>
      <div class="component-row" style="margin-top: 1rem;">
        <CfButton variant="primary" size="sm">Small (28px)</CfButton>
        <CfButton variant="primary" size="md">Medium (36px)</CfButton>
        <CfButton variant="primary" size="lg">Large (42px)</CfButton>
        <CfButton variant="outline" size="sm" icon="clipboard">Copy</CfButton>
        <CfButton variant="secondary" size="sm" disabled>Disabled</CfButton>
      </div>
    </CfCard>

    <!-- Section 2: Badges & Status Pills -->
    <CfCard title="2. Badges & Status Indicators (CfBadge)" subtitle="Colored badges, status dots, and product pills" class="showcase-section">
      <div class="component-row">
        <CfBadge variant="neutral">Default</CfBadge>
        <CfBadge variant="success" dot>Active</CfBadge>
        <CfBadge variant="warning" dot>Warning</CfBadge>
        <CfBadge variant="danger" dot>Exceeded Limit</CfBadge>
        <CfBadge variant="brand">Cloudflare Brand</CfBadge>
        <CfBadge variant="info">DNS Only</CfBadge>
        <CfBadge variant="beta">Beta Feature</CfBadge>
        <CfBadge variant="success" icon="check">Deployed</CfBadge>
      </div>
    </CfCard>

    <!-- Section 3: Stat Cards & Sparklines -->
    <CfCard title="3. Stat Cards & Analytics KPIs (CfStatCard)" subtitle="Cloudflare tabular value display with delta indicators and sparklines" class="showcase-section">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1rem;">
        <CfStatCard 
          label="Total Requests" 
          value="48.2M" 
          change="+14.2%" 
          trend="up" 
          subtext="vs previous 24h"
          tooltip="Total edge requests handled by global Anycast network"
        />
        <CfStatCard 
          label="Edge Bandwidth" 
          value="1.84 TB" 
          change="-3.8%" 
          trend="down" 
          subtext="vs previous 24h"
          tooltip="Egress bytes delivered cached from Cloudflare PoPs"
        />
        <CfStatCard 
          label="Median CPU Time" 
          value="4.2 ms" 
          change="0.0%" 
          trend="neutral" 
          subtext="V8 Isolate runtime"
          tooltip="p50 execution wall-clock time across all workers"
        />
        <CfStatCard 
          label="Threats Mitigated" 
          value="129.4k" 
          change="+28.4%" 
          trend="up" 
          subtext="WAF & DDoS blocks"
          tooltip="Malicious packets stopped before hitting origin servers"
        />
      </div>
    </CfCard>

    <!-- Section 4: Segmented Pill Tabs -->
    <CfCard title="4. Segmented Tabs (CfTabs)" subtitle="Exact 36px Cloudflare pill switcher" class="showcase-section">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <div>
          <span style="font-size: 13px; color: #71717a; margin-bottom: 0.5rem; display: block;">Active Tab: <strong>{{ activeTab }}</strong></span>
          <CfTabs v-model="activeTab" :tabs="['Overview', 'Metrics', 'Deployments', 'Settings']" />
        </div>
        <div>
          <CfTabs v-model="activeSubTab" :tabs="tabObjects" />
        </div>
      </div>
    </CfCard>

    <!-- Section 5: Form Controls & Inputs -->
    <CfCard title="5. Form Inputs & Form Controls (CfInput, CfSelect, CfSwitch, CfCheckbox, CfRadio, CfTextarea)" subtitle="High-contrast inputs, selects, toggles, checkboxes, and radio buttons" class="showcase-section">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1.5rem;">
        <CfInput v-model="inputText" label="Worker Name" placeholder="e.g. my-awesome-worker" prefix-icon="server" clearable hint="Must be unique across your subdomain" />
        <CfInput v-model="searchQuery" label="Search Query" placeholder="Search applications..." prefix-icon="search" clearable />
        <CfSelect v-model="selectedEnv" label="Deployment Environment" :options="envOptions" hint="Target edge routing group" />
        <CfInput v-model="errorInput" label="Custom Domain" placeholder="api.example.com" error="Domain must be delegated to Cloudflare DNS" />
        
        <div style="display: flex; flex-direction: column; gap: 0.75rem;">
          <CfSwitch v-model="switchVal1" label="Enable Live Tail Observability" />
          <CfSwitch v-model="switchVal2" label="Smart Placement (Edge Routing)" />
          <CfSwitch v-model="switchValDisabled" label="Enterprise Mutual TLS (Disabled)" disabled />
        </div>

        <div style="display: flex; flex-direction: column; gap: 0.75rem;">
          <CfCheckbox v-model="check1" label="Enable R2 Object Versioning" description="Retain previous copies of modified objects" />
          <CfCheckbox v-model="check2" label="Auto-purge CDN cache on update" />
        </div>

        <div style="display: flex; flex-direction: column; gap: 0.75rem;">
          <span style="font-size: 13px; font-weight: 500; color: #18181b;">Storage Class:</span>
          <CfRadio v-model="radioClass" value="standard" label="Standard Storage" description="Frequent access with sub-millisecond retrieval" />
          <CfRadio v-model="radioClass" value="infrequent" label="Infrequent Access" description="Cost-optimized cold archival tier" />
        </div>

        <div style="grid-column: 1 / -1;">
          <CfTextarea v-model="sampleTextarea" label="Wrangler Configuration (wrangler.toml)" placeholder="Enter TOML configuration..." :rows="3" hint="Parsed locally before isolate deployment" :maxlength="500" />
        </div>
      </div>
    </CfCard>

    <!-- Section 6: Popovers, Tooltips & Context Menus -->
    <CfCard title="6. Tooltips & Dropdown Menus (CfTooltip, CfDropdownMenu)" subtitle="Hover tooltips with arrow pointer and action context menus" class="showcase-section">
      <div class="component-row" style="gap: 1.5rem;">
        <CfTooltip text="Copy connection string to clipboard" position="top">
          <CfButton variant="outline" icon="clipboard">Hover Me (Top Tooltip)</CfButton>
        </CfTooltip>

        <CfTooltip text="Zone is configured for Cloudflare SSL Universal" position="bottom">
          <CfBadge variant="success" dot>Hover Me (Bottom)</CfBadge>
        </CfTooltip>

        <CfDropdownMenu :items="sampleMenuItems" align="left" @select="handleMenuSelect">
          <template #trigger="{ isOpen }">
            <CfButton variant="secondary" icon="more-horizontal">
              Actions Menu
              <CfIcon name="chevron-down" size="14" />
            </CfButton>
          </template>
        </CfDropdownMenu>

        <CfButton variant="primary" @click="triggerToast">
          Trigger Floating Toast
        </CfButton>
      </div>
    </CfCard>

    <!-- Section 7: Progress Bars & Avatars -->
    <CfCard title="7. Progress Indicators & Avatars (CfProgress, CfAvatar)" subtitle="Linear progress bars with striped animations, user avatars with status indicators" class="showcase-section">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; align-items: center;">
        <div style="display: flex; flex-direction: column; gap: 1rem;">
          <CfProgress :value="72" show-label variant="brand">
            <template #label>R2 Bucket Storage Quota</template>
          </CfProgress>
          <CfProgress :value="45" show-label variant="success" striped animated>
            <template #label>Deployment Propagation (330/330 PoPs)</template>
          </CfProgress>
          <CfProgress :value="88" show-label variant="danger">
            <template #label>Monthly Workers Free Limit</template>
          </CfProgress>
        </div>

        <div style="display: flex; align-items: center; gap: 1rem; flex-wrap: wrap;">
          <CfAvatar name="Sarah Connor" size="xl" status="online" />
          <CfAvatar name="Alex Rivera" size="lg" status="away" />
          <CfAvatar name="Matthew Prince" size="md" status="busy" />
          <CfAvatar name="Dev Ops" size="sm" status="offline" />
          <CfAvatar name="CDN Edge" size="xs" />
        </div>
      </div>
    </CfCard>

    <!-- Section 8: Code Block & Accordion -->
    <CfCard title="8. Code Blocks & Accordions (CfCodeBlock, CfAccordion)" subtitle="Syntax code blocks with 1-click copy and collapsible disclosure panels" class="showcase-section">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <CfCodeBlock 
          :code="sampleCodeSnippet" 
          language="typescript" 
          title="index.ts — Cloudflare Worker Handler" 
        />

        <CfAccordion title="How does Cloudflare Workers V8 Isolate security work?" default-open>
          Cloudflare Workers uses V8 isolates rather than containers or virtual machines. Each isolate provides a secure, lightweight runtime sandbox that starts in under 5 milliseconds with negligible memory overhead compared to traditional serverless architectures.
        </CfAccordion>
        <CfAccordion title="What are the limits of R2 Object Storage?">
          Cloudflare R2 provides zero egress fees across all storage tiers. Individual objects can be up to 5TB in size, with unlimited total storage capacity and automatic geographic replication.
        </CfAccordion>
      </div>
    </CfCard>

    <!-- Section 9: Empty States & Pagination -->
    <CfCard title="9. Empty States & Pagination (CfEmptyState, CfPagination)" subtitle="Cloudflare dashed cloud illustration and segmented page switchers" class="showcase-section">
      <div style="display: flex; flex-direction: column; gap: 1.5rem;">
        <CfEmptyState
          title="No D1 SQL databases found"
          description="Create serverless SQL databases replicated globally at the edge with zero egress fees."
          action-label="Create Database"
          action-icon="plus"
          @action="triggerToast"
        />

        <div style="display: flex; justify-content: flex-end;">
          <CfPagination 
            v-model:current-page="currentPage" 
            :total="128" 
            :page-size="10" 
          />
        </div>
      </div>
    </CfCard>

    <!-- Section 10: Alert Banners -->
    <CfCard title="10. Notification Banners (CfAlert)" subtitle="Used for edge notices, system warnings, and error banners" class="showcase-section">
      <div style="display: flex; flex-direction: column; gap: 0.75rem;">
        <CfAlert variant="info" title="Zero Packet Loss Notice">
          Anycast edge routing is operating with zero packet loss across all 330+ locations globally.
        </CfAlert>
        <CfAlert variant="success" title="Worker Successfully Deployed">
          V8 isolate bundle propagated to 330 PoP edge datacenters in 280ms.
        </CfAlert>
        <CfAlert variant="warning" title="SSL Certificate Revalidation">
          Edge certificate validation is completing in the background. Traffic remains encrypted.
        </CfAlert>
        <CfAlert variant="danger" title="Rate Limit Warning" dismissible>
          Daily free tier quota has reached 85%. Consider upgrading to Workers Paid.
        </CfAlert>
      </div>
    </CfCard>

    <!-- Section 11: Data Table & Shimmer Skeletons -->
    <CfCard title="11. Data Table & Skeleton Loaders (CfTable, CfSkeleton, CfLoader)" subtitle="Production tabular component with selection, custom cells, and loading state" class="showcase-section">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
        <div style="display: flex; align-items: center; gap: 1rem;">
          <span style="font-size: 13px; color: #71717a;">Selected: <strong>{{ selectedTableKeys.length }}</strong> rows</span>
          <CfButton variant="outline" size="sm" @click="tableLoading = !tableLoading">
            Toggle Shimmer Skeletons
          </CfButton>
        </div>
        <div style="display: flex; align-items: center; gap: 0.5rem;">
          <CfLoader size="20" />
          <span style="font-size: 12px; color: #71717a;">Active Edge Stream</span>
        </div>
      </div>

      <CfTable
        :columns="sampleColumns"
        :data="sampleData"
        row-key="id"
        :loading="tableLoading"
        selectable
        :selected-keys="selectedTableKeys"
        @toggle-select="handleRowSelect"
        @toggle-select-all="handleSelectAll"
      >
        <template #col-status="{ val }">
          <CfBadge :variant="val === 'Active' ? 'success' : 'warning'" dot>{{ val }}</CfBadge>
        </template>
        <template #col-domain="{ val }">
          <span class="font-mono text-xs" style="color: #0051c3;">{{ val }}</span>
        </template>
        <template #actions="{ row }">
          <CfButton variant="ghost" size="sm" icon="clipboard" icon-only :title="`Copy ${row.name}`" @click="copyRowName(row.name)" />
        </template>
      </CfTable>
    </CfCard>

    <!-- Sample Interactive Modal -->
    <CfModal v-model="showSampleModal" title="Deploy New Serverless Worker" size="md">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <p style="color: #71717a; font-size: 13px; margin: 0;">
          Create a production-ready Cloudflare Worker running on V8 isolates with zero cold start latency.
        </p>
        <CfInput v-model="modalWorkerName" label="Worker Name" placeholder="e.g. auth-gateway" prefix-icon="server" required />
        <CfSwitch v-model="modalDeployInstant" label="Deploy immediately to production edge routes" />
      </div>
      <template #footer>
        <CfButton variant="outline" @click="showSampleModal = false">Cancel</CfButton>
        <CfButton variant="primary" icon="check" @click="completeModalDeploy">Deploy Application</CfButton>
      </template>
    </CfModal>

    <!-- Sample Slide-Over Drawer -->
    <CfDrawer v-model="showDrawer" title="Application Settings" subtitle="Configure runtime variables and bindings" width="480px">
      <div style="display: flex; flex-direction: column; gap: 1.25rem;">
        <CfInput label="Production Branch" model-value="main" prefix-icon="git-branch" />
        <CfSelect label="Node.js Compatibility Flag" :options="['nodejs_compat_v2', 'nodejs_compat', 'none']" model-value="nodejs_compat_v2" />
        <CfTextarea label="Environment Variables" placeholder="KEY=VALUE" :rows="4" />
        <CfSwitch label="Enable Tail Workers logging" :model-value="true" />
      </div>
      <template #footer="{ close }">
        <CfButton variant="outline" @click="close">Cancel</CfButton>
        <CfButton variant="primary" @click="close">Save Changes</CfButton>
      </template>
    </CfDrawer>

    <!-- Global Toast Container -->
    <CfToast />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'

const toast = useToast()

const showSampleModal = ref(false)
const showDrawer = ref(false)
const btnLoading = ref(false)
const inputText = ref('frontend-auth-worker')
const searchQuery = ref('')
const errorInput = ref('bad-domain')
const selectedEnv = ref('production')
const switchVal1 = ref(true)
const switchVal2 = ref(true)
const switchValDisabled = ref(false)
const check1 = ref(true)
const check2 = ref(false)
const radioClass = ref('standard')
const sampleTextarea = ref('name = "cloudflare-auth-worker"\nmain = "src/index.ts"\ncompatibility_date = "2024-04-01"')
const activeTab = ref('Overview')
const activeSubTab = ref('objects')
const modalWorkerName = ref('api-rate-limiter')
const modalDeployInstant = ref(true)
const tableLoading = ref(false)
const selectedTableKeys = ref([1])
const currentPage = ref(1)

const envOptions = [
  { value: 'production', label: 'Production (Global Anycast)' },
  { value: 'staging', label: 'Staging (Preview Edge)' },
  { value: 'development', label: 'Development (Local Isolated)' }
]

const tabObjects = [
  { label: 'Objects', value: 'objects', icon: 'folder', count: 397 },
  { label: 'Metrics', value: 'metrics', icon: 'bar-chart-2' },
  { label: 'Settings', value: 'settings', icon: 'settings' }
]

const sampleMenuItems = [
  { id: 'view', label: 'View Analytics', icon: 'bar-chart-2' },
  { id: 'duplicate', label: 'Duplicate Worker', icon: 'copy', shortcut: '⌘D' },
  { id: 'export', label: 'Export TOML', icon: 'download' },
  { divider: true },
  { id: 'delete', label: 'Delete Application', icon: 'trash', danger: true }
]

const sampleCodeSnippet = `export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    if (url.pathname === "/health") {
      return new Response(JSON.stringify({ status: "healthy", pop: request.cf?.colo }), {
        headers: { "content-type": "application/json" }
      });
    }
    return new Response("Cloudflare Worker Active");
  }
};`

const sampleColumns = [
  { key: 'name', label: 'Application Name', width: '220px' },
  { key: 'status', label: 'Status', width: '130px' },
  { key: 'domain', label: 'Primary Route', width: '260px' },
  { key: 'requests', label: 'Requests (24h)', align: 'right', width: '140px' },
  { key: 'cpu', label: 'CPU Median', align: 'right', width: '120px' }
]

const sampleData = [
  { id: 1, name: 'code-collection', status: 'Active', domain: 'codecollection.northlab.dev', requests: '12', cpu: '8ms' },
  { id: 2, name: 'auth-gateway', status: 'Active', domain: 'auth.northlab.dev', requests: '4,890', cpu: '3ms' },
  { id: 3, name: 'image-optimizer', status: 'Deploying', domain: 'images.northlab.dev', requests: '210', cpu: '14ms' }
]

const simulateBtnLoad = () => {
  btnLoading.value = true
  setTimeout(() => { btnLoading.value = false }, 1500)
}

const triggerToast = () => {
  toast.success('Configuration saved and distributed across 330 Cloudflare PoPs in 180ms.')
}

const handleMenuSelect = (item) => {
  toast.info(`Selected action: ${item.label}`)
}

const resetStates = () => {
  inputText.value = 'frontend-auth-worker'
  activeTab.value = 'Overview'
  switchVal1.value = true
  switchVal2.value = true
  toast.info('Components showcase reset to defaults.')
}

const handleRowSelect = (id) => {
  const idx = selectedTableKeys.value.indexOf(id)
  if (idx > -1) selectedTableKeys.value.splice(idx, 1)
  else selectedTableKeys.value.push(id)
}

const handleSelectAll = () => {
  if (selectedTableKeys.value.length === sampleData.length) {
    selectedTableKeys.value = []
  } else {
    selectedTableKeys.value = sampleData.map(d => d.id)
  }
}

const copyRowName = (name) => {
  toast.success(`Copied "${name}" to clipboard.`)
}

const completeModalDeploy = () => {
  showSampleModal.value = false
  toast.success(`Worker "${modalWorkerName.value}" deployed successfully.`)
}
</script>

<style scoped>
.kumo-components-showcase {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 0 4rem 0;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
}

.showcase-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--color-kumo-line, #e2e8f0);
}

.showcase-title {
  font-size: 24px;
  font-weight: 600;
  color: #18181b;
  margin: 0;
  letter-spacing: -0.015em;
}

[data-theme="dark"] .showcase-title {
  color: #f8fafc;
}

.showcase-subtitle {
  font-size: 14px;
  color: #71717a;
  margin: 0;
  max-width: 780px;
}

.showcase-section {
  width: 100%;
}

.component-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
}
</style>
